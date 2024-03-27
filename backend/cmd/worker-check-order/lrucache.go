package main

import "container/list"

type order struct {
	key   string
	order []byte
}

type LRU struct {
	capacity int
	items    map[string]*list.Element
	queue    *list.List
}

func New(cap int) *LRU {
	return &LRU{
		capacity: cap,
		items:    make(map[string]*list.Element),
		queue:    list.New(),
	}
}

func (l *LRU) Get(key string) ([]byte, bool) {
	elem, isExist := l.items[key]

	if isExist == false {
		return nil, false
	}

	l.queue.MoveToFront(elem)
	return elem.Value.(*order).order, true
}

func (l *LRU) Set(key string, value []byte) {
	if elem, isExist := l.items[key]; isExist == true {
		l.queue.MoveToFront(elem)
		elem.Value.(*order).order = value
	}

	if l.queue.Len() == l.capacity {
		l.purge()
	}

	item := &order{
		key:   key,
		order: value,
	}

	elem := l.queue.PushFront(item)
	l.items[key] = elem
}

func (l *LRU) Clear() {
	l.queue = list.New()
}

func (l *LRU) purge() {
	if elem := l.queue.Back(); elem != nil {
		deletedOrder := l.queue.Remove(elem).(*order)
		delete(l.items, deletedOrder.key)
	}
}
