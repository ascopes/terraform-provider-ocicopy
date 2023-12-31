package duration_type

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Custom type that allows using strings to denote arbitrary durations, such
// as "5s", "12ns", "1.5m", "2h", "3m45s", etc.
type DurationType struct {
	basetypes.StringType
}

// Return a plain-text description of the type.
func (duration DurationType) Description(context.Context) string {
	return "String value that represents an arbitrary time duration. " +
		"For example, '5m', '10m20s', '2h44ms', etc."
}

// Determine if another attribute type is compatible with this type.
func (duration DurationType) Equal(other attr.Type) bool {
	_, ok := other.(DurationType)
	return ok
}

// Return a markdown description of this type.
func (duration DurationType) MarkdownDescription(context.Context) string {
	return "String value that represents an arbitrary time duration. " +
		"The value should take the form of one or more `<number><unit>` pairs, " +
		"such as `5m`, `2m30s`, `-35m`, etc.\n" +
		"\n" +
		"Valid units are:\n" +
		"\n" +
		"| Symbol    | Corresponding unit |" +
		"|----------:|:-------------------|" +
		"| `ns`      | Nanoseconds        |" +
		"| `us`      | Microseconds       |" +
		"| `ms`      | Milliseconds       |" +
		"| `s`       | Seconds            |" +
		"| `m`       | Minutes            |" +
		"| `h`       | Hours              |"
}

// Get a string description of this type.
func (duration DurationType) String() string {
	return "Duration"
}

// Parse a Terraform raw value into an instance of this type.
func (duration DurationType) ValueFromTerraform(ctx context.Context, value tftypes.Value) (attr.Value, error) {
	attrValue, err := duration.StringType.ValueFromTerraform(ctx, value)
	if err != nil {
		return nil, err
	}

	stringValue := attrValue.(basetypes.StringValue)
	rawValue := stringValue.ValueString()
	// Remove accidental whitespace (e.g. "5h 3m" should be the same as "5h3m" semantically)
	rawDuration, err := time.ParseDuration(rawValue)

	if err != nil {
		return nil, err
	}

	durationValue := DurationValue{StringValue: stringValue, durationValue: rawDuration}
	return durationValue, nil
}

// Get a new value holder for this type.
func (DurationType) ValueType(context.Context) attr.Value {
	return DurationValue{}
}
