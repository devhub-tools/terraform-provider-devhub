resource "devhub_workflow" "example" {
  name = "terraform_test"

  inputs = [
    {
      key         = "user_id"
      description = "User ID"
      type        = "string"
    }
  ]

  steps = [
    {
      name = "api-step"

      api_action = {
        endpoint             = "https://api.example.com/endpoint"
        method               = "GET"
        expected_status_code = 200
        include_devhub_jwt   = true
        headers = [
          { key = "content-type", value = "application/json" }
        ]
        body = "{\"user_id\": \"$${user_id}\"}"
      }
    },
    {
      name = "slack-step"

      slack_action = {
        slack_channel = "my-channel"
        message       = "Hello, world!"
        link_text     = "Click here"
      }
    },
    {
      name = "query-step"

      query_action = {
        query         = "SELECT * FROM users WHERE id = '$${user_id}'"
        credential_id = "crd_xxx"
        timeout       = 10
      }
    },
    {
      name = "approval-step"

      approval_action = {
        required_approvals = 1
      }
    },
    {
      name = "slack-reply-step"

      slack_reply_action = {
        reply_to_step_name = "slack-step"
        message            = "Workflow completed successfully"
      }
    }
  ]
}
