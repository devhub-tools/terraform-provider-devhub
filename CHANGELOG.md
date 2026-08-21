## 0.1.0 (Unreleased)

FEATURES:

* resource/devhub_terradesk_workspace: add `command` to replace the image's entrypoint on the init, plan and apply containers
* resource/devhub_terradesk_workspace: add `init_containers`, which run after the repository checkout and before `terraform init`
* resource/devhub_terradesk_workspace: add `raw` to `secrets`, exposing the value under its literal name instead of `TF_VAR_<name>`
* resource/devhub_terradesk_workspace: add `slack_channel` and `slack_notify_on`, posting to slack when a plan is waiting on a human
