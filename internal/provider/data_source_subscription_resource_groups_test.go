package provider

import (
	"context"
	"reflect"
	"testing"

	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func TestFlattenResourceGroups(t *testing.T) {
	ctx := context.Background()
	resourceGroup := "shared_coconut"

	var results = []*wiz.GraphSearchResult{}

	expected := []interface{}{
		map[string]interface{}{
			"id":   "8f137cbc-0810-55ff-acd6-f7574eb0d071",
			"name": resourceGroup,
		},
	}

	graphEntity := &wiz.GraphEntity{
		ID:   "8f137cbc-0810-55ff-acd6-f7574eb0d071",
		Name: resourceGroup,
	}

	var entities = []wiz.GraphEntity{}
	entities = append(entities, *graphEntity)

	var searchresult = &wiz.GraphSearchResult{
		Entities: entities,
	}
	results = append(results, searchresult)

	flattened := flattenResourceGroups(ctx, &results)

	if !reflect.DeepEqual(flattened, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			flattened,
			expected,
		)
	}

}
