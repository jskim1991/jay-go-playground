package main

import (
	"context"
	"fmt"
	"time"
)

type processor struct {
	outA chan string
	outB chan string
	outC chan string
	inC  chan aggregatedResult
	errs chan error
}

type request struct {
	A string
	B string
}

type aggregatedResult struct {
	A string
	B string
}

func (p *processor) launch(ctx context.Context, data request) {
	go func() {
		result, err := getResultA(ctx, data.A)
		if err != nil {
			p.errs <- err
			return
		}
		p.outA <- result
	}()

	go func() {
		result, err := getResultB(ctx, data.B)
		if err != nil {
			p.errs <- err
			return
		}
		p.outB <- result
	}()

	go func() {
		select {
		case <-ctx.Done():
			return
		case inputC := <-p.inC:
			result, err := getResultC(ctx, inputC)
			if err != nil {
				p.errs <- err
				return
			}
			p.outC <- result
		}
	}()
}

func (p *processor) waitForAB(ctx context.Context) (aggregatedResult, error) {
	var input aggregatedResult
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
			return aggregatedResult{}, err
		case <-ctx.Done():
			return aggregatedResult{}, ctx.Err()
		}
	}
	return input, nil
}

func (p *processor) waitForC(ctx context.Context) (string, error) {
	select {
	case out := <-p.outC:
		return out, nil
	case err := <-p.errs:
		return "", err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func getResultA(ctx context.Context, input string) (string, error) {
	fmt.Println("Invoking A")
	time.Sleep(1 * time.Second)
	return "results from A", nil
}

func getResultB(ctx context.Context, input string) (string, error) {
	fmt.Println("Invoking B")
	time.Sleep(2 * time.Second)
	return "results from B", nil
}

func getResultC(ctx context.Context, input aggregatedResult) (string, error) {
	fmt.Println("Invoking C")
	return "results from C", nil
}

func main() {
	p := processor{
		outA: make(chan string),
		outB: make(chan string),
		outC: make(chan string),
		inC:  make(chan aggregatedResult),
		errs: make(chan error),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data := request{
		A: "inputA",
		B: "inputB",
	}
	p.launch(ctx, data)

	ab, err := p.waitForAB(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Gathered data from A and B:", ab)

	p.inC <- ab
	c, err := p.waitForC(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Request handled successfully:", c)
}
