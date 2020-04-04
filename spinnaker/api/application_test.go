package api

import (
	"strings"
	"testing"
)

func TestValidateApplicationCloudProviders(t *testing.T) {
	tcs := map[string]struct {
		cloudProviders string
		appName        string
		shouldPass     bool
	}{
		"pass":                                          {"kubernetes", "test-1", true},
		"pass with multiple provider":                   {"kubernetes,gce", "test1", true},
		"fail with single invalid cloudProvider":        {"mercari", "test-1", false},
		"fail with valid and invalid mix cloudProvider": {"kubernetes,mercari", "test-1", false},
		"fail with over max lenght":                     {"kubernetes", "very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-long-app", false},
	}

	for n, tc := range tcs {
		t.Run(n, func(t *testing.T) {
			for _, p := range strings.Split(tc.cloudProviders, ",") {
				err := validateSpinnakerApplicationNameByCloudProvider(tc.appName, p)
				if err != nil && tc.shouldPass {
					t.Fatalf("failed: %v", err)
				}
			}
		})
	}
}
