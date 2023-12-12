package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

var extraConfigErrorSummary = "Invalid extra configuration"

func TestAddFieldError(t *testing.T) {
	diags := diag.Diagnostics{}
	fieldName := "foo"
	keyName := "bar"

	// Test case 1: Invalid extra configuration
	expectedDiags := diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "An issue was encountered while processing the `extraConfig` field.",
			Detail:   "missing or invalid foo field in bar",
		},
	}

	extraConfigErrorSummary = "Invalid extra configuration"
	actualDiags := addFieldError(diags, fieldName, keyName)

	if len(actualDiags) != len(expectedDiags) {
		t.Errorf("Expected %d diagnostics, but got %d", len(expectedDiags), len(actualDiags))
	}

	for i, actualDiag := range actualDiags {
		expectedDiag := expectedDiags[i]

		if actualDiag.Severity != expectedDiag.Severity {
			t.Errorf("Expected severity %v, but got %v", expectedDiag.Severity, actualDiag.Severity)
		}

		if actualDiag.Summary != expectedDiag.Summary {
			t.Errorf("Expected summary %q, but got %q", expectedDiag.Summary, actualDiag.Summary)
		}

		if actualDiag.Detail != expectedDiag.Detail {
			t.Errorf("Expected detail %q, but got %q", expectedDiag.Detail, actualDiag.Detail)
		}
	}
}
