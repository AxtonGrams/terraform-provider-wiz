package provider

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"wiz.io/hashicorp/terraform-provider-wiz/internal"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/vendor"
)

func TestFlattenAutomationActionParamsServiceNowTicket(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		map[string]interface{}{
			"base_url":      "https://example.com",
			"client_id":     "a8f5a149-f563-4742-ab81-80e29984325b",
			"client_secret": "8d8b11b7-3706-41bb-a8fc-22c6ac3c17d4",
			"password":      "d974d737-7735-4308-90ca-537fc2dc5bcd",
			"user":          "b243cf01-f38a-4ba0-8c56-2500cedb5639",
			"ticket_fields": []interface{}{
				map[string]interface{}{
					"attach_evidence_csv": utils.ConvertBoolToPointer(false),
					"custom_fields":       `{"testp":"testv"}`,
					"description":         "test description",
					"summary":             "summary",
					"table_name":          "c7c2bba1-9cf2-47f4-9cb6-a213890d1ee7",
				},
			},
		},
	}
	var expanded = &vendor.ServiceNowAutomationActionParams{
		BaseURL:      "https://example.com",
		ClientID:     "a8f5a149-f563-4742-ab81-80e29984325b",
		ClientSecret: "__secret_content__",
		Password:     "__secret_content__",
		User:         "b243cf01-f38a-4ba0-8c56-2500cedb5639",
		TicketFields: vendor.ServiceNowTicketFields{
			AttachEvidenceCSV: utils.ConvertBoolToPointer(false),
			CustomFields:      json.RawMessage(`{"testp":"testv"}`),
			Description:       "test description",
			Summary:           "summary",
			TableName:         "c7c2bba1-9cf2-47f4-9cb6-a213890d1ee7",
		},
	}
	// the provider reads the state for secret information because the api does not return secrets
	// we need to pass the sensitive values to emulate the provider logic
	stateParamsSet := schema.NewSet(
		schema.HashResource(
			&schema.Resource{
				Schema: map[string]*schema.Schema{
					"password": {
						Type:      schema.TypeString,
						Required:  true,
						Sensitive: true,
					},
					"client_secret": {
						Type:      schema.TypeString,
						Optional:  true,
						Sensitive: true,
					},
				},
			},
		),
		[]interface{}{
			map[string]interface{}{
				"password":      "d974d737-7735-4308-90ca-537fc2dc5bcd",
				"client_secret": "8d8b11b7-3706-41bb-a8fc-22c6ac3c17d4",
			},
		},
	)

	automationActionParams := flattenAutomationActionParams(ctx, stateParamsSet, "SERVICENOW_TICKET", expanded)
	if !reflect.DeepEqual(automationActionParams, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			automationActionParams,
			expected,
		)
	}
}

func TestFlattenAutomationActionParamsServiceNowUpdateTicket(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		map[string]interface{}{
			"base_url":      "https://example.com",
			"client_id":     "a8f5a149-f563-4742-ab81-80e29984325b",
			"client_secret": "8d8b11b7-3706-41bb-a8fc-22c6ac3c17d4",
			"fields":        `{"testp":"testv"}`,
			"password":      "d974d737-7735-4308-90ca-537fc2dc5bcd",
			"table_name":    "c7c2bba1-9cf2-47f4-9cb6-a213890d1ee7",
			"user":          "b243cf01-f38a-4ba0-8c56-2500cedb5639",
		},
	}
	var expanded = &vendor.ServiceNowUpdateTicketAutomationActionParams{
		BaseURL:      "https://example.com",
		ClientID:     "a8f5a149-f563-4742-ab81-80e29984325b",
		ClientSecret: "__secret_content__",
		Fields:       json.RawMessage(`{"testp":"testv"}`),
		Password:     "__secret_content__",
		TableName:    "c7c2bba1-9cf2-47f4-9cb6-a213890d1ee7",
		User:         "b243cf01-f38a-4ba0-8c56-2500cedb5639",
	}

	// the provider reads the state for secret information because the api does not return secrets
	// we need to pass the sensitive values to emulate the provider logic
	stateParamsSet := schema.NewSet(
		schema.HashResource(
			&schema.Resource{
				Schema: map[string]*schema.Schema{
					"password": {
						Type:      schema.TypeString,
						Required:  true,
						Sensitive: true,
					},
					"client_secret": {
						Type:      schema.TypeString,
						Optional:  true,
						Sensitive: true,
					},
				},
			},
		),
		[]interface{}{
			map[string]interface{}{
				"password":      "d974d737-7735-4308-90ca-537fc2dc5bcd",
				"client_secret": "8d8b11b7-3706-41bb-a8fc-22c6ac3c17d4",
			},
		},
	)

	automationActionParams := flattenAutomationActionParams(ctx, stateParamsSet, "SERVICENOW_UPDATE_TICKET", expanded)
	if !reflect.DeepEqual(automationActionParams, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			automationActionParams,
			expected,
		)
	}
}

