package set

import "fmt"

// Interface for a set implementation.
//
// The type parameter T may have additional constraints in any implementation.
//
// There is zero guarantee about any kind of thread-safety or safety between
// Goroutines. In addition, concurrent modification is not supported, including
// during iteration.
type Set[T any] interface {
	fmt.Stringer
	fmt.GoStringer

	// Add an item to the set.
	//
	// Return true if the item already existed and was replaced, or false if it was a new item.
	Add(item T) bool

	// Determine if the set contains the current item.
	//
	// Return true if it was present or false if it was not present.
	Contains(item T) bool

	// Create an iterator across this set that can be used in `for` expressions.
	//
	// A channel that will emit items in this set in an arbitrary order.
	Iterator() chan T

	// Get the length of the set.
	//
	// Return the number of items in the set.
	Len() int

	// Remove an item from the set.
	//
	// Return true if the item was present and was removed, or false if it was not present.
	Remove(item T) bool
}
