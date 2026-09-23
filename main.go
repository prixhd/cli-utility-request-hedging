package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"
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

	var timeout int

	flag.IntVar(&timeout, "t", 15, "таймаут в секундах")
	flag.IntVar(&timeout, "timeout", 15, "таймаут в секундах")

	flag.Parse()

	urls := flag.Args()

	ch := make(chan result, len(urls))

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(timeout)*time.Second,
	)
	defer cancel()

	for _, url := range urls {
		go getStatusToUrl(url, ch, ctx)
	}

	for i := 0; i < len(urls); i++ {
		select {
		case currentResult := <-ch:
			if currentResult.err != nil {
				fmt.Println("Ошибка: ", currentResult.err)
				continue
			}

			fmt.Println("Первым ответил: ", currentResult.url)
			fmt.Println("Status: ", currentResult.resp.Status)

			for key, values := range currentResult.resp.Header {
				for _, value := range values {
					fmt.Printf("%s: %s\n", key, value)
				}
			}

			fmt.Println()

			_, err := io.Copy(os.Stdout, currentResult.resp.Body)
			if err != nil {
				fmt.Println("Ошибка при чтении body: ", err)
			}

			cancel()
			currentResult.resp.Body.Close()
			return

		case <-ctx.Done():
			fmt.Println("Время закончилось")
			return
		}
	}

	fmt.Println("Все запросы завершились ошибкой")
}
