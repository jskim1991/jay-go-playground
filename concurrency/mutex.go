package main

import (
	"fmt"
	"sync"
)

var (
	mutex   sync.Mutex
	balance int
)

func deposit(value int, wg *sync.WaitGroup) {
	mutex.Lock()
	fmt.Println("Depositing", value)
	balance += value
	fmt.Printf("New balance %d\n\n", balance)
	mutex.Unlock()
	wg.Done()
}

func withdraw(value int, wg *sync.WaitGroup) {
	mutex.Lock()
	fmt.Println("Withdrawing", value)
	balance -= value
	fmt.Printf("New balance %d\n\n", balance)
	mutex.Unlock()
	wg.Done()
}

func main() {
	fmt.Println("Go Mutex tutorial")

	wg := sync.WaitGroup{}
	wg.Add(2)

	balance = 1000
	fmt.Printf("Starting balance %d\n\n", balance)
	go withdraw(750, &wg)
	go deposit(500, &wg)

	wg.Wait()
}
