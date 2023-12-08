package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccProvisionResource_Create(t *testing.T) {
	config := `
provider "provisioner" {
  sns_topic = "arn:aws:sns:us-east-2:809674927168:PublishDNS"
  region    = "us-east-2"
}

resource "provisioner_provision" "test" {
  name           = "testhost"
  instance_id    = "i-040c13fa937c0a51c"
  private_ip     = "172.31.4.211"
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
		},
	})
}

func TestAccProvisionResource_Update(t *testing.T) {
	config := `
provider "provisioner" {
  sns_topic = "arn:aws:sns:us-east-2:809674927168:PublishDNS"
  region    = "us-east-2"
}

resource "provisioner_provision" "test" {
  name           = "testhost"
  instance_id    = "i-040c13fa937c0a51c"
  private_ip     = "172.31.4.211"
}
`

	config_update := `
provider "provisioner" {
  sns_topic = "arn:aws:sns:us-east-2:809674927168:PublishDNS"
  region    = "us-east-2"
}

resource "provisioner_provision" "test" {
  name           = "testhost1"
  instance_id    = "i-040c13fa937c0a51c"
  private_ip     = "172.31.4.211"
}
`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				Config: config_update,
			},
		},
	})
}