func TestFlattenAutomationActionParamsWebhookBasic(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		map[string]interface{}{
			"url":                "https://example.com",
			"body":               "test123",
			"client_certificate": "204b37d9-0c96-4bc4-9895-e19b467632b6",
			"auth_username":      "1bdea88d-7628-45ff-968b-0672519f9bfd",
			"auth_password":      "a6903d88-0bc7-43fc-a651-7a67b481c54b",
		},
	}
	var expanded = &vendor.WebhookAutomationActionParams{
		URL:               "https://example.com",
		Body:              "test123",
		ClientCertificate: "204b37d9-0c96-4bc4-9895-e19b467632b6",
		AuthenticationType: internal.EnumType{
			Type: "WebhookAutomationActionAuthenticationBasic",
		},
		Authentication: vendor.WebhookAutomationActionAuthenticationBasic{
			Username: "1bdea88d-7628-45ff-968b-0672519f9bfd",
			Password: "__secret_content__",
		},
	}
	// the provider reads the state for secret information because the api does not return secrets
	// we need to pass the sensitive values to emulate the provider logic
	stateParamsSet := schema.NewSet(
		schema.HashResource(
			&schema.Resource{
				Schema: map[string]*schema.Schema{
					"auth_password": {
						Type:      schema.TypeString,
						Required:  true,
						Sensitive: true,
					},
				},
			},
		),
		[]interface{}{
			map[string]interface{}{
				"auth_password": "a6903d88-0bc7-43fc-a651-7a67b481c54b",
			},
		},
	)

	automationActionParams := flattenAutomationActionParams(ctx, stateParamsSet, "WEBHOOK", expanded)
	if !reflect.DeepEqual(automationActionParams, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			automationActionParams,
			expected,
		)
	}
}

func TestFlattenAutomationActionParamsWebhookTokenBearer(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		map[string]interface{}{
			"url":                "https://example.com",
			"body":               "test123",
			"client_certificate": "204b37d9-0c96-4bc4-9895-e19b467632b6",
			"auth_token":         "e8f33c93-bd56-404a-85a4-d1094be4562d",
		},
	}
	var expanded = &vendor.WebhookAutomationActionParams{
		URL:               "https://example.com",
		Body:              "test123",
		ClientCertificate: "204b37d9-0c96-4bc4-9895-e19b467632b6",
		AuthenticationType: internal.EnumType{
			Type: "WebhookAutomationActionAuthenticationTokenBearer",
		},
		Authentication: vendor.WebhookAutomationActionAuthenticationTokenBearer{
			Token: "__secret_content__",
		},
	}
	// we need to pass schema.Set for the sensitive values
	// the provider reads the state for secret information because the api does not return secrets
	// we need to pass the sensitive values to emulate the provider logic
	stateParamsSet := schema.NewSet(
		schema.HashResource(
			&schema.Resource{
				Schema: map[string]*schema.Schema{
					"auth_token": {
						Type:      schema.TypeString,
						Required:  true,
						Sensitive: true,
					},
				},
			},
		),
		[]interface{}{
			map[string]interface{}{
				"auth_token": "e8f33c93-bd56-404a-85a4-d1094be4562d",
			},
		},
	)

	automationActionParams := flattenAutomationActionParams(ctx, stateParamsSet, "WEBHOOK", expanded)
	if !reflect.DeepEqual(automationActionParams, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			automationActionParams,
			expected,
		)
	}
}

func TestFlattenAutomationActionParamsWebhookNoAuth(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		map[string]interface{}{
			"url":                "https://example.com",
			"body":               "test123",
			"client_certificate": "204b37d9-0c96-4bc4-9895-e19b467632b6",
		},
	}
	var expanded = &vendor.WebhookAutomationActionParams{
		URL:               "https://example.com",
		Body:              "test123",
		ClientCertificate: "204b37d9-0c96-4bc4-9895-e19b467632b6",
	}
	// we need to pass schema.Set for the sensitive values
	var localSchemas []interface{}
	stateParamsSet := schema.NewSet(schema.HashString, localSchemas)
	automationActionParams := flattenAutomationActionParams(ctx, stateParamsSet, "WEBHOOK", expanded)
	if !reflect.DeepEqual(automationActionParams, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			automationActionParams,
			expected,
		)
	}
}

