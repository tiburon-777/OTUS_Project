package cache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, l.Len(), 0)
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10)
		require.Equal(t, []int{10}, transpond(l))
		require.Equal(t, 1, l.Len())

		l.Remove(l.Front())
		require.Equal(t, []int{}, transpond(l))
		require.Equal(t, 0, l.Len())

		l.PushFront(10)
		require.Equal(t, []int{10}, transpond(l))
		require.Equal(t, 1, l.Len())

		l.PushBack(20)
		require.Equal(t, []int{20, 10}, transpond(l))
		require.Equal(t, 2, l.Len())

		l.PushBack(30)
		require.Equal(t, []int{30, 20, 10}, transpond(l))
		require.Equal(t, 3, l.Len())

		middle := l.Back().Next // 20
		l.Remove(middle)
		require.Equal(t, []int{30, 10}, transpond(l))
		require.Equal(t, 2, l.Len())

		l.PushFront(20)
		require.Equal(t, []int{30, 10, 20}, transpond(l))
		require.Equal(t, 3, l.Len())

		l.Remove(l.Back())
		require.Equal(t, []int{10, 20}, transpond(l))
		require.Equal(t, 2, l.Len())

		l.PushBack(30)
		require.Equal(t, []int{30, 10, 20}, transpond(l))
		require.Equal(t, 3, l.Len())

		l.Remove(l.Front())
		require.Equal(t, []int{30, 10}, transpond(l))
		require.Equal(t, 2, l.Len())

		l.PushFront(20)
		require.Equal(t, []int{30, 10, 20}, transpond(l))
		require.Equal(t, 3, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		}
		require.Equal(t, []int{70, 50, 30, 10, 20, 40, 60, 80}, transpond(l))
		require.Equal(t, 8, l.Len())

		l.MoveToFront(l.Front())
		require.Equal(t, []int{70, 50, 30, 10, 20, 40, 60, 80}, transpond(l))
		require.Equal(t, 8, l.Len())

		l.MoveToFront(l.Back().Next.Next) // 30
		require.Equal(t, []int{70, 50, 10, 20, 40, 60, 80, 30}, transpond(l))
		require.Equal(t, 8, l.Len())

		l.MoveToFront(l.Back())
		require.Equal(t, []int{50, 10, 20, 40, 60, 80, 30, 70}, transpond(l))
		require.Equal(t, 8, l.Len())
	})
}

func transpond(l *List) []int {
	elems := make([]int, 0, l.Len())
	for i := l.Back(); i != nil; i = i.Next {
		elems = append(elems, i.Value.(int))
	}
	return elems
}
