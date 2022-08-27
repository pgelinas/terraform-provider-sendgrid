package sendgrid_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sendgrid "github.com/taharah/terraform-provider-sendgrid/sdk"
)

func TestAccSendgridAPIKeyBasic(t *testing.T) {
	name := "terraform-api-key-" + acctest.RandString(10)
	scopes := []string{"mail.send"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSendgridAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSendgridAPIKeyConfigBasic(name, scopes),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridAPIKeyExists("sendgrid_api_key.this", name),
				),
			},
		},
	})
}

func TestAccSendgridAPIKeyUpdateScopesAndName(t *testing.T) {
	name := "terraform-api-key-" + acctest.RandString(10)
	newName := "terraform-api-key-" + acctest.RandString(10)
	scopes := []string{"mail.send"}
	newScopes := []string{"mail.send", "sender_verification_eligible"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSendgridAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSendgridAPIKeyConfigBasic(name, scopes),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridAPIKeyExists("sendgrid_api_key.this", name),
				),
			},
			{
				Config: testAccCheckSendgridAPIKeyConfigBasic(name, newScopes),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridAPIKeyExists("sendgrid_api_key.this", name),
				),
			},
			{
				Config: testAccCheckSendgridAPIKeyConfigBasic(newName, newScopes),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridAPIKeyExists("sendgrid_api_key.this", newName),
				),
			},
		},
	})
}

func TestAccSendgridAPIKeyUpdateName(t *testing.T) {
	name := "terraform-api-key-" + acctest.RandString(10)
	newName := "terraform-api-key-" + acctest.RandString(10)
	scopes := []string{"mail.send"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSendgridAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSendgridAPIKeyConfigBasic(name, scopes),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridAPIKeyExists("sendgrid_api_key.this", name),
				),
			},
			{
				Config: testAccCheckSendgridAPIKeyConfigBasic(newName, scopes),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridAPIKeyExists("sendgrid_api_key.this", newName),
				),
			},
		},
	})
}

func testAccCheckSendgridAPIKeyDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*sendgrid.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sendgrid_api_key" {
			continue
		}

		apiKeyID := rs.Primary.ID

		if _, err := c.DeleteAPIKey(apiKeyID); err != nil {
			return err.Err
		}
	}

	return nil
}

func testAccCheckSendgridAPIKeyConfigBasic(name string, scopes []string) string {
	return fmt.Sprintf(`
resource "sendgrid_api_key" "this" {
  name = %q
  scopes = %s
}`, name, formatResourceList(scopes))
}

func testAccCheckSendgridAPIKeyExists(resource, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]

		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No apiKeyID set")
		}

		if rs.Primary.Attributes["name"] != name {
			return fmt.Errorf("Name incorrect: %s != %s", rs.Primary.Attributes["name"], name)
		}
		return nil
	}
}
