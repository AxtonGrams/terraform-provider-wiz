package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceWizReportGraphQuery_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix(ResourcePrefix)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TestCase(TcProject)) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceWizReportGraphQueryBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"wiz_report_graph_query.foo",
						"name",
						rName,
					),
					resource.TestCheckResourceAttr(
						"wiz_report_graph_query.foo",
						"project_id",
						"2c38b8fa-c315-57ea-9de4-e3a19592d796",
					),
					resource.TestCheckResourceAttr(
						"wiz_report_graph_query.foo",
						"query",
						"{\"select\": true, \"type\": [\"CONTAINER_IMAGE\"], \"where\": {\"name\": {\"CONTAINS\": [\"atlantis\"]}}}",
					),
					resource.TestCheckResourceAttr(
						"wiz_report_graph_query.foo",
						"run_interval_hours",
						"48",
					),
					resource.TestCheckResourceAttr(
						"wiz_report_graph_query.foo",
						"run_starts_at",
						"2023-12-06 16:30:00 +0000 UTC",
					),
				),
			},
		},
	})
}

func testResourceWizReportGraphQueryBasic(rName string) string {
	return fmt.Sprintf(`
resource "wiz_report_graph_query" "foo" {
  name = "%s"
  project_id = "6c3858fa-c807-57ea-9de4-d3e19536d796"
  run_interval_hours = 48
  run_starts_at = "2023-12-06 16:30:00 +0000 UTC"
  query = "{\"select\": true, \"type\": [\"CONTAINER_IMAGE\"], \"where\": {\"name\": {\"CONTAINS\": [\"foo\"]}}}"
}
`, rName)
}
