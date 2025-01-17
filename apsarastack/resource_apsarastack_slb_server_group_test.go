package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccApsaraStackSlbServerGroup_vpc(t *testing.T) {
	var v *slb.DescribeVServerGroupAttributeResponse
	resourceId := "apsarastack_slb_server_group.default"
	ra := resourceAttrInit(resourceId, serverGroupMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccSlbServerGroupVpc")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbServerGroupDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${apsarastack_slb.default.id}",
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${apsarastack_instance.default.0.id}", "${apsarastack_instance.default.1.id}"},
							"port":       "100",
							"weight":     "10",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "tf-testAccSlbServerGroupVpcUpdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccSlbServerGroupVpcUpdate",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${apsarastack_instance.default.0.id}", "${apsarastack_instance.default.1.id}"},
							"port":       "100",
							"weight":     "10",
						},
						{
							"server_ids": []string{"${apsarastack_network_interface.default.0.id}"},
							"port":       "70",
							"weight":     "10",
							"type":       "eni",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#": "2",
					}),
				),
				ExpectNonEmptyPlan: true,
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}",
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${apsarastack_instance.default.0.id}", "${apsarastack_instance.default.1.id}"},
							"port":       "100",
							"weight":     "10",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":      "tf-testAccSlbServerGroupVpc",
						"servers.#": "1",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccApsaraStackSlbServerGroup_multi_vpc(t *testing.T) {
	var v *slb.DescribeVServerGroupAttributeResponse
	resourceId := "apsarastack_slb_server_group.default.1"
	ra := resourceAttrInit(resourceId, serverGroupMultiClassicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccSlbServerGroupVpc")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbServerGroupMultiVpcDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${apsarastack_slb.default.id}",
					"count":            "2",
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${apsarastack_instance.default.0.id}", "${apsarastack_instance.default.1.id}"},
							"port":       "100",
							"weight":     "10",
						},
						{
							"server_ids": "${apsarastack_instance.default.*.id}",
							"port":       "80",
							"weight":     "100",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccApsaraStackSlbServerGroup_classic(t *testing.T) {
	var v *slb.DescribeVServerGroupAttributeResponse
	resourceId := "apsarastack_slb_server_group.default"
	ra := resourceAttrInit(resourceId, serverGroupMultiClassicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccSlbServerGroupClassic")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceServerGroupClassicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.SlbClassicNoSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${apsarastack_slb.default.id}",
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${apsarastack_instance.default.0.id}", "${apsarastack_instance.default.1.id}"},
							"port":       "100",
							"weight":     "10",
						},
						{
							"server_ids": "${apsarastack_instance.default.*.id}",
							"port":       "80",
							"weight":     "100",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "tf-testAccSlbServerGroupClassicUpdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccSlbServerGroupClassicUpdate",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${apsarastack_instance.default.0.id}", "${apsarastack_instance.default.1.id}"},
							"port":       "100",
							"weight":     "10",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#": "1",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}",
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${apsarastack_instance.default.0.id}", "${apsarastack_instance.default.1.id}"},
							"port":       "100",
							"weight":     "10",
						},
						{
							"server_ids": "${apsarastack_instance.default.*.id}",
							"port":       "80",
							"weight":     "100",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":      "tf-testAccSlbServerGroupClassic",
						"servers.#": "2",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func resourceSlbServerGroupDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "apsarastack_instance_types" "default" {
  cpu_core_count    = 1
  memory_size       = 2
}
data "apsarastack_instance_types" "new" {
	eni_amount = 2
}
data "apsarastack_images" "default" {
  name_regex = "^ubuntu_18.*64"
  most_recent = true
  owners = "system"
}
resource "apsarastack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
  vpc_id = "${apsarastack_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"
  name = "${var.name}"
}
resource "apsarastack_security_group" "default" {
  name = "${var.name}"
  vpc_id = "${apsarastack_vpc.default.id}"
}
resource "apsarastack_network_interface" "default" {
    count = 1
    name = "${var.name}"
    vswitch_id = "${apsarastack_vswitch.default.id}"
    security_groups = [ "${apsarastack_security_group.default.id}" ]
}
resource "apsarastack_instance" "default" {
  image_id = "${data.apsarastack_images.default.images.0.id}"
  instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
  instance_name = "${var.name}"
  count = "2"
  security_groups = "${apsarastack_security_group.default.*.id}"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${apsarastack_vswitch.default.id}"
}
resource "apsarastack_instance" "new" {
  image_id = "${data.apsarastack_images.default.images.0.id}"
  instance_type = "${data.apsarastack_instance_types.new.instance_types.0.id}"
  instance_name = "${var.name}"
  count = "1"
  security_groups = "${apsarastack_security_group.default.*.id}"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.apsarastack_instance_types.new.instance_types.0.availability_zones.0}"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${apsarastack_vswitch.default.id}"
}
resource "apsarastack_network_interface_attachment" "default" {
	count = 1
    instance_id = "${apsarastack_instance.new.0.id}"
    network_interface_id = "${element(apsarastack_network_interface.default.*.id, count.index)}"
}
resource "apsarastack_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${apsarastack_vswitch.default.id}"
}
`, name)
}

func resourceServerGroupClassicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "apsarastack_instance_types" "default" {
  cpu_core_count    = 1
  memory_size       = 2
}
data "apsarastack_images" "default" {
  name_regex = "^ubuntu_18.*64"
  most_recent = true
  owners = "system"
}
resource "apsarastack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
  vpc_id = "${apsarastack_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"
  name = "${var.name}"
}
resource "apsarastack_security_group" "default" {
  name = "${var.name}"
  vpc_id = "${apsarastack_vpc.default.id}"
}
resource "apsarastack_instance" "default" {
  image_id = "${data.apsarastack_images.default.images.0.id}"
  instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
  instance_name = "${var.name}"
  count = "2"
  security_groups = "${apsarastack_security_group.default.*.id}"

  internet_max_bandwidth_out = "10"
  availability_zone = "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${apsarastack_vswitch.default.id}"
}
resource "apsarastack_slb" "default" {
  name = "${var.name}"
}
`, name)

}

func resourceSlbServerGroupMultiVpcDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSlbServerGroupVpc"
}
data "apsarastack_instance_types" "default" {
  cpu_core_count    = 1
  memory_size       = 2
}
data "apsarastack_images" "default" {
        name_regex = "^ubuntu_18.*64"
  most_recent = true
  owners = "system"
}
resource "apsarastack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
  vpc_id = "${apsarastack_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"
  name = "${var.name}"
}
resource "apsarastack_security_group" "default" {
  name = "${var.name}"
  vpc_id = "${apsarastack_vpc.default.id}"
}
resource "apsarastack_instance" "default" {
  image_id = "${data.apsarastack_images.default.images.0.id}"
  instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
  instance_name = "${var.name}"
  count = "2"
  security_groups = "${apsarastack_security_group.default.*.id}"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${apsarastack_vswitch.default.id}"
}
resource "apsarastack_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${apsarastack_vswitch.default.id}"
}
`)

}

var serverGroupMap = map[string]string{
	"name":      "tf-server-group",
	"servers.#": "1",
}

var serverGroupMultiClassicMap = map[string]string{
	"servers.#": "2",
}

var serversMap = []map[string]interface{}{
	{
		"server_ids": []string{"${apsarastack_instance.default.0.id}"},
		"port":       "1",
		"weight":     "10",
	},
	{
		"server_ids": []string{"${apsarastack_instance.default.1.id}"},
		"port":       "2",
		"weight":     "10",
	},
	{
		"server_ids": []string{"${apsarastack_instance.default.2.id}"},
		"port":       "3",
		"weight":     "10",
	},
	//{
	//	"server_ids": []string{"${apsarastack_instance.default.3.id}"},
	//	"port":       "4",
	//	"weight":     "10",
	//},
	//{
	//	"server_ids": []string{"${apsarastack_instance.default.4.id}"},
	//	"port":       "5",
	//	"weight":     "10",
	//},
	//{
	//	"server_ids": []string{"${apsarastack_instance.default.5.id}"},
	//	"port":       "6",
	//	"weight":     "10",
	//},
	//{
	//	"server_ids": []string{"${apsarastack_instance.default.6.id}"},
	//	"port":       "7",
	//	"weight":     "10",
	//},
	//{
	//	"server_ids": []string{"${apsarastack_instance.default.7.id}"},
	//	"port":       "8",
	//	"weight":     "10",
	//},
	//{
	//	"server_ids": []string{"${apsarastack_instance.default.8.id}"},
	//	"port":       "9",
	//	"weight":     "10",
	//},
	//{
	//	"server_ids": []string{"${apsarastack_instance.default.9.id}"},
	//	"port":       "10",
	//	"weight":     "10",
	//},
	//{
	//	"server_ids": []string{"${apsarastack_instance.default.10.id}"},
	//	"port":       "11",
	//	"weight":     "10",
	//},
	//{
	//	"server_ids": []string{"${apsarastack_instance.default.11.id}"},
	//	"port":       "12",
	//	"weight":     "10",
	//},
	//{
	//	"server_ids": []string{"${apsarastack_instance.default.12.id}"},
	//	"port":       "13",
	//	"weight":     "10",
	//},
	//{
	//	"server_ids": []string{"${apsarastack_instance.default.13.id}"},
	//	"port":       "14",
	//	"weight":     "10",
	//},
}
