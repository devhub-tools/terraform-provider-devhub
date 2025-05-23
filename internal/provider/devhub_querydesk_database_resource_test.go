package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatabaseResource(t *testing.T) {
	name := fmt.Sprintf("database_%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDatabaseResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "name", name),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "adapter", "POSTGRES"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "hostname", "localhost"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "ssl", "false"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "restrict_access", "true"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.0.username", "postgres"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.0.password", "password"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.0.reviews_required", "0"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.0.default_credential", "true"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.1.username", "another"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.1.password", "password2"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.1.reviews_required", "1"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.1.default_credential", "false"),
					resource.TestCheckResourceAttrSet("devhub_querydesk_database.test", "credential_ids.postgres"),
					resource.TestCheckResourceAttrSet("devhub_querydesk_database.test", "credential_ids.another"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "devhub_querydesk_database.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credentials.0.password", "credentials.1.password"},
			},
			// Update and Read testing
			{
				Config: testAccDatabaseResourceConfig(name + "_updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "name", name+"_updated"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "adapter", "POSTGRES"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "hostname", "localhost"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "ssl", "false"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "restrict_access", "true"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.0.username", "postgres"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.0.password", "password"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.0.reviews_required", "0"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.0.default_credential", "true"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.1.username", "another"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.1.password", "password2"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.1.reviews_required", "1"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.1.default_credential", "false"),
					resource.TestCheckResourceAttrSet("devhub_querydesk_database.test", "credential_ids.postgres"),
					resource.TestCheckResourceAttrSet("devhub_querydesk_database.test", "credential_ids.another"),
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

	credentials = [
		{
			username = "postgres"
			password = "password"
			reviews_required = 0
			default_credential = true
		},
		{
			username = "another"
			password = "password2"
			reviews_required = 1
	}
	]
}
`, name)
}
