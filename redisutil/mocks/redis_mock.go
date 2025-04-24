package mocks

import (
	"errors"
	"fmt"
	"time"
)

type MockCache struct {
	store       map[string]any
	GetCalls    []string
	SetCalls    []string
	DeleteCalls []string
}

func NewMockCache() *MockCache {
	return &MockCache{
		store: map[string]any{},
	}
}

func (m *MockCache) Set(key string, val any, exp time.Duration) (err error) {
	m.SetCalls = append(m.SetCalls, key)
	m.store[key] = val
	return nil
}

func (m *MockCache) Get(key string) (val string, err error) {
	m.GetCalls = append(m.GetCalls, key)
	out, ok := m.store[key]
	if !ok {
		return "", errors.New("val in key not found")
	}

	val, ok = out.(string)
	if !ok {
		return "", fmt.Errorf("value for key %s is not string", out)
	}
	return val, nil
}

func (m *MockCache) Delete(key string) (delCount int64, err error) {
	m.DeleteCalls = append(m.DeleteCalls, key)
	delete(m.store, key)
	return 0, nil
}
