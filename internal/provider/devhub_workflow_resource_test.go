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
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.0.api_action.endpoint", "https://api.example.com/endpoint"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.0.api_action.method", "GET"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.0.api_action.expected_status_code", "200"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.0.api_action.include_devhub_jwt", "true"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.0.api_action.headers.0.key", "content-type"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.0.api_action.headers.0.value", "application/json"),
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
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.0.api_action.endpoint", "https://api.example.com/endpoint"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.0.api_action.method", "GET"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.0.api_action.expected_status_code", "200"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.0.api_action.include_devhub_jwt", "true"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.0.api_action.headers.0.key", "content-type"),
					resource.TestCheckResourceAttr("devhub_workflow.test", "steps.0.api_action.headers.0.value", "application/json"),
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

			api_action = {
				endpoint = "https://api.example.com/endpoint"
				method = "GET"
				expected_status_code = 200
				include_devhub_jwt = true
				headers = [
					{ key = "content-type", value = "application/json" }
				]
			}
		}
	]
}
`, name)
}
