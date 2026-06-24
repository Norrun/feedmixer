package datautils

import (
	"errors"
	"sync"
)

func NewSimpleCache() *SimpleCace {
	return &SimpleCace{mu: sync.RWMutex{}, m: make(map[TypeKey]any)}
}

type SimpleCace struct {
	mu sync.RWMutex
	m  map[TypeKey]any
}

func (receiver *SimpleCace) Get(key TypeKey) (any, bool) {
	receiver.mu.RLock()
	v, ok := receiver.m[key]
	receiver.mu.RUnlock()
	return v, ok
}

func (receiver *SimpleCace) Set(key TypeKey, value any) error {
	if !key.Validate(value) {
		return errors.New("Invalid type")
	}
	receiver.mu.Lock()
	receiver.m[key] = value
	receiver.mu.Unlock()
	return nil
}

func (receiver *SimpleCace) Delete(key TypeKey) {
	receiver.mu.Lock()
	delete(receiver.m, key)
	receiver.mu.Unlock()
}

func (receiver *SimpleCace) BulkGet(keys ...TypeKey) []Option[any] {
	res := make([]Option[any], 0)
	receiver.mu.RLock()
	for _, v := range keys {
		v, ok := receiver.m[v]
		res = append(res, Option[any]{v, ok})
	}
	receiver.mu.RUnlock()
	return res
}

type KeyValue struct {
	Key   TypeKey
	Value any
}

func (receiver *SimpleCace) BulkSet(entries ...KeyValue) error {
	errs := make([]error, 0)
	for _, v := range entries {
		if !v.Key.Validate(v.Value) {
			errs = append(errs, errors.New("Invalid type"))
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	receiver.mu.Lock()
	for _, v := range entries {
		receiver.m[v.Key] = v.Value
	}
	receiver.mu.Unlock()

	return nil
}

func (receiver *SimpleCace) BulkDelete(keys ...TypeKey) {
	receiver.mu.Lock()
	for _, key := range keys {
		delete(receiver.m, key)
	}
	receiver.mu.Unlock()
}

func (receiver *SimpleCace) Clear() {
	receiver.mu.Lock()
	clear(receiver.m)
	receiver.mu.Unlock()
}

func (receiver *SimpleCace) Has(key TypeKey) bool {
	receiver.mu.RLock()
	_, ok := receiver.m[key]
	receiver.mu.RUnlock()
	return ok
}
