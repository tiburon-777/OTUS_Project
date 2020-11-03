package main

import (
	"context"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestIntegrationPositive(t *testing.T) {
	wg := sync.WaitGroup{}
	server := &http.Server{Addr: ":3000", Handler: http.FileServer(http.Dir("../examples"))}
	go server.ListenAndServe()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	go func(ctx context.Context) {
		main()
	}(ctx)

	// Реализовать тесты логики приложения (ресайзы по разным требованиям):
	wg.Add(2)
	t.Run("remote server return jpeg", func(t *testing.T) {
		defer wg.Done()
		body, _, err := request("http://localhost/fill/1024/504/localhost:3000/gopher_original_1024x504.jpg", 15*time.Second)
		require.NoError(t, err)
		require.NotNil(t, body)
		//require.Equal(t, 200, resp.StatusCode)
	})
	t.Run("found pic in cache", func(t *testing.T) {
		defer wg.Done()
	})

	// Закрыть сервер и приложение
	wg.Wait()
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("can't stop publishing test static")
	}
}

func TestIntegrationNegative(t *testing.T) {
	// Развернуть веб сервис со статическими картинками
	wg := sync.WaitGroup{}
	server := &http.Server{Addr: ":3000", Handler: http.FileServer(http.Dir("../examples"))}
	go server.ListenAndServe()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// Запустить наше приложение
	go func(ctx context.Context) {
		main()
	}(ctx)

	// Реализовать тесты отказа:
	wg.Add(5)
	t.Run("remote server not exist", func(t *testing.T) {
		defer wg.Done()
	})
	t.Run("remote server exists, but pic not found (404 Not Found)", func(t *testing.T) {
		defer wg.Done()
	})
	t.Run("remote server exists, but pic is not pic", func(t *testing.T) {
		defer wg.Done()
	})
	t.Run("remote server return ISE (500)", func(t *testing.T) {
		defer wg.Done()
	})
	t.Run("remote server return plain html or test", func(t *testing.T) {
		defer wg.Done()
	})

	// Закрыть сервер и приложение
	wg.Wait()
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("can't stop publishing test static")
	}
}

func request(addr string, timeout time.Duration) ([]byte, *http.Response, error) {
	client := &http.Client{}
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	req, err := http.NewRequestWithContext(ctx, "GET", addr, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Close = true
	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}
	if err = res.Body.Close(); err != nil {
		return nil, nil, err
	}
	return body, res, nil
}