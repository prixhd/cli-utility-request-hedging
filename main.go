package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	urls := os.Args[1:]

	for _, url := range urls {

		resp, err := http.Get(url)

		// Обработка ошибки
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("Status:", resp.Status)

		for key, values := range resp.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", key, value)
			}
		}

		resp.Body.Close()
	}
}
