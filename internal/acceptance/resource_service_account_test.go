package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func resourceWizServiceAccountCheckHelper(rName string, rType string) resource.TestCheckFunc {
	commonChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			// check that type is set to the provided type
			"wiz_service_account.test",
			"type",
			rType,
		),
		resource.TestCheckResourceAttr(
			// check that name is set to the provided name
			"wiz_service_account.test",
			"name",
			rName,
		),
		resource.TestMatchResourceAttr(
			// check that client_id is set to a non-empty string
			"wiz_service_account.test",
			"client_id",
			regexp.MustCompile(`\w`),
		),
		resource.TestMatchResourceAttr(
			// check that client_secret is set to a non-empty string
			"wiz_service_account.test",
			"client_secret",
			regexp.MustCompile(`\w`),
		),
		resource.TestMatchResourceAttr(
			// check that created_at matches UTC datetime format
			"wiz_service_account.test",
			"created_at",
			regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{1,}Z`),
		),
		resource.TestMatchResourceAttr(
			// check that last_rotated_at matches UTC datetime format
			"wiz_service_account.test",
			"last_rotated_at",
			regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{1,}Z`),
		),
		resource.TestMatchResourceAttr(
			// check that recreate_if_rotated is a boolean
			"wiz_service_account.test",
			"recreate_if_rotated",
			regexp.MustCompile(`true|false`),
		),
	}

	// Include additional checks if type is THIRD_PARTY
	if rType == "THIRD_PARTY" {
		commonChecks = append(commonChecks,
			resource.TestCheckResourceAttrSet(
				// check for a set assigned project id
				"wiz_service_account.test",
				"assigned_projects.0",
			),
			resource.TestCheckResourceAttrSet(
				// check for a 2nd value in the list of scopes
				"wiz_service_account.test",
				"scopes.1",
			))
	}

	return resource.ComposeTestCheckFunc(commonChecks...)
}

func TestAccResourceWizServiceAccount_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix(ResourcePrefix)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TestCase(TcCommon)) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceWizServiceAccountBasic(rName, "THIRD_PARTY"), // the default type for GRAPHQL service account are THIRD_PARTY
				Check:  resourceWizServiceAccountCheckHelper(rName, "THIRD_PARTY"),
			},
			{
				Config: testResourceWizServiceAccountBasic(rName, "KUBERNETES_ADMISSION_CONTROLLER"),
				Check:  resourceWizServiceAccountCheckHelper(rName, "KUBERNETES_ADMISSION_CONTROLLER"),
			},
			{
				Config: testResourceWizServiceAccountBasic(rName, "BROKER"),
				Check:  resourceWizServiceAccountCheckHelper(rName, "BROKER"),
			},
		},
	})
}

func testResourceWizServiceAccountBasic(rName string, rType string) string {
	switch rType {
	// THIRD_PARTY service accounts require scopes and can accept assigned_projects
	case "THIRD_PARTY":
		project := uuid.New().String()
		return fmt.Sprintf(`
			resource "wiz_service_account" "test" {
				name                = "%s"
				scopes              = ["read:service_accounts", "create:service_accounts"]
				assigned_projects = [ "%s" ]
			  }

		`, rName, project)
	default:
		return fmt.Sprintf(`
				resource "wiz_service_account" "test" {
					name                = "%s"
					type                = "%s"
				  }

			`, rName, rType)
	}
}
