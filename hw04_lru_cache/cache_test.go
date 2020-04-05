package hw04_lru_cache //nolint:golint,stylecheck

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

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})
	t.Run("empty cashe#2", func(t *testing.T) {
		c := NewCache(0)
		wasInCache := c.Set("first", 0)
		require.False(t, wasInCache)
		c2 := NewCache(-1)
		wasInCache = c2.Set("first", -1)
		require.False(t, wasInCache)
	})
	t.Run("simple", func(t *testing.T) {

		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})
	t.Run("clear cache", func(t *testing.T) {
		c := NewCache(5)
		require.Error(t, c.Clear())
		c.Set("first", 1)
		require.NoError(t, c.Clear())
	})
	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(5) //емкость 5 для 10 элементов
		//set
		c.Set("Albania", ".al")
		c.Set("Australia", ".com.au")
		c.Set("Austria", ".at")
		c.Set("Belgium", ".be")
		c.Set("Italy", ".it")
		c.Get("Italy")
		c.Set("Japan", ".jp")
		c.Set("Russian Federation", ".ru")
		c.Set("Singapore", ".sg")
		c.Set("Taiwan", ".tw")
		c.Set("USA", ".us")
		//get/set
		c.Get("Russian Federation")
		c.Get("Japan")
		c.Get("Italy")
		c.Get("USA")
		c.Get("USA")
		c.Get("USA")
		c.Get("Japan")
		c.Get("Russian Federation")
		c.Get("Russian Federation")
		c.Set("South Korea", ".kr")
		c.Set("Romania", ".ro")
		c.Get("USA")
		require.Equal(t, []Key{"Japan", "Russian Federation", "South Korea", "Romania", "USA"}, printCash(c))
	})
}
func printCash(l Cache) []Key {
	sl := []Key{}
	elem := l.(*lruCache).queue.Back()
	for i := 0; i < l.(*lruCache).queue.Len(); i++ {
		if item, ok := elem.Value.(*cacheItem); ok {
			sl = append(sl, item.key)
		}
		elem = elem.Next
	}
	return sl
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
