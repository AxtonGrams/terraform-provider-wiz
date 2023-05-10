package provider

import (
	"context"
	"reflect"
	"testing"

	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func TestFlattenCloudAccounts(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		map[string]interface{}{
			"id":             "3bc53af8-6661-4a57-a2ee-69a4f2853bba",
			"external_id":    "04b1b4a5-755f-41a9-94ab-6e12173c9b3c",
			"name":           "a6250069-ef06-46ff-bdfe-aea00557f41d",
			"cloud_provider": "0767b7a3-d540-4b9c-8afd-a018aa7da0fb",
			"status":         "9b6e7ae9-e0f6-4748-8171-a6b7a8f385ec",
			"linked_project_ids": []interface{}{
				"3d9ef88a-84f9-4a84-9a67-e5cdd28ad35f",
				"55e9138d-e48f-4155-a2ac-364eb00005db",
			},
			"source_connector_ids": []interface{}{
				"7ac2f620-3882-4c35-91f0-7631eef430c6",
				"dc303b7d-d03d-47f5-9d40-cf16906fb542",
			},
		},
	}

	// Create sample data
	accs := &ReadCloudAccounts{
		CloudAccounts: wiz.CloudAccountConnection{
			Nodes: []*wiz.CloudAccount{
				{
					ID:            "3bc53af8-6661-4a57-a2ee-69a4f2853bba",
					ExternalID:    "04b1b4a5-755f-41a9-94ab-6e12173c9b3c",
					Name:          "a6250069-ef06-46ff-bdfe-aea00557f41d",
					CloudProvider: "0767b7a3-d540-4b9c-8afd-a018aa7da0fb",
					Status:        "9b6e7ae9-e0f6-4748-8171-a6b7a8f385ec",
					LinkedProjects: []*wiz.Project{
						{
							ID: "55e9138d-e48f-4155-a2ac-364eb00005db",
						},
						{
							ID: "3d9ef88a-84f9-4a84-9a67-e5cdd28ad35f",
						},
					},
					SourceConnectors: []wiz.Connector{
						{
							ID: "7ac2f620-3882-4c35-91f0-7631eef430c6",
						},
						{
							ID: "dc303b7d-d03d-47f5-9d40-cf16906fb542",
						},
					},
				},
			}},
	}
	cloudAccounts := make([]interface{}, 0)
	cloudAccounts = append(cloudAccounts, accs)

	flattened := flattenCloudAccounts(ctx, cloudAccounts)

	if !reflect.DeepEqual(flattened, expected) {
		t.Errorf("Unexpected result. Expected: %v, but got: %v", expected, flattened)
	}

}

func TestFlattenProjectIDs(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		"225697ef-8d21-42e1-8195-46d29b285ee6",
		"b0a03462-697e-4ef8-af52-0e8122c6eb7f",
		"d24f22fb-088d-4586-ba8a-9524260f7427",
	}

	var projects = &[]*wiz.Project{
		{
			ID: "b0a03462-697e-4ef8-af52-0e8122c6eb7f",
		},
		{
			ID: "225697ef-8d21-42e1-8195-46d29b285ee6",
		},
		{
			ID: "d24f22fb-088d-4586-ba8a-9524260f7427",
		},
	}

	flattened := flattenProjectIDs(ctx, projects)

	if !reflect.DeepEqual(flattened, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			flattened,
			expected,
		)
	}
}

func TestFlattenSourceConnectorIDs(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		"317d6352-69e0-47e0-a280-76dd3e2e9659",
		"796def2c-70c6-4dc6-85a1-991616a98f4a",
		"d84b87ad-a38f-4ff1-9ee3-761521fbbaab",
	}

	var connectors = &[]wiz.Connector{
		{
			ID: "d84b87ad-a38f-4ff1-9ee3-761521fbbaab",
		},
		{
			ID: "796def2c-70c6-4dc6-85a1-991616a98f4a",
		},
		{
			ID: "317d6352-69e0-47e0-a280-76dd3e2e9659",
		},
	}

	flattened := flattenSourceConnectorIDs(ctx, connectors)

	if !reflect.DeepEqual(flattened, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			flattened,
			expected,
		)
	}
}
