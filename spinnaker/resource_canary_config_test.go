package spinnaker

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mercari/terraform-provider-spinnaker/spinnaker/api"
)

const (
	testDesc = "Test description"
)

func TestAccResourceSourceSpinnakerCanaryConfig_basic(t *testing.T) {
	resourceName := "spinnaker_canary_config.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpinnakerCanaryConfigConfigDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccSpinnakerCanaryConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpinnakerCanaryConfigExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", testDesc),
					resource.TestCheckResourceAttr(resourceName, "applications.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "metric.name", "CPU"),
					resource.TestCheckResourceAttr(resourceName, "metric.groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "metric.query.type", "stackdriver"),
					resource.TestCheckResourceAttr(resourceName, "metric.query.service_type", "stackdriver"),
					resource.TestCheckResourceAttr(resourceName, "metric.query.resource_type", "k8s_node"),
					resource.TestCheckResourceAttr(resourceName, "metric.query.metric_type", "kubernetes.io/anthos/gkeconnect_dialer_connection_attempts_total"),
				),
			},
		},
	})
}

func testAccCheckSpinnakerCanaryConfigConfigDestroy(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Canary Config not found, canary config: %s", n)
		}

		id := rs.Primary.ID
		if id == "" {
			return fmt.Errorf("No canary config ID is set")
		}

		client := testAccProvider.Meta().(gateConfig).client
		cfg := &canaryConfigRead{}

		retry := 5
		for {
			if err := api.GetCanaryConfig(client, id, cfg); err != nil {
				if strings.Contains(err.Error(), "not found") {
					return nil
				}

				return err
			}

			if cfg == nil {
				return nil
			}

			if cfg != nil {
				retry--
				log.Printf("[INFO] Retring CheckDestroy in 1 seconds, retry count: %v", 5-retry)
				time.Sleep(1 * time.Second)
			}

			if retry <= 0 {
				break
			}
		}

		return fmt.Errorf("Spinnaker canary config still exists, canary config: %s", id)
	}
}

func testAccCheckSpinnakerCanaryConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Canary config not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Canary Config ID is set")
		}
		client := testAccProvider.Meta().(gateConfig).client
		err := resource.Retry(1*time.Minute, func() *resource.RetryError {
			_, resp, err := client.V2CanaryConfigControllerApi.GetCanaryConfigUsingGET(client.Context, rs.Primary.ID, nil)
			if resp != nil {
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					return resource.RetryableError(fmt.Errorf("canary config does not exit"))
				} else if resp.StatusCode != http.StatusOK {
					return resource.NonRetryableError(fmt.Errorf("encountered an error getting canary config, status code: %d", resp.StatusCode))
				}
			}
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Unable to find Canary config after retries: %s", err)
		}
		return nil
	}
}

func testAccSpinnakerCanaryConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "spinnaker_canary_config" "test" {
  name         = "%s"
  description  = "%s"
  applications = ["keke-test"]

  metric {
    name   = "CPU"
    groups = ["Group 1"]

    query {
      type          = "stackdriver"
      service_type  = "stackdriver"
      resource_type = "k8s_node"
      metric_type   = "kubernetes.io/anthos/gkeconnect_dialer_connection_attempts_total"
    }
  }

  classifier {
    group_weights = {
      "Group 1" = 100
    }
  }
}
`, rName, testDesc)
}
