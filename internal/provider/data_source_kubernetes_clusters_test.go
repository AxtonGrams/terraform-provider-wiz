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
					"name":           "AWS",
				},
			},
		},
		map[string]interface{}{
			"id":   "8f137cbc-0810-55ff-acd6-f7574eb0d072",
			"name": "24x7-prod",
			"cloud_account": []interface{}{
				map[string]interface{}{
					"cloud_provider": "AZURE",
					"external_id":    "31e5304c-baca-54fa-bff3-aaf493eaeae2",
					"id":             "31e5304c-baca-54fa-bff3-aaf493eaeae0",
					"name":           "AZURE",
				},
			},
		},
	}

	clusters := &ReadKubernetesClusters{
		KubernetesClusters: wiz.KubernetesClusterConnection{
			Nodes: []*wiz.KubernetesCluster{
				{
					ID:   "8f137cbc-0810-55ff-acd6-f7574eb0d071",
					Name: "24x7-dev",
					CloudAccount: wiz.CloudAccount{
						ID:            "31e5304c-baca-54fa-bff3-aaf493eaeae0",
						ExternalID:    "151668690081",
						CloudProvider: "AWS",
						Name:          "AWS",
					},
				},
				{
					ID:   "8f137cbc-0810-55ff-acd6-f7574eb0d072",
					Name: "24x7-prod",
					CloudAccount: wiz.CloudAccount{
						ID:            "31e5304c-baca-54fa-bff3-aaf493eaeae0",
						ExternalID:    "31e5304c-baca-54fa-bff3-aaf493eaeae2",
						CloudProvider: "AZURE",
						Name:          "AZURE",
					},
				},
			},
		}}

	clusterLinks := make([]interface{}, 0)
	clusterLinks = append(clusterLinks, clusters)

	flattened := flattenClusters(ctx, clusterLinks)

	if !reflect.DeepEqual(flattened, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			flattened,
			expected,
		)
	}
}
