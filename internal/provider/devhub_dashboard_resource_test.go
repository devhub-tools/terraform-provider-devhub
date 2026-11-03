package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDashboardResource(t *testing.T) {
	name := fmt.Sprintf("dashboard_%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDashboardResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_dashboard.test", "name", name),
					resource.TestCheckResourceAttr("devhub_dashboard.test", "panels.0.title", "Users"),
					resource.TestCheckResourceAttr("devhub_dashboard.test", "panels.0.inputs.0.key", "user_id"),
					resource.TestCheckResourceAttr("devhub_dashboard.test", "panels.0.inputs.0.description", "User ID"),
					resource.TestCheckResourceAttr("devhub_dashboard.test", "panels.0.query_details.query", "SELECT * FROM users WHERE id = '$${user_id}'"),
					resource.TestCheckResourceAttrPair("devhub_dashboard.test", "panels.0.query_details.credential_id", "devhub_querydesk_database.test", "credentials.0.id"),
					resource.TestCheckResourceAttr("devhub_dashboard.test", "panels.0.actions.0.label", "Disable user"),
					resource.TestCheckResourceAttrPair("devhub_dashboard.test", "panels.0.actions.0.workflow_id", "devhub_workflow.test", "id"),
					resource.TestCheckResourceAttr("devhub_dashboard.test", "panels.0.actions.0.input_mappings.0.workflow_input_key", "user_id"),
					resource.TestCheckResourceAttr("devhub_dashboard.test", "panels.0.actions.0.input_mappings.0.column", "id"),
					resource.TestCheckResourceAttr("devhub_dashboard.test", "panels.0.actions.0.input_mappings.1.workflow_input_key", "reason"),
					resource.TestCheckResourceAttr("devhub_dashboard.test", "panels.0.actions.0.input_mappings.1.prompt_user", "true"),
					resource.TestCheckResourceAttr("devhub_dashboard.test", "panels.0.actions.0.input_mappings.1.prompt", "Why are you disabling this user?"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "devhub_dashboard.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDashboardResourceConfig(name + "_updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_dashboard.test", "name", name+"_updated"),
					resource.TestCheckResourceAttr("devhub_dashboard.test", "panels.0.title", "Users"),
					resource.TestCheckResourceAttr("devhub_dashboard.test", "panels.0.inputs.0.key", "user_id"),
					resource.TestCheckResourceAttr("devhub_dashboard.test", "panels.0.inputs.0.description", "User ID"),
					resource.TestCheckResourceAttr("devhub_dashboard.test", "panels.0.query_details.query", "SELECT * FROM users WHERE id = '$${user_id}'"),
					resource.TestCheckResourceAttrPair("devhub_dashboard.test", "panels.0.query_details.credential_id", "devhub_querydesk_database.test", "credentials.0.id"),
					resource.TestCheckResourceAttr("devhub_dashboard.test", "panels.0.actions.0.label", "Disable user"),
					resource.TestCheckResourceAttrPair("devhub_dashboard.test", "panels.0.actions.0.workflow_id", "devhub_workflow.test", "id"),
					resource.TestCheckResourceAttr("devhub_dashboard.test", "panels.0.actions.0.input_mappings.0.workflow_input_key", "user_id"),
					resource.TestCheckResourceAttr("devhub_dashboard.test", "panels.0.actions.0.input_mappings.0.column", "id"),
					resource.TestCheckResourceAttr("devhub_dashboard.test", "panels.0.actions.0.input_mappings.1.workflow_input_key", "reason"),
					resource.TestCheckResourceAttr("devhub_dashboard.test", "panels.0.actions.0.input_mappings.1.prompt_user", "true"),
					resource.TestCheckResourceAttr("devhub_dashboard.test", "panels.0.actions.0.input_mappings.1.prompt", "Why are you disabling this user?"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDashboardResourceConfig(name string) string {
	return providerConfig + fmt.Sprintf(`
# panels' credentials and actions' workflows have to exist in the organization
resource "devhub_querydesk_database" "test" {
  name     = "%[1]s_db"
  adapter  = "POSTGRES"
  hostname = "localhost"
  database = "mydb"

  credentials = [
    {
      username           = "postgres"
      password           = "password"
      reviews_required   = 0
      default_credential = true
    }
  ]
}

resource "devhub_workflow" "test" {
  name = "%[1]s_workflow"

  inputs = [
    { key = "user_id", type = "string" },
    { key = "reason", type = "string" }
  ]

  steps = [
    {
      name = "condition-step"
      condition_action = {
        condition  = "true"
        when_false = "succeeded"
      }
    }
  ]
}

resource "devhub_dashboard" "test" {
  name     = %[1]q
  panels = [
    {
      title = "Users"

			inputs = [
				{
					key = "user_id"
					description = "User ID"
				}
			]

			actions = [
				{
					label = "Disable user"
					workflow_id = devhub_workflow.test.id

					input_mappings = [
						{
							workflow_input_key = "user_id"
							column = "id"
						},
						{
							workflow_input_key = "reason"
							prompt_user = true
							prompt = "Why are you disabling this user?"
						}
					]
				}
			]

			query_details = {
				query = "SELECT * FROM users WHERE id = '$$${user_id}'"
				credential_id = devhub_querydesk_database.test.credentials[0].id
			}
		},
    {
      title = "Users"

			query_details = {
				query = "SELECT * FROM users WHERE id = '$${user_id}'"
				credential_id = devhub_querydesk_database.test.credentials[0].id
			}
		}
	]
}
`, name)
}
