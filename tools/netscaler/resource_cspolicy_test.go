/*
Copyright 2016 Citrix Systems, Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package netscaler

import (
	"fmt"
	"github.com/chiradeep/go-nitro/netscaler"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccCspolicy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCspolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCspolicy_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCspolicyExist("netscaler_cspolicy.foo", nil),

					resource.TestCheckResourceAttr(
						"netscaler_cspolicy.foo", "lbvserver", "test-lb-vserver"),
					resource.TestCheckResourceAttr(
						"netscaler_cspolicy.foo", "policyname", "test_policy"),
					resource.TestCheckResourceAttr(
						"netscaler_cspolicy.foo", "rule", "CLIENT.IP.SRC.SUBNET(24).EQ(10.217.84.0)"),
				),
			},
		},
	})
}

func testAccCheckCspolicyExist(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No lb vserver name is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed!")
			}

			*id = rs.Primary.ID
		}

		nsClient := testAccProvider.Meta().(*NetScalerNitroClient).client
		data, err := nsClient.FindResource(netscaler.Cspolicy.Type(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if data == nil {
			return fmt.Errorf("LB vserver %s not found", n)
		}

		return nil
	}
}

func testAccCheckCspolicyDestroy(s *terraform.State) error {
	nsClient := testAccProvider.Meta().(*NetScalerNitroClient).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netscaler_cspolicy" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No name is set")
		}

		_, err := nsClient.FindResource(netscaler.Cspolicy.Type(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("LB vserver %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

const testAccCspolicy_basic = `

resource "netscaler_csvserver" "foo" {
  
  ipv46 = "10.202.11.11"
  name = "terraform-cs"
  port = 443
  servicetype = "SSL"
}


resource "netscaler_cspolicy" "foo" {
  
  lbvserver = "test-lb-vserver"
  policyname = "test_policy"
  rule = "CLIENT.IP.SRC.SUBNET(24).EQ(10.217.84.0)"

  depends_on = ["netscaler_csvserver.foo"]

}
`
