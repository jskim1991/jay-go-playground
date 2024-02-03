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

func (p *processorV2) launchV2(ctx context.Context, data requestV2) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		result, err := getResultAV2(ctx, data.A)
		wg.Done()
		if err != nil {
			p.errs <- err
			return
		}
		p.outA <- result
	}()

	go func() {
		result, err := getResultBV2(ctx, data.B)
		wg.Done()
		if err != nil {
			p.errs <- err
			return
		}
		p.outB <- result
	}()

	wg.Wait()
}

func (p *processorV2) waitForABV2(ctx context.Context) (aggregatedResultV2, error) {
	var input aggregatedResultV2
	count := 0
	for count < 2 {
		select {
		case a := <-p.outA:
			input.A = a
			count++
		case b := <-p.outB:
			input.B = b
			count++
		case err := <-p.errs:
			return aggregatedResultV2{}, err
		case <-ctx.Done():
			return aggregatedResultV2{}, ctx.Err()
		}
	}
	return input, nil
}

func getResultAV2(ctx context.Context, input string) (string, error) {
	fmt.Println("Invoking A")
	time.Sleep(1 * time.Second)
	return "results from A", nil
}

func getResultBV2(ctx context.Context, input string) (string, error) {
	fmt.Println("Invoking B")
	time.Sleep(2 * time.Second)
	return "results from B", nil
}

func getResultCV2(ctx context.Context, input aggregatedResultV2) (string, error) {
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
	p.launchV2(ctx, data)

	ab, err := p.waitForABV2(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Gathered data from A and B:", ab)

	c, err := getResultCV2(ctx, ab)
	if err != nil {
		p.errs <- err
		return
	}

	fmt.Println("Request handled successfully:", c)
}
