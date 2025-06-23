package queue

import (
	"fmt"
	"strings"

	"github.com/nikolaihg/deadlink-scraper-go/linktype"
)

type Queue struct {
	totalQueued int
	number      int
	elements    []linktype.Link
}

func New() *Queue {
	return &Queue{
		elements: make([]linktype.Link, 0),
	}
}

func (q *Queue) Enqueue(link linktype.Link) {
	q.elements = append(q.elements, link)
	q.totalQueued++
	q.number++
}

func (q *Queue) Dequeue() linktype.Link {
	link := q.elements[0]
	q.elements = q.elements[1:]
	q.number--
	return link
}

func (q *Queue) Size() int {
	return q.number
}

func (q *Queue) IsEmpty() bool {
	return q.Size() == 0
}

func (q *Queue) Print() {
	fmt.Println(q.String())
}

func (q *Queue) String() string {
	var sb strings.Builder
	sb.WriteString("Queue contents:\n")
	for i, link := range q.elements {
		sb.WriteString(fmt.Sprintf("%d: URL: %s, Type: %d\n", i+1, link.URL, link.Type))
	}
	return sb.String()
}
