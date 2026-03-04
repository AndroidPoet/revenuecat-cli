package paywalls

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/AndroidPoet/revenuecat-cli/internal/api"
	"github.com/AndroidPoet/revenuecat-cli/internal/cli"
	"github.com/AndroidPoet/revenuecat-cli/internal/output"
)

var (
	offeringID string
)

// PaywallsCmd manages paywalls
var PaywallsCmd = &cobra.Command{
	Use:   "paywalls",
	Short: "Manage paywalls",
	Long:  `Create paywalls for offerings in your RevenueCat project.`,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a paywall for an offering",
	RunE:  runCreate,
}

func init() {
	createCmd.Flags().StringVar(&offeringID, "offering-id", "", "offering ID to create paywall for")
	createCmd.MarkFlagRequired("offering-id")

	PaywallsCmd.AddCommand(createCmd)
}

func parseTimeout() time.Duration {
	d, err := time.ParseDuration(cli.GetTimeout())
	if err != nil {
		return 60 * time.Second
	}
	return d
}

func runCreate(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would create paywall for offering '%s'", offeringID)
		return nil
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}

	ctx, cancel := client.Context()
	defer cancel()

	var result map[string]interface{}
	path := fmt.Sprintf("/projects/%s/offerings/%s/paywall", client.GetProjectID(), offeringID)
	if err := client.Post(ctx, path, map[string]interface{}{}, &result); err != nil {
		return err
	}

	output.PrintSuccess("Paywall created for offering '%s'", offeringID)
	return output.Print(result)
}
