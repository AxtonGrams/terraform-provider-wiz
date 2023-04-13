package provider

import (
	"context"
	"reflect"
	"testing"

	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func TestFlattenAssignedProjectIDs(t *testing.T) {
	ctx := context.Background()

	expected := []interface{}{
		"2dc9a5ee-b52e-41a2-a13f-75c57d466acf",
		"bc0dc093-e74e-4eea-9734-e3e5cfe1ecab",
	}

	var expanded = []wiz.Project{
		{
			ID: "2dc9a5ee-b52e-41a2-a13f-75c57d466acf",
		},
		{
			ID: "bc0dc093-e74e-4eea-9734-e3e5cfe1ecab",
		},
	}

	assignedProjectIDs := flattenAssignedProjectIDs(ctx, expanded)

	if !reflect.DeepEqual(assignedProjectIDs, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			assignedProjectIDs,
			expected,
		)
	}
}
