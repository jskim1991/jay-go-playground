package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type processorV2 struct {
	outA chan string
	outB chan string
	outC chan string
	inC  chan aggregatedResultV2
	errs chan error
}

type requestV2 struct {
	A string
	B string
}

type aggregatedResultV2 struct {
	A string
	B string
}

func v2GetResultA(ctx context.Context, input string) (string, error) {
	fmt.Println("Invoking A")
	time.Sleep(1 * time.Second)
	return "results from A", nil
}

func v2GetResultB(ctx context.Context, input string) (string, error) {
	fmt.Println("Invoking B")
	time.Sleep(2 * time.Second)
	return "results from B", nil
}

func v2GetResultC(ctx context.Context, input aggregatedResultV2) (string, error) {
	fmt.Println("Invoking C")
	return "results from C", nil
}

func main() {
	p := processorV2{
		outA: make(chan string),
		outB: make(chan string),
		outC: make(chan string),
		inC:  make(chan aggregatedResultV2),
		errs: make(chan error),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data := requestV2{
		A: "inputA",
		B: "inputB",
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		result, err := v2GetResultA(ctx, data.A)
		wg.Done()
		if err != nil {
			p.errs <- err
			return
		}
		p.outA <- result
	}()

	go func() {
		result, err := v2GetResultB(ctx, data.B)
		wg.Done()
		if err != nil {
			p.errs <- err
			return
		}
		p.outB <- result
	}()

	wg.Wait()

	ab := aggregatedResultV2{
		A: <-p.outA,
		B: <-p.outB,
	}
	fmt.Println("Gathered data from A and B:", ab)

	c, err := v2GetResultC(ctx, ab)
	if err != nil {
		p.errs <- err
		return
	}

	fmt.Println("Request handled successfully:", c)
}
