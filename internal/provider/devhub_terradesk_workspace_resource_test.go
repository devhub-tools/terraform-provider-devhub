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
				Config: testAccWorkspaceResourceConfig("my_workspace"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "name", "my_workspace"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "repository", "devhub-tools/devhub"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "path", "terraform"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "docker_image", "hashicorp/terraform:1.10"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "env_vars.0.name", "ENV_VAR"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "env_vars.0.value", "env-var-value"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "secrets.0.name", "my_secret"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "secrets.0.value", "secret-value"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "devhub_terradesk_workspace.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secrets.0.value"},
			},
			// Update and Read testing
			{
				Config: testAccWorkspaceResourceConfig("another_workspace"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "name", "another_workspace"),
					// everything else should be the same
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "repository", "devhub-tools/devhub"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "path", "terraform"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "docker_image", "hashicorp/terraform:1.10"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "env_vars.0.name", "ENV_VAR"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "env_vars.0.value", "env-var-value"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "secrets.0.name", "my_secret"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "secrets.0.value", "secret-value"),
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

func TestAccWorkspaceWithWorkloadIdentityResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccWorkspaceWithWorkloadIdentityResourceConfig("devhub@google.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "name", "my_workspace"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "workload_identity.enabled", "true"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "workload_identity.service_account_email", "devhub@google.com"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "workload_identity.provider", "projects/123456789/locations/global/workloadIdentityPools/pools/devhub/providers/devhub"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "devhub_terradesk_workspace.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secrets.0.value"},
			},
			// Update and Read testing
			{
				Config: testAccWorkspaceWithWorkloadIdentityResourceConfig("serviceaccount@google.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "name", "my_workspace"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "workload_identity.enabled", "true"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "workload_identity.service_account_email", "serviceaccount@google.com"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "workload_identity.provider", "projects/123456789/locations/global/workloadIdentityPools/pools/devhub/providers/devhub"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccWorkspaceWithWorkloadIdentityResourceConfig(email string) string {
	return providerConfig + fmt.Sprintf(`
resource "devhub_terradesk_workspace" "test" {
  name     		 = "my_workspace"
  repository   = "devhub-tools/devhub"
	path 				 = "terraform"
	docker_image = "hashicorp/terraform:1.10"

	workload_identity = {
		enabled             = true
		service_account_email = %[1]q
		provider            = "projects/123456789/locations/global/workloadIdentityPools/pools/devhub/providers/devhub"
	}
}
`, email)
}
