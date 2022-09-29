resource "wiz_security_framework" "test" {
  name        = "terraform-test-security-framework1"
  description = "test description"
  enabled     = true
  category {
    name        = "AM Asset Management"
    description = "test am description"
    sub_category {
      title = "AM-1 Track asset inventory and their risks"
    }
  }
  category {
    name        = "test category 2"
    description = "test description 2"
    sub_category {
      title       = "test subcategory"
      description = "bad stuff now"
    }
    sub_category {
      title       = "test subcategory 2"
      description = "bad stuff could happen"
    }
  }
}
