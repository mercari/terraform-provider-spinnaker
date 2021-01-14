package spinnaker

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSpinnakerApplication_basic(t *testing.T) {
	resourceName := "spinnaker_application.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpinnakerApplicatioDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccSpinnakerApplication_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpinnakerApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "email", "acceptance@test.com"),
					resource.TestCheckResourceAttr(resourceName, "instance_port", strconv.Itoa(defaultInstancePort)),
				),
			},
		},
	})
}

func TestAccDataSourceSpinnakerApplication_instancePort(t *testing.T) {
	resourceName := "spinnaker_application.test"
	rand.Seed(time.Now().UnixNano())
	rName := acctest.RandomWithPrefix("tf-acc-test")
	rPort := rand.Intn(8000) + 1 // avoid 0
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpinnakerApplicatioDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccSpinnakerApplication_instancePort(rName, rPort),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpinnakerApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "email", "acceptance@test.com"),
					resource.TestCheckResourceAttr(resourceName, "instance_port", strconv.Itoa(rPort)),
				),
			},
		},
	})
}

func TestAccDataSourceSpinnakerApplication_cloudProviders(t *testing.T) {
	resourceName := "spinnaker_application.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	cloudProvider := "kubernetes"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpinnakerApplicatioDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccSpinnakerApplication_cloudProvider(rName, cloudProvider),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpinnakerApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "email", "acceptance@test.com"),
					resource.TestCheckResourceAttr(resourceName, "instance_port", strconv.Itoa(defaultInstancePort)),
					resource.TestCheckResourceAttr(resourceName, "cloud_providers", cloudProvider),
				),
			},
		},
	})
}
