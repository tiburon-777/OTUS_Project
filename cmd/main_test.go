package main

import (
	"context"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"
)

const testPortBase = 3000

func TestIntegrationPositive(t *testing.T) {
	testPort := strconv.Itoa(testPortBase + 1)
	wg := sync.WaitGroup{}
	server := &http.Server{Addr: "localhost:" + testPort, Handler: http.FileServer(http.Dir("../test/data"))}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Println(err.Error())
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	go func(ctx context.Context) {
		main()
	}(ctx)
	time.Sleep(3*time.Second)

	// Реализовать тесты логики приложения (ресайзы по разным требованиям):
	wg.Add(18)
	t.Run("test static", func(t *testing.T) {
		defer wg.Done()
		body, resp, err := request("http://localhost:"+testPort+"/gopher_original_1024x504.jpg", 15*time.Second)
		require.NoError(t, err)
		require.NotNil(t, body)
		require.Equal(t, 200, resp.StatusCode)
	})
	t.Run("remote server return jpeg in original size", func(t *testing.T) {
		defer wg.Done()
		body, resp, err := request("http://localhost:8080/fill/1024/504/localhost:"+testPort+"/gopher_original_1024x504.jpg", 15*time.Second)
		require.NoError(t, err)
		fSize := 63488
		require.InDelta(t, len(body),fSize, float64(fSize/100)*2, "File size should be about "+strconv.Itoa(fSize/1024)+"Kb~2%")
		require.Equal(t, 200, resp.StatusCode)
		require.Equal(t,resp.Header.Get("X-From-Appcache"), "")
	})
	t.Run("remote server return png in original size", func(t *testing.T) {
		defer wg.Done()
		body, resp, err := request("http://localhost:8080/fill/972/603/localhost:"+testPort+"/test.png", 15*time.Second)
		require.NoError(t, err)
		fSize := 433896
		require.InDelta(t, len(body),fSize, float64(fSize/100)*2, "File size should be about "+strconv.Itoa(fSize/1024)+"Kb~2%")
		require.Equal(t, 200, resp.StatusCode)
		require.Equal(t,resp.Header.Get("X-From-Appcache"), "")
	})
	t.Run("remote server return gif in original size", func(t *testing.T) {
		defer wg.Done()
		body, resp, err := request("http://localhost:8080/fill/480/320/localhost:"+testPort+"/test.gif", 15*time.Second)
		require.NoError(t, err)
		fSize := 34508
		require.InDelta(t, len(body),fSize, float64(fSize/100)*2, "File size should be about "+strconv.Itoa(fSize/1024)+"Kb~2%")
		require.Equal(t, 200, resp.StatusCode)
		require.Equal(t,resp.Header.Get("X-From-Appcache"), "")
	})
	t.Run("found pic in cache", func(t *testing.T) {
		defer wg.Done()
		body, resp, err := request("http://localhost:8080/fill/1024/504/localhost:"+testPort+"/gopher_original_1024x504.jpg", 15*time.Second)
		require.NoError(t, err)
		fSize := 63488
		require.InDelta(t, len(body),fSize, float64(fSize/100)*2, "File size should be about "+strconv.Itoa(fSize/1024)+"Kb~2%")
		require.Equal(t, 200, resp.StatusCode)
		require.Equal(t,resp.Header.Get("X-From-Appcache"), "true")
	})
	t.Run("resize PNG to 400x400", func(t *testing.T) {
		defer wg.Done()
		body, resp, err := request("http://localhost:8080/fill/400/400/localhost:"+testPort+"/test.png", 15*time.Second)
		require.NoError(t, err)
		fSize := 161317
		require.InDelta(t, len(body),fSize, float64(fSize/100)*2, "File size should be about "+strconv.Itoa(fSize/1024)+"Kb~2%")
		require.Equal(t, 200, resp.StatusCode)
		require.Equal(t,resp.Header.Get("X-From-Appcache"), "")
	})
	t.Run("resize GIF to 200x200", func(t *testing.T) {
		defer wg.Done()
		body, resp, err := request("http://localhost:8080/fill/200/200/localhost:"+testPort+"/test.gif", 15*time.Second)
		require.NoError(t, err)
		fSize := 11913
		require.InDelta(t, len(body),fSize, float64(fSize/100)*2, "File size should be about "+strconv.Itoa(fSize/1024)+"Kb~2%")
		require.Equal(t, 200, resp.StatusCode)
		require.Equal(t,resp.Header.Get("X-From-Appcache"), "")
	})
	t.Run("resize JPEG to 50x50", func(t *testing.T) {
		defer wg.Done()
		body, resp, err := request("http://localhost:8080/fill/50/50/localhost:"+testPort+"/gopher_original_1024x504.jpg", 15*time.Second)
		require.NoError(t, err)
		fSize := 1437
		require.InDelta(t, len(body),fSize, float64(fSize/100)*2, "File size should be about "+strconv.Itoa(fSize/1024)+"Kb~2%")
		require.Equal(t, 200, resp.StatusCode)
		require.Equal(t,resp.Header.Get("X-From-Appcache"), "")
	})
	t.Run("resize JPEG to 200x70", func(t *testing.T) {
		defer wg.Done()
		body, resp, err := request("http://localhost:8080/fill/200/70/localhost:"+testPort+"/gopher_original_1024x504.jpg", 15*time.Second)
		require.NoError(t, err)
		fSize := 3875
		require.InDelta(t, len(body),fSize, float64(fSize/100)*2, "File size should be about "+strconv.Itoa(fSize/1024)+"Kb~2%")
		require.Equal(t, 200, resp.StatusCode)
		require.Equal(t,resp.Header.Get("X-From-Appcache"), "")
	})
	t.Run("resize JPEG to 256x126", func(t *testing.T) {
		defer wg.Done()
		body, resp, err := request("http://localhost:8080/fill/256/126/localhost:"+testPort+"/gopher_original_1024x504.jpg", 15*time.Second)
		require.NoError(t, err)
		fSize := 6803
		require.InDelta(t, len(body),fSize, float64(fSize/100)*2, "File size should be about "+strconv.Itoa(fSize/1024)+"Kb~2%")
		require.Equal(t, 200, resp.StatusCode)
		require.Equal(t,resp.Header.Get("X-From-Appcache"), "")
	})
	t.Run("resize JPEG to 333x666", func(t *testing.T) {
		defer wg.Done()
		body, resp, err := request("http://localhost:8080/fill/333/666/localhost:"+testPort+"/gopher_original_1024x504.jpg", 15*time.Second)
		require.NoError(t, err)
		fSize := 28749
		require.InDelta(t, len(body),fSize, float64(fSize/100)*2, "File size should be about "+strconv.Itoa(fSize/1024)+"Kb~2%")
		require.Equal(t, 200, resp.StatusCode)
		require.Equal(t,resp.Header.Get("X-From-Appcache"), "")
	})
	t.Run("resize JPEG to 500x500", func(t *testing.T) {
		defer wg.Done()
		body, resp, err := request("http://localhost:8080/fill/500/500/localhost:"+testPort+"/gopher_original_1024x504.jpg", 15*time.Second)
		require.NoError(t, err)
		fSize := 32606
		require.InDelta(t, len(body),fSize, float64(fSize/100)*2, "File size should be about "+strconv.Itoa(fSize/1024)+"Kb~2%")
		require.Equal(t, 200, resp.StatusCode)
		require.Equal(t,resp.Header.Get("X-From-Appcache"), "")
	})
	t.Run("resize JPEG to 1024x252", func(t *testing.T) {
		defer wg.Done()
		body, resp, err := request("http://localhost:8080/fill/1024/252/localhost:"+testPort+"/gopher_original_1024x504.jpg", 15*time.Second)
		require.NoError(t, err)
		fSize := 30356
		require.InDelta(t, len(body),fSize, float64(fSize/100)*2, "File size should be about "+strconv.Itoa(fSize/1024)+"Kb~2%")
		require.Equal(t, 200, resp.StatusCode)
		require.Equal(t,resp.Header.Get("X-From-Appcache"), "")
	})
	t.Run("resize JPEG to 2000x1000", func(t *testing.T) {
		defer wg.Done()
		body, resp, err := request("http://localhost:8080/fill/2000/1000/localhost:"+testPort+"/gopher_original_1024x504.jpg", 15*time.Second)
		require.NoError(t, err)
		fSize := 151996
		require.InDelta(t, len(body),fSize, float64(fSize/100)*2, "File size should be about "+strconv.Itoa(fSize/1024)+"Kb~2%")
		require.Equal(t, 200, resp.StatusCode)
		require.Equal(t,resp.Header.Get("X-From-Appcache"), "")
	})
	t.Run("remote server not exist (502 Bad request)", func(t *testing.T) {
		defer wg.Done()
		_, resp, err := request("http://localhost:8080/fill/1024/252/abracadabra/fakepic.jpg", 15*time.Second)
		require.NoError(t, err)
		require.Equal(t, 502, resp.StatusCode)
		require.Equal(t,resp.Header.Get("X-From-Appcache"), "")
	})
	t.Run("remote server exists, but pic not found (404 Not Found)", func(t *testing.T) {
		defer wg.Done()
		_, resp, err := request("http://localhost:8080/fill/1024/252/localhost:"+testPort+"/fakepic.jpg", 15*time.Second)
		require.NoError(t, err)
		require.Equal(t, 404, resp.StatusCode)
		require.Equal(t,resp.Header.Get("X-From-Appcache"), "")
	})
	t.Run("remote server exists, but pic is not pic (500 Internal Server Error)", func(t *testing.T) {
		defer wg.Done()
		_, resp, err := request("http://localhost:8080/fill/1024/252/localhost:"+testPort+"/test.exe", 15*time.Second)
		require.NoError(t, err)
		require.Equal(t, 500, resp.StatusCode)
		require.Equal(t,resp.Header.Get("X-From-Appcache"), "")
	})
	t.Run("remote server return plain html or texst", func(t *testing.T) {
		defer wg.Done()
		_, resp, err := request("http://localhost:8080/fill/1024/252/localhost:"+testPort+"/test.html", 15*time.Second)
		require.NoError(t, err)
		require.Equal(t, 500, resp.StatusCode)
		require.Equal(t,resp.Header.Get("X-From-Appcache"), "")
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
