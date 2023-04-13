package provider

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func TestGetSecuritySubCategories(t *testing.T) {
	ctx := context.Background()

	var expectedSubCategory1 = wiz.SecuritySubCategoryInput{
		ID:          "efd96f85-293a-462c-8f38-d2731d93db3d",
		Title:       "03371981-bceb-42f2-8bb7-61a4b1408b64",
		Description: "5c53a6be-75ff-453f-9ba6-3464de7fab42",
	}

	var expectedSubCategory2 = wiz.SecuritySubCategoryInput{
		ID:          "7622db38-8b5c-44ba-a766-2d8d48a63eb6",
		Title:       "337b63d5-b553-4802-80da-f20ea10de803",
		Description: "31d44de7-ce27-4361-8c3f-5329bb9a8231",
	}

	var expected = []wiz.SecuritySubCategoryInput{}

	expected = append(expected, expectedSubCategory1)
	expected = append(expected, expectedSubCategory2)

	d := schema.NewSet(
		schema.HashResource(
			&schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:        schema.TypeString,
						Description: "Internal identifier for the security subcategory. Specify an existing identifier to use an existing subcategory. If not provided, a new subcategory will be created.",
						Optional:    true,
						Computed:    true,
					},
					"title": {
						Type:        schema.TypeString,
						Description: "Title of the security subcategory.",
						Required:    true,
					},
					"description": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Description of the security subcategory.",
					},
				},
			},
		),
		[]interface{}{
			map[string]interface{}{
				"id":          "efd96f85-293a-462c-8f38-d2731d93db3d",
				"title":       "03371981-bceb-42f2-8bb7-61a4b1408b64",
				"description": "5c53a6be-75ff-453f-9ba6-3464de7fab42",
			},
			map[string]interface{}{
				"id":          "7622db38-8b5c-44ba-a766-2d8d48a63eb6",
				"title":       "337b63d5-b553-4802-80da-f20ea10de803",
				"description": "31d44de7-ce27-4361-8c3f-5329bb9a8231",
			},
		},
	)

	subCategories := getSecuritySubCategories(ctx, d)

	if !reflect.DeepEqual(expected, subCategories) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			utils.PrettyPrint(subCategories),
			utils.PrettyPrint(expected),
		)
	}
}

func testGetSecurityCategories(t *testing.T) {
	ctx := context.Background()

	var expectedCategory1 = wiz.SecurityCategoryInput{
		ID:          "547e6ca8-8caf-4229-a232-e128939dec52",
		Name:        "9872a0e0-3861-47bb-8faf-8093507a5cab",
		Description: "465e9993-dd6e-4902-8e1b-c0c54373623b",
		SubCategories: []wiz.SecuritySubCategoryInput{
			{
				ID:          "efd96f85-293a-462c-8f38-d2731d93db3d",
				Title:       "03371981-bceb-42f2-8bb7-61a4b1408b64",
				Description: "5c53a6be-75ff-453f-9ba6-3464de7fab42",
			},
			{
				ID:          "7622db38-8b5c-44ba-a766-2d8d48a63eb6",
				Title:       "337b63d5-b553-4802-80da-f20ea10de803",
				Description: "31d44de7-ce27-4361-8c3f-5329bb9a8231",
			},
		},
	}

	var expectedCategory2 = wiz.SecurityCategoryInput{
		ID:          "db64775d-66d5-404d-bf1e-3dc3fbf52b6c",
		Name:        "4d6fa680-96e3-489e-8e41-8603cbf9902a",
		Description: "f70f4eae-f8e0-4771-bbef-e700c0b2a394",
		SubCategories: []wiz.SecuritySubCategoryInput{
			{
				ID:    "9a505144-4819-423a-9375-25d569404c4f",
				Title: "d5785980-30ca-4ae6-9eed-93f78d0ec826",
			},
		},
	}

	var expected = []wiz.SecurityCategoryInput{}
	expected = append(expected, expectedCategory1)
	expected = append(expected, expectedCategory2)

	d := schema.TestResourceDataRaw(
		t,
		resourceWizSecurityFramework().Schema,
		map[string]interface{}{
			"name":        "20615763-2732-4eaf-aa28-9fe1baf008eb",
			"description": "db477e0e-e8c2-498f-991d-6f7ee690a971",
			"enabled":     false,
			"category": []interface{}{
				map[string]interface{}{
					"id":   "547e6ca8-8caf-4229-a232-e128939dec52",
					"name": "9872a0e0-3861-47bb-8faf-8093507a5cab",
					"sub_category": []interface{}{
						map[string]interface{}{
							"id":          "efd96f85-293a-462c-8f38-d2731d93db3d",
							"title":       "03371981-bceb-42f2-8bb7-61a4b1408b64",
							"description": "5c53a6be-75ff-453f-9ba6-3464de7fab42",
						},
						map[string]interface{}{
							"id":          "7622db38-8b5c-44ba-a766-2d8d48a63eb6",
							"title":       "337b63d5-b553-4802-80da-f20ea10de803",
							"description": "31d44de7-ce27-4361-8c3f-5329bb9a8231",
						},
					},
				},
				map[string]interface{}{
					"name": "f9480c26-b70a-440c-9143-a22f9b01254d",
					"sub_category": []interface{}{
						map[string]interface{}{
							"id":    "9a505144-4819-423a-9375-25d569404c4f",
							"title": "d5785980-30ca-4ae6-9eed-93f78d0ec826",
						},
					},
				},
			},
		},
	)

	categories := getSecurityCategories(ctx, d)

	if !reflect.DeepEqual(expected, categories) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			utils.PrettyPrint(categories),
			utils.PrettyPrint(expected),
		)
	}
}

