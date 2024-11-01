terraform {
  required_providers {
    devhub = {
      source = "registry.terraform.io/devhub/devhub"
    }
  }
}

provider "devhub" {
  api_key = "dh_b3JnXzAxSjIyR0M5NUpRS0JGRzI0WE5LQjRNQjUyOr0dxLoeRj8hSUvpCViuHdEHXM-hid12JKckd_4DWFv1"
  host    = "http://localhost:4000"
}
