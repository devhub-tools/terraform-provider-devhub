// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	providerConfig = `
provider "devhub" {
	host    = "http://localhost:4000"
	api_key = "dh_b3JnXzAxSjIyR0M5NUpRS0JGRzI0WE5LQjRNQjUyOr0dxLoeRj8hSUvpCViuHdEHXM-hid12JKckd_4DWFv1"
}
`
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"devhub": providerserver.NewProtocol6WithError(New("test")()),
}
