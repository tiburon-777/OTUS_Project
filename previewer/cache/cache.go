package cache

import (
	"log"
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                              // Очистить кэш
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

func (l *lruCache) Set(key Key, value interface{}) bool {
	if _, exists := l.items[key]; exists {
		l.mx.RLock()
		l.items[key].Value = Item{Value: value, Key: key}
		l.queue.MoveToFront(l.items[key])
		l.mx.RUnlock()
		return exists
	}
	if l.queue.Len() == l.capacity {
		l.mx.RLock()
		k, ok := l.queue.Back().Value.(Item)
		if !ok {
			log.Fatal("Ошибка приведения типов")
		}
		delete(l.items, k.Key)
		l.queue.Remove(l.queue.Back())
		l.mx.RUnlock()
	}
	l.mx.RLock()
	l.items[key] = l.queue.PushFront(Item{Value: value, Key: key})
	l.mx.RUnlock()
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mx.Lock()
	defer l.mx.Unlock()
	if l.items[key] == nil {
		return nil, false
	}
	l.queue.MoveToFront(l.items[key])
	s, ok := l.items[key].Value.(Item)
	if !ok {
		log.Fatal("Ошибка приведения типов")
	}
	return s.Value, true
}

func (l *lruCache) Clear() {
	l.mx.Lock()
	l.items = nil
	l.queue.len = 0
	l.queue.Info = ListItem{}
	l.mx.Unlock()
}
