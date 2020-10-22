package models

import "net/url"

type Config struct {
	LruCapasity int
}

type Query struct {
	Height int
	Width int
	URL *url.URL
}