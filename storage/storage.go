package storage

import "errors"

type MemoryCache struct {
	Cache map[string]interface{}
}

var (
	errKeyHold = errors.New("Key holding the wrong kind of value")
)

//Create new MemoryCache instance
func New() *MemoryCache {
	return &MemoryCache{Cache: make(map[string]interface{})}
}

// Get the value of a key
func (m *MemoryCache) Get(key string) (value string, err error) {
	switch v := m.Cache[key].(type) {
	case string:
		value = v

	case nil:
		value = "nil"

	default:
		err = errKeyHold
	}

	return
}

// Sets the value at the specified key
func (m *MemoryCache) Set(key string, value string) (err error) {
	m.Cache[key] = value
	return
}

// Set the string value of the field
func (m *MemoryCache) HSet(key string, field string, value string) (err error) {
	switch v := m.Cache[key].(type) {
	case map[string]string:
		v[field] = value
		m.Cache[key] = v

	case nil:
		m.Cache[key] = map[string]string{
			field: value}

	default:
		err = errKeyHold
	}

	return
}

// Get the value of a hash field stored at specified key
func (m *MemoryCache) HGet(key string, field string) (value string, err error) {
	switch v := m.Cache[key].(type) {
	case map[string]string:
		v2, ok := v[field]
		if !ok {
			value = "nil"
			break
		}
		value = v2

	case nil:
		value = "nil"

	default:
		err = errKeyHold
	}
	return
}

// LPush prepend one or multiple values to a list
func (m *MemoryCache) LPush(key string, values ...string) (err error) {
	var list []string
	switch v := m.Cache[key].(type) {
	case []string:
		list = v

	case nil:
		list = make([]string, len(values))

	default:
		err = errKeyHold
		return
	}

	for _, v := range values {
		list = append([]string{v}, list...)
	}

	m.Cache[key] = list
	return

}
