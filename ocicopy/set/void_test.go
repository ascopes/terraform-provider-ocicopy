package set

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func Test_Void_HasZeroSize(t *testing.T) {
	// When
	size := int(unsafe.Sizeof(voidValue))

	// Then
	assert.Equal(t, 0, size)
}
