package mapping

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
)

// Check that the value is neither null nor unknown, and then consume it.
func IfPresent[T attr.Value](value T, then func(T)) {
	if !value.IsNull() && !value.IsUnknown() {
		then(value)
	}
}

// Check that the value is not nil, then consume it.
func IfNotNil[T any](value *T, then func(*T)) {
	if value != nil {
		then(value)
	}
}
