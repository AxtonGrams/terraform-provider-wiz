package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceWizReportGraphQuery_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix(ResourcePrefix)
	projectID := os.Getenv("WIZ_PROJECT_ID")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TestCase(TcReportGraphQuery)) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceWizReportGraphQueryBasic(rName, projectID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"wiz_report_graph_query.foo",
						"name",
						rName,
					),
					resource.TestCheckResourceAttr(
						"wiz_report_graph_query.foo",
						"project_id",
						projectID,
					),
					resource.TestCheckResourceAttr(
						"wiz_report_graph_query.foo",
						"query",
						"{\"select\": true, \"type\": [\"CONTAINER_IMAGE\"], \"where\": {\"name\": {\"CONTAINS\": [\"foo\"]}}}",
					),
					resource.TestCheckResourceAttr(
						"wiz_report_graph_query.foo",
						"run_interval_hours",
						"48",
					),
					resource.TestCheckResourceAttr(
						"wiz_report_graph_query.foo",
						"run_starts_at",
						"2023-12-06 16:00:00 +0000 UTC",
					),
				),
			},
		},
	})
}

func testResourceWizReportGraphQueryBasic(rName, projectID string) string {
	return fmt.Sprintf(`
resource "wiz_report_graph_query" "foo" {
  name = "%s"
  project_id = "%s"
  run_interval_hours = 48
  run_starts_at = "2023-12-06 16:00:00 +0000 UTC"
  query = "{\"select\": true, \"type\": [\"CONTAINER_IMAGE\"], \"where\": {\"name\": {\"CONTAINS\": [\"foo\"]}}}"
}
`, rName, projectID)
}
