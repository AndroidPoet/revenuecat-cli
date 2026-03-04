package metrics

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/AndroidPoet/revenuecat-cli/internal/api"
	"github.com/AndroidPoet/revenuecat-cli/internal/cli"
	"github.com/AndroidPoet/revenuecat-cli/internal/output"
)

// MetricsCmd manages metrics
var MetricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "View metrics",
	Long:  `View metrics and analytics for your RevenueCat project.`,
}

var overviewCmd = &cobra.Command{
	Use:   "overview",
	Short: "Get metrics overview",
	RunE:  runOverview,
}

func init() {
	MetricsCmd.AddCommand(overviewCmd)
}

type MetricsOverview struct {
	ActiveSubscribers int     `json:"active_subscribers"`
	ActiveTrials      int     `json:"active_trials"`
	MRR               float64 `json:"mrr"`
	Revenue           float64 `json:"revenue"`
}

func parseTimeout() time.Duration {
	d, err := time.ParseDuration(cli.GetTimeout())
	if err != nil {
		return 60 * time.Second
	}
	return d
}

func runOverview(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}

	ctx, cancel := client.Context()
	defer cancel()

	var overview MetricsOverview
	path := fmt.Sprintf("/projects/%s/metrics/overview", client.GetProjectID())
	if err := client.Get(ctx, path, &overview); err != nil {
		return err
	}

	return output.Print(overview)
}
