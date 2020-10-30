package application

import (
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestBuildQuery(t *testing.T) {

	urlParcer := func(u string) *url.URL {
		res, _ := url.Parse(u)
		return res
	}

	table := []struct {
		url       *url.URL
		expWidth  int
		expHeight int
		expURL    *url.URL
		err       bool
		msg       string
	}{
		{
			url: urlParcer("/fill/10/10/domain.me/some/pic.jpg"), expWidth: 10, expHeight: 10, expURL: urlParcer("http://domain.me/some/pic.jpg"), err: false, msg: "Normal request",
		},
		{
			url: urlParcer("/fill/10/10/pic.jpg"), expWidth: 10, expHeight: 10, expURL: urlParcer("http://pic.jpg"), err: false, msg: "Short URL",
		},
		{
			url: urlParcer("/fill/10"), expWidth: 0, expHeight: 0, expURL: nil, err: true, msg: "Only width",
		},
		{
			url: urlParcer("/fill/10/10"), expWidth: 0, expHeight: 0, expURL: nil, err: true, msg: "Only dimensions",
		},
		{
			url: urlParcer("/fill/qwew/qwew/domain.me/some/pic.jpg"), expWidth: 0, expHeight: 0, expURL: nil, err: true, msg: "Strings in dimensions",
		},
		{
			url: urlParcer("/fill/domain.me/some/pic.jpg"), expWidth: 0, expHeight: 0, expURL: nil, err: true, msg: "No dimensions",
		},
	}

	for _, dat := range table {
		t.Run(dat.msg, func(t *testing.T) {
			i := false
			query, err := buildQuery(dat.url)
			if err != nil {
				i = true
			}
			require.Equal(t, dat.err, i, dat.msg)
			require.Equal(t, dat.expWidth, query.Width, dat.msg)
			require.Equal(t, dat.expHeight, query.Height, dat.msg)
			require.Equal(t, dat.expURL, query.URL, dat.msg)
		})
	}
}
