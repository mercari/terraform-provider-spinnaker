package spinnaker

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mercari/terraform-provider-spinnaker/spinnaker/api"
)

func TestAccResourceSourceSpinnakerPipeline_basic(t *testing.T) {
	resourceName := acctest.RandomWithPrefix("tf-acc-test")
	application := acctest.RandomWithPrefix("tf-acc-test")
	pipeline := `{
  		"keepWaitingPipelines": false,
  		"limitConcurrent": true,
  		"notifications": [],
  		"stages": [],
  		"triggers": []
		}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpinnakerPipelineDestroy("spinnaker_pipeline.test", application),
		Steps: []resource.TestStep{
			{
				Config: testAccSpinnakerPipeline_basic(resourceName, application, pipeline),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpinnakerPipelineExists("spinnaker_pipeline.test", application),
					resource.TestCheckResourceAttr("spinnaker_pipeline.test", "name", resourceName),
					resource.TestCheckResourceAttr("spinnaker_pipeline.test", "application", application),
				),
			},
			{
				ResourceName:      "spinnaker_pipeline.test",
				ImportStateId:     fmt.Sprintf("%s%s", application, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSpinnakerPipelineDestroy(resourceName string, applicationName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		pipelineName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(gateConfig).client
		pipeline := &pipelineRead{}

		if !ok {
			return fmt.Errorf("Pipeline not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Pipeline ID is set")
		}

		if _, err := api.GetPipeline(client, applicationName, pipelineName, pipeline); err != nil {
			if err.Error() == api.ErrCodeNoSuchEntityException {
				return nil
			}
			return err
		}

		return fmt.Errorf("Pipeline still exists: %s", pipelineName)
	}
}

func testAccCheckSpinnakerPipelineExists(resourceName string, applicationName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		pipelineName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(gateConfig).client
		pipeline := &pipelineRead{}

		if !ok {
			return fmt.Errorf("Pipeline not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Pipeline ID is set")
		}

		if _, err := api.GetPipeline(client, applicationName, pipelineName, pipeline); err != nil {
			return err
		}

		return nil
	}
}

func testAccSpinnakerPipeline_basic(rName string, application string, pipeline string) string {
	return fmt.Sprintf(`
resource "spinnaker_application" "test" {
	name = %q
	email =  "acceptance@test.com"
}

resource "spinnaker_pipeline" "test" {
	name  = %q
	application = %q
	pipeline = %q
}
`, application, rName, application, pipeline)
}
