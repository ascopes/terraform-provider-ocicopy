package set

import (
	"fmt"
)

func NewMapHashSet[T comparable]() Set[T] {
	return &mapHashSet[T]{
		data: make(map[T]void, 16),
	}
}

type mapHashSet[T comparable] struct {
	data map[T]void
}

func (set *mapHashSet[T]) Add(item T) bool {
	_, exists := set.data[item]
	set.data[item] = voidValue
	return exists
}

func (set *mapHashSet[T]) Contains(item T) bool {
	_, exists := set.data[item]
	return exists
}

func (set *mapHashSet[T]) GoString() string {
	return fmt.Sprintf("mapHashSet{data: %#v}", set.data)
}

func (set *mapHashSet[T]) Iterator() chan T {
	c := make(chan T)
	go func() {
		defer close(c)
		for item := range set.data {
			c <- item
		}
	}()
	return c
}

func (set *mapHashSet[T]) Len() int {
	return len(set.data)
}

func (set *mapHashSet[T]) Remove(item T) bool {
	_, exists := set.data[item]
	if exists {
		delete(set.data, item)
	}
	return exists
}

func (set *mapHashSet[T]) String() string {
	return set.GoString()
}