func TestGetServicenowParams(t *testing.T) {
	ctx := context.Background()

	var expected = &vendor.CreateServiceNowAutomationActionParamInput{
		BaseURL:      "https://example.com",
		User:         "e0fe3fd4-0ad5-4e83-9814-17b7907aa860",
		Password:     "b6c58df3-8e35-4532-914b-99e9a332c252",
		ClientID:     "d3b04157-738a-412c-9a82-c2d398b7fd0b",
		ClientSecret: "7d286055-c62a-432a-9637-c34511fb90a7",
		TicketFields: vendor.CreateServiceNowFieldsInput{
			TableName:         "b2bc2894-b8a4-4e95-8ffe-ad06c9e1398b",
			Summary:           "Test Summary",
			Description:       "Description",
			AttachEvidenceCSV: utils.ConvertBoolToPointer(true),
			CustomFields:      json.RawMessage(`{"test":"test"}`),
		},
	}

	d := schema.TestResourceDataRaw(
		t,
		resourceWizAutomationAction().Schema,
		map[string]interface{}{
			"name":                          "test",
			"is_accessible_to_all_projects": true,
			"servicenow_params": []interface{}{
				map[string]interface{}{
					"base_url":      "https://example.com",
					"user":          "e0fe3fd4-0ad5-4e83-9814-17b7907aa860",
					"password":      "b6c58df3-8e35-4532-914b-99e9a332c252",
					"client_id":     "d3b04157-738a-412c-9a82-c2d398b7fd0b",
					"client_secret": "7d286055-c62a-432a-9637-c34511fb90a7",
					"ticket_fields": []interface{}{
						map[string]interface{}{
							"table_name":          "b2bc2894-b8a4-4e95-8ffe-ad06c9e1398b",
							"summary":             "Test Summary",
							"description":         "Description",
							"attach_evidence_csv": true,
							"custom_fields":       `{"test":"test"}`,
						},
					},
				},
			},
		},
	)

	automationActionParams := getServicenowParams(ctx, d.Get("servicenow_params"))

	if !reflect.DeepEqual(automationActionParams, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			automationActionParams,
			expected,
		)
	}
}

func TestGetServicenowUpdateTicketParams(t *testing.T) {
	ctx := context.Background()

	var expected = &vendor.CreateServiceNowUpdateTicketAutomationActionParamInput{
		BaseURL:      "https://example.com",
		User:         "e0fe3fd4-0ad5-4e83-9814-17b7907aa860",
		Password:     "b6c58df3-8e35-4532-914b-99e9a332c252",
		TableName:    "b2bc2894-b8a4-4e95-8ffe-ad06c9e1398b",
		Fields:       json.RawMessage(`{"test":"test"}`),
		ClientID:     "d3b04157-738a-412c-9a82-c2d398b7fd0b",
		ClientSecret: "7d286055-c62a-432a-9637-c34511fb90a7",
	}

	d := schema.TestResourceDataRaw(
		t,
		resourceWizAutomationAction().Schema,
		map[string]interface{}{
			"name":                          "test",
			"is_accessible_to_all_projects": true,
			"servicenow_update_ticket_params": []interface{}{
				map[string]interface{}{
					"base_url":      "https://example.com",
					"user":          "e0fe3fd4-0ad5-4e83-9814-17b7907aa860",
					"password":      "b6c58df3-8e35-4532-914b-99e9a332c252",
					"table_name":    "b2bc2894-b8a4-4e95-8ffe-ad06c9e1398b",
					"fields":        `{"test":"test"}`,
					"client_id":     "d3b04157-738a-412c-9a82-c2d398b7fd0b",
					"client_secret": "7d286055-c62a-432a-9637-c34511fb90a7",
				},
			},
		},
	)

	automationActionParams := getServicenowUpdateTicketParams(ctx, d.Get("servicenow_update_ticket_params"))

	if !reflect.DeepEqual(automationActionParams, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			automationActionParams,
			expected,
		)
	}
}

func TestGetWebhookAutomationActionParams(t *testing.T) {
	ctx := context.Background()

	var expected = &vendor.CreateWebhookAutomationActionParamsInput{
		URL:               "https://example.com",
		Body:              "326252d8-2191-4ccb-b260-9f6918f7c522",
		ClientCertificate: "57178574-db59-4a50-8da5-63373de2711a",
		AuthUsername:      "e0fe3fd4-0ad5-4e83-9814-17b7907aa860",
		AuthPassword:      "854ff4cd-dbdb-48d0-b876-fab1ef02e778",
		AuthToken:         "b6c58df3-8e35-4532-914b-99e9a332c252",
	}

	d := schema.TestResourceDataRaw(
		t,
		resourceWizAutomationAction().Schema,
		map[string]interface{}{
			"name":                          "test",
			"is_accessible_to_all_projects": true,
			"webhook_params": []interface{}{
				map[string]interface{}{
					"url":                "https://example.com",
					"body":               "326252d8-2191-4ccb-b260-9f6918f7c522",
					"client_certificate": "57178574-db59-4a50-8da5-63373de2711a",
					"auth_username":      "e0fe3fd4-0ad5-4e83-9814-17b7907aa860",
					"auth_password":      "854ff4cd-dbdb-48d0-b876-fab1ef02e778",
					"auth_token":         "b6c58df3-8e35-4532-914b-99e9a332c252",
				},
			},
		},
	)

	automationActionParams := getWebhookAutomationActionParams(ctx, d.Get("webhook_params"))

	if !reflect.DeepEqual(automationActionParams, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			automationActionParams,
			expected,
		)
	}
}
