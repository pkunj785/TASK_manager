package kv

import (
	"fmt"
	"sync"
)

/*
	In-memory KV store 
*/

type Storer[K comparable, V any] interface {
	Put(K, V) error
	Get(K) (V, error)
	Update(K, V) error
	Delete(K) error
}

///////////////////////////////

type KVStore[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

//////////////////////////////////////

func NewKVStore[K comparable, V any]() *KVStore[K, V] {
	return &KVStore[K, V]{
		data: make(map[K]V),
	}
}

////////////////////////////////////

func (s *KVStore[K, V]) Put(key K, val V) error {

	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = val

	return nil
}

////////////////////////////////////////////////

func (s *KVStore[K, V]) Get(key K) (V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.data[key]

	if !ok {
		return val, fmt.Errorf("the key (%v) does not exist", key)
	}

	return val, nil
}

/////////////////////////////////////////////////

func (s *KVStore[K, V]) Update(key K, val V) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Has(key) {
		return fmt.Errorf("the key (%v) does not exist", key)
	}

	s.data[key] = val

	return nil
}

/////////////////////////////////////////////////////

func (s *KVStore[K, V]) Delete(key K) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[key]
	if !ok {
		return fmt.Errorf("the key (%v) does not exits", key)
	}

	delete(s.data, key)

	fmt.Printf("the key (%v) successfully deleted\n", key)
	return nil
}

///////////////////////////////////////////

func (s *KVStore[K, V]) Has(key K) bool {
	_, ok := s.data[key]

	return ok
}