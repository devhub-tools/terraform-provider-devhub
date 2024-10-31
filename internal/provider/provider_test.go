// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	providerConfig = `
provider "devhub" {
	host    = "http://localhost:4000"
	api_key = "test"
}
`
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"devhub": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
}
