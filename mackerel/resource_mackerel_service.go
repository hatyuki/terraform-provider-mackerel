package mackerel

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelService() *schema.Resource {
	return &schema.Resource{
		Create:   resourceMackerelServiceCreate,
		Read:     resourceMackerelServiceRead,
		Update:   nil,
		Delete:   resourceMackerelServiceDelete,
		Exists:   resourceMackerelServiceExists,
		Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough},

		SchemaVersion: 0,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"memo": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  `Managed by Terraform`,
			},
			"roles": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func resourceMackerelServiceCreate(d *schema.ResourceData, m interface{}) error {
	mkr := m.(*mackerel.Client)

	param := &mackerel.CreateServiceParam{
		Name: d.Get("name").(string),
		Memo: d.Get("memo").(string),
	}

	service, err := mkr.CreateService(param)
	if err != nil {
		return err
	}

	d.SetId(service.Name)

	return setMackerelServiceAttr(d, service)
}

func resourceMackerelServiceRead(d *schema.ResourceData, m interface{}) error {
	mkr := m.(*mackerel.Client)

	services, err := mkr.FindServices()
	if err != nil {
		return err
	}

	for _, service := range services {
		if service.Name == d.Id() {
			return setMackerelServiceAttr(d, service)
		}
	}

	return fmt.Errorf(`service "%s" not found`, d.Id())
}

func resourceMackerelServiceDelete(d *schema.ResourceData, m interface{}) error {
	mkr := m.(*mackerel.Client)
	_, err := mkr.DeleteService(d.Id())
	return err
}

func resourceMackerelServiceExists(d *schema.ResourceData, m interface{}) (bool, error) {
	mkr := m.(*mackerel.Client)

	if _, err := mkr.GetServiceMetaDataNameSpaces(d.Id()); err != nil {
		if err.(*mackerel.APIError).StatusCode != http.StatusNotFound {
			return false, err
		}
		return false, nil
	}

	return true, nil
}

func setMackerelServiceAttr(d *schema.ResourceData, service *mackerel.Service) error {
	if err := d.Set("name", service.Name); err != nil {
		return err
	}
	if err := d.Set("memo", service.Memo); err != nil {
		return err
	}
	if err := d.Set("roles", service.Roles); err != nil {
		return err
	}

	return nil
}
