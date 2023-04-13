package provider

import (
	"context"
	"reflect"
	"testing"

	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func TestFlattenTargetPlatformIDs(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		"02ef6af4-f2fe-45ea-a119-a4e4150ffb6c",
		"b488610a-6846-404c-8e5b-28ee04846dda",
		"cdd2c255-921d-4ea9-b348-5660a7b9d459",
	}

	var plats = []wiz.Technology{
		{
			ID: "cdd2c255-921d-4ea9-b348-5660a7b9d459",
		},
		{
			ID: "02ef6af4-f2fe-45ea-a119-a4e4150ffb6c",
		},
		{
			ID: "b488610a-6846-404c-8e5b-28ee04846dda",
		},
	}

	flattened := flattenTargetPlatformIDs(ctx, plats)

	if !reflect.DeepEqual(flattened, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			flattened,
			expected,
		)
	}
}
