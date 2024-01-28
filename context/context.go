package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

const RequestId = "request_id"

func enrichContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, RequestId, rand.Int())
}

func doSomething(ctx context.Context) {
	rID := ctx.Value(RequestId)
	fmt.Println("Request ID: ", rID)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("timed out")
			return
		default:
			fmt.Println("processing...")
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	fmt.Println("Go context tutorial")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	ctx = enrichContext(ctx)

	go doSomething(ctx)

	select {
	case <-ctx.Done():
		fmt.Println(fmt.Sprintf("main(): timed out: %s", ctx.Err()))
	}
	time.Sleep(2 * time.Second)
}
