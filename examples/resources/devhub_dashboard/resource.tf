resource "devhub_dashboard" "example" {
  name = "terraform_test"

  panels = [
    {
      title = "Users"

      inputs = [
        {
          key         = "user_id"
          description = "User ID"
          type        = "string"
        }
      ]

      query_details = {
        query         = "SELECT * FROM users WHERE id = '$${user_id}'"
        credential_id = "crd_xxx"
        timeout       = 10
      }
    }
  ]
}
