package asana

import (
	"github.com/alisdair/asana/v1"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"strconv"
)

var (
	projectColors = []string{
		"dark-pink", "dark-green", "dark-blue", "dark-red",
		"dark-teal", "dark-brown", "dark-orange", "dark-purple",
		"dark-warm-gray", "light-pink", "light-green", "light-blue",
		"light-red", "light-teal", "light-yellow", "light-orange",
		"light-purple", "light-warm-gray",
	}
	projectLayouts = []string{string(asana.ListLayout), string(asana.BoardLayout)}
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the project. This is generally a short sentence fragment that fits on a line in the UI for maximum readability.",
			},
			"notes": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "More detailed, free-form textual information associated with the project.",
			},
			"color": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				Description:  "Color of the project. Must be either null or one of: dark-pink, dark-green, dark-blue, dark-red, dark-teal, dark-brown, dark-orange, dark-purple, dark-warm-gray, light-pink, light-green, light-blue, light-red, light-teal, light-yellow, light-orange, light-purple, light-warm-gray.",
				ValidateFunc: validation.StringInSlice(projectColors, false),
			},
			"layout": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      asana.BoardLayout,
				Description:  "The layout (board or list view) of the project.",
				ValidateFunc: validation.StringInSlice(projectLayouts, false),
			},
			"workspace": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The workspace or organization that the project belongs to.",
			},
			"public": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "True if the project is public to the organization.",
			},
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	client, err := asana.NewClient()
	if err != nil {
		return err
	}

	request := &asana.ProjectRequest{
		Name:                 d.Get("name").(string),
		Notes:                d.Get("notes").(string),
		Color:                d.Get("color").(string),
		Layout:               asana.Layout(d.Get("layout").(string)),
		Workspace:            d.Get("workspace").(string),
		PublicToOrganization: d.Get("public").(bool),
	}

	project, err := client.CreateProject(request)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(project.ID, 10))

	return nil
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	client, err := asana.NewClient()
	if err != nil {
		return err
	}

	project, err := client.FindProjectByID(d.Id())

	// If update returns a 404, the project has been deleted
	if httpErr, ok := err.(*asana.HTTPError); ok && httpErr.Code() == 404 {
		d.SetId("")
		return nil
	} else if err != nil {
		return err
	}

	d.Set("name", project.Name)
	d.Set("notes", project.Notes)
	d.Set("color", project.Color)
	// Public is not returned by the API
	d.Set("workspace", strconv.FormatInt(project.Workspace.ID, 10))

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	client, err := asana.NewClient()
	if err != nil {
		return err
	}

	request := &asana.ProjectRequest{
		ProjectID:            d.Id(),
		Name:                 d.Get("name").(string),
		Notes:                d.Get("notes").(string),
		Color:                d.Get("color").(string),
		PublicToOrganization: d.Get("public").(bool),
	}
	_, err = client.UpdateProject(request)

	// If update returns a 404, the project has been deleted
	if httpErr, ok := err.(*asana.HTTPError); ok && httpErr.Code() == 404 {
		d.SetId("")
		return nil
	}

	return err
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	client, err := asana.NewClient()
	if err != nil {
		return err
	}

	err = client.DeleteProjectByID(d.Id())

	// Consider 404 a success
	if httpErr, ok := err.(*asana.HTTPError); ok && httpErr.Code() == 404 {
		return nil
	}

	return err
}
