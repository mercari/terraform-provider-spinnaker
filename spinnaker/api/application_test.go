package api

import (
	"testing"
)

func TestValidateApplicationCloudProviders(t *testing.T) {
	tcs := map[string]struct {
		cloudProviders []string
		appName        string
		shouldPass     bool
	}{
		"pass":                                          {[]string{"kubernetes"}, "test-1", true},
		"pass with multiple provider":                   {[]string{"kubernetes", "gce"}, "test1", true},
		"fail with single invalid cloudProvider":        {[]string{"mercari"}, "test-1", false},
		"fail with valid and invalid mix cloudProvider": {[]string{"kubernetes", "mercari"}, "test-1", false},
		"fail with over max lenght":                     {[]string{"kubernetes"}, "very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-long-app", false},
	}

	for n, tc := range tcs {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			for _, p := range tc.cloudProviders {
				err := validateSpinnakerApplicationNameByCloudProvider(tc.appName, p)
				if err != nil && tc.shouldPass {
					t.Fatalf("failed: %v", err)
				}
			}
		})
	}
}

func TestValidateSpinnakerApplicationAccess(t *testing.T) {
	tcs := map[string]struct {
		accesses   []string
		shouldPass bool
	}{
		"pass":                      {[]string{"WRITE"}, true},
		"pass with multiple access": {SupportedAccesses, true},
		"not supported access":      {[]string{"MERCARI", "KEKE"}, false},
		"mixture of supported access and not supported access": {[]string{"WRITE", "KEKE", "BLOG"}, false},
	}

	for n, tc := range tcs {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			if err := validateSpinnakerApplicationAccess(tc.accesses); err != nil && tc.shouldPass {
				t.Fatalf("failed: %v", err)
			}
		})
	}
}
