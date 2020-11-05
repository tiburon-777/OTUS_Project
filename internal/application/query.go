package application

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Query struct {
	Width  int
	Height int
	URL    *url.URL
}

func buildQuery(u *url.URL) (q Query, err error) {
	t := strings.Split(u.Path, "/")
	if len(t) < 5 {
		return Query{}, errors.New("need more params")
	}
	q.Width, err = strconv.Atoi(t[2])
	if err != nil {
		return Query{}, errors.New("width must be an integer")
	}
	q.Height, err = strconv.Atoi(t[3])
	if err != nil {
		return Query{}, errors.New("height must be an integer")
	}
	tn := "http://" + strings.Join(t[4:], "/")
	q.URL, err = q.URL.Parse(tn)
	if err != nil {
		return Query{}, errors.New("not valid url")
	}
	return q, nil
}

func (q Query) id() string {
	return strconv.Itoa(q.Width) + "/" + strconv.Itoa(q.Height) + "/" + q.URL.Path
}

func (q Query) fromOrigin(timeout time.Duration) ([]byte, *http.Response, error) {
	client := &http.Client{}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://"+q.URL.Host+q.URL.Path, nil)
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
