package mackerel

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mackerelio/mackerel-client-go"
)

func TestAccMackerelRole_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(checkMackerelServiceDestroy, checkMackerelRoleDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMackerelRoleConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					checkMackerelRoleExists,
					resource.TestCheckResourceAttr("mackerel_role.example", "service_name", "ExampleService"),
					resource.TestCheckResourceAttr("mackerel_role.example", "name", "ExampleRole"),
					resource.TestCheckResourceAttr("mackerel_role.example", "memo", "This is an example."),
				),
			},
		},
	})
}

func TestAccMackerelRole_Import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(checkMackerelServiceDestroy, checkMackerelRoleDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMackerelRoleConfigImport,
			},
			{
				ResourceName:      "mackerel_role.example",
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mackerel_role.example", "memo", "Managed by Terraform"),
				),
			},
		},
	})
}

func checkMackerelRoleExists(s *terraform.State) error {
	mkr := testAccProvider.Meta().(*mackerel.Client)

	for _, r := range s.RootModule().Resources {
		if r.Type != "mackerel_role" {
			continue
		}

		if ok, err := mackerelRoleExistsHelper(mkr, r.Primary.ID); err != nil {
			return fmt.Errorf("received an error retrieving role: %s", err)
		} else if !ok {
			return fmt.Errorf("the specified role does not exist")
		}
	}

	return nil
}

func checkMackerelRoleDestroy(s *terraform.State) error {
	mkr := testAccProvider.Meta().(*mackerel.Client)

	for _, r := range s.RootModule().Resources {
		if r.Type != "mackerel_role" {
			continue
		}

		if ok, err := mackerelRoleExistsHelper(mkr, r.Primary.ID); err != nil {
			return fmt.Errorf("received an error retrieving role: %s", err)
		} else if ok {
			return fmt.Errorf("the specified role still exists")
		}
	}

	return nil
}

func mackerelRoleExistsHelper(mkr *mackerel.Client, id string) (bool, error) {
	name := strings.Split(id, ":")
	if _, err := mkr.GetRoleMetaDataNameSpaces(name[0], name[1]); err != nil {
		if err.(*mackerel.APIError).StatusCode != http.StatusNotFound {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

const testAccCheckMackerelRoleConfigBasic = `
resource "mackerel_service" "example" {
  name = "ExampleService"
}

resource "mackerel_role" "example" {
  service_name = "${mackerel_service.example.name}"
  name         = "ExampleRole"
  memo         = "This is an example."
}
`

const testAccCheckMackerelRoleConfigImport = `
resource "mackerel_service" "example" {
  name = "ExampleService"
}

resource "mackerel_role" "example" {
  service_name = "${mackerel_service.example.name}"
  name         = "ExampleRole"
}
`
