package provider

import (
	"context"
	"reflect"
	"testing"

	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func TestFlattenKubernetesClusters(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		map[string]interface{}{
			"id":   "8f137cbc-0810-55ff-acd6-f7574eb0d071",
			"name": "24x7-dev",
			"cloud_account": []interface{}{
				map[string]interface{}{
					"cloud_provider": "AWS",
					"external_id":    "151668690081",
					"id":             "31e5304c-baca-54fa-bff3-aaf493eaeae0",
					"name":           "MYK8S_TENANT_STAGE",
				},
			},
		},
	}

	var clusterLinks = &[]*wiz.KubernetesCluster{
		{
			ID:   "8f137cbc-0810-55ff-acd6-f7574eb0d071",
			Name: "24x7-dev",
			CloudAccount: *&wiz.CloudAccount{
				ID:            "31e5304c-baca-54fa-bff3-aaf493eaeae0",
				ExternalID:    "151668690081",
				CloudProvider: "AWS",
				Name:          "MYK8S_TENANT_STAGE",
			},
		},
	}

	flattened := flattenClusters(ctx, clusterLinks)

	if !reflect.DeepEqual(flattened, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			flattened,
			expected,
		)
	}
}
