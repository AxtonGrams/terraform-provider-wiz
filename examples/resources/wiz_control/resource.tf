resource "wiz_control" "test" {
  name                      = "test control 2"
  enabled                   = false
  description               = "test control 2 description"
  project_id                = "*"
  severity                  = "LOW"
  resolution_recommendation = "fix it"
  security_sub_categories = [
    "wsct-id-8",
  ]
  query = jsonencode(
    {
      "relationships" : [
        {
          "type" : [
            {
              "reverse" : true,
              "type" : "CONTAINS"
            }
          ],
          "with" : {
            "select" : true,
            "type" : [
              "SUBSCRIPTION"
            ]
          }
        }
      ]
    }
  )
  scope_query = jsonencode(
    {
      "type" : [
        "SUBSCRIPTION"
      ]
    }
  )
}
