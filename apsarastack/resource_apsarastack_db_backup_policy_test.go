package apsarastack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func testAccCheckDBBackupPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "apsarastack_db_backup_policy" {
			continue
		}
		request := rds.CreateDescribeBackupPolicyRequest()
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "rds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.DBInstanceId = rs.Primary.ID
		_, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DescribeBackupPolicy(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
				continue
			}
			return WrapError(err)
		}
	}
	return nil
}

func TestAccApsaraStackDBBackupPolicy_mysql(t *testing.T) {
	var v *rds.DescribeBackupPolicyResponse
	resourceId := "apsarastack_db_backup_policy.default"
	serverFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeBackupPolicy")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBbackuppolicy"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBBackupPolicyMysqlConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDBBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":                 "${apsarastack_db_instance.default.id}",
					"enable_backup_log":           "true",
					"local_log_retention_hours":   "18",
					"high_space_usage_protection": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_period": []string{"Wednesday", "Monday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_time": "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_time": "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_retention_period": "800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_retention_period": "800",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup_retention_period": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup_retention_period": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_log_retention_hours": "24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_log_retention_hours": "24",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_log_retention_space": "35",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_log_retention_space": "35",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"high_space_usage_protection": "Disable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"high_space_usage_protection": "Disable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"compress_type": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compress_type": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"archive_backup_retention_period": "50",
					"archive_backup_keep_count":       "3",
					"archive_backup_keep_policy":      "ByWeek",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"archive_backup_retention_period": "50",
						"archive_backup_keep_count":       "3",
						"archive_backup_keep_policy":      "ByWeek",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"archive_backup_keep_policy": "KeepAll",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"archive_backup_keep_policy": "KeepAll",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_backup_log": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_backup_log": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":                     "${apsarastack_db_instance.default.id}",
					"preferred_backup_period":         []string{"Tuesday", "Monday", "Wednesday"},
					"preferred_backup_time":           "13:00Z-14:00Z",
					"backup_retention_period":         "900",
					"enable_backup_log":               "true",
					"log_backup_retention_period":     "7",
					"local_log_retention_hours":       "48",
					"high_space_usage_protection":     "Enable",
					"archive_backup_retention_period": "150",
					"archive_backup_keep_count":       "2",
					"archive_backup_keep_policy":      "ByMonth",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#":       "3",
						"preferred_backup_time":           "13:00Z-14:00Z",
						"backup_retention_period":         "900",
						"enable_backup_log":               "true",
						"log_backup_retention_period":     "7",
						"local_log_retention_hours":       "48",
						"high_space_usage_protection":     "Enable",
						"archive_backup_retention_period": "150",
						"archive_backup_keep_count":       "2",
						"archive_backup_keep_policy":      "ByMonth",
					}),
				),
			}},
	})
}

func resourceDBBackupPolicyMysqlConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "apsarastack_zones" "default" {
  available_resource_creation = "Rds"
}

	data "apsarastack_db_instance_engines" "default" {
  		instance_charge_type = "PostPaid"
  		engine               = "MySQL"
  		engine_version       = "5.6"
	}

	data "apsarastack_db_instance_classes" "default" {
 	 	engine = "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine}"
		engine_version = "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine_version}"
	}

resource "apsarastack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
  vpc_id            = "${apsarastack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.apsarastack_db_instance_classes.default.instance_classes.0.zone_ids.0.sub_zone_ids.0}"
  name              = "${var.name}"
}
resource "apsarastack_db_instance" "default" {
  	vswitch_id       = "${apsarastack_vswitch.default.id}"
  	instance_name    = "${var.name}"
  	engine 			 = "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine}"
	engine_version   = "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine_version}"
	instance_type    = "${data.apsarastack_db_instance_classes.default.instance_classes.0.instance_class}"
  	instance_storage = "${data.apsarastack_db_instance_classes.default.instance_classes.0.storage_range.min}"
}`, name)
}

func TestAccApsaraStackDBBackupPolicy_pgdb(t *testing.T) {
	var v *rds.DescribeBackupPolicyResponse
	resourceId := "apsarastack_db_backup_policy.default"
	serverFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeBackupPolicy")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBbackuppolicy"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBBackupPolicyPostgreSQLConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDBBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":                 "${apsarastack_db_instance.default.id}",
					"enable_backup_log":           "true",
					"local_log_retention_hours":   "1",
					"high_space_usage_protection": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_period": []string{"Monday", "Wednesday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_time": "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_time": "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_retention_period": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_retention_period": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup_retention_period": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup_retention_period": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_log_retention_hours": "24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_log_retention_hours": "24",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_log_retention_space": "35",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_log_retention_space": "35",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"high_space_usage_protection": "Disable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"high_space_usage_protection": "Disable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"compress_type": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compress_type": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_backup_log": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_backup_log": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_period":     []string{"Tuesday", "Wednesday", "Monday"},
					"preferred_backup_time":       "10:00Z-11:00Z",
					"backup_retention_period":     "20",
					"enable_backup_log":           "true",
					"log_backup_retention_period": "15",
					"local_log_retention_hours":   "48",
					"high_space_usage_protection": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#":   "3",
						"preferred_backup_time":       "10:00Z-11:00Z",
						"backup_retention_period":     "20",
						"enable_backup_log":           "true",
						"log_backup_retention_period": "15",
						"local_log_retention_hours":   "48",
						"high_space_usage_protection": "Enable",
					}),
				),
			}},
	})
}

func resourceDBBackupPolicyPostgreSQLConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "apsarastack_db_instance_engines" "default" {
	engine               = "PostgreSQL"
	engine_version       = "10.0"
	instance_charge_type = "PostPaid"
}

data "apsarastack_db_instance_classes" "default" {
	engine               = "PostgreSQL"
	engine_version       = "10.0"
	instance_charge_type = "PostPaid"
}

resource "apsarastack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
  vpc_id            = "${apsarastack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.apsarastack_db_instance_classes.default.instance_classes.0.zone_ids.0.id}"
  name              = "${var.name}"
}
resource "apsarastack_db_instance" "default" {
  	vswitch_id       = "${apsarastack_vswitch.default.id}"
  	instance_name    = "${var.name}"
  	engine 			 = "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine}"
	engine_version   = "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine_version}"
	instance_type    = "${data.apsarastack_db_instance_classes.default.instance_classes.0.instance_class}"
  	instance_storage = "${data.apsarastack_db_instance_classes.default.instance_classes.0.storage_range.min}"
	zone_id          = "${data.apsarastack_db_instance_classes.default.instance_classes.0.zone_ids.0.id}"
}`, name)
}

