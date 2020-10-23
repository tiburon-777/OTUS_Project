package application

import (
	"errors"
	"github.com/tiburon-777/OTUS_Project/previewer/models"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func buildQuery(u *url.URL) (q models.Query, err error) {
	t := strings.Split(u.Path, "/")
	q.Width,err=strconv.Atoi(t[2])
	if err!=nil {
		return models.Query{}, errors.New("width must be an integer")
	}
	q.Height,err=strconv.Atoi(t[3])
	if err!=nil {
		return models.Query{}, errors.New("height must be an integer")
	}
	tn := "http://"+strings.Join(t[4:],"/")
	q.URL,err=q.URL.Parse(tn)
	if err!=nil {
		return models.Query{}, errors.New("not valid url")
	}
	return q,nil
}

func getPic(q models.Query) ([]byte,http.Header,error) {
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