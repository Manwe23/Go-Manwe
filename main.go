// main
package main

import (
	"errors"
	"fmt"
	"time"
)

type queue_item struct {
	val  int
	tail *queue_item
}

type queue struct {
	head *queue_item
	size int
}

var mutex bool = false

func fib(n int) int {
	if n < 2 {
		return n
	}

	return fib(n-1) + fib(n-2)
}

func golibroda(c <-chan bool) {
	for {
		fmt.Println("Golibroda rozpoczyna prace...")
		time.Sleep(5 * time.Second)
		fmt.Println("Golibroda konczy prace.")
	}

}

func (q queue) is_empty() bool {
	return q.size == 0
}

func (q *queue) pop() (int, error) {
	if q.is_empty() {
		return 0, errors.New("queue is empty !")
	}
	var tmp int = q.head.val
	q.size--
	q.head = q.head.tail
	return tmp, nil
}

func (q *queue) push(item int) {
	var i *queue_item
	i = new(queue_item)
	i.val = item
	i.tail = q.head
	q.head = i
	q.size++
}

func main() {
	fmt.Println("Hello World!")

	var i int
	var q queue //= &queue{}
	q.size = 0
	q.push(7)
	q.push(18)
	q.push(29)
	if x, err := q.pop(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(x)
	}
	fmt.Println(q.pop())
	c := make(chan bool)
	go golibroda(c)
	fmt.Scanf("%d", &i)
	//fmt.Println(fib(30))
}
