package main

import (
	"fmt"
	"time"
)

func compute(value int, c chan bool) {
	for i := 0; i < value; i++ {
		time.Sleep(time.Second)
		fmt.Println(i)
	}
	fmt.Println("compute finished")
	c <- true
}

func main() {
	fmt.Println("Go concurrency tutorial")

	c := make(chan bool)

	go compute(3, c)
	go compute(3, c)

	fmt.Println("main finished")

	// main function completes before goroutine completes
	// any goroutines that have yet to complete by the end of main function call are terminated
	<-c
	<-c
}
