package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceWizUser_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix(ResourcePrefix)
	smtpDomain := os.Getenv("WIZ_SMTP_DOMAIN")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TestCase(TcUser)) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceWizUserBasic(rName, smtpDomain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						// check that name is correctly set
						"wiz_user.foo",
						"name",
						rName,
					),
					resource.TestCheckResourceAttr(
						// check that email is correctly set
						"wiz_user.foo",
						"email",
						fmt.Sprintf("%s@%s", rName, smtpDomain),
					),
					resource.TestCheckResourceAttr(
						// check that role is set to specified value
						"wiz_user.foo",
						"role",
						"PROJECT_MEMBER",
					),
					resource.TestCheckResourceAttrSet(
						// check for a set assigned project id
						"wiz_user.foo",
						"assigned_project_ids.0",
					),
					resource.TestCheckResourceAttrSet(
						// check that send_email_invite is set (can be null)
						"wiz_user.foo",
						"send_email_invite",
					),
				),
			},
		},
	})
}

func testResourceWizUserBasic(rName string, smtpDomain string) string {
	project := uuid.New().String()
	return fmt.Sprintf(`
	resource "wiz_user" "foo" {
		name                 = "%[1]s"
		email                = "%[1]s@%s"
		role                 = "PROJECT_MEMBER"
		assigned_project_ids = [ "%s" ]
	  } 
`, rName, smtpDomain, project)
}
