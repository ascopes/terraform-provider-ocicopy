package internal

type void = struct{}

var voidValue void

// Make a set from the given slice.
func sliceToSet[T comparable](slice []T) map[T]void {
	mapping := make(map[T]struct{}, len(slice))
	for _, item := range slice {
		mapping[item] = voidValue
	}
	return mapping
}
