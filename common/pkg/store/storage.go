// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package store is an in-memory collection keyed by Identity(), safe for concurrent use.
package store

import (
	"slices"
	"sync"
)

// Identifiable is anything the store can key. Identity must be stable for the lifetime of the value.
type Identifiable interface {
	Identity() string
}

// OrgEnvAware is a request scoped to an organization and environment. Used to pick a tenant bucket.
type OrgEnvAware interface {
	GetOrgId() string
	GetEnvId() string
}

type pageRange struct {
	page int
	size int
}

func (s pageRange) fromIndex() int {
	if s.page < 1 {
		return 0
	}
	return (s.page - 1) * s.size
}

func (s pageRange) toIndex(sliceSize int) int {
	if s.size < 1 {
		s.size = 10
	}
	from := s.fromIndex()
	return min(from+s.size, sliceSize)
}

// Store holds values of T by Identity(). Insertion order is kept. Callers must not copy a Store.
type Store[T Identifiable] struct {
	data        map[string]T
	identifiers []string
	mutex       sync.Mutex
}

// NewStoreWithData returns a store that already contains identifiables. Duplicate identities upsert in order.
func NewStoreWithData[T Identifiable](identifiables ...T) *Store[T] {
	s := NewStore[T]()
	for _, keyed := range identifiables {
		s.Put(keyed)
	}
	return s
}

// NewStore returns an empty store.
func NewStore[T Identifiable]() *Store[T] {
	s := &Store[T]{
		data:        make(map[string]T),
		identifiers: make([]string, 0),
		mutex:       sync.Mutex{},
	}
	return s
}

// GetAll returns a snapshot of all values in insertion order. The slice is independent of later mutations.
func (s *Store[T]) GetAll() []T {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	result := make([]T, 0, len(s.identifiers))
	for _, id := range s.identifiers {
		result = append(result, s.data[id])
	}
	return result
}

// GetPage returns a snapshot of one page in insertion order. page < 1 is treated as the first page.
// size < 1 defaults to 10. A page past the end returns nil, not a panic.
func (s *Store[T]) GetPage(page int, size int) []T {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	r := make([]T, 0)

	pr := pageRange{page: page, size: size}
	from := pr.fromIndex()
	to := pr.toIndex(len(s.identifiers))

	if from >= len(s.identifiers) {
		return r
	}

	for _, k := range s.identifiers[from:to] {
		r = append(r, s.data[k])
	}
	return r
}

// Get returns the value for key. ok is false when the key is absent.
func (s *Store[T]) Get(key string) (T, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	identifiable, ok := s.data[key]
	if !ok {
		return *new(T), false
	}
	return identifiable, true
}

// Put inserts or replaces by Identity(). A new key is appended; an existing key keeps its position.
func (s *Store[T]) Put(identified T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	id := identified.Identity()
	if _, ok := s.data[id]; !ok {
		s.identifiers = append(s.identifiers, id)
	}
	s.data[id] = identified
}

// DeleteByKey removes the key if present. Missing keys are a no-op.
func (s *Store[T]) DeleteByKey(identifier string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.data[identifier]; ok {
		index := slices.Index(s.identifiers, identifier)
		if index >= 0 {
			s.identifiers = slices.Delete(s.identifiers, index, index+1)
		}
	}
	delete(s.data, identifier)
}

// Delete removes identified by Identity(). Missing keys are a no-op.
func (s *Store[T]) Delete(identified T) {
	s.DeleteByKey(identified.Identity())
}
