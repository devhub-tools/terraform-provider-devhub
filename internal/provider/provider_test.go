// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var providerConfig = fmt.Sprintf(`
provider "devhub" {
	host = %q
}
`, os.Getenv("DEVHUB_HOST"))

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"devhub": providerserver.NewProtocol6WithError(New("test")()),
}
