package firecore

import (
	"fmt"
	"github.com/streamingfast/logging"
	"go.uber.org/zap"
	"time"

	"github.com/spf13/cobra"
	"github.com/streamingfast/bstream"
	"github.com/streamingfast/bstream/transform"
	"github.com/streamingfast/dstore"
	pbbstream "github.com/streamingfast/pbgo/sf/bstream/v1"
	"google.golang.org/protobuf/proto"
)

var UnsafePayloadKind pbbstream.Protocol = pbbstream.Protocol_UNKNOWN

// Block represents the chain-specific Protobuf block. Chain specific's block
// model must implement this interface so that Firehose core is able to properly
// marshal/unmarshal your block into/to the Firehose block envelope binary format.
//
// All the methods are prefixed with `GetFirehoseBlock` to avoid any potential
// conflicts with the fields/getters of your chain's block model that would
// prevent you from implementing this interface.
//
// Consumer of your chain's protobuf block model don't need to be aware of those
// details, they are internal Firehose core information that are required to function
// properly.
//
// The value you return for each of those methods must be done respecting Firehose rules
// which are enumarated in the documentation of each method.
type Block interface {
	proto.Message

	// GetFirehoseBlockID returns the block ID as a string, usually in the representation
	// used by your chain (hex, base58, base64, etc.). The block ID must be unique across
	// all blocks that will ever exist on your chain.
	GetFirehoseBlockID() string

	// GetFirehoseBlockNumber returns the block number as an unsigned integer. The block
	// number could be shared by multiple blocks in which case one is the canonical one
	// and the others are forks (resolution of forks is handled by Firehose core later in the
	// block processing pipeline).
	//
	// The value should be sequentially ordered which means that a block with block number 10
	// has come before block 11. Firehose core will deal with block skips without problem though
	// (e.g. block 1, is produced then block 3 where block 3's parent is block 1).
	GetFirehoseBlockNumber() uint64

	// GetFirehoseBlockPreviousID returns the block ID of the parent block as a string. All blocks
	// ever produced must have a parent block ID except for the genesis block which is the first
	// one. The value must be the same as the one returned by GetFirehoseBlockID() of the parent.
	//
	// If it's the genesis block, return an empty string.
	GetFirehoseBlockParentID() string

	// GetFirehoseBlockTime returns the block timestamp as a time.Time of when the block was
	// produced. This should the consensus agreed time of the block.
	GetFirehoseBlockTime() time.Time

	// GetFirehoseBlockLIBNum returns the last irreversible block number as an unsigned integer
	// of this block. This is one of the most important piece of information for Firehose core.
	// as it determines when "forks" are now stalled and should be removed from memory and it
	// drives a bunch of important write processes that will write the block to disk only when the
	// block is now irreversible.
	//
	// The value returned should be the oldest block that should turned to be irreversible when this
	// block was produced. Assume for example the current block is 100. If finality rule of a chain
	// is that a block become irreversible after 12 blocks has been produced, then the value returned
	// in this case should be 88 (100 - 12) which means that when block 100 was produced, block 88
	// can now be considered irreversible.
	//
	// Irreversibility is chain specific and how the value here is returned depends on the chain. On
	// probabilistic irreversible chains, like Bitcoin, the value returned here is usually the current
	// block number - <threshold> where <threshold> is choosen to be safe enough in all situations (ensure
	// that is block number < <threshold>, then you properly cap to 0).
	//
	// On deterministic irreversible chains, usually the last irreversible block number if part of the
	// consensus and as such should be part of the Protobuf block model somewhere. In those cases, this
	// value should be returned here.
	GetFirehoseBlockLIBNum() uint64
}

// BlockEncoder is the interface of an object that is going to a chain specific
// block implementing [Block] interface that will be encoded into [bstream.Block]
// type which is the type used by Firehose core to "envelope" the block.
type BlockEncoder interface {
	Encode(block Block) (blk *bstream.Block, err error)
}

type BlockEncoderFunc func(block Block) (blk *bstream.Block, err error)

func (f BlockEncoderFunc) Encode(block Block) (blk *bstream.Block, err error) {
	return f(block)
}

type CommandExecutor func(cmd *cobra.Command, args []string) (err error)

func NewGenericBlockEncoder(protocolVersion int32) BlockEncoder {
	return BlockEncoderFunc(func(block Block) (blk *bstream.Block, err error) {
		return EncodeBlock(protocolVersion, block)
	})
}

var zlog, _ = logging.PackageLogger("types", "github.com/streamingfast/firehose-core/types")

func EncodeBlock(protocolVersion int32, b Block) (blk *bstream.Block, err error) {
	content, err := proto.Marshal(b)
	if err != nil {
		zlog.Error("failed to marshal to binary", zap.Error(err), zap.Any("block", b))
		return nil, fmt.Errorf("unable to marshal to binary form: %s", err)
	}

	bstreamBlock := &bstream.Block{
		Id:             b.GetFirehoseBlockID(),
		Number:         b.GetFirehoseBlockNumber(),
		PreviousId:     b.GetFirehoseBlockParentID(),
		Timestamp:      b.GetFirehoseBlockTime(),
		LibNum:         b.GetFirehoseBlockLIBNum(),
		PayloadVersion: protocolVersion,

		// PayloadKind is not actually used anymore and should be left to UNKNOWN
		PayloadKind: UnsafePayloadKind,
	}

	return bstream.GetBlockPayloadSetter(bstreamBlock, content)
}

type BlockIndexerFactory[B Block] func(indexStore dstore.Store, indexSize uint64) (BlockIndexer[B], error)

type BlockIndexer[B Block] interface {
	ProcessBlock(block B) error
}

// BlockTransformerFactory is a bit convoluted, but yes it's a function acting as a factory that returns itself
// a factory. The reason for this is that the factory needs to be able to access the index store and the index
// size to be able to create the actual factory.
//
// In the context of `firehose-core` transform registration, this function will be called exactly once
// for the overall process. The returns [transform.Factory] will be used multiple times (one per request
// requesting this transform).
type BlockTransformerFactory func(indexStore dstore.Store, indexPossibleSizes []uint64) (*transform.Factory, error)
