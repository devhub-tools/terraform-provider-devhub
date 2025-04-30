package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWorkflowResource(t *testing.T) {
	name := fmt.Sprintf("workflow_%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccWorkflowResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_workflow.test", "name", name),
					resource.TestCheckResourceAttr("devhub_workflow.test", "inputs.0.key", "user_id"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "inputs.0.description", "User ID"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "inputs.0.type", "string"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.0.name", "approval-step"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.0.approval_action.reviews_required", "1"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.1.name", "api-step"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.1.api_action.endpoint", "https://api.example.com/endpoint"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.1.api_action.method", "GET"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.1.api_action.expected_status_code", "200"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.1.api_action.include_devhub_jwt", "true"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.1.api_action.headers.0.key", "content-type"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.1.api_action.headers.0.value", "application/json"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.2.name", "query-step"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.2.query_action.query", "SELECT * FROM users WHERE id = '${user_id}'"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.2.query_action.timeout", "10"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.2.query_action.credential_id", "crd_123"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.3.name", "slack-step"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.3.slack_action.slack_channel", "#general"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.3.slack_action.message", "Hello, world!"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.3.slack_action.link_text", "Click here"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.4.name", "slack-reply-step"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.4.slack_reply_action.reply_to_step_name", "slack-step"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.4.slack_reply_action.message", "Hello, world!"),
					resource.TestCheckResourceAttrSet("devhub_workflow.test", "id"),
					resource.TestCheckResourceAttrSet("devhub_workflow.test", "steps.0.id"),
					resource.TestCheckResourceAttrSet("devhub_workflow.test", "steps.1.id"),
					resource.TestCheckResourceAttrSet("devhub_workflow.test", "steps.2.id"),
					resource.TestCheckResourceAttrSet("devhub_workflow.test", "steps.3.id"),
					resource.TestCheckResourceAttrSet("devhub_workflow.test", "steps.4.id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "devhub_workflow.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccWorkflowResourceConfig(name + "_updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_workflow.test", "name", name+"_updated"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "inputs.0.key", "user_id"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "inputs.0.description", "User ID"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "inputs.0.type", "string"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.0.name", "approval-step"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.0.approval_action.reviews_required", "1"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.1.name", "api-step"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.1.api_action.endpoint", "https://api.example.com/endpoint"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.1.api_action.method", "GET"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.1.api_action.expected_status_code", "200"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.1.api_action.include_devhub_jwt", "true"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.1.api_action.headers.0.key", "content-type"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.1.api_action.headers.0.value", "application/json"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.2.name", "query-step"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.2.query_action.query", "SELECT * FROM users WHERE id = '${user_id}'"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.2.query_action.timeout", "10"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.2.query_action.credential_id", "crd_123"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.3.name", "slack-step"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.3.slack_action.slack_channel", "#general"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.3.slack_action.message", "Hello, world!"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.3.slack_action.link_text", "Click here"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.4.name", "slack-reply-step"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.4.slack_reply_action.reply_to_step_name", "slack-step"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.4.slack_reply_action.message", "Hello, world!"),
					resource.TestCheckResourceAttrSet("devhub_workflow.test", "id"),
					resource.TestCheckResourceAttrSet("devhub_workflow.test", "steps.0.id"),
					resource.TestCheckResourceAttrSet("devhub_workflow.test", "steps.1.id"),
					resource.TestCheckResourceAttrSet("devhub_workflow.test", "steps.2.id"),
					resource.TestCheckResourceAttrSet("devhub_workflow.test", "steps.3.id"),
					resource.TestCheckResourceAttrSet("devhub_workflow.test", "steps.4.id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccWorkflowResourceConfig(name string) string {
	return providerConfig + fmt.Sprintf(`
resource "devhub_workflow" "test" {
  name     = %[1]q
  inputs = [
    {
      key = "user_id"
      description = "User ID"
      type = "string"
    }
  ]
  steps = [
		{
			name = "approval-step"
			approval_action = {
				reviews_required = 1
			}
		},
		{
      name = "api-step"
			api_action = {
				endpoint = "https://api.example.com/endpoint"
				method = "GET"
				expected_status_code = 200
				include_devhub_jwt = true
				headers = [
					{ key = "content-type", value = "application/json" }
				]
			}
		},
		{
			name = "query-step"
			query_action = {
				query = "SELECT * FROM users WHERE id = '$${user_id}'"
				credential_id = "crd_123"
				timeout = 10
			}
		},
		{
			name = "slack-step"
			slack_action = {
				slack_channel = "#general"
				message = "Hello, world!"
				link_text = "Click here"
			}
		},
		{
			name = "slack-reply-step"
			slack_reply_action = {
				reply_to_step_name = "slack-step"
				message = "Hello, world!"
			}
		}
	]
}
`, name)
}
