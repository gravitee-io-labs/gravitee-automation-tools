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

package store

import (
	"slices"
	"sync"
)

type Identifiable interface {
	Identity() string
}

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

type Store[T Identifiable] struct {
	data        map[string]T
	identifiers []string
	mutex       sync.Mutex
}

func NewStoreWithData[T Identifiable](identifiables ...T) *Store[T] {
	s := NewStore[T]()
	for _, keyed := range identifiables {
		s.Put(keyed)
	}
	return s
}

func NewStore[T Identifiable]() *Store[T] {
	s := &Store[T]{
		data:        make(map[string]T),
		identifiers: make([]string, 0),
		mutex:       sync.Mutex{},
	}
	return s
}

func (s *Store[T]) GetAll() []T {
	s.mutex.Lock()
	n := len(s.identifiers)
	s.mutex.Unlock()
	return s.GetPage(1, n)
}

func (s *Store[T]) GetPage(page int, size int) []T {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	r := make([]T, 0)

	pr := pageRange{page: page, size: size}
	from := pr.fromIndex()
	to := pr.toIndex(len(s.identifiers))

	if from >= len(s.identifiers) {
		return nil
	}

	for _, k := range s.identifiers[from:to] {
		r = append(r, s.data[k])
	}
	return r
}

func (s *Store[T]) Get(key string) (T, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	identifiable, ok := s.data[key]
	if !ok {
		identifiable = *new(T)
	}
	return identifiable, identifiable.Identity() != ""
}

func (s *Store[T]) Put(identified T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	id := identified.Identity()
	if _, ok := s.data[id]; !ok {
		s.identifiers = append(s.identifiers, id)
	}
	s.data[id] = identified
}

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

func (s *Store[T]) Delete(identified T) {
	s.DeleteByKey(identified.Identity())
}
