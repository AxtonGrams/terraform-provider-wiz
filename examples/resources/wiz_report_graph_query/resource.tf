# A simple example
resource "wiz_report_graph_query" "foo" {
  name       = "foo"
  project_id = "2c38b8fa-c315-57ea-9de4-e3a19592d796"
  query      = <<EOF
{
  "select": true,
  "type": [
    "CONTAINER_IMAGE"
  ],
  "where": {
    "name": {
      "CONTAINS": [
        "foo"
      ]
    }
  }
}
EOF
}

# Scheduling enabled (both run_interval_hours and run_starts_at must be set)
resource "wiz_report_graph_query" "foo" {
  name               = "foo"
  project_id         = "2c38b8fa-c315-57ea-9de4-e3a19592d796"
  run_interval_hours = 24
  run_starts_at      = "2023-12-06 16:00:00 +0000 UTC"
  query              = <<EOF
{
  "select": true,
  "type": [
    "CONTAINER_IMAGE"
  ],
  "where": {
    "name": {
      "CONTAINS": [
        "foo"
      ]
    }
  }
}
EOF
}
