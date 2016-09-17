package storage

import (
	"errors"
	"time"
)

const (
	// Default expiration time
	DefaultExpiration time.Duration = -1
)

var (
	errKeyHold   = errors.New("Key holding the wrong kind of value")
	ErrNullValue = errors.New("Null value")
)

type Item struct {
	Object     interface{}
	Expiration int64
}

type MemoryCache struct {
	items   map[string]Item
	expires map[string]bool
	//mu    sync.RWMutex
}

//Create new MemoryCache instance
func New() *MemoryCache {
	return &MemoryCache{items: make(map[string]Item), expires: make(map[string]bool)}
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
		err = ErrNullValue

	default:
		err = errKeyHold
	}

	return
}

// Remove the specified keys.
func (m *MemoryCache) Del(key string) bool {
	_, ok := m.items[key]
	delete(m.items, key)
	delete(m.expires, key)
	return ok
}

// Set expiration time for specified key
func (m *MemoryCache) SetTTL(key string, d time.Duration) error {

	if _, ok := m.items[key]; !ok {
		return ErrNullValue
	}
	e := int64(DefaultExpiration)
	if d > 0 {
		e = time.Now().Add(d).Unix()
	}
	value := m.items[key]
	value.Expiration = int64(e)
	m.items[key] = value
	m.expires[key] = true
	return nil
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
			err = ErrNullValue
			break
		}
		value = v2

	case nil:
		err = ErrNullValue

	default:
		err = errKeyHold
	}
	return
}

// Get all the fields and values stored in a hash at specified key
func (m *MemoryCache) HGetAll(key string) (values []string, err error) {
	switch v := m.items[key].Object.(type) {
	case map[string]string:
		for k := range v {
			values = append(values, k, v[k])
		}
	case nil:
		err = ErrNullValue

	default:
		err = errKeyHold
	}
	return
}

func (m *MemoryCache) HDel(key string, fields ...string) (n int, err error) {
	switch v := m.items[key].Object.(type) {
	case map[string]string:
		for _, f := range fields {
			_, ok := v[f]
			if ok {
				n++
				delete(v, f)
			}
		}
		m.Set(key, v)

	case nil:
		err = ErrNullValue

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
func (m *MemoryCache) Lindex(key string, i int) (value string, err error) {
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
		err = ErrNullValue

	case nil:
		err = ErrNullValue

	default:
		err = errKeyHold

	}
	return
}

// Get the length of the list stored at key
func (m *MemoryCache) Llen(key string) (n int, err error) {
	switch v := m.items[key].Object.(type) {
	case []string:
		return len(v), nil

	case nil:
		err = ErrNullValue

	default:
		err = errKeyHold
	}

	return
}

// Remove and get the first element in a list
func (m *MemoryCache) LPop(key string) (value string, err error) {
	switch v := m.items[key].Object.(type) {
	case []string:
		if len(v) < 1 {
			err = ErrNullValue
			break
		}
		value = v[0]
		m.Set(key, v[1:])

	case nil:
		err = ErrNullValue

	default:
		err = errKeyHold

	}
	return
}

// Check for key expire
func (m *MemoryCache) IsExpire(key string) bool {
	now := time.Now().Unix()
	_, ok := m.expires[key]
	return ok && (now > m.items[key].Expiration)
}

// Get expire key list
func (m *MemoryCache) ExpireList() (list []string) {
	for key := range m.expires {
		list = append(list, key)
	}
	return
}
