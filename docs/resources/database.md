---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "querydesk_database Resource - terraform-provider-querydesk"
subcategory: ""
description: |-
  Database resource
---

# querydesk_database (Resource)

Database resource

## Example Usage

```terraform
resource "querydesk_database" "example" {
  name     = "terraform_test"
  adapter  = "postgres"
  hostname = "localhost"
  database = "mydb"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `adapter` (String) The adapter to use to establish the connection. Currently only `postgres` and `mysql` are supported, but  sql server is on the roadmap.
- `database` (String) The name of the database to connect to.
- `hostname` (String) The hostname for connecting to the database, either an ip or url.
- `name` (String) The name for users to use to identity the database.

### Optional

- `cacertfile` (String, Sensitive) The server ca cert to use with ssl connections, `ssl` must be set to `true`.
- `certfile` (String, Sensitive) The client cert to use with ssl connections, `ssl` must be set to `true`.
- `keyfile` (String, Sensitive) The client key to use with ssl connections, `ssl` must be set to `true`.
- `restrict_access` (Boolean) Whether access to this databases should be explicitly granted to users or if any authenticated user can access it.
- `ssl` (Boolean) Set to `true` to turn on ssl connections for this database.

### Read-Only

- `id` (String) Database id.