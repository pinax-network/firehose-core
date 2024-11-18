package rpc

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/go-multierror"
)

var ErrorNoMoreClient = errors.New("no more clients")

type Clients[C any] struct {
	clients               []C
	maxBlockFetchDuration time.Duration
	rollingStrategy       RollingStrategy[C]
}

func NewClients[C any](maxBlockFetchDuration time.Duration, rollingStrategy RollingStrategy[C]) *Clients[C] {
	return &Clients[C]{
		maxBlockFetchDuration: maxBlockFetchDuration,
		rollingStrategy:       rollingStrategy,
	}
}

func (c *Clients[C]) Add(client C) {
	c.clients = append(c.clients, client)
}

func WithClients[C any, V any](clients *Clients[C], f func(context.Context, C) (v V, err error)) (v V, err error) {
	var errs error

	clients.rollingStrategy.reset()
	client, err := clients.rollingStrategy.next(clients)
	if err != nil {
		errs = multierror.Append(errs, err)
		return v, errs
	}

	for {

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, clients.maxBlockFetchDuration)

		v, err := f(ctx, client)
		cancel()

		if err != nil {
			errs = multierror.Append(errs, err)
			client, err = clients.rollingStrategy.next(clients)
			if err != nil {
				errs = multierror.Append(errs, err)
				return v, errs
			}

			continue
		}
		return v, nil
	}
}

type RollingStrategy[C any] interface {
	reset()
	next(clients *Clients[C]) (C, error)
}

type RollingStrategyRoundRobin[C any] struct {
	fistCallToNewClient bool
	usedClientCount     int
	nextClientIndex     int
}

func NewRollingStrategyRoundRobin[C any]() RollingStrategy[C] {
	return &RollingStrategyRoundRobin[C]{
		fistCallToNewClient: true,
	}
}

func (s *RollingStrategyRoundRobin[C]) reset() {
	s.usedClientCount = 0
}
func (s *RollingStrategyRoundRobin[C]) next(clients *Clients[C]) (client C, err error) {
	if len(clients.clients) == s.usedClientCount {
		return client, ErrorNoMoreClient
	}

	if s.fistCallToNewClient {
		s.fistCallToNewClient = false
		client = clients.clients[0]
		s.usedClientCount = s.usedClientCount + 1
		s.nextClientIndex = s.nextClientIndex + 1
		return client, nil
	}

	if s.nextClientIndex == len(clients.clients) { //roll to 1st client
		s.nextClientIndex = 0
	}

	if s.usedClientCount == 0 { //just been reset
		s.nextClientIndex = s.prevIndex(clients)
		client = clients.clients[s.nextClientIndex]
		s.usedClientCount = s.usedClientCount + 1
		s.nextClientIndex = s.nextClientIndex + 1
		return client, nil
	}

	if s.nextClientIndex == len(clients.clients) { //roll to 1st client
		client = clients.clients[0]
		s.usedClientCount = s.usedClientCount + 1
		return client, nil
	}

	client = clients.clients[s.nextClientIndex]
	s.usedClientCount = s.usedClientCount + 1
	s.nextClientIndex = s.nextClientIndex + 1
	return client, nil
}

func (s *RollingStrategyRoundRobin[C]) prevIndex(clients *Clients[C]) int {
	if s.nextClientIndex == 0 {
		return len(clients.clients) - 1
	}
	return s.nextClientIndex - 1
}

type RollingStrategyAlwaysUseFirst[C any] struct {
	nextIndex int
}

func NewRollingStrategyAlwaysUseFirst[C any]() *RollingStrategyAlwaysUseFirst[C] {
	return &RollingStrategyAlwaysUseFirst[C]{}
}

func (s *RollingStrategyAlwaysUseFirst[C]) reset() {
	s.nextIndex = 0
}

func (s *RollingStrategyAlwaysUseFirst[C]) next(c *Clients[C]) (client C, err error) {
	if len(c.clients) <= s.nextIndex {
		return client, ErrorNoMoreClient
	}
	client = c.clients[s.nextIndex]
	s.nextIndex++
	return client, nil

}
