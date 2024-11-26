package rpc

import (
	"context"
	"sort"
)

type SortValueFetcher interface {
	fetchSortValue(ctx context.Context) (sortValue uint64, err error)
}

type SortDirection int

const (
	SortDirectionAscending SortDirection = iota
	SortDirectionDescending
)

func Sort[C any](ctx context.Context, clients *Clients[C], direction SortDirection) error {
	type sortable struct {
		clientIndex int
		sortValue   uint64
	}
	var sortableValues []sortable
	for i, client := range clients.clients {
		var v uint64
		var err error
		if s, ok := any(client).(SortValueFetcher); ok {
			v, err = s.fetchSortValue(ctx)
			if err != nil {
				//do nothing
			}
		}
		sortableValues = append(sortableValues, sortable{i, v})
	}

	sort.Slice(sortableValues, func(i, j int) bool {
		if direction == SortDirectionAscending {
			return sortableValues[i].sortValue < sortableValues[j].sortValue
		}
		return sortableValues[i].sortValue > sortableValues[j].sortValue
	})

	var sorted []C
	for _, v := range sortableValues {
		sorted = append(sorted, clients.clients[v.clientIndex])
	}

	clients.lock.Lock()
	defer clients.lock.Unlock()
	clients.clients = sorted

	return nil
}
