package blockpoller

import (
	pbbstream "github.com/streamingfast/bstream/pb/sf/bstream/v1"
)

type BlockFetcher[C any] interface {
	IsBlockAvailable(requestedSlot uint64) bool
	Fetch(client C, blkNum uint64) (b *pbbstream.Block, skipped bool, err error)
}
