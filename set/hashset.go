package set

import (
	"fmt"
	"strings"
)

func NewHashSet[T comparable]() Set[T] {
	return &hashSet[T]{
		data: make(map[T]void, 16),
	}
}

type hashSet[T comparable] struct {
	data map[T]void
}

func (set *hashSet[T]) Add(item T) bool {
	_, exists := set.data[item]
	set.data[item] = voidValue
	return exists
}

func (set *hashSet[T]) Contains(item T) bool {
	_, exists := set.data[item]
	return exists
}

func (set *hashSet[T]) GoString() string {
	return fmt.Sprintf("hashSet{data: %#v}", set.data)
}

func (set *hashSet[T]) Iterator() chan T {
	c := make(chan T)
	go func() {
		defer close(c)
		for item := range set.data {
			c <- item
		}
	}()
	return c
}

func (set *hashSet[T]) Len() int {
	return len(set.data)
}

func (set *hashSet[T]) Remove(item T) bool {
	_, exists := set.data[item]
	if exists {
		delete(set.data, item)
	}
	return exists
}

func (set *hashSet[T]) String() string {
	bldr := strings.Builder{}
	bldr.WriteRune('{')
	i := 0
	for item := range set.Iterator() {
		if i > 0 {
			bldr.WriteString(", ")
		}
		bldr.WriteString(fmt.Sprintf("%#v", item))
		i++
	}
	bldr.WriteRune('}')
	return bldr.String()
}
