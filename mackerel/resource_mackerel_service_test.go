package mackerel

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mackerelio/mackerel-client-go"
)

func TestAccMackerelService_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkMackerelServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMackerelServiceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					checkMackerelServiceExists,
					resource.TestCheckResourceAttr("mackerel_service.example", "name", "ExampleService"),
					resource.TestCheckResourceAttr("mackerel_service.example", "memo", "This is an example."),
					resource.TestCheckResourceAttr("mackerel_service.example", "roles.#", "0"),
				),
			},
		},
	})
}

func TestAccMackerelService_Import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkMackerelServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMackerelServiceConfigImport,
			},
			{
				ResourceName:      "mackerel_service.example",
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mackerel_service.example", "memo", "Managed by Terraform"),
				),
			},
		},
	})
}

func checkMackerelServiceExists(s *terraform.State) error {
	mkr := testAccProvider.Meta().(*mackerel.Client)

	for _, r := range s.RootModule().Resources {
		if ok, err := mackerelServiceExistsHelper(mkr, r.Primary.ID); err != nil {
			return fmt.Errorf("received an error retrieving service: %s", err)
		} else if !ok {
			return fmt.Errorf("the specified service does not exist")
		}
	}

	return nil
}

func checkMackerelServiceDestroy(s *terraform.State) error {
	mkr := testAccProvider.Meta().(*mackerel.Client)

	for _, r := range s.RootModule().Resources {
		if ok, err := mackerelServiceExistsHelper(mkr, r.Primary.ID); err != nil {
			return fmt.Errorf("received an error retrieving service: %s", err)
		} else if ok {
			return fmt.Errorf("the specified service still exists")
		}
	}

	return nil
}

func mackerelServiceExistsHelper(mkr *mackerel.Client, name string) (bool, error) {
	if _, err := mkr.GetServiceMetaDataNameSpaces(name); err != nil {
		if err.(*mackerel.APIError).StatusCode != http.StatusNotFound {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

const testAccCheckMackerelServiceConfigBasic = `
resource "mackerel_service" "example" {
  name = "ExampleService"
  memo = "This is an example."
}
`

const testAccCheckMackerelServiceConfigImport = `
resource "mackerel_service" "example" {
  name = "ExampleService"
}
`
