package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWorkspaceResource(t *testing.T) {
	name := fmt.Sprintf("workspace_%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccWorkspaceResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "name", name),
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
				Config: testAccWorkspaceResourceConfig(name + "_updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "name", name+"_updated"),
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

func TestAccWorkspaceWithInitContainersResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccWorkspaceWithInitContainersResourceConfig("1password/op:2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "command", "/home/terradesk/bin/op run -- terraform"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "init_containers.0.name", "install-op"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "init_containers.0.image", "1password/op:2"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "init_containers.0.command", `sh -c "cp /usr/local/bin/op /home/terradesk/bin/"`),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "secrets.0.name", "OP_SERVICE_ACCOUNT_TOKEN"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "secrets.0.raw", "true"),
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
				Config: testAccWorkspaceWithInitContainersResourceConfig("1password/op:2.30"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "init_containers.0.image", "1password/op:2.30"),
					// everything else should be the same
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "command", "/home/terradesk/bin/op run -- terraform"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "init_containers.0.name", "install-op"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "secrets.0.raw", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccWorkspaceWithInitContainersResourceConfig(image string) string {
	return providerConfig + fmt.Sprintf(`
resource "devhub_terradesk_workspace" "test" {
  name     		 = "my_workspace"
  repository   = "devhub-tools/devhub"
	path 				 = "terraform"
	docker_image = "hashicorp/terraform:1.10"
	command      = "/home/terradesk/bin/op run -- terraform"

	init_containers = [
		{
			name = "install-op"
			image = %[1]q
			command = "sh -c \"cp /usr/local/bin/op /home/terradesk/bin/\""
		}
	]

	secrets = [
		{
			name = "OP_SERVICE_ACCOUNT_TOKEN"
			value = "ops_token"
			raw = true
		}
	]
}
`, image)
}

func TestAccWorkspaceWithSlackNotificationsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccWorkspaceWithSlackNotificationsResourceConfig("all_plans"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "slack_channel", "#terraform"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "slack_notify_on", "all_plans"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "devhub_terradesk_workspace.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccWorkspaceWithSlackNotificationsResourceConfig("default_branch"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "slack_notify_on", "default_branch"),
					// everything else should be the same
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "slack_channel", "#terraform"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccWorkspaceWithSlackNotificationsResourceConfig(notifyOn string) string {
	return providerConfig + fmt.Sprintf(`
resource "devhub_terradesk_workspace" "test" {
  name     		 = "my_workspace"
  repository   = "devhub-tools/devhub"
	path 				 = "terraform"
	docker_image = "hashicorp/terraform:1.10"

	slack_channel   = "#terraform"
	slack_notify_on = %[1]q
}
`, notifyOn)
}
