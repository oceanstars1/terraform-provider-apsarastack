package apsarastack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("apsarastack_key_pair", &resource.Sweeper{
		Name: "apsarastack_key_pair",
		F:    testSweepKeyPairs,
	})
}

func testSweepKeyPairs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting ApsaraStack client: %s", err)
	}
	client := rawClient.(*connectivity.ApsaraStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"testAcc",
		"terraform-test-",
	}

	var pairs []ecs.KeyPair
	req := ecs.CreateDescribeKeyPairsRequest()
	req.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		req.Scheme = "https"
	} else {
		req.Scheme = "http"
	}
	req.Headers = map[string]string{"RegionId": client.RegionId}
	req.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeKeyPairs(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving Key Pairs: %s", err)
		}
		resp, _ := raw.(*ecs.DescribeKeyPairsResponse)
		if resp == nil || len(resp.KeyPairs.KeyPair) < 1 {
			break
		}
		pairs = append(pairs, resp.KeyPairs.KeyPair...)

		if len(resp.KeyPairs.KeyPair) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	for _, v := range pairs {
		name := v.KeyPairName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Key Pair: %s", name)
			continue
		}
		log.Printf("[INFO] Deleting Key Pair: %s", name)
		req := ecs.CreateDeleteKeyPairsRequest()
		if strings.ToLower(client.Config.Protocol) == "https" {
			req.Scheme = "https"
		} else {
			req.Scheme = "http"
		}
		req.Headers = map[string]string{"RegionId": client.RegionId}
		req.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		req.KeyPairNames = convertListToJsonString(append(make([]interface{}, 0, 1), name))
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteKeyPairs(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Key Pair (%s): %s", name, err)
		}
	}
	return nil
}

// this method is referenced by other file, so hold it temporarily.
func testAccCheckKeyPairExists(n string, keypair *ecs.KeyPair) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Key Pair ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
		ecsService := EcsService{client}

		response, err := ecsService.DescribeKeyPair(rs.Primary.ID)

		log.Printf("[WARN] disk ids %#v", rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Finding Key Pair %#v got an error: %#v.", rs.Primary.ID, err)
		}
		*keypair = response
		return nil
	}
}

func testAccCheckKeyPairDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "apsarastack_key_pair" {
			continue
		}

		// Try to find the Disk
		client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
		ecsService := EcsService{client}

		response, err := ecsService.DescribeKeyPair(rs.Primary.ID)
		os.Remove(rs.Primary.Attributes["key_file"])

		if err != nil {
			// Verify the error is what we want
			if NotFoundError(err) {
				continue
			}
			return err
		}
		if response.KeyPairName != "" {
			return fmt.Errorf("Error Key Pair still exist")
		}
	}

	return nil
}

func TestAccApsaraStackKeyPairBasic(t *testing.T) {
	var v ecs.KeyPair
	resourceId := "apsarastack_key_pair.default"
	ra := resourceAttrInit(resourceId, testAccCheckKeyPairBasicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKeyPairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyPairConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccKeyPairConfig_public_key(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"public_key": "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg",
					}),
				),
			},
			{
				Config: testAccKeyPairConfig_tag(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccKeyPairConfig_key_name(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_name": fmt.Sprintf("tf-testAccKeyPairConfig%d", rand),
					}),
				),
			},
		},
	})

}

func TestAccApsaraStackKeyPairMulti(t *testing.T) {
	var v ecs.KeyPair
	resourceId := "apsarastack_key_pair.default.9"
	ra := resourceAttrInit(resourceId, testAccCheckKeyPairBasicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKeyPairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyPairConfigMulti(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

var testAccCheckKeyPairBasicMap = map[string]string{
	"finger_print": CHECKSET,
	"key_name":     CHECKSET,
}

func testAccKeyPairConfigBasic(rand int) string {
	return fmt.Sprintf(`
resource "apsarastack_key_pair" "default" {
	key_name ="tf-testAccKeyPairConfig%d"
}
`, rand)
}

func testAccKeyPairConfig_public_key(rand int) string {
	return fmt.Sprintf(`
resource "apsarastack_key_pair" "default" {
	key_name ="tf-testAccKeyPairConfig%d"
	public_key = "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg"
    
}
`, rand)
}
func testAccKeyPairConfig_tag(rand int) string {
	return fmt.Sprintf(`
resource "apsarastack_key_pair" "default" {
	key_name ="tf-testAccKeyPairConfig%d"
	public_key = "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg"
    
}
`, rand)
}

func testAccKeyPairConfig_key_name(rand int) string {
	return fmt.Sprintf(`
resource "apsarastack_key_pair" "default" {
	key_name  = "tf-testAccKeyPairConfig%d"
	public_key = "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg"
    
}
`, rand)
}

func testAccKeyPairConfigMulti(rand int) string {
	return fmt.Sprintf(`
resource "apsarastack_key_pair" "default" {
	count = 10
	key_name =  "tf-testAccKeyPairConfig%d${count.index}"
}
`, rand)
}
