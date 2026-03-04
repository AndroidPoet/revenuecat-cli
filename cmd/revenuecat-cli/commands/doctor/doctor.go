package doctor

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/AndroidPoet/revenuecat-cli/internal/api"
	"github.com/AndroidPoet/revenuecat-cli/internal/cli"
	"github.com/AndroidPoet/revenuecat-cli/internal/config"
	"github.com/AndroidPoet/revenuecat-cli/internal/output"
)

// DoctorCmd runs diagnostic checks
var DoctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Run diagnostic checks",
	Long: `Run diagnostic checks to verify your revenuecat-cli setup.

Checks:
  1. Configuration file loads correctly
  2. API key is configured
  3. Project ID is set
  4. RevenueCat API is reachable`,
	RunE: runDoctor,
}

type checkResult struct {
	Check  string `json:"check"`
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

func runDoctor(cmd *cobra.Command, args []string) error {
	results := make([]checkResult, 0, 4)

	// Check 1: Configuration
	configPath := config.GetConfigPath()
	if configPath != "" {
		results = append(results, checkResult{
			Check:  "Configuration",
			Status: "OK",
			Detail: configPath,
		})
	} else {
		results = append(results, checkResult{
			Check:  "Configuration",
			Status: "WARN",
			Detail: "No config file found. Run 'rc auth login' to configure.",
		})
	}

	// Check 2: API Key
	apiKey, err := config.GetAPIKey()
	if err != nil {
		results = append(results, checkResult{
			Check:  "API Key",
			Status: "FAIL",
			Detail: err.Error(),
		})
	} else {
		masked := apiKey[:6] + "..." + apiKey[len(apiKey)-4:]
		results = append(results, checkResult{
			Check:  "API Key",
			Status: "OK",
			Detail: masked,
		})
	}

	// Check 3: Project ID
	projectID := cli.GetProjectID()
	if projectID == "" {
		results = append(results, checkResult{
			Check:  "Project ID",
			Status: "WARN",
			Detail: "Not set. Use --project flag or RC_PROJECT env var.",
		})
	} else {
		results = append(results, checkResult{
			Check:  "Project ID",
			Status: "OK",
			Detail: projectID,
		})
	}

	// Check 4: API connectivity (only if we have key + project)
	if apiKey != "" && projectID != "" {
		client, err := api.NewClient(projectID, 10*time.Second)
		if err != nil {
			results = append(results, checkResult{
				Check:  "API Connectivity",
				Status: "FAIL",
				Detail: err.Error(),
			})
		} else {
			ctx, cancel := client.Context()
			defer cancel()

			var project map[string]interface{}
			path := fmt.Sprintf("/projects/%s", projectID)
			if err := client.Get(ctx, path, &project); err != nil {
				results = append(results, checkResult{
					Check:  "API Connectivity",
					Status: "FAIL",
					Detail: err.Error(),
				})
			} else {
				results = append(results, checkResult{
					Check:  "API Connectivity",
					Status: "OK",
					Detail: "Successfully connected to RevenueCat API",
				})
			}
		}
	} else {
		results = append(results, checkResult{
			Check:  "API Connectivity",
			Status: "SKIP",
			Detail: "Requires API key and project ID",
		})
	}

	// Print summary
	allOK := true
	for _, r := range results {
		icon := "✓"
		if r.Status == "FAIL" {
			icon = "✗"
			allOK = false
		} else if r.Status == "WARN" || r.Status == "SKIP" {
			icon = "!"
		}
		output.PrintInfo("%s %s: %s (%s)", icon, r.Check, r.Status, r.Detail)
	}

	if allOK {
		output.PrintSuccess("\nAll checks passed!")
	}

	return output.Print(results)
}
