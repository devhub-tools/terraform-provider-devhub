package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWorkspaceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccWorkspaceResourceConfig("my_database"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "name", "my_database"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "adapter", "POSTGRES"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "hostname", "localhost"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "ssl", "false"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "restrict_access", "true"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "enable_data_protection", "false"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "credentials.0.username", "postgres"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "credentials.0.password", "password"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "credentials.0.reviews_required", "0"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "credentials.0.default_credential", "true"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "credentials.1.username", "another"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "credentials.1.password", "password2"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "credentials.1.reviews_required", "1"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "credentials.1.default_credential", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "devhub_terradesk_workspace.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credentials.0.password", "credentials.1.password"},
			},
			// Update and Read testing
			{
				Config: testAccDatabaseResourceConfig("another_database"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "name", "another_database"),
					// everything else should be the same
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "adapter", "POSTGRES"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "hostname", "localhost"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "ssl", "false"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "restrict_access", "true"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "enable_data_protection", "false"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "credentials.0.username", "postgres"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "credentials.0.password", "password"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "credentials.0.reviews_required", "0"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "credentials.0.default_credential", "true"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "credentials.1.username", "another"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "credentials.1.password", "password2"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "credentials.1.reviews_required", "1"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "credentials.1.default_credential", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccWorkspaceResourceConfig(name string) string {
	return providerConfig + fmt.Sprintf(`
resource "devhub_terradesk_workspace" "test" {
  name     		 = %[1]q
  repository   = "devhub-tools/devhub"
	path 				 = "terraform"
	docker_image = "hashicorp/terraform:1.10"

	env_vars = [
		{
			name = "ENV_VAR"
			value = "env-var-value"
		}
	]

	secrets = [
		{
			name = "my_secret"
			value = "secret-value"
		}
	]
}
`, name)
}
