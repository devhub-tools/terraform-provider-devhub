resource "devhub_querydesk_database" "example" {
  name     = "terraform_test"
  adapter  = "POSTGRES"
  hostname = "localhost"
  database = "mydb"

  credentials = [
    {
      username           = "postgres"
      password           = "postgres"
      reviews_required   = 0
      default_credential = true
    }
  ]
}
