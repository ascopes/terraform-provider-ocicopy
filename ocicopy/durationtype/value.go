package durationtype

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

func (duration DurationValue) Equal(other attr.Value) bool {
	that, ok := other.(DurationValue)
	if !ok {
		return false
	}

	return duration.durationValue.String() == that.durationValue.String()
}

func (duration DurationValue) Type(context.Context) attr.Type {
	return DurationType{}
}

func (duration DurationValue) ValueDuration() time.Duration {
	return duration.durationValue
}
