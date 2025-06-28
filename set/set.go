package set

import (
	"sync"

	"github.com/nikolaihg/deadlink-scraper-go/linktype"
)

type Set struct {
	mu   sync.RWMutex
	data map[string]linktype.Link
}

func New() *Set {
	return &Set{
		data: make(map[string]linktype.Link),
	}
}

func (s *Set) Add(link linktype.Link) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[link.URL] = link
}

func (s *Set) Contains(link linktype.Link) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.data[link.URL]
	return exists
}

func (s *Set) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

func (s *Set) Values() []linktype.Link {
	s.mu.RLock()
	defer s.mu.RUnlock()
	values := make([]linktype.Link, 0, len(s.data))
	for _, v := range s.data {
		values = append(values, v)
	}
	return values
}
