package cache

import (
	"fmt"
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) (bool, error) // Добавить значение в кэш по ключу
	Get(key Key) (interface{}, bool, error)       // Получить значение из кэша по ключу
	Clear()                                       // Очистить кэш
}

type lruCache struct {
	capacity int
	queue    *List
	items    map[Key]*ListItem
	mx       sync.RWMutex
}

type Item struct {
	Key   Key
	Value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem),
	}
}

func (l *lruCache) Set(key Key, value interface{}) (bool, error) {
	l.mx.Lock()
	defer l.mx.Unlock()
	if _, exists := l.items[key]; exists {
		l.items[key].Value = Item{Value: value, Key: key}
		l.queue.MoveToFront(l.items[key])
		return exists, nil
	}
	if l.queue.Len() == l.capacity {
		k, ok := l.queue.Back().Value.(Item)
		if !ok {
			return false, fmt.Errorf("can't cast type")
		}
		delete(l.items, k.Key)
		l.queue.Remove(l.queue.Back())
	}
	l.items[key] = l.queue.PushFront(Item{Value: value, Key: key})
	return false, nil
}

func (l *lruCache) Get(key Key) (interface{}, bool, error) {
	l.mx.Lock()
	defer l.mx.Unlock()
	if l.items[key] == nil {
		return nil, false, nil
	}
	l.queue.MoveToFront(l.items[key])
	s, ok := l.items[key].Value.(Item)
	if !ok {
		return nil, false, fmt.Errorf("can't cast type")
	}
	return s.Value, true, nil
}

func (l *lruCache) Clear() {
	l.mx.Lock()
	defer l.mx.Unlock()
	l.items = nil
	l.queue.len = 0
	l.queue.Info = ListItem{}
}
