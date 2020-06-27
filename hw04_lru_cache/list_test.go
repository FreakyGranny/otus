package hw04_lru_cache //nolint:golint,stylecheck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, l.Len(), 7)
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}

func TestCustomList(t *testing.T) {
	t.Run("first-front-push", func(t *testing.T) {
		l := NewList()

		l.PushFront(15)
		require.NotNil(t, l.Back())
		require.Equal(t, 15, l.Back().Value)
	})
	t.Run("first-back-push", func(t *testing.T) {
		l := NewList()

		l.PushBack(25)
		require.NotNil(t, l.Front())
		require.Equal(t, 25, l.Front().Value)
	})
	t.Run("one-side-push", func(t *testing.T) {
		l := NewList()
		l.PushFront(36)
		l.PushFront(26)
		require.Equal(t, 26, l.Front().Value)
		require.Equal(t, 36, l.Back().Value)
	})
	t.Run("double-side-push", func(t *testing.T) {
		l := NewList()
		l.PushFront(35)
		l.PushBack(25)
		require.Equal(t, 35, l.Front().Value)
		require.Equal(t, 25, l.Back().Value)
	})
	t.Run("one-side-push-and-move", func(t *testing.T) {
		l := NewList()
		p := l.PushFront(37)
		l.PushFront(27)
		l.PushFront(17)
		l.MoveToFront(p)
		require.Equal(t, 37, l.Front().Value)
		require.Equal(t, 27, l.Back().Value)
	})
}
