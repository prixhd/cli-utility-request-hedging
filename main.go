package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

type result struct {
	url  string
	resp *http.Response
	err  error
}

func getStatusToUrl(url string, ch chan result, ctx context.Context) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		ch <- result{
			url: url,
			err: err,
		}
		return
	}

	resp, err := http.DefaultClient.Do(req)
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	urls := os.Args[1:]

	for _, url := range urls {
		go getStatusToUrl(url, ch, ctx)
	}

	for i := 0; i < len(urls); i++ {
		currentResult := <-ch

		if currentResult.err != nil {
			fmt.Println("Ошибка: ", currentResult.err)
			continue
		}

		fmt.Println("Первым ответил: ", currentResult.url)
		fmt.Println("Status: ", currentResult.resp.Status)

		cancel()

		currentResult.resp.Body.Close()
		return
	}

	fmt.Println("Все запросы завершились ошибкой")
}
