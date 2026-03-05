package watch

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/AndroidPoet/revenuecat-cli/internal/api"
	"github.com/AndroidPoet/revenuecat-cli/internal/cli"
)

var interval string

// WatchCmd watches resources in real-time
var WatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch resources in real-time",
	Long:  `Watch RevenueCat resources with live-refreshing output.`,
}

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Watch metrics in real-time",
	RunE:  runWatchMetrics,
}

func init() {
	metricsCmd.Flags().StringVar(&interval, "interval", "5s", "refresh interval (e.g. 5s, 1m)")

	WatchCmd.AddCommand(metricsCmd)
}

func parseTimeout() time.Duration {
	d, err := time.ParseDuration(cli.GetTimeout())
	if err != nil {
		return 60 * time.Second
	}
	return d
}

func runWatchMetrics(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	dur, err := time.ParseDuration(interval)
	if err != nil {
		return fmt.Errorf("invalid interval %q: %w", interval, err)
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}

	pid := client.GetProjectID()
	path := fmt.Sprintf("/projects/%s/metrics/overview", pid)

	for {
		// Clear screen
		fmt.Print("\033[H\033[2J")

		fmt.Printf("RevenueCat Metrics — %s\n", time.Now().Format("2006-01-02 15:04:05"))
		fmt.Println("─────────────────────────────────────")

		ctx, cancel := client.Context()

		var metrics struct {
			ActiveSubscribers int     `json:"active_subscribers"`
			ActiveTrials      int     `json:"active_trials"`
			MRR               float64 `json:"mrr"`
			Revenue           float64 `json:"revenue"`
		}

		if err := client.Get(ctx, path, &metrics); err != nil {
			cancel()
			fmt.Printf("Error: %s\n", err)
		} else {
			cancel()
			fmt.Printf("  Active Subscribers:  %d\n", metrics.ActiveSubscribers)
			fmt.Printf("  Active Trials:       %d\n", metrics.ActiveTrials)
			fmt.Printf("  MRR:                 $%.2f\n", metrics.MRR)
			fmt.Printf("  Revenue:             $%.2f\n", metrics.Revenue)
		}

		fmt.Println()
		fmt.Printf("Press Ctrl+C to stop. Refreshing every %s...\n", dur)
		time.Sleep(dur)
	}
}
