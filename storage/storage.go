package storage

import (
	"errors"
	"sync"
	"time"
)

const (
	// Default expiration time
	DefaultExpiration time.Duration = -1
	// nil value
	NIL = "(nil)"
)

var (
	errKeyHold = errors.New("Key holding the wrong kind of value")
)

type Item struct {
	Object     interface{}
	Expiration int64
}

type MemoryCache struct {
	items map[string]Item
	mu    sync.RWMutex
}

//Create new MemoryCache instance
func New() *MemoryCache {
	return &MemoryCache{items: make(map[string]Item)}
}

// Sets the value at the specified key
func (m *MemoryCache) Set(key string, value interface{}) (err error) {
	m.items[key] = Item{
		Object:     value,
		Expiration: int64(DefaultExpiration),
	}

	return
}

// Get the value of a key
func (m *MemoryCache) Get(key string) (value string, err error) {
	switch v := m.items[key].Object.(type) {
	case string:
		value = v

	case nil:
		value = NIL

	default:
		err = errKeyHold
	}

	return
}

// Remove the specified keys.
func (m *MemoryCache) Remove(key string) {
	delete(m.items, key)
	return
}

// Set expiration time for specified key
func (m *MemoryCache) SetTTL(key string, d time.Duration) {

	e := int64(DefaultExpiration)
	if d > 0 {
		e = time.Now().Add(d).Unix()
	}
	value := m.items[key]
	value.Expiration = int64(e)
	m.items[key] = value
	return
}

// Set the string value of the field
func (m *MemoryCache) HSet(key string, field string, value string) (err error) {
	switch v := m.items[key].Object.(type) {
	case map[string]string:
		v[field] = value
		m.Set(key, v)

	case nil:
		m.Set(key, map[string]string{
			field: value})

	default:
		err = errKeyHold
	}

	return
}

// Get the value of a hash field stored at specified key
func (m *MemoryCache) HGet(key string, field string) (value string, err error) {
	switch v := m.items[key].Object.(type) {
	case map[string]string:
		v2, ok := v[field]
		if !ok {
			value = NIL
			break
		}
		value = v2

	case nil:
		value = NIL

	default:
		err = errKeyHold
	}
	return
}

// LPush prepend one or multiple values to a list
func (m *MemoryCache) LPush(key string, values ...string) (err error) {
	var list []string
	switch v := m.items[key].Object.(type) {
	case []string:
		list = v

	case nil:
		list = make([]string, 0)

	default:
		err = errKeyHold
		return
	}

	for _, v := range values {
		list = append([]string{v}, list...)
	}

	m.Set(key, list)
	return

}

// Get element from a list by its index
func (m *MemoryCache) Lindex(key string, i int) (value string) {
	switch v := m.items[key].Object.(type) {
	case []string:
		n := len(v)
		if i >= 0 && i < n {
			value = v[i]
			break
		}
		if i >= -n && i < 0 {
			value = v[n+i]
			break
		}
		value = NIL

	default:
		value = NIL

	}
	return
}
