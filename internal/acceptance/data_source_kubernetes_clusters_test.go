package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDatasourceWizKubernetesClusters_basic tests the basic functionality of the datasource
// wiz_kubernetes_clusters. The assumption is that at least two clusters exist in the Wiz tenant in order
// to validate pagination functionality
func TestAccDatasourceWizKubernetesClusters_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TestCase(TcCommon)) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceWizKubernetesClustersBasic(1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						// check first kubernetes cluster has an id that matches the UUID regex
						"data.wiz_kubernetes_clusters.foo",
						"kubernetes_clusters.0.id",
						regexp.MustCompile(UUIDPattern),
					),
					resource.TestMatchResourceAttr(
						// check first kubernetes cluster has a name that is set to a non-empty string
						"data.wiz_kubernetes_clusters.foo",
						"kubernetes_clusters.0.name",
						regexp.MustCompile(`\w`),
					),
					resource.TestMatchResourceAttr(
						// check cloud_account_block has id that matches the UUID regex
						"data.wiz_kubernetes_clusters.foo",
						"kubernetes_clusters.0.cloud_account.0.id",
						regexp.MustCompile(UUIDPattern),
					),
					resource.TestMatchResourceAttr(
						// check cloud_account_block has external_id set to a non-empty string, different cloud providers have different formats
						"data.wiz_kubernetes_clusters.foo",
						"kubernetes_clusters.0.cloud_account.0.external_id",
						regexp.MustCompile(`\w`),
					),
					resource.TestMatchResourceAttr(
						// check cloud_account_block has name set to a non-empty string
						"data.wiz_kubernetes_clusters.foo",
						"kubernetes_clusters.0.cloud_account.0.name",
						regexp.MustCompile(`\w`),
					),
					resource.TestMatchResourceAttr(
						// check cloud_account_block has cloud_provider set to a non-empty string
						"data.wiz_kubernetes_clusters.foo",
						"kubernetes_clusters.0.cloud_account.0.cloud_provider",
						regexp.MustCompile(`\w`),
					),
				),
			},
			{
				Config: testAccDatasourceWizKubernetesClustersBasic(2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						// check that the second kubernetes cluster has an id that matches the UUID regex
						"data.wiz_kubernetes_clusters.foo",
						"kubernetes_clusters.1.id",
						regexp.MustCompile(`\w`),
					),
				),
			},
		},
	})
}

func testAccDatasourceWizKubernetesClustersBasic(maxPages int) string {
	return fmt.Sprintf(`
	data "wiz_kubernetes_clusters" "foo" {
		first  = 1
		max_pages = %d
		search = ""
	  }
	  
`, maxPages)
}
