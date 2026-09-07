package store

import (
	"context"
	"slices"
)

type Identifiable interface {
	Identity() string
}

type opType byte

const (
	noop = iota
	opRead
	opList
	opWrite
	opDelete
)

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

type operation[T Identifiable] struct {
	identifier string
	page       pageRange
	opType     opType
	identified T
}

type Store[T Identifiable] struct {
	data           map[string]T
	identifiers    []string
	semaphore      chan operation[T]
	singleSupplier chan T
	manySupplier   chan []T
}

func NewStoreWithData[T Identifiable](ctx context.Context, keyed ...T) *Store[T] {
	s := NewStore[T](ctx)
	for _, keyed := range keyed {
		s.Put(keyed)
	}
	return s
}

func NewStore[T Identifiable](ctx context.Context) *Store[T] {
	s := &Store[T]{
		data:           make(map[string]T),
		identifiers:    make([]string, 0),
		semaphore:      make(chan operation[T]),
		singleSupplier: make(chan T),
		manySupplier:   make(chan []T),
	}

	go s.start(ctx)

	return s
}

func (s *Store[T]) start(ctx context.Context) {
	for {
		select {
		case action := <-s.semaphore:
			switch action.opType {
			case opRead:
				keyed, ok := s.data[action.identifier]
				if !ok {
					keyed = *new(T)
				}
				go func() {
					s.singleSupplier <- keyed
				}()
			case opList:
				go func() {
					r := make([]T, 0)
					from := action.page.fromIndex()
					to := action.page.toIndex(len(s.identifiers))
					for _, k := range s.identifiers[from:to] {
						r = append(r, s.data[k])
					}
					s.manySupplier <- r
				}()
			case opWrite:
				s.identifiers = append(s.identifiers, action.identified.Identity())
				s.data[action.identified.Identity()] = action.identified
			case opDelete:
				if _, ok := s.data[action.identifier]; ok {
					index := slices.Index(s.identifiers, action.identifier)
					s.identifiers = slices.Delete(s.identifiers, index, index+1)
				}
				delete(s.data, action.identifier)
			case noop:
				// no op!
			}

		case <-ctx.Done():
			return
		}
	}
}

func (s *Store[T]) GetAll() []T {
	s.do(operation[T]{
		opType: opList,
		page:   pageRange{page: 0, size: len(s.identifiers)},
	})
	return <-s.manySupplier
}

func (s *Store[T]) GetPage(page int, size int) []T {
	s.do(operation[T]{
		opType: opList,
		page:   pageRange{page: page, size: size},
	})
	return <-s.manySupplier
}

func (s *Store[T]) Get(key string) (T, bool) {
	s.do(operation[T]{
		opType:     opRead,
		identifier: key})
	data := <-s.singleSupplier
	return data, data.Identity() != ""
}

func (s *Store[T]) Put(identified T) {
	s.do(operation[T]{
		opType:     opWrite,
		identified: identified,
	})
}

func (s *Store[T]) DeleteByKey(identifier string) {
	s.do(operation[T]{
		opType:     opDelete,
		identifier: identifier})
}

func (s *Store[T]) Delete(identied T) {
	s.DeleteByKey(identied.Identity())
}

func (s *Store[T]) do(op operation[T]) {
	s.semaphore <- op
}
