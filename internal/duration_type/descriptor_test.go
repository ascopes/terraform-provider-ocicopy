package duration_type

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"gotest.tools/v3/assert"
)

func TestDurationType_Description(t *testing.T) {
	input := DurationType{}
	len := len(input.Description(context.TODO()))
	assert.Check(t, len > 0, "Empty string returned")
}

func TestDurationType_Equal(t *testing.T) {
	tests := []struct {
		input attr.Type
		want  bool
	}{
		{input: types.StringType, want: false},
		{input: types.Int64Type, want: false},
		{input: types.Float64Type, want: false},
		{input: types.BoolType, want: false},
		{input: types.SetType{}, want: false},
		{input: types.ListType{}, want: false},
		{input: types.MapType{}, want: false},
		{input: types.ObjectType{}, want: false},
		{input: DurationType{}, want: true},
	}

	for _, test := range tests {
		assert.Check(
			t,
			(DurationType{}).Equal(test.input) == test.want,
			"Expected duration_type{}.Equal(%T) to be %v, got %v",
			test.input,
			test.want,
			!test.want,
		)
	}
}

func TestDurationType_MarkdownDescription(t *testing.T) {
	input := DurationType{}
	len := len(input.MarkdownDescription(context.TODO()))
	assert.Check(t, len > 0, "Empty string returned")
}

func TestDurationType_String(t *testing.T) {
	input := DurationType{}
	want := "Duration"
	got := input.String()
	assert.Equal(t, want, got)
}

func TestDurationType_ValueFromTerraform(t *testing.T) {
	ctx := context.TODO()

	t.Run("Fails if value is not a convertable type", func(t *testing.T) {
		// Given
		value := tftypes.NewValue(tftypes.Number, 123)

		// When
		_, err := DurationType{}.ValueFromTerraform(ctx, value)

		// Then
		assert.ErrorContains(t, err, "can't unmarshal tftypes.Number into *string, expected string")
	})

	t.Run("Fails if value is not a valid duration", func(t *testing.T) {
		// Given
		value := tftypes.NewValue(tftypes.String, "500x")

		// When
		_, err := DurationType{}.ValueFromTerraform(ctx, value)

		// Then
		assert.ErrorContains(t, err, "time: unknown unit \"x\" in duration \"500x\"")
	})

	t.Run("Parses valid durations successfully", func(t *testing.T) {
		// Given
		value := tftypes.NewValue(tftypes.String, "0h45m3s")

		// When
		result, err := DurationType{}.ValueFromTerraform(ctx, value)
		castResult := result.(DurationValue)

		// Then
		assert.NilError(t, err, "Expected no error")
		assert.Equal(t, "45m3s", castResult.ValueDuration().String())
	})
}

func TestDurationType_ValueType(t *testing.T) {
	input := DurationType{}
	want := DurationValue{}
	got := input.ValueType(context.TODO())
	assert.Equal(t, want, got)
}
