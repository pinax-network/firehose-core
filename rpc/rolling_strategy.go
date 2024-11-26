package rpc

type RollingStrategy[C any] interface {
	reset()
	next(clients *Clients[C]) (C, error)
}

type StickyRollingStrategy[C any] struct {
	fistCallToNewClient bool
	usedClientCount     int
	nextClientIndex     int
}

func NewStickyRollingStrategy[C any]() *StickyRollingStrategy[C] {
	return &StickyRollingStrategy[C]{
		fistCallToNewClient: true,
	}
}

func (s *StickyRollingStrategy[C]) reset() {
	s.usedClientCount = 0
}
func (s *StickyRollingStrategy[C]) next(clients *Clients[C]) (client C, err error) {
	clients.lock.Lock()
	defer clients.lock.Unlock()

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

func (s *StickyRollingStrategy[C]) prevIndex(clients *Clients[C]) int {
	clients.lock.Lock()
	defer clients.lock.Unlock()

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
	c.lock.Lock()
	defer c.lock.Unlock()

	if len(c.clients) <= s.nextIndex {
		return client, ErrorNoMoreClient
	}
	client = c.clients[s.nextIndex]
	s.nextIndex++
	return client, nil

}
