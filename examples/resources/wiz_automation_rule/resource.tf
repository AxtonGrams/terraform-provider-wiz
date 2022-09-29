resource "wiz_automation_rule" "rule1" {
  enabled = true
  name    = "test_recreate"
  filters = jsonencode(
    {
      "relatedEntity" : {
        "cloudPlatform" : [
          "AWS"
        ]
      }
    }
  )
  action_id      = wiz_automation_action.jira_transition.id
  trigger_source = "ISSUES"
  trigger_type = [
    "UPDATED",
  ]
}
