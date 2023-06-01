package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceWizSAMLIdp_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix(ResourcePrefix)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TestCase(TcCommon)) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceWizSAMLIdpBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"wiz_saml_idp.test",
						"issuer_url",
						"https://ping.example.com/idp/SSO.saml2",
					),
					resource.TestCheckResourceAttr(
						"wiz_saml_idp.test",
						"login_url",
						"https://ping.example.com/idp/SSO.saml2",
					),
					resource.TestCheckResourceAttr(
						"wiz_saml_idp.test",
						"logout_url",
						"https://ping.example.com/idp/SLO.saml2",
					),
					resource.TestCheckResourceAttr(
						"wiz_saml_idp.test",
						"use_provider_managed_roles",
						"true",
					),
					resource.TestCheckResourceAttr(
						"wiz_saml_idp.test",
						"allow_manual_role_override",
						"false",
					),
					resource.TestCheckResourceAttr(
						"wiz_saml_idp.test",
						"merge_groups_mapping_by_role",
						"false",
					),
				),
			},
		},
	})
}

func testResourceWizSAMLIdpBasic(rName string) string {
	return fmt.Sprintf(`
resource "wiz_saml_idp" "test" {
  name                         = "%s"
  issuer_url                   = "https://ping.example.com/idp/SSO.saml2"
  login_url                    = "https://ping.example.com/idp/SSO.saml2"
  logout_url                   = "https://ping.example.com/idp/SLO.saml2"
  use_provider_managed_roles   = true
  allow_manual_role_override   = false
  merge_groups_mapping_by_role = false
  certificate                  = <<-EOT
-----BEGIN CERTIFICATE-----
MIIFpzCCA4+gAwIBAgIJAKY0mQyPWs1eMA0GCSqGSIb3DQEBCwUAMGoxCzAJBgNV
BAYTAlVTMRAwDgYDVQQIDAdOb3doZXJlMRwwGgYDVQQHDBNOb3RoaW5nIHRvIHNl
ZSBoZXJlMRwwGgYDVQQKDBNEZWZhdWx0IENvbXBhbnkgTHRkMQ0wCwYDVQQDDARw
aW5nMB4XDTIyMDYyNDE2MDU1M1oXDTMyMDYyMTE2MDU1M1owajELMAkGA1UEBhMC
VVMxEDAOBgNVBAgMB05vd2hlcmUxHDAaBgNVBAcME05vdGhpbmcgdG8gc2VlIGhl
cmUxHDAaBgNVBAoME0RlZmF1bHQgQ29tcGFueSBMdGQxDTALBgNVBAMMBHBpbmcw
ggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDwkVOUv+IpjkU9AyWgpUlC
wtidRa79Twq6z93fIzuStVt7ITKM0aoIwrtkINlZ3Y78zQb+k3/43QNjO89fj9Bc
UJaXHBLbRrp1T06n0gbl2pQ2dMQ+GfhhIe82o/4XtbnZFzLDVhnVWxHeZsaqQ934
VFo/uz6zrfWjaqHhn66TK6ItYM3xRxxb24WVVMjXDpSiboGQl49t266TTidnn/fq
RuPBc/8lAYcaye4U4r0ExIOw/VT24S0W9X+OBzKJNOI5Moz8c/wPnPGHRlrp/toa
iS3k3sKxMHTbQaP6EYYvgyEHR0aK6SJU39Gf6jm91EJguJpJVCSOW0XC0Vo0O7gH
9mybgZ6hOtDKIrk8DlfDjpOypZVq542oVTMSjIsaFb1Y5LIaVrKD51j8KrZmyE3d
EGwatw6VsLerF9yi5wv6wi/oGgLLM5IeLJAa4sk1+eqGxg2nAaZd3EgW36BzlL9P
GjcbX1Zaiwe4MWxorR3iMADV8a5JPmd8oC+WXXCqq8ddxVsl2nkJ47jUwHPzY2Uu
Hd9XyLoqlY9GTFsz+I6pd/NU9YD/xSuw2jVhaynRjFEboAecZp9kA8fwD7qwKW71
NIru/KgoKbHl2yQ513LgpzunlT7xMtSId5bvdb7U4pIUH1GZYHZ+wiTi1qU600aH
HrLI4jSsA5YhnaARnpzQowIDAQABo1AwTjAdBgNVHQ4EFgQUN6E/ZVbo9vAPlFZy
phVOm9zZ4kwwHwYDVR0jBBgwFoAUN6E/ZVbo9vAPlFZyphVOm9zZ4kwwDAYDVR0T
BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAgEA6wt/A7kxhjt6RRBBP8WRE2w8oPtx
sh+4DQWwtHIIiaLd0nuqXYgmlnIBQcuayKgVWvKiMZwlCg7M8aJHb1DV68jn21/J
5S3C+YHwImP6BlY8158ID280ZU3JoZfyIH+MToumdbfd4wMlKNNOvHbMoUtgq4rS
Cd6vnZ8RF/5MykD4CPhC4zc3hI+xfrtcd2gnSKwDbU5wY94uB0f2QK/qbGnVjz7P
eUFAdPKOiFIRBCVzNJj6JxSUq1CABqW5UDmAK9bDpuWybnRKrZrhg4PA0O7ZIgyq
D5WaGPD+U+zozrN7YXTu97Ey/S4HFEjJBZaOkPFSTKG6u+l2sjrD6dgrq3ae7t4W
fxI7pamr7cd1t316c/8zowl7JUkHZMvu2kE3CrCyMKgwZ22EuVSEfSKss5fPelWs
ZGNFLHcI8Xmk+RL13bAh/41bxEt80WQrVJRg3X5mFhzTFec6Ox9v5loh2sEd7jh7
dKouC0o4KxPVfiAK0FJL9aaB/K/rrSb6ddzal4hZ8t91AKRsUXTw/Iu8+nJ0tbU5
gu3BmRbPJ5DphvXRY1yy3GERGpQIHWSn7XxvH/OlXO+mHHWNNYa5SW7V4RY4UcSD
Z8lCchNPFJqIlyvk9LSEorFq4tT21t/pgVOFgw0yJaTyBZ/IvIimjwNHJBnIeBQ2
GfRTgIAGAQ8ZFfQ=
-----END CERTIFICATE-----
EOT
}
`, rName)
}
