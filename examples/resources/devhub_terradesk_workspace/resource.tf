resource "devhub_terradesk_workspace" "example" {
  name         = "default"
  repository   = "devhub-tools/devhub"
  path         = "terraform"
  docker_image = "hashicorp/terraform:1.10"

  env_vars = [
    {
      name  = "ENV_VAR"
      value = "env-var-value"
    }
  ]

  secrets = [
    {
      name  = "my_secret"
      value = "secret-value"
    }
  ]

  workload_identity = {
    enabled               = true
    service_account_email = "devhub@my-project.iam.gserviceaccount.com"
    provider              = google_iam_workload_identity_pool_provider.devhub.name
  }
}

resource "google_iam_workload_identity_pool" "devhub" {
  project                   = google_project.default.project_id
  workload_identity_pool_id = "devhub"
  display_name              = "devhub"
}

resource "google_iam_workload_identity_pool_provider" "devhub" {
  project                            = google_project.default.project_id
  workload_identity_pool_id          = google_iam_workload_identity_pool.devhub.workload_identity_pool_id
  workload_identity_pool_provider_id = "devhub"

  attribute_mapping = {
    "google.subject"                   = "assertion.sub"
    "attribute.terradesk_workspace_id" = "assertion.terradesk_workspace_id"
  }

  oidc {
    issuer_uri = "https://devhub.example.com"
  }
}