func TestAccApsaraStackDBBackupPolicy_SQLServer(t *testing.T) {
	var v *rds.DescribeBackupPolicyResponse
	resourceId := "apsarastack_db_backup_policy.default"
	serverFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeBackupPolicy")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBbackuppolicy"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBBackupPolicySQLServerConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDBBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${apsarastack_db_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_period": []string{"Wednesday", "Monday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_time": "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_time": "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_retention_period": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_retention_period": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup_frequency": "LogInterval",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup_frequency": "LogInterval",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"compress_type": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compress_type": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_period": []string{"Wednesday", "Tuesday"},
					"preferred_backup_time":   "11:00Z-12:00Z",
					"backup_retention_period": "13",
					"log_backup_frequency":    "LogInterval",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "2",
						"preferred_backup_time":     "11:00Z-12:00Z",
						"backup_retention_period":   "13",
						"log_backup_frequency":      "LogInterval",
					}),
				),
			}},
	})
}

func resourceDBBackupPolicySQLServerConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "apsarastack_db_instance_engines" "default" {
	engine               = "SQLServer"
	engine_version       = "2012"
}

data "apsarastack_db_instance_classes" "default" {
	engine = "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine}"
	engine_version = "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine_version}"
}
resource "apsarastack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
  vpc_id            = "${apsarastack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.apsarastack_db_instance_classes.default.instance_classes.0.zone_ids.0.sub_zone_ids.0}"
  name              = "${var.name}"
}
resource "apsarastack_db_instance" "default" {
  	vswitch_id       = "${apsarastack_vswitch.default.id}"
  	instance_name    = "${var.name}"
  	engine 			 = "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine}"
	engine_version   = "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine_version}"
	instance_type    = "${data.apsarastack_db_instance_classes.default.instance_classes.0.instance_class}"
  	instance_storage = "${data.apsarastack_db_instance_classes.default.instance_classes.0.storage_range.min}"
}`, name)
}

// Unknown current resource exists
func TestAccApsaraStackDBBackupPolicy_PPAS(t *testing.T) {
	var v *rds.DescribeBackupPolicyResponse
	resourceId := "apsarastack_db_backup_policy.default"
	serverFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeBackupPolicy")
	name := "tf-testAccDBbackuppolicy"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBBackupPolicyPPASConfigDependence)
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.RdsPPASNoSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDBBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":                 "${apsarastack_db_instance.default.id}",
					"high_space_usage_protection": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_period": []string{"Wednesday", "Monday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_time": "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_time": "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_retention_period": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_retention_period": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup_retention_period": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup_retention_period": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_log_retention_hours": "24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_log_retention_hours": "24",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_log_retention_space": "35",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_log_retention_space": "35",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"high_space_usage_protection": "Disable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"high_space_usage_protection": "Disable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"compress_type": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compress_type": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_backup_log": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_backup_log": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_period":     []string{"Wednesday", "Monday", "Tuesday"},
					"preferred_backup_time":       "10:00Z-11:00Z",
					"backup_retention_period":     "20",
					"enable_backup_log":           "true",
					"log_backup_retention_period": "15",
					"local_log_retention_hours":   "48",
					"high_space_usage_protection": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#":   "3",
						"preferred_backup_time":       "10:00Z-11:00Z",
						"backup_retention_period":     "20",
						"enable_backup_log":           "true",
						"log_backup_retention_period": "15",
						"local_log_retention_hours":   "48",
						"high_space_usage_protection": "Enable",
					}),
				),
			}},
	})
}

func resourceDBBackupPolicyPPASConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "apsarastack_zones" "default" {
  available_resource_creation = "Rds"
}
data "apsarastack_db_instance_engines" "default" {
	engine               = "PPAS"
	engine_version       = "10.0"
    instance_charge_type = "PostPaid"
    multi_zone           = true
}

data "apsarastack_db_instance_classes" "default" {
	engine               = "PPAS"
	engine_version       = "10.0"
    instance_charge_type = "PostPaid"
    multi_zone           = true
}
resource "apsarastack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
  vpc_id            = "${apsarastack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.apsarastack_db_instance_classes.default.instance_classes.0.zone_ids.0.sub_zone_ids.0}"
  name              = "${var.name}"
}
resource "apsarastack_db_instance" "default" {
  	vswitch_id       = "${apsarastack_vswitch.default.id}"
  	instance_name    = "${var.name}"
  	engine 			 = "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine}"
	engine_version   = "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine_version}"
	instance_type    = "${data.apsarastack_db_instance_classes.default.instance_classes.0.instance_class}"
  	instance_storage = "${data.apsarastack_db_instance_classes.default.instance_classes.0.storage_range.min}"
	zone_id          = "${data.apsarastack_db_instance_classes.default.instance_classes.0.zone_ids.0.id}"
}`, name)
}
