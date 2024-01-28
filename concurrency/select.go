package main

import (
	"fmt"
	"time"
)

func fibonacci(c chan int, quit chan int) {
	x := 0
	y := 1

	for {
		select {
		case c <- x:
			fmt.Println("updating")
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func main() {
	c := make(chan int)
	quit := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
			time.Sleep(200 * time.Millisecond)
		}
		quit <- 0
	}()

	fibonacci(c, quit)
	fmt.Println("main finished")
}