func TestFlattenSecurityCategories(t *testing.T) {
	ctx := context.Background()

	expected := []interface{}{
		map[string]interface{}{
			"id":          "57195246-327a-4266-a5f6-8af81bc5b752",
			"name":        "e6e1b491-ad79-4d9f-9ea2-cae34ab6befa",
			"description": "a9649710-556f-490c-b383-b1324dc21711",
			"sub_category": []interface{}{
				map[string]interface{}{
					"id":          "8453e1b3-cc94-4966-8d72-baec7331e7df",
					"title":       "2ab89df0-c150-49b7-aad3-1a4706a4e613",
					"description": "10a2649b-1539-4465-9bae-ff2b35b9f524",
				},
			},
		},
	}

	var expanded = wiz.SecurityFramework{
		Name:        "5c8268cc-2c41-4f14-a331-5a4dae67b836",
		Description: "abb8b08a-964a-4a7d-8a33-b1d0f2caea35",
		Enabled:     true,
		ID:          "721b490c-31fc-4af1-ba9b-5713c37a569a",
		Categories: []wiz.SecurityCategory{
			{
				ID:          "57195246-327a-4266-a5f6-8af81bc5b752",
				Name:        "e6e1b491-ad79-4d9f-9ea2-cae34ab6befa",
				Description: "a9649710-556f-490c-b383-b1324dc21711",
				SubCategories: []wiz.SecuritySubCategory{
					{
						ID:          "8453e1b3-cc94-4966-8d72-baec7331e7df",
						Title:       "2ab89df0-c150-49b7-aad3-1a4706a4e613",
						Description: "10a2649b-1539-4465-9bae-ff2b35b9f524",
					},
				},
			},
		},
	}

	securityCategory := flattenSecurityCategories(ctx, expanded)

	if !reflect.DeepEqual(securityCategory, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			securityCategory,
			expected,
		)
	}
}

func TestFlattenSecuritySubCategories(t *testing.T) {
	ctx := context.Background()

	expected := []interface{}{
		map[string]interface{}{
			"id":          "8453e1b3-cc94-4966-8d72-baec7331e7df",
			"title":       "2ab89df0-c150-49b7-aad3-1a4706a4e613",
			"description": "10a2649b-1539-4465-9bae-ff2b35b9f524",
		},
	}

	var expanded = []wiz.SecuritySubCategory{
		{
			ID:          "8453e1b3-cc94-4966-8d72-baec7331e7df",
			Title:       "2ab89df0-c150-49b7-aad3-1a4706a4e613",
			Description: "10a2649b-1539-4465-9bae-ff2b35b9f524",
		},
	}

	securitySubCategory := flattenSecuritySubCategories(ctx, expanded)

	if !reflect.DeepEqual(securitySubCategory, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			securitySubCategory,
			expected,
		)
	}
}
