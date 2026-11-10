package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// fakeDevhub is an in-memory Devhub API that keeps every request body, so a test
// can see what the provider sent for values the state never holds.
type fakeDevhub struct {
	mu       sync.Mutex
	nextId   int
	objects  map[string]map[string]any
	requests []map[string]any
}

func newFakeDevhub(t *testing.T) (*fakeDevhub, string) {
	fake := &fakeDevhub{objects: map[string]map[string]any{}}
	server := httptest.NewServer(http.HandlerFunc(fake.handle))
	t.Cleanup(server.Close)

	return fake, server.URL
}

func (f *fakeDevhub) handle(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	segments := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/"), "/")
	id := ""
	if len(segments) == 3 {
		id = segments[2]
	}

	switch r.Method {
	case http.MethodGet:
		object, ok := f.objects[id]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		_ = json.NewEncoder(w).Encode(object)

	case http.MethodDelete:
		delete(f.objects, id)

	case http.MethodPost, http.MethodPatch:
		raw, _ := io.ReadAll(r.Body)

		var body, response map[string]any
		_ = json.Unmarshal(raw, &body)
		_ = json.Unmarshal(raw, &response)
		f.requests = append(f.requests, body)

		if id == "" {
			f.nextId++
			id = fmt.Sprintf("id_%d", f.nextId)
		}

		response["id"] = id

		for _, key := range []string{"credentials", "secrets", "env_vars"} {
			items, _ := response[key].([]any)
			for index, item := range items {
				fields := item.(map[string]any)
				// the API never returns these
				delete(fields, "password")
				delete(fields, "value")
				fields["id"] = fmt.Sprintf("%s_%s_%d", id, key, index)
			}
		}

		f.objects[id] = response
		_ = json.NewEncoder(w).Encode(response)
	}
}

// sent is the field of the nested key's first element in the last create or update.
func (f *fakeDevhub) sent(t *testing.T, key, field string) string {
	f.mu.Lock()
	defer f.mu.Unlock()

	if len(f.requests) == 0 {
		t.Fatal("no create or update reached the API")
	}

	last := f.requests[len(f.requests)-1]
	if key == "" {
		value, _ := last[field].(string)
		return value
	}

	items := last[key].([]any)
	value, _ := items[0].(map[string]any)[field].(string)
	return value
}

func (f *fakeDevhub) requestCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()

	return len(f.requests)
}

func checkSent(t *testing.T, fake *fakeDevhub, key, field, want string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if got := fake.sent(t, key, field); got != want {
			return fmt.Errorf("sent %s %s %q, want %q", key, field, got, want)
		}

		return nil
	}
}

func checkRequestCount(fake *fakeDevhub, want int) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if got := fake.requestCount(); got != want {
			return fmt.Errorf("API got %d creates and updates, want %d", got, want)
		}

		return nil
	}
}

func fakeProviderConfig(host string) string {
	return fmt.Sprintf(`
provider "devhub" {
  host    = %q
  api_key = "test"
}
`, host)
}

func TestWorkspaceSecretValueWo(t *testing.T) {
	fake, host := newFakeDevhub(t)

	config := func(value string, version int) string {
		return fakeProviderConfig(host) + fmt.Sprintf(`
resource "devhub_terradesk_workspace" "test" {
  name         = "write_only"
  repository   = "devhub-tools/devhub"
  docker_image = "hashicorp/terraform:1.10"

  secrets = [
    {
      name             = "my_secret"
      value_wo         = %q
      value_wo_version = %d
    }
  ]
}
`, value, version)
	}

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config("first", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkSent(t, fake, "secrets", "value", "first"),
					resource.TestCheckNoResourceAttr("devhub_terradesk_workspace.test", "secrets.0.value_wo"),
					resource.TestCheckNoResourceAttr("devhub_terradesk_workspace.test", "secrets.0.value"),
					resource.TestCheckResourceAttr("devhub_terradesk_workspace.test", "secrets.0.value_wo_version", "1"),
				),
			},
			// a new value alone isn't a change Terraform can see
			{
				Config: config("second", 1),
				Check:  checkRequestCount(fake, 1),
			},
			// bumping the version sends it
			{
				Config: config("second", 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkRequestCount(fake, 2),
					checkSent(t, fake, "secrets", "value", "second"),
					resource.TestCheckNoResourceAttr("devhub_terradesk_workspace.test", "secrets.0.value_wo"),
				),
			},
		},
	})
}

func TestDatabasePasswordWo(t *testing.T) {
	fake, host := newFakeDevhub(t)

	config := func(password string, version int) string {
		return fakeProviderConfig(host) + fmt.Sprintf(`
resource "devhub_querydesk_database" "test" {
  name     = "write_only"
  adapter  = "POSTGRES"
  hostname = "localhost"
  database = "postgres"

  keyfile_wo         = "key-%[1]s"
  keyfile_wo_version = %[2]d

  credentials = [
    {
      username            = "postgres"
      password_wo         = %[1]q
      password_wo_version = %[2]d
      reviews_required    = 0
    }
  ]
}
`, password, version)
	}

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config("first", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkSent(t, fake, "credentials", "password", "first"),
					checkSent(t, fake, "", "keyfile", "key-first"),
					resource.TestCheckNoResourceAttr("devhub_querydesk_database.test", "credentials.0.password_wo"),
					resource.TestCheckNoResourceAttr("devhub_querydesk_database.test", "credentials.0.password"),
					resource.TestCheckNoResourceAttr("devhub_querydesk_database.test", "keyfile_wo"),
					resource.TestCheckNoResourceAttr("devhub_querydesk_database.test", "keyfile"),
				),
			},
			{
				Config: config("second", 1),
				Check:  checkRequestCount(fake, 1),
			},
			{
				Config: config("second", 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkRequestCount(fake, 2),
					checkSent(t, fake, "credentials", "password", "second"),
					checkSent(t, fake, "", "keyfile", "key-second"),
				),
			},
		},
	})
}

func TestWorkspaceSecretValueValidation(t *testing.T) {
	_, host := newFakeDevhub(t)

	config := func(secrets string) string {
		return fakeProviderConfig(host) + fmt.Sprintf(`
resource "devhub_terradesk_workspace" "test" {
  name         = "write_only"
  repository   = "devhub-tools/devhub"
  docker_image = "hashicorp/terraform:1.10"
  %s
}
`, secrets)
	}

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config(`secrets = [{ name = "s", value = "a", value_wo = "b" }]`),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
			{
				Config:      config(`secrets = [{ name = "s" }]`),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
			{
				Config:      config(`secrets = [{ name = "s", value = "a", value_wo_version = 1 }]`),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
			// with no default, a workspace without secrets reads back as null, not []
			{
				Config: config(""),
				Check:  resource.TestCheckNoResourceAttr("devhub_terradesk_workspace.test", "secrets.#"),
			},
			{
				Config:   config(""),
				PlanOnly: true,
			},
		},
	})
}
