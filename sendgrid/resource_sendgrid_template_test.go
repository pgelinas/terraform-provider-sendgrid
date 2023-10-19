package sendgrid_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	sendgrid "github.com/taharah/terraform-provider-sendgrid/sdk"
)

func TestAccSendgridTemplateBasic(t *testing.T) {
	name := "terraform-template-" + acctest.RandString(10)
	generation := "dynamic"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSendgridTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSendgridTemplateConfigBasic(name, generation),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridTemplateExists("sendgrid_template.this", name),
				),
			},
		},
	})
}

func testAccCheckSendgridTemplateDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*sendgrid.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sendgrid_template" {
			continue
		}

		templateID := rs.Primary.ID

		if _, err := c.DeleteTemplate(context.Background(), templateID); err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckSendgridTemplateConfigBasic(name, generation string) string {
	return fmt.Sprintf(`
resource "sendgrid_template" "this" {
  name = %q
  generation = %q
}`, name, generation)
}

func testAccCheckSendgridTemplateExists(resource, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No template ID set: %s", resource)
		}

		if rs.Primary.Attributes["name"] != name {
			return fmt.Errorf("Name incorrect: %s != %s", rs.Primary.Attributes["name"], name)
		}

		return nil
	}
}
