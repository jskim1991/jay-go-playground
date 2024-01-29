package main

import (
	"errors"
	"fmt"
)

type Queue struct {
	items []int
}

func (q *Queue) Enqueue(item int) {
	q.items = append(q.items, item)
}

func (q *Queue) Dequeue() (int, error) {
	if len(q.items) == 0 {
		return 0, errors.New("empty queue")
	}

	item := q.items[0]
	q.items = q.items[1:]

	return item, nil
}

func main() {
	fmt.Println("Queue with slices")

	q := Queue{}
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	fmt.Println(q)

	for i := 0; i < 4; i++ {
		item, err := q.Dequeue()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(item)
		}
	}
}
