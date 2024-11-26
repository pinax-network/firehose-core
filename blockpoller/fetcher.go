package blockpoller

import (
	"context"

	pbbstream "github.com/streamingfast/bstream/pb/sf/bstream/v1"
)

type BlockFetcher[C any] interface {
	IsBlockAvailable(requestedSlot uint64) bool
	Fetch(ctx context.Context, client C, blkNum uint64) (b *pbbstream.Block, skipped bool, err error)
}

type HeadBlockNumberFetcher[C any] interface {
	FetchHeadBlockNumber(ctx context.Context, client C) (uint64, error)
}
