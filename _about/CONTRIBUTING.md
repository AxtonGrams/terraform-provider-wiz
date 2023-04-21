# Welcome

_Please Note: This documentation is intended for Terraform Provider code developers. Typical operators writing and applying Terraform configurations do not need to read or understand this material._

## Contribute

Please follow the following steps to ensure your contribution goes smoothly.

### 1. Configure Development Environment

Install Terraform and Go. Clone the repository, compile the provider, and set up testing.

### 2. Change Code

### 3. Write Tests

Changes must be covered by acceptance tests for all contributions.

### 4. Create a Pull Request

When your contribution is ready, Create a Pull Request in the Wiz provider repository.

Include the output from the acceptance tests for the resource you created or altered.  Acceptance tests can be targeted to the specific resources as follows:

```
$ TF_ACC=1 go test ./internal/acceptance/... -v -run='TestAccResourceWizSAMLIdp_basic'
=== RUN   TestAccResourceWizSAMLIdp_basic
2023/04/20 16:09:44 [DEBUG] POST https://auth.app.wiz.io/oauth/token
2023/04/20 16:09:45 [DEBUG] POST https://auth.app.wiz.io/oauth/token
2023/04/20 16:09:46 [DEBUG] POST https://auth.app.wiz.io/oauth/token
2023/04/20 16:09:47 [DEBUG] POST https://api.us8.app.wiz.io/graphql
2023/04/20 16:09:48 [DEBUG] POST https://api.us8.app.wiz.io/graphql
2023/04/20 16:09:50 [DEBUG] POST https://auth.app.wiz.io/oauth/token
2023/04/20 16:09:51 [DEBUG] POST https://auth.app.wiz.io/oauth/token
2023/04/20 16:09:51 [DEBUG] POST https://api.us8.app.wiz.io/graphql
2023/04/20 16:09:52 [DEBUG] POST https://auth.app.wiz.io/oauth/token
2023/04/20 16:09:53 [DEBUG] POST https://auth.app.wiz.io/oauth/token
2023/04/20 16:09:54 [DEBUG] POST https://auth.app.wiz.io/oauth/token
2023/04/20 16:09:55 [DEBUG] POST https://api.us8.app.wiz.io/graphql
--- PASS: TestAccResourceWizSAMLIdp_basic (11.93s)
PASS
ok      wiz.io/hashicorp/terraform-provider-wiz/internal/acceptance     11.950s
```
