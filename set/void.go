package set

// Type that has zero size that the compiler may be able to just optimise out.
type void struct{}

// A raw void value that has zero size and can be reused globally.
var voidValue = void{}
