package set

import (
	"fmt"
	"reflect"
	"sort"
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

	// Obtain a sorted slice as map key iteration order can differ
	// between calls, which makes this output difficult to predict
	// and even harder to actually test.
	keys := reflect.ValueOf(myMap).MapKeys()
	sort.Slice(keys, func(a int, b int) bool {
		return keys[a].Interface().(T) < keys[b].Interface().(T)
	})

	for i, item := range keys {
		if i > 0 {
			writer.WriteString(", ")
		}
		writer.WriteString(fmt.Sprintf("%#v", item))
	}
	
	bldr.WriteRune('}')
	return bldr.String()
}
