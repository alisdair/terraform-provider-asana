package asana

import (
	"encoding/json"
	"github.com/alisdair/asana/v1"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strconv"
)

func resourceTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceTaskCreate,
		Read:   resourceTaskRead,
		Delete: resourceTaskDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the task. This is generally a short sentence fragment that fits on a line in the UI for maximum readability. However, it can be longer.",
				ForceNew:    true,
			},
			"notes": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "More detailed, free-form textual information associated with the task.",
				ForceNew:    true,
			},
			"workspace": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "The workspace this task is associated with. Once created, task cannot be moved to a different workspace. This attribute can only be specified at creation time. Not needed if projects are specified.",
				ForceNew:    true,
			},
			"projects": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Array of projects this task is associated with. At task creation time, this array can be used to add the task to many projects at once.",
				ForceNew:    true,
			},
		},
	}
}

func resourceTaskCreate(d *schema.ResourceData, m interface{}) error {
	client, err := asana.NewClient()
	if err != nil {
		return err
	}

	request := &asana.TaskRequest{
		Name:      d.Get("name").(string),
		Notes:     d.Get("notes").(string),
		Workspace: d.Get("workspace").(string),
	}

	if v, ok := d.GetOk("projects"); ok {
		projects := make([]int64, len(v.([]interface{})))
		for i, project := range v.([]interface{}) {
			id, err := strconv.ParseInt(project.(string), 10, 64)
			if err != nil {
				return err
			}
			projects[i] = id
		}
		request.Projects = projects
	}
	log.Printf("task request: %v\n", request)
	jsonRequest, _ := json.Marshal(request)
	log.Printf("task request json: %s\n", jsonRequest)

	task, err := client.CreateTask(request)
	if err != nil {
		log.Printf("create task error: %v\n", err)
		return err
	}

	d.SetId(strconv.FormatInt(task.ID, 10))

	return nil
}

func resourceTaskRead(d *schema.ResourceData, m interface{}) error {
	client, err := asana.NewClient()
	if err != nil {
		return err
	}

	task, err := client.FindTaskByID(d.Id())

	// If update returns a 404, the task has been deleted
	if httpErr, ok := err.(*asana.HTTPError); ok && httpErr.Code() == 404 {
		d.SetId("")
		return nil
	} else if err != nil {
		return err
	}

	d.Set("name", task.Name)
	d.Set("notes", task.Notes)
	d.Set("workspace", strconv.FormatInt(task.Workspace.ID, 10))

	var projects []string
	for _, project := range task.Projects {
		projects = append(projects, strconv.FormatInt(project.ID, 10))
	}
	d.Set("projects", projects)

	return nil
}

func resourceTaskDelete(d *schema.ResourceData, m interface{}) error {
	client, err := asana.NewClient()
	if err != nil {
		return err
	}

	err = client.DeleteTask(d.Id())

	// Consider 404 a success
	if httpErr, ok := err.(*asana.HTTPError); ok && httpErr.Code() == 404 {
		return nil
	}

	return err
}
