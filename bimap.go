package bimap

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

var ErrKeyValExists = errors.New("key or value exists")

type BiMap[T, U comparable] struct {
	rwLock sync.RWMutex
	front  map[T]U
	back   map[U]T
}

type option[T, U comparable] interface {
	apply(*BiMap[T, U])
}

type initialOption[T, U comparable] map[T]U

func (io initialOption[T, U]) apply(m *BiMap[T, U]) {
	for k, v := range map[T]U(io) {
		m.front[k] = v
		m.back[v] = k
	}
}

// WithInitialMap returns an initialOption object that implements the option interface
func WithInitialMap[T, U comparable](m map[T]U) option[T, U] {
	return initialOption[T, U](m)
}

// New returns a BiMap object
func New[T, U comparable](options ...option[T, U]) *BiMap[T, U] {
	m := &BiMap[T, U]{
		front: make(map[T]U),
		back:  make(map[U]T),
	}
	for _, opt := range options {
		opt.apply(m)
	}
	return m
}

// GetFront returns the value and its existence by the given key in front map
func (m *BiMap[T, U]) GetFront(key T) (U, bool) {
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()
	v, ok := m.front[key]
	return v, ok
}

// GetBack returns the value and its existence by the given key in back map
func (m *BiMap[T, U]) GetBack(key U) (T, bool) {
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()
	v, ok := m.back[key]
	return v, ok
}

// SetFront sets the value with corresponding key in the front map, it will return an error if either key or value exist
func (m *BiMap[T, U]) SetFront(key T, val U) error {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	var ok bool
	if _, ok = m.front[key]; !ok {
		_, ok = m.back[val]
	}
	if ok {
		return ErrKeyValExists
	}
	m.front[key] = val
	m.back[val] = key
	return nil
}

// SetBack sets the value with corresponding key in the back map, it will return an error if either key or value exist
func (m *BiMap[T, U]) SetBack(key U, val T) error {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	var ok bool
	if _, ok = m.back[key]; !ok {
		_, ok = m.front[val]
	}
	if ok {
		return ErrKeyValExists
	}
	m.back[key] = val
	m.front[val] = key
	return nil
}

// DeleteFront deletes the value of the given key in front map
func (m *BiMap[T, _]) DeleteFront(key T) {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	v, ok := m.front[key]
	if !ok {
		return
	}
	delete(m.front, key)
	delete(m.back, v)
}

// DeleteBack deletes the value of the given key in back map
func (m *BiMap[_, U]) DeleteBack(key U) {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	v, ok := m.back[key]
	if !ok {
		return
	}
	delete(m.back, key)
	delete(m.front, v)
}

// Front returns a new map object that contains all key-value pairs in front map
func (m *BiMap[T, U]) Front() map[T]U {
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()
	nm := make(map[T]U, len(m.front))
	for k, v := range m.front {
		nm[k] = v
	}
	return nm
}

// Back returns a new map object that contains all key-value pairs in back map
func (m *BiMap[T, U]) Back() map[U]T {
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()
	nm := make(map[U]T, len(m.back))
	for k, v := range m.back {
		nm[k] = v
	}
	return nm
}

// Len returns the length of the BiMap object
func (m *BiMap[_, _]) Len() int {
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()
	return len(m.front)
}

// For iterate over the map for the given function
func (m *BiMap[T, U]) For(fn func(f T, b U)) {
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()
	for f, b := range m.front {
		fn(f, b)
	}
}

// String returns a string representation the BiMap object
func (m *BiMap[T, U]) String() string {
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()
	pairs := make([]string, 0, len(m.front))
	for f, b := range m.front {
		pairs = append(pairs, fmt.Sprintf("%v:%v", f, b))
	}
	return "map[" + strings.Join(pairs, " ") + "]"
}
