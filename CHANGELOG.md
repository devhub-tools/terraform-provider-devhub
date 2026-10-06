## 2.44.0 (October 5, 2026)

BREAKING CHANGES:

* resource/devhub_querydesk_database: rename `agent_id` to `tunnel_id`
* resource/devhub_terradesk_workspace: rename `agent_id` to `tunnel_id`

FEATURES:

* resource/devhub_querydesk_database: add `connection_parameters`, the names and values a Postgres database sends when a connection opens, such as `pgdog.role` = `replica`
* resource/devhub_terradesk_workspace: add `drift_detection_enabled`, `drift_cron_expression` and `drift_slack_channel` to run drift detection on a schedule
* resource/devhub_terradesk_workspace: add `trigger_paths`, additional repository paths that also start a plan

## 2.39.1 (August 20, 2026)

FEATURES:

* resource/devhub_terradesk_workspace: add `slack_channel` and `slack_notify_on`, posting to slack when a plan is waiting on a human

## 2.39.0 (August 19, 2026)

FEATURES:

* resource/devhub_querydesk_database: add `review_bypass_allowed`, letting a super admin run their own query that hasn't met `reviews_required`
* resource/devhub_querydesk_database: add `verify_server_hostname` to verify the server certificate matches the hostname
* resource/devhub_terradesk_workspace: add `command` to replace the image's entrypoint on the init, plan and apply containers
* resource/devhub_terradesk_workspace: add `init_containers`, which run after the repository checkout and before `terraform init`
* resource/devhub_terradesk_workspace: add `raw` to `secrets`, exposing the value under its literal name instead of `TF_VAR_<name>`

## 2.34.0 (July 27, 2026)

BREAKING CHANGES:

* resource/devhub_querydesk_database: remove `hostname` from `credentials`

FEATURES:

* resource/devhub_dashboard: add `actions` to `panels`, row actions that trigger a workflow

## 2.31.1 (July 8, 2026)

BUG FIXES:

* provider: prefer provider block attributes over environment variables

## 2.31.0 (July 1, 2026)

FEATURES:

* resource/devhub_querydesk_database: add `ai_enabled` and `ai_max_rows`
* resource/devhub_querydesk_database: add `ai_allowed` and `timeout` to `credentials`

## 2.24.1 (March 31, 2026)

BREAKING CHANGES:

* resource/devhub_workflow: remove the read-only `id` from `steps`

## 2.23.0 (March 9, 2026)

FEATURES:

* resource/devhub_querydesk_database: support the `SQLSERVER` and `ORACLE` adapters
* resource/devhub_workflow: add `group` to organize workflows into folders
* resource/devhub_workflow: add `required` to `inputs`

BUG FIXES:

* resource/devhub_workflow: fix a workflow with no inputs

## 2.10.3 (July 9, 2025)

BUG FIXES:

* resource/devhub_querydesk_database: send an unset `port` as null

## 2.10.2 (June 23, 2025)

BUG FIXES:

* resource/devhub_workflow: read `cron_schedule` back correctly

## 2.10.0 (June 20, 2025)

NOTES:

* Version numbers now follow the Devhub release they shipped with.

FEATURES:

* resource/devhub_workflow: add `cron_schedule` to trigger a workflow on a schedule
* resource/devhub_workflow: add `condition` to `steps`, and a `condition_action` step

## 1.4.0 (June 16, 2025)

FEATURES:

* resource/devhub_dashboard: add `restricted_access`
* resource/devhub_querydesk_database: add `port`
* resource/devhub_querydesk_database: add `hostname` to `credentials`

## 1.3.1 (June 11, 2025)

FEATURES:

* resource/devhub_querydesk_database: support the `CLICKHOUSE` adapter

## 1.3.0 (May 28, 2025)

FEATURES:

* **New Data Source:** `devhub_role`
* **New Data Source:** `devhub_user`
* resource/devhub_workflow: add `permissions` to `approval_action`

## 1.2.5 (April 30, 2025)

BUG FIXES:

* resource/devhub_workflow: allow a step's `name` to be empty

## 1.2.4 (April 30, 2025)

BUG FIXES:

* resource/devhub_workflow: handle an empty `description`

## 1.2.3 (April 30, 2025)

BREAKING CHANGES:

* resource/devhub_workflow: rename `required_approvals` to `reviews_required` on `approval_action`

FEATURES:

* resource/devhub_querydesk_database: add the read-only `credential_ids`, a map of credential ids by username

## 1.2.2 (April 28, 2025)

FEATURES:

* resource/devhub_workflow: add `trigger_linear_label_name`

## 1.2.1 (April 28, 2025)

BUG FIXES:

* resource/devhub_dashboard, resource/devhub_workflow: loosen attribute validation

## 1.2.0 (April 28, 2025)

FEATURES:

* **New Resource:** `devhub_dashboard`
* **New Resource:** `devhub_workflow`

## 1.1.1 (March 21, 2025)

BREAKING CHANGES:

* resource/devhub_querydesk_database: remove `enable_data_protection`

## 1.1.0 (February 17, 2025)

BREAKING CHANGES:

* resource/devhub_querydesk_database: remove `slack_webhook_url`

FEATURES:

* **New Resource:** `devhub_terradesk_workspace`

## 1.0.1 (November 4, 2024)

FEATURES:

* resource/devhub_querydesk_database: add `agent_id`

## 1.0.0 (November 1, 2024)

BREAKING CHANGES:

* The provider now targets Devhub. `database` is renamed to `devhub_querydesk_database`, and `database_user` is replaced by its `credentials` list.

FEATURES:

* provider: `api_key` can be set with the `DEVHUB_API_KEY` environment variable
* resource/devhub_querydesk_database: add `group` and `slack_channel`

## 0.2.0 (June 30, 2023)

FEATURES:

* **New Resource:** `database_user`

## 0.1.0 (June 29, 2023)

FEATURES:

* **New Resource:** `database`
