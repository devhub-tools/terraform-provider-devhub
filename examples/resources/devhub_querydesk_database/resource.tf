resource "devhub_querydesk_database" "example" {
  name        = "terraform_test"
  adapter     = "POSTGRES"
  hostname    = "localhost"
  database    = "mydb"
  ai_enabled  = true
  ai_max_rows = 500

  credentials = [
    {
      username           = "postgres"
      password           = "postgres"
      reviews_required   = 0
      default_credential = true
      ai_allowed         = true
      timeout            = 30
    }
  ]
}
