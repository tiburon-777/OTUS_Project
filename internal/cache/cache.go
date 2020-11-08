package cache

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) (bool, error) // Добавить значение в кэш по ключу
	Get(key Key) (interface{}, bool, error)       // Получить значение из кэша по ключу
	Clear() error                                 // Очистить кэш
}

type lruCache struct {
	capacity int
	path     string
	queue    *List
	items    map[Key]*ListItem
	mx       sync.RWMutex
}

type Item struct {
	Key   Key
	Value interface{}
}

func NewCache(capacity int, path string) Cache {
	if _, err := ioutil.ReadDir(path); err != nil {
		log.Printf("cache directory %s not exists. Try to create.\n", path)
		err = os.MkdirAll(path, 0777)
		if err != nil {
			log.Fatalf("can't create cache directory %s: %s", path, err.Error())
		}
	}
	return &lruCache{
		capacity: capacity,
		path:     path,
		queue:    NewList(),
		items:    make(map[Key]*ListItem),
	}
}

func (l *lruCache) Set(key Key, value interface{}) (bool, error) {
	l.mx.Lock()
	defer l.mx.Unlock()
	if _, exists := l.items[key]; exists {
		err := l.loadOut(key, value)
		if err != nil {
			return false, fmt.Errorf("can't replace file %s: %w", path.Join([]string{l.path, string(key)}...), err)
		}
		l.items[key].Value = Item{Value: value, Key: key}
		l.queue.MoveToFront(l.items[key])
		return exists, nil
	}
	if l.queue.Len() == l.capacity {
		k, ok := l.queue.Back().Value.(Item)
		if !ok {
			return false, fmt.Errorf("can't cast type")
		}
		err := l.remove(k.Value.(Key))
		if err != nil {
			return false, fmt.Errorf("can't delete file %s: %w", path.Join([]string{l.path, k.Value.(string)}...), err)
		}
		delete(l.items, k.Key)
		l.queue.Remove(l.queue.Back())
	}
	err := l.loadOut(key, value)
	if err != nil {
		return false, fmt.Errorf("can't save file %s: %w", path.Join([]string{l.path, string(key)}...), err)
	}
	if l.items == nil {
		l.items = make(map[Key]*ListItem)
	}
	l.items[key] = l.queue.PushFront(Item{Value: key, Key: key})
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
	pic, err := l.loadIn(s.Key)
	if err != nil {
		return nil, false, fmt.Errorf("can't load file %s: %w", path.Join([]string{l.path, string(s.Key)}...), err)
	}
	return pic, true, nil
}

func (l *lruCache) Clear() error {
	l.mx.Lock()
	defer l.mx.Unlock()
	err := l.drop()
	if err != nil {
		return fmt.Errorf("can't remove files from %s: %w", l.path, err)
	}
	l.items = nil
	l.queue.len = 0
	l.queue.Info = ListItem{}
	return nil
}

func (l *lruCache) loadOut(name Key, pic interface{}) error {
	filename := path.Join([]string{l.path, string(name)}...)
	err := ioutil.WriteFile(filename, pic.([]byte), 0600)
	if err != nil {
		return fmt.Errorf("can't create or write file %s: %w", filename, err)
	}
	return nil
}

func (l *lruCache) loadIn(name Key) ([]byte, error) {
	filename := path.Join([]string{l.path, string(name)}...)
	f, err := os.Open(filename)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return nil, fmt.Errorf("can't open file %s: %w", filename, err)
	}
	res, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("can't read file %s: %w", filename, err)
	}
	return res, nil
}

func (l *lruCache) remove(name Key) error {
	filename := path.Join([]string{l.path, string(name)}...)
	err := os.RemoveAll(filename)
	if err != nil {
		return fmt.Errorf("can't remove file %s: %w", filename, err)
	}
	return nil
}

func (l *lruCache) drop() error {
	dir, err := ioutil.ReadDir(l.path)
	if err != nil {
		return fmt.Errorf("can't read directory %s: %w", l.path, err)
	}
	for _, d := range dir {
		if d.Name() != "nofile" {
			err := os.Remove(path.Join([]string{l.path, d.Name()}...))
			if err != nil {
				return fmt.Errorf("can't remove file %s/%s: %w", l.path, d.Name(), err)
			}
		}
	}
	return nil
}
