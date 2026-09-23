package main

import (
	"fmt"
	"net/http"
	"os"
)

func getStatusToUrl(url string, ch chan string) {

	resp, err := http.Get(url)

	// Обработка ошибки
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	ch <- url

}

func main() {
	ch := make(chan string)

	urls := os.Args[1:]

	for _, url := range urls {
		go getStatusToUrl(url, ch)

	}
	firstURL := <-ch
	fmt.Println("Первым ответил: ", firstURL)
}
