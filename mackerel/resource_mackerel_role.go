package mackerel

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelRole() *schema.Resource {
	return &schema.Resource{
		Create:   resourceMackerelRoleCreate,
		Read:     resourceMackerelRoleRead,
		Update:   nil,
		Delete:   resourceMackerelRoleDelete,
		Exists:   resourceMackerelRoleExists,
		Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough},

		SchemaVersion: 0,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
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
		},
	}
}

func resourceMackerelRoleCreate(d *schema.ResourceData, m interface{}) error {
	mkr := m.(*mackerel.Client)

	serviceName := d.Get("service_name").(string)
	param := &mackerel.CreateRoleParam{
		Name: d.Get("name").(string),
		Memo: d.Get("memo").(string),
	}

	role, err := mkr.CreateRole(serviceName, param)
	if err != nil {
		return err
	}

	d.SetId(serviceName + ":" + role.Name)
	return setMackerelRoleAttr(d, role)
}

func resourceMackerelRoleRead(d *schema.ResourceData, m interface{}) error {
	mkr := m.(*mackerel.Client)

	name := strings.Split(d.Id(), ":")
	serviceName, roleName := name[0], name[1]

	roles, err := mkr.FindRoles(serviceName)
	if err != nil {
		return err
	}

	for _, role := range roles {
		if role.Name == roleName {
			if err := d.Set("service_name", serviceName); err != nil {
				return err
			}
			return setMackerelRoleAttr(d, role)
		}
	}

	return fmt.Errorf(`role "%s" not found`, roleName)
}

func resourceMackerelRoleDelete(d *schema.ResourceData, m interface{}) error {
	mkr := m.(*mackerel.Client)

	name := strings.Split(d.Id(), ":")
	serviceName, roleName := name[0], name[1]
	if _, err := mkr.DeleteRole(serviceName, roleName); err != nil {
		return err
	}
	return nil
}

func resourceMackerelRoleExists(d *schema.ResourceData, m interface{}) (bool, error) {
	mkr := m.(*mackerel.Client)

	name := strings.Split(d.Id(), ":")
	serviceName, roleName := name[0], name[1]
	if _, err := mkr.GetRoleMetaDataNameSpaces(serviceName, roleName); err != nil {
		if err.(*mackerel.APIError).StatusCode != http.StatusNotFound {
			return false, err
		}
		return false, nil
	}

	return true, nil
}

func setMackerelRoleAttr(d *schema.ResourceData, role *mackerel.Role) error {
	if err := d.Set("name", role.Name); err != nil {
		return err
	}
	if err := d.Set("memo", role.Memo); err != nil {
		return err
	}

	return nil
}
