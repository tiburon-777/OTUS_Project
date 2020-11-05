package cache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok, err := c.Get("aaa")
		require.NoError(t, err)
		require.False(t, ok)

		_, ok, err = c.Get("bbb")
		require.NoError(t, err)
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache, err := c.Set("aaa", 100)
		require.NoError(t, err)
		require.False(t, wasInCache)

		wasInCache, err = c.Set("bbb", 200)
		require.NoError(t, err)
		require.False(t, wasInCache)

		val, ok, err := c.Get("aaa")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok, err = c.Get("bbb")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache, err = c.Set("aaa", 300)
		require.NoError(t, err)
		require.True(t, wasInCache)

		val, ok, err = c.Get("aaa")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok, err = c.Get("ccc")
		require.NoError(t, err)
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)

		wasInCache, err := c.Set("aaa", 100)
		require.NoError(t, err)
		require.False(t, wasInCache)

		wasInCache, err = c.Set("bbb", 200)
		require.NoError(t, err)
		require.False(t, wasInCache)

		wasInCache, err = c.Set("ccc", 300)
		require.NoError(t, err)
		require.False(t, wasInCache)

		_, ok, err := c.Get("bbb")
		require.NoError(t, err)
		require.True(t, ok)

		_, ok, err = c.Get("aaa")
		require.NoError(t, err)
		require.True(t, ok)

		wasInCache, err = c.Set("ddd", 400)
		require.NoError(t, err)
		require.False(t, wasInCache)

		_, ok, err = c.Get("ddd")
		require.NoError(t, err)
		require.True(t, ok)

		_, ok, err = c.Get("ccc")
		require.NoError(t, err)
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
