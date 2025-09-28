package map_utils

import (
	"sync"
)

type LockMapMutateActionType = int

const (
	LockMapMutateActionType_Noop LockMapMutateActionType = iota
	LockMapMutateActionType_Replace
	LockMapMutateActionType_Delete
)

type LockMapMutateAction[T comparable, V any] struct {
	Action LockMapMutateActionType
	Key    T
	Value  V
}

type LockMap[T comparable, V any] struct {
	Mutex sync.RWMutex
	Map   map[T]V
}

func NewLockMap[T comparable, V any]() *LockMap[T, V] {
	return &LockMap[T, V]{
		Mutex: sync.RWMutex{},
		Map:   map[T]V{},
	}
}

func (m *LockMap[T, V]) Get(key T) (V, bool) {
	m.Mutex.RLock()
	defer m.Mutex.RUnlock()
	value, has_value := m.Map[key]
	return value, has_value
}

func (m *LockMap[T, V]) Set(key T, value V) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.Map[key] = value
}

func (m *LockMap[T, V]) ForEach(callback func(value V, key T) bool) {
	m.Mutex.RLock()
	defer m.Mutex.RUnlock()
	for key, value := range m.Map {
		if !callback(value, key) {
			break
		}
	}
}

func (m *LockMap[T, V]) Mutate(callback func(value V, key T) LockMapMutateAction[T, V]) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	for key, value := range m.Map {
		action := callback(value, key)
		if action.Action == LockMapMutateActionType_Delete {
			delete(m.Map, action.Key)
		} else if action.Action == LockMapMutateActionType_Replace {
			m.Map[action.Key] = action.Value
		}
	}
}
