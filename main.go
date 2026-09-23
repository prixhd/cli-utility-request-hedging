package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
)

func getStatusToUrl(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(url)

	// Обработка ошибки
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("URL: ", url)
	fmt.Println("Status: ", resp.Status)

	for key, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

}

func main() {
	var wg sync.WaitGroup
	urls := os.Args[1:]
	wg.Add(len(urls))

	for _, url := range urls {
		go getStatusToUrl(url, &wg)

	}

	wg.Wait()
}
