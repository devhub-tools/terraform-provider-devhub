package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrderResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDatabaseResourceConfig("my_database"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "name", "my_database"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "ssl", "false"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "restrict_access", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "devhub_querydesk_database.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDatabaseResourceConfig("another_database"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "name", "another_database"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDatabaseResourceConfig(name string) string {
	return providerConfig + fmt.Sprintf(`
resource "devhub_querydesk_database" "test" {
  name     = %[1]q
  adapter  = "POSTGRES"
  hostname = "localhost"
  database = "mydb"
}
`, name)
}
