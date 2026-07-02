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
				Config: testAccDatabaseResourceConfig(name, true, 100, true, 30),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "name", name),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "adapter", "POSTGRES"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "hostname", "localhost"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "ssl", "false"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "restrict_access", "true"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "ai_enabled", "true"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "ai_max_rows", "100"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.0.username", "postgres"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.0.password", "password"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.0.reviews_required", "0"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.0.default_credential", "true"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.0.ai_allowed", "true"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.0.timeout", "30"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.1.username", "another"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.1.password", "password2"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.1.reviews_required", "1"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.1.default_credential", "false"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.1.ai_allowed", "false"),
					resource.TestCheckNoResourceAttr("devhub_querydesk_database.test", "credentials.1.timeout"),
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
				Config: testAccDatabaseResourceConfig(name+"_updated", false, 250, false, 60),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "name", name+"_updated"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "adapter", "POSTGRES"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "hostname", "localhost"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "ssl", "false"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "restrict_access", "true"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "ai_enabled", "false"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "ai_max_rows", "250"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.0.username", "postgres"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.0.password", "password"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.0.reviews_required", "0"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.0.default_credential", "true"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.0.ai_allowed", "false"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.0.timeout", "60"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.1.username", "another"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.1.password", "password2"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.1.reviews_required", "1"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.1.default_credential", "false"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.1.ai_allowed", "false"),
					resource.TestCheckResourceAttrSet("devhub_querydesk_database.test", "credential_ids.postgres"),
					resource.TestCheckResourceAttrSet("devhub_querydesk_database.test", "credential_ids.another"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDatabaseResourceConfig(name string, aiEnabled bool, aiMaxRows int, aiAllowed bool, timeout int) string {
	return providerConfig + fmt.Sprintf(`
resource "devhub_querydesk_database" "test" {
  name        = %[1]q
  adapter     = "POSTGRES"
  hostname    = "localhost"
  database    = "mydb"
  ai_enabled  = %[2]t
  ai_max_rows = %[3]d

	credentials = [
		{
			username = "postgres"
			password = "password"
			reviews_required = 0
			default_credential = true
			ai_allowed = %[4]t
			timeout = %[5]d
		},
		{
			username = "another"
			password = "password2"
			reviews_required = 1
	}
	]
}
`, name, aiEnabled, aiMaxRows, aiAllowed, timeout)
}
