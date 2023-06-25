package internal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type apiError struct {
	error
	detail string
}

func newApiError(err error, detail string, args ...any) apiError {
	var parsedDetail string
	if len(args) > 0 {
		parsedDetail = fmt.Sprintf(detail, args...)
	} else {
		parsedDetail = detail
	}
	return apiError{detail: parsedDetail, error: err}
}

func singleApiError(err error, summary string, args ...any) []apiError {
	return []apiError{newApiError(err, summary, args...)}
}

func processApiErrors(diagnostics *diag.Diagnostics, summary string, errors ...apiError) bool {
	for _, err := range errors {
		diagnostics.AddError(summary, err.detail)
	}

	return len(errors) > 0
}
