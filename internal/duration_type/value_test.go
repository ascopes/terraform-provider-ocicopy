package duration_type

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"gotest.tools/v3/assert"
)

func TestNewDurationValue(t *testing.T) {
	// Given
	duration := durationOf("5m32s")

	// When
	durationValue := NewDurationValue(duration)

	// Then
	assert.Equal(t, "5m32s", durationValue.ValueString())
	assert.Equal(t, duration, durationValue.ValueDuration())
}

func TestDurationValue_Equal(t *testing.T) {
	tests := []struct {
		first  attr.Value
		second attr.Value
		want   bool
	}{
		{
			first:  NewDurationValue(durationOf("5m23s")),
			second: basetypes.NewStringValue("foobar"),
			want:   false,
		},
		{
			first:  NewDurationValue(durationOf("5m23s")),
			second: NewDurationValue(durationOf("5m23s")),
			want:   true,
		},
		{
			first:  NewDurationValue(durationOf("5m23s")),
			second: NewDurationValue(durationOf("0h5m23s")),
			want:   true,
		},
		{
			first:  NewDurationValue(durationOf("5m23s")),
			second: NewDurationValue(durationOf("0h5m26s")),
			want:   false,
		},
	}

	for _, test := range tests {
		assert.Check(
			t,
			test.first.Equal(test.second) == test.want,
			"Expected %#v{}.Equal(%#v) to be %v, got %v",
			test.first,
			test.second,
			test.want,
			!test.want,
		)
	}
}

func TestDurationValue_Type(t *testing.T) {
	// Given
	duration := durationOf("5m32s")

	// When
	durationValue := NewDurationValue(duration)

	// Then
	assert.Equal(t, DurationType{}, durationValue.Type(context.TODO()))
}

func TestDurationValue_ValueDuration(t *testing.T) {
	// Given
	duration := durationOf("5h2m30s")

	// When
	durationValue := NewDurationValue(duration)

	// Then
	assert.Equal(t, duration, durationValue.ValueDuration())
}

func durationOf(value string) time.Duration {
	parsed, err := time.ParseDuration(value)
	if err != nil {
		panic(err)
	}
	return parsed
}
