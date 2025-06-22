package queue

import "github.com/nikolaihg/deadlink-scraper-go/linktype"

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
