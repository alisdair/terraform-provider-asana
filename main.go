package main

import (
	"github.com/alisdair/terraform-provider-asana/asana"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return asana.Provider()
		},
	})
}
