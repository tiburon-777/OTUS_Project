package application

import (
	"bytes"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"errors"
	"image/jpeg"
	"github.com/nfnt/resize"
)

type Query struct {
	Height int
	Width int
	URL *url.URL
}

func BuildQuery(u *url.URL) (q Query, err error) {
	t := strings.Split(u.Path, "/")
	q.Width,err=strconv.Atoi(t[2])
	if err!=nil {
		return Query{}, errors.New("width must be an integer")
	}
	q.Height,err=strconv.Atoi(t[3])
	if err!=nil {
		return Query{}, errors.New("height must be an integer")
	}
	tn := "http://"+strings.Join(t[4:],"/")
	q.URL,err=q.URL.Parse(tn)
	if err!=nil {
		return Query{}, errors.New("not valid url")
	}
	return q,nil
}

func (q Query) id() string {
	return strconv.Itoa(q.Width)+"/"+strconv.Itoa(q.Height)+"/"+q.URL.Path
}

func (q Query) fromOrigin() ([]byte, http.Header, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://" + q.URL.Host + "/" + q.URL.Path, nil)
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
	return body, res.Header, nil
}

func (q Query) resize(b []byte) ([]byte, error) {
	i,_,err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	log.Println("ресайзим")
	m := resize.Resize(uint(q.Width), uint(q.Height), i, resize.Bicubic)
	var g []byte
	s := bytes.NewBuffer(g)
	if err = jpeg.Encode(s,m,nil); err != nil {
		return nil, err
	}
	return g, nil
}