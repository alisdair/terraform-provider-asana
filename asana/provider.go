package asana

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// Provider creates the Asana provider
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"asana_project": resourceProject(),
		},
	}
}
