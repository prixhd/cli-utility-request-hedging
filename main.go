package main

import (
	"fmt"
	"net/http"
	"os"
)

type result struct {
	url  string
	resp *http.Response
	err  error
}

func getStatusToUrl(url string, ch chan result) {

	resp, err := http.Get(url)

	// Обработка ошибки
	if err != nil {
		ch <- result{
			url: url,
			err: err,
		}
		return
	}

	ch <- result{
		url:  url,
		resp: resp,
	}

}

func main() {
	ch := make(chan result)

	urls := os.Args[1:]

	for _, url := range urls {
		go getStatusToUrl(url, ch)
	}

	for i := 0; i < len(urls); i++ {
		result := <-ch

		if result.err != nil {
			fmt.Println("Ошибка: ", result.err)
			continue
		}

		fmt.Println("Первым ответил: ", result.url)
		fmt.Println("Status: ", result.resp.Status)

		result.resp.Body.Close()
		return
	}

}
