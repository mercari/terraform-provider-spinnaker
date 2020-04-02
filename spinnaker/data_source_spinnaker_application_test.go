package spinnaker

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
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
					resource.TestCheckResourceAttr(resourceName, "instance_port", "80"),
				),
			},
		},
	})
}
