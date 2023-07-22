package config_test

import (
	"testing"

	"github.com/ascopes/terraform-provider-ocicopy/internal/config"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestDurationTypeEqual(t *testing.T) {
	type test struct {
		input attr.Type
		want  bool
	}

	tests := []test{
		{input: types.StringType, want: false},
		{input: types.Int64Type, want: false},
		{input: types.Float64Type, want: false},
		{input: types.BoolType, want: false},
		{input: types.SetType{}, want: false},
		{input: types.ListType{}, want: false},
		{input: types.MapType{}, want: false},
		{input: types.ObjectType{}, want: false},
		{input: config.DurationType{}, want: true},
	}

	for _, tst := range tests {
		if (config.DurationType{}).Equal(tst.input) != tst.want {
			t.Fatalf("Expected DurationType{}.Equal(%T) to be %v, got %v", tst.input, tst.want, !tst.want)
		}
	}
}

func TestDurationTypeType(t *testing.T) {
	input := config.DurationType{}
	want := "Duration"
	got := input.String()

	if got != want {
		t.Fatalf("Expected DurationType.String() to be equal to '%s', got '%s'", got, want)
	}
}
