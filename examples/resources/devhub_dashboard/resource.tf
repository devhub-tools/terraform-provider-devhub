resource "devhub_dashboard" "example" {
  name = "terraform_test"

  panels = [
    {
      title = "Users"

      inputs = [
        {
          key         = "user_id"
          description = "User ID"
        }
      ]

      query_details = {
        query         = "SELECT * FROM users WHERE id = '$${user_id}'"
        credential_id = "crd_xxx"
      }

      actions = [
        {
          label       = "Disable user"
          workflow_id = "wf_xxx"

          # One mapping per input of the workflow being triggered. Either pull
          # the value from a column of the clicked row, or prompt the user for
          # it when they click the button.
          input_mappings = [
            {
              workflow_input_key = "user_id"
              column             = "id"
            },
            {
              workflow_input_key = "reason"
              prompt_user        = true
              prompt             = "Why are you disabling this user?"
            }
          ]
        }
      ]
    }
  ]
}
