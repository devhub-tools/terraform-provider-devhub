resource "devhub_terradesk_workspace" "example" {
  name         = "default"
  repository   = "devhub-tools/devhub"
  path         = "terraform"
  docker_image = "hashicorp/terraform:1.10"

  # wrap every terraform command in `op run` so 1Password resolves op:// references
  command = "/home/terradesk/bin/op run -- terraform"

  # post to slack when a plan is waiting on a human, main branch only
  slack_channel   = "#terraform"
  slack_notify_on = "default_branch"

  # nightly refresh-only run reporting anything changed outside terraform
  drift_detection_enabled = true
  drift_cron_expression   = "0 2 * * *"
  drift_slack_channel     = "#infra-alerts"

  # stage the op binary onto the shared volume before terraform init runs
  init_containers = [
    {
      name    = "install-op"
      image   = "1password/op:2"
      command = "sh -c \"cp /usr/local/bin/op /home/terradesk/bin/\""
    }
  ]

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
    },
    {
      # raw secrets skip the TF_VAR_ prefix, so `op` can read its own token
      name  = "OP_SERVICE_ACCOUNT_TOKEN"
      value = var.op_service_account_token
      raw   = true
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
