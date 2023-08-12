package set

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
)

func Test_NewMapHashSet_ReturnsNewMapHashSet(t *testing.T) {
	// When
	set := NewMapHashSet[string]()

	// Then
	assert.Equal(t, 0, set.Len())
}

func Test_MapHashSet_Add_AddsItemToSetWhenNotPresent(t *testing.T) {
	// Given
	set := NewMapHashSet[string]()
	set.Add("foo")
	set.Add("bar")
	set.Add("baz")

	// When
	exists := set.Add("foobarbaz")

	// Then
	assert.Assert(t, !exists, "Expected set.Add(...) to return false, but returned true")
}

func Test_MapHashSet_Add_AddsItemToSetWhenPresent(t *testing.T) {
	// Given
	set := NewMapHashSet[string]()
	set.Add("foo")
	set.Add("bar")
	set.Add("baz")
	set.Add("foobarbaz")

	// When
	exists := set.Add("foobarbaz")

	// Then
	assert.Assert(t, exists, "Expected set.Add(...) to return true, but returned false")
}

func Test_MapHashSet_Contains_ReturnsFalseWhenNotPresent(t *testing.T) {
	// Given
	set := NewMapHashSet[string]()
	set.Add("foo")
	set.Add("bar")
	set.Add("baz")

	// When
	exists := set.Contains("foobarbaz")

	// Then
	assert.Assert(t, !exists, "Expected set.Contains(...) to return false but returned true")
}

func Test_MapHashSet_Contains_ReturnsTrueWhenPresent(t *testing.T) {
	// Given
	set := NewMapHashSet[string]()
	set.Add("foo")
	set.Add("bar")
	set.Add("baz")
	set.Add("foobarbaz")

	// When
	exists := set.Contains("foobarbaz")

	// Then
	assert.Assert(t, exists, "Expected set.Contains(...) to return true but returned false")
}

func Test_MapHashSet_GoString_ReturnsInnerMapGoString(t *testing.T) {
	// Given
	set := NewMapHashSet[string]()
	set.Add("foo")
	set.Add("bar")
	set.Add("baz")
	mapHashSet := set.(*mapHashSet[string])

	// When
	actualStr := set.GoString()

	// Then
	assert.Equal(t, fmt.Sprintf("mapHashSet{data: %#v}", mapHashSet.data), actualStr)
}

func Test_MapHashSet_Iterator_YieldsItemsInSet(t *testing.T) {
	// Given
	set := NewMapHashSet[string]()
	set.Add("foo")
	set.Add("bar")
	set.Add("baz")

	expectedItems := []string{"foo", "bar", "baz"}

	// When
	c := set.Iterator()
	actualItems := make([]string, 0)
	for item := range c {
		actualItems = append(actualItems, item)
	}

	// Then
	assert.Equal(t, len(expectedItems), len(actualItems))
	for _, actualItem := range actualItems {
		found := false
		for _, expectedItem := range expectedItems {
			if expectedItem == actualItem {
				found = true
				break
			}
		}

		assert.Assert(t, found, "Expected %#v to be yielded by iterator, but it was not", actualItem)
	}
}

func Test_MapHashSet_Len_ReturnsExpectedValue(t *testing.T) {
	for i := 0; i < 15; i++ {
		// Given
		set := NewMapHashSet[int]()

		for j := 0; j < i; j++ {
			set.Add(j)
		}

		// Then
		assert.Equal(t, set.Len(), i)
	}
}

func Test_MapHashSet_Remove_ReturnsFalseWhenItemDoesNotExist(t *testing.T) {
	// Given
	set := NewMapHashSet[string]()
	set.Add("foo")
	set.Add("bar")
	set.Add("baz")

	// When
	existed := set.Remove("bork")

	// Then
	assert.Assert(t, !existed, "Expected .Remove(...) to return false but returned true")
	assert.Assert(t, !set.Contains("bork"), "Expected .Contains(\"bork\") to return false but returned true")
}

func Test_MapHashSet_Remove_ReturnsTrueWhenItemExists(t *testing.T) {
	// Given
	set := NewMapHashSet[string]()
	set.Add("foo")
	set.Add("bar")
	set.Add("baz")
	set.Add("bork")

	// When
	existed := set.Remove("bork")

	// Then
	assert.Assert(t, existed, "Expected .Remove(...) to return true but returned false")
	assert.Assert(t, !set.Contains("bork"), "Expected .Contains(\"bork\") to return false but returned true")
}

func Test_MapHashSet_String_ReturnsExpectedValue(t *testing.T) {
	// Given
	set := NewMapHashSet[int]()
	set.Add(1366)
	set.Add(768)
	set.Add(1266)

	// When
	actualStr := set.String()
	wantedStr := set.GoString()

	// Then
	assert.Equal(t, wantedStr, actualStr)
}
