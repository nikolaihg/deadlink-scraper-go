package set

import "github.com/nikolaihg/deadlink-scraper-go/linktype"

type Set struct {
	data map[string]linktype.Link
}

func New() *Set {
	return &Set{
		data: make(map[string]linktype.Link),
	}
}

func (s *Set) Add(link linktype.Link) {
	s.data[link.URL] = link
}

func (s *Set) Contains(link linktype.Link) bool {
	_, exists := s.data[link.URL]
	return exists
}

func (s *Set) Size() int {
	return len(s.data)
}

func (s *Set) Values() []linktype.Link {
	values := make([]linktype.Link, 0, len(s.data))
	for _, v := range s.data {
		values = append(values, v)
	}
	return values
}
