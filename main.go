// main
package main

import (
	"errors"
	"fmt"
	"math/rand"
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

var mutex bool = false
var sleep = false

func fib(n int) int {
	if n < 2 {
		return n
	}

	return fib(n-1) + fib(n-2)
}

func golibroda(c <-chan bool, q *queue) {
	for {
		if x, err := q.pop(); err != nil {
			fmt.Println("Golibroda obsługuje klienta nr:", x)
			time.Sleep(5 * time.Second)
			fmt.Println("Golibroda konczy prace.")
		} else {
			break
		}
	}
	sleep = true
	fmt.Println("Golibroda idzie spać. Zzzz...")
	for {
		if !sleep {
			break
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Println("Golibroda budzi sie!")
	golibroda(c, q)
}

func poczekalnia(c <-chan bool, q *queue) {
	w8_time_generator := rand.New(rand.NewSource(99))
	var w8_time int = w8_time_generator.Intn(10)
	var client_nr int = 0
	for {
		client_nr++
		fmt.Println("Przyszedl klient nr:", client_nr)
		fmt.Println("Nastepny klient za:", w8_time)
		time.Sleep(time.Duration(w8_time) * time.Second)
		q.push(client_nr)
		if sleep {
			sleep = false
		}
		w8_time = w8_time_generator.Intn(10)
	}
}

func main() {
	fmt.Println("Witaj swiat!")

	var i int
	var q queue //= &queue{}
	q.size = 0
	c := make(chan bool)
	c2 := make(chan bool)
	go golibroda(c, &q)
	go poczekalnia(c2, &q)
	fmt.Scanf("%d", &i)
	//fmt.Println(fib(30))
}
