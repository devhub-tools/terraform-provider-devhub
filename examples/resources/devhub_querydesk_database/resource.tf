resource "devhub_querydesk_database" "example" {
  name                  = "terraform_test"
  adapter               = "POSTGRES"
  hostname              = "localhost"
  database              = "mydb"
  ai_enabled            = true
  ai_max_rows           = 500
  ai_timeout            = 60
  review_bypass_allowed = true

  # Lets members ask to run one query unredacted, after a peer review. Leave it out to keep the database's current value.
  unredacted_runs_allowed = true

  # Sent to the server when a connection opens. Here pgdog routes the database to the replica.
  connection_parameters = [
    {
      name  = "pgdog.role"
      value = "replica"
    }
  ]

  credentials = [
    {
      username           = "postgres"
      password           = "postgres"
      reviews_required   = 1
      default_credential = true
      ai_allowed         = true
      timeout            = 30
    }
  ]
}
