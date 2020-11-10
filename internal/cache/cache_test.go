package cache

import (
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		cacheDir, err := ioutil.TempDir("", "cache_.")
		require.NoError(t, err, err)
		defer os.RemoveAll(cacheDir)
		c, err := NewCache(10, cacheDir)
		require.NoError(t, err, err)
		err = c.Clear()
		require.NoError(t, err, err)

		_, ok, err := c.Get("aaa")
		require.NoError(t, err)
		require.False(t, ok)

		_, ok, err = c.Get("bbb")
		require.NoError(t, err)
		require.False(t, ok)

		err = c.Clear()
		require.NoError(t, err, err)
	})

	t.Run("simple", func(t *testing.T) {
		cacheDir, err := ioutil.TempDir("", "cache_.")
		require.NoError(t, err, err)
		defer os.RemoveAll(cacheDir)
		c, err := NewCache(5, cacheDir)
		require.NoError(t, err, err)
		err = c.Clear()
		require.NoError(t, err, err)

		wasInCache, err := c.Set("aaa", []byte("pic #1111"))
		require.NoError(t, err)
		require.False(t, wasInCache)

		wasInCache, err = c.Set("bbb", []byte("pic #2222"))
		require.NoError(t, err)
		require.False(t, wasInCache)

		val, ok, err := c.Get("aaa")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, []byte("pic #1111"), val)

		val, ok, err = c.Get("bbb")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, []byte("pic #2222"), val)

		wasInCache, err = c.Set("aaa", []byte("pic #3333"))
		require.NoError(t, err)
		require.True(t, wasInCache)

		val, ok, err = c.Get("aaa")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, []byte("pic #3333"), val)

		val, ok, err = c.Get("ccc")
		require.NoError(t, err)
		require.False(t, ok)
		require.Nil(t, val)

		err = c.Clear()
		require.NoError(t, err, err)
	})

	t.Run("purge logic", func(t *testing.T) {
		cacheDir, err := ioutil.TempDir("", "cache_.")
		require.NoError(t, err, err)
		defer os.RemoveAll(cacheDir)
		c, err := NewCache(3, cacheDir)
		require.NoError(t, err, err)
		err = c.Clear()
		require.NoError(t, err, err)

		wasInCache, err := c.Set("aaa", []byte("pic #1111"))
		require.NoError(t, err)
		require.False(t, wasInCache)

		wasInCache, err = c.Set("bbb", []byte("pic #2222"))
		require.NoError(t, err)
		require.False(t, wasInCache)

		wasInCache, err = c.Set("ccc", []byte("pic #3333"))
		require.NoError(t, err)
		require.False(t, wasInCache)

		_, ok, err := c.Get("bbb")
		require.NoError(t, err)
		require.True(t, ok)

		_, ok, err = c.Get("aaa")
		require.NoError(t, err)
		require.True(t, ok)

		wasInCache, err = c.Set("ddd", []byte("pic #4444"))
		require.NoError(t, err)
		require.False(t, wasInCache)

		_, ok, err = c.Get("ddd")
		require.NoError(t, err)
		require.True(t, ok)

		_, ok, err = c.Get("ccc")
		require.NoError(t, err)
		require.False(t, ok)

		err = c.Clear()
		require.NoError(t, err, err)
	})
}

func TestCacheMultithreading(t *testing.T) {
	cacheDir, err := ioutil.TempDir("", "cache_.")
	require.NoError(t, err, err)
	defer os.RemoveAll(cacheDir)
	c, err := NewCache(10, cacheDir)
	require.NoError(t, err, err)
	err = c.Clear()
	require.NoError(t, err, err)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 10_000; i++ {
			itm := strconv.Itoa(i)
			_, err := c.Set(Key(itm), []byte(itm))
			require.NoError(t, err, err)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 10_000; i++ {
			itm := strconv.Itoa(rand.Intn(10_000))
			b, s, err := c.Get(Key(itm))
			require.NoError(t, err, err)
			if s {
				require.Equal(t, itm, string(b.([]uint8)))
			}

		}
	}()

	wg.Wait()

	err = c.Clear()
	require.NoError(t, err, err)
}
