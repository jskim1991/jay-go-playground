package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var urls = []string{
	"https://www.google.com",
	"https://naver.com",
}

func fetchStatus(writer http.ResponseWriter, request *http.Request) {
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			resp, err := http.Get(url)
			if err != nil {
				fmt.Fprintf(writer, "%+v\n", err)
			}
			fmt.Fprintf(writer, "%s: %+v\n", url, resp.Status)
			wg.Done()
		}(url)
	}
	wg.Wait()
}

func main() {
	fmt.Println("Go WaitGroup tutorial")

	http.HandleFunc("/", fetchStatus)
	log.Fatal(http.ListenAndServe(":8080", nil))

	fmt.Println("main finished")
}
