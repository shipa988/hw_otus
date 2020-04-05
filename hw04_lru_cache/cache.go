package hw04_lru_cache //nolint:golint,stylecheck
import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
)

type Key string

var _ Cache = (*lruCache)(nil)

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear() error                        // Очистить кэш
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*cacheItem
	mux      *sync.Mutex
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	if l.capacity == 0 {
		return false
	}
	l.mux.Lock()
	defer l.mux.Unlock()
	if item, ok := l.items[key]; ok { // элемент присутствует в словаре
		item.value = value
		//l.queue.PushFront(item)
		if err := l.queue.MoveToFront(item.qvalue); err == nil { //переместить элемент в начало очереди
			item.qvalue = l.queue.Front()
			return ok
		}
	} else { //элемента нет в словаре
		item = &cacheItem{
			key:   key,
			value: value,
		}
		l.queue.PushFront(item)
		item.qvalue = l.queue.Front() //для идентификации cacheItem внутри  queue )
		l.items[key] = item
		if l.queue.Len() > l.capacity { //размер очереди больше ёмкости кэша
			if err := l.Clear(); err != nil {
				fmt.Printf("%v\r\n", errors.Wrap(err, "can't Set element "+string(key)))
			}
		}
		return ok
	}
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mux.Lock()
	defer l.mux.Unlock()
	if item, ok := l.items[key]; ok { // элемент присутствует в словаре
		if err := l.queue.MoveToFront(item.qvalue); err == nil { //переместить элемент в начало очереди
			item.qvalue = l.queue.Front()
			return item.value, ok
		}
	} else { //элемента нет в словаре
		return nil, ok
	}
	return nil, false
}

func (l *lruCache) Clear() error {
	lastItem := l.queue.Back() //последний элемент из очереди
	if lastItem == nil {
		return errors.New("queue is empty")
	} else if item, ok := lastItem.Value.(*cacheItem); ok {
		delete(l.items, item.key) // удалить его значение из словаря
		// удалить последний элемент из очереди
		if err := l.queue.Remove(lastItem); err != nil {
			return errors.Wrap(err, "can't clear lruCache")
		}
	}
	return nil
}

type cacheItem struct {
	key    Key
	value  interface{}
	qvalue *listItem //адрес для идентификации cacheItem внутри  queue
}

func NewCache(capacity int) Cache {
	if capacity < 0 {
		capacity = 0
	}
	cash := &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    map[Key]*cacheItem{},
		mux:      &sync.Mutex{},
	}
	return cash
}
