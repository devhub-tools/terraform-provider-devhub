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
					resource.TestCheckResourceAttr("devhub_dashboard.test", "panels.0.query_details.query", "SELECT * FROM users WHERE id = '${user_id}'"),
					resource.TestCheckResourceAttr("devhub_dashboard.test", "panels.0.query_details.credential_id", "123"),
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
					resource.TestCheckResourceAttr("devhub_dashboard.test", "panels.0.query_details.query", "SELECT * FROM users WHERE id = '${user_id}'"),
					resource.TestCheckResourceAttr("devhub_dashboard.test", "panels.0.query_details.credential_id", "123"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDashboardResourceConfig(name string) string {
	return providerConfig + fmt.Sprintf(`
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

			query_details = {
				query = "SELECT * FROM users WHERE id = '$${user_id}'"
				credential_id = "123"
			}
		},
    {
      title = "Users"

			query_details = {
				query = "SELECT * FROM users WHERE id = '$${user_id}'"
				credential_id = "123"
			}
		}
	]
}
`, name)
}
