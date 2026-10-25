package provider

import (
	"fmt"
	"regexp"
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
				Config: testAccDatabaseResourceConfig(name, true, 100, true, 30, true, `ai_timeout = 120
  connection_parameters = [
    { name = "pgdog.role", value = "replica" },
    { name = "options", value = "-c search_path=reporting" },
  ]`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "name", name),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "adapter", "POSTGRES"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "hostname", "localhost"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "ssl", "false"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "restrict_access", "true"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "review_bypass_allowed", "true"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "ai_enabled", "true"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "ai_max_rows", "100"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "ai_timeout", "120"),
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
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "connection_parameters.#", "2"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "connection_parameters.0.name", "pgdog.role"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "connection_parameters.0.value", "replica"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "connection_parameters.1.name", "options"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "connection_parameters.1.value", "-c search_path=reporting"),
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
				Config: testAccDatabaseResourceConfig(name+"_updated", false, 250, false, 60, false, `ai_timeout = 45
  connection_parameters = [{ name = "pgdog.role", value = "primary" }]`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "name", name+"_updated"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "adapter", "POSTGRES"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "hostname", "localhost"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "ssl", "false"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "restrict_access", "true"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "review_bypass_allowed", "false"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "ai_enabled", "false"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "ai_max_rows", "250"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "ai_timeout", "45"),
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
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "connection_parameters.#", "1"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "connection_parameters.0.name", "pgdog.role"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "connection_parameters.0.value", "primary"),
				),
			},
			// Removing the attribute clears the database's connection parameters
			{
				Config: testAccDatabaseResourceConfig(name+"_updated", false, 250, false, 60, false, "ai_timeout = 45"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("devhub_querydesk_database.test", "connection_parameters.#"),
				),
			},
			// Leaving ai_timeout out sets the database back to 10
			{
				Config: testAccDatabaseResourceConfig(name+"_updated", false, 250, false, 60, false, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "ai_timeout", "10"),
				),
			},
			// An empty list is accepted and plans clean
			{
				Config: testAccDatabaseResourceConfig(name+"_updated", false, 250, false, 60, false, "connection_parameters = []"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "connection_parameters.#", "0"),
				),
			},
			// A name the server reserves is refused with its reason
			{
				Config:      testAccDatabaseResourceConfig(name+"_updated", false, 250, false, 60, false, `connection_parameters = [{ name = "user", value = "postgres" }]`),
				ExpectError: regexp.MustCompile(`can't be a connection\s+parameter`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDatabaseResourceUnredactedRunsAllowed(t *testing.T) {
	name := fmt.Sprintf("database_%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// A configuration that sets it on creates the database with it on
			{
				Config: testAccDatabaseUnredactedRunsConfig(name, "unredacted_runs_allowed = true"),
				Check:  resource.TestCheckResourceAttr("devhub_querydesk_database.test", "unredacted_runs_allowed", "true"),
			},
			// Import reads the database's current value
			{
				ResourceName:            "devhub_querydesk_database.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credentials.0.password"},
			},
			// Changing it to off updates the database
			{
				Config: testAccDatabaseUnredactedRunsConfig(name, "unredacted_runs_allowed = false"),
				Check:  resource.TestCheckResourceAttr("devhub_querydesk_database.test", "unredacted_runs_allowed", "false"),
			},
			{
				Config: testAccDatabaseUnredactedRunsConfig(name, "unredacted_runs_allowed = true"),
				Check:  resource.TestCheckResourceAttr("devhub_querydesk_database.test", "unredacted_runs_allowed", "true"),
			},
			// Removing it from a configuration that had it on leaves the database on, which is
			// what a database that predates the setting sees. Each step fails on a plan that
			// is not empty after its apply, so this also proves no change is planned.
			{
				Config: testAccDatabaseUnredactedRunsConfig(name, ""),
				Check:  resource.TestCheckResourceAttr("devhub_querydesk_database.test", "unredacted_runs_allowed", "true"),
			},
			// An unrelated change still sends nothing for it
			{
				Config: testAccDatabaseUnredactedRunsConfig(name+"_renamed", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "name", name+"_renamed"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "unredacted_runs_allowed", "true"),
				),
			},
		},
	})
}

func TestAccDatabaseResourceUnredactedRunsAllowedOmitted(t *testing.T) {
	name := fmt.Sprintf("database_%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// A configuration that omits it creates the database with it off
			{
				Config: testAccDatabaseUnredactedRunsConfig(name, ""),
				Check:  resource.TestCheckResourceAttr("devhub_querydesk_database.test", "unredacted_runs_allowed", "false"),
			},
			// and plans no change afterwards
			{
				Config:   testAccDatabaseUnredactedRunsConfig(name, ""),
				PlanOnly: true,
			},
		},
	})
}

func testAccDatabaseUnredactedRunsConfig(name string, unredactedRunsAllowed string) string {
	return providerConfig + fmt.Sprintf(`
resource "devhub_querydesk_database" "test" {
  name     = %[1]q
  adapter  = "POSTGRES"
  hostname = "localhost"
  database = "mydb"
  %[2]s

  credentials = [
    {
      username           = "postgres"
      password           = "password"
      reviews_required   = 0
      default_credential = true
    }
  ]
}
`, name, unredactedRunsAllowed)
}

func testAccDatabaseResourceConfig(name string, aiEnabled bool, aiMaxRows int, aiAllowed bool, timeout int, reviewBypassAllowed bool, extraAttributes string) string {
	return providerConfig + fmt.Sprintf(`
resource "devhub_querydesk_database" "test" {
  name        = %[1]q
  adapter     = "POSTGRES"
  hostname    = "localhost"
  database    = "mydb"
  ai_enabled  = %[2]t
  ai_max_rows = %[3]d
  review_bypass_allowed = %[6]t
  %[7]s

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
`, name, aiEnabled, aiMaxRows, aiAllowed, timeout, reviewBypassAllowed, extraAttributes)
}

// A configuration that never mentions ai_timeout creates the database at 10, and the
// plan after the apply is empty, which every step asserts unless told otherwise.
func TestAccDatabaseResourceAiTimeoutDefault(t *testing.T) {
	name := fmt.Sprintf("database_%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseResourceConfig(name, true, 100, true, 30, false, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "ai_timeout", "10"),
					resource.TestCheckResourceAttr("devhub_querydesk_database.test", "credentials.0.timeout", "30"),
				),
			},
			{
				Config:   testAccDatabaseResourceConfig(name, true, 100, true, 30, false, ""),
				PlanOnly: true,
			},
		},
	})
}
