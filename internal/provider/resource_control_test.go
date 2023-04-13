package provider

import (
	"context"
	"reflect"
	"testing"

	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func TestFlattenControlSecuritySubCategories(t *testing.T) {
	ctx := context.Background()

	expected := []interface{}{
		"6b5d5a05-1186-4f70-ae0c-bde55cc9e6aa",
		"33bf37f5-9d7e-4e0e-a081-ca362a2223b5",
	}

	var expanded = []*wiz.SecuritySubCategory{
		{
			ID: "6b5d5a05-1186-4f70-ae0c-bde55cc9e6aa",
		},
		{
			ID: "33bf37f5-9d7e-4e0e-a081-ca362a2223b5",
		},
	}

	ssc := flattenControlSecuritySubCategories(ctx, expanded)

	if !reflect.DeepEqual(ssc, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			ssc,
			expected,
		)
	}
}
