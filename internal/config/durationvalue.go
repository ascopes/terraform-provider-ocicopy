package config

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type DurationValue struct {
	basetypes.StringValue

	durationValue time.Duration
}

func (v DurationValue) Equal(other attr.Value) bool {
	that, ok := other.(DurationValue)
	if !ok {
		return false
	}

	return v.durationValue.String() == (that.durationValue.String())
}

func (DurationValue) Type(context.Context) attr.Type {
	return DurationType{}
}

func (v DurationValue) ValueDuration() time.Duration {
	return v.durationValue
}
