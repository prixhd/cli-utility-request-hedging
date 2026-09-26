package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type result struct {
	url    string
	status string
	header http.Header
	body   []byte
	err    error
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
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- result{
			url: url,
			err: err,
		}
		return
	}

	ch <- result{
		url:    url,
		status: resp.Status,
		header: resp.Header,
		body:   body,
	}
}

func printHelp() {
	fmt.Print(`Использование:
hedgedcurl [options] [urls...]

Опции:
-t, --timeout SECONDS   таймаут запросов (по умолчанию 15 секунд)
-h, --help              показать помощь
`)
}

func printRes(res result) {
	fmt.Println("Первым ответил:", res.url)
	fmt.Println("Status:", res.status)

	for key, values := range res.header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

	fmt.Println()

	_, err := os.Stdout.Write(res.body)
	if err != nil {
		fmt.Println("Ошибка при выводе body:", err)
	}
}

func main() {

	var timeout int
	var help bool

	flag.IntVar(&timeout, "t", 15, "таймаут в секундах")
	flag.IntVar(&timeout, "timeout", 15, "таймаут в секундах")

	flag.BoolVar(&help, "h", false, "показать помощь")
	flag.BoolVar(&help, "help", false, "показать помощь")

	flag.Parse()

	if help {
		printHelp()
		return
	}

	if timeout <= 0 {
		fmt.Println("Ошибка! Таймаут должен быть больше 0 секунд")
		os.Exit(1)
	}

	urls := flag.Args()

	if len(urls) == 0 {
		fmt.Println("Ошибка! Укажите хотя бы один URL")
		os.Exit(1)
	}

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

			cancel()
			printRes(currentResult)
			return

		case <-ctx.Done():
			// После того как сработал таймаут, мы проверяем полностью канал на результат - был ли там поставлен успешный ответ в канал перед таймаутом или нет
			for {
				select {
				case currentResult := <-ch:
					if currentResult.err == nil {
						cancel()
						printRes(currentResult)
						return
					}
				default:
					fmt.Println("Время закончилось")
					cancel()
					os.Exit(228)
				}
			}

		}
	}

	fmt.Println("Все запросы завершились c ошибкой")
	cancel()
	os.Exit(1)
}
