package status

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/AndroidPoet/revenuecat-cli/internal/api"
	"github.com/AndroidPoet/revenuecat-cli/internal/cli"
	"github.com/AndroidPoet/revenuecat-cli/internal/output"
)

// StatusCmd displays a project status dashboard
var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Project status dashboard",
	Long:  `Fetch apps, metrics, entitlements, and offerings in parallel and display a combined project summary.`,
	RunE:  runStatus,
}

type StatusSummary struct {
	ProjectID         string  `json:"project_id"`
	Apps              int     `json:"apps"`
	Entitlements      int     `json:"entitlements"`
	Offerings         int     `json:"offerings"`
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

func runStatus(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}

	ctx, cancel := client.Context()
	defer cancel()

	pid := client.GetProjectID()

	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		errs    []error
		summary = StatusSummary{ProjectID: pid}
	)

	type listCount struct {
		Items json.RawMessage `json:"items"`
	}

	// Fetch apps count
	wg.Add(1)
	go func() {
		defer wg.Done()
		var resp listCount
		path := fmt.Sprintf("/projects/%s/apps?limit=5", pid)
		if err := client.Get(ctx, path, &resp); err != nil {
			mu.Lock()
			errs = append(errs, fmt.Errorf("apps: %w", err))
			mu.Unlock()
			return
		}
		var items []json.RawMessage
		if err := json.Unmarshal(resp.Items, &items); err != nil {
			mu.Lock()
			errs = append(errs, fmt.Errorf("apps parse: %w", err))
			mu.Unlock()
			return
		}
		mu.Lock()
		summary.Apps = len(items)
		mu.Unlock()
	}()

	// Fetch metrics overview
	wg.Add(1)
	go func() {
		defer wg.Done()
		var metrics struct {
			ActiveSubscribers int     `json:"active_subscribers"`
			ActiveTrials      int     `json:"active_trials"`
			MRR               float64 `json:"mrr"`
			Revenue           float64 `json:"revenue"`
		}
		path := fmt.Sprintf("/projects/%s/metrics/overview", pid)
		if err := client.Get(ctx, path, &metrics); err != nil {
			// Metrics endpoint may 404; treat as zero
			return
		}
		mu.Lock()
		summary.ActiveSubscribers = metrics.ActiveSubscribers
		summary.ActiveTrials = metrics.ActiveTrials
		summary.MRR = metrics.MRR
		summary.Revenue = metrics.Revenue
		mu.Unlock()
	}()

	// Fetch entitlements count
	wg.Add(1)
	go func() {
		defer wg.Done()
		var resp listCount
		path := fmt.Sprintf("/projects/%s/entitlements?limit=5", pid)
		if err := client.Get(ctx, path, &resp); err != nil {
			mu.Lock()
			errs = append(errs, fmt.Errorf("entitlements: %w", err))
			mu.Unlock()
			return
		}
		var items []json.RawMessage
		if err := json.Unmarshal(resp.Items, &items); err != nil {
			mu.Lock()
			errs = append(errs, fmt.Errorf("entitlements parse: %w", err))
			mu.Unlock()
			return
		}
		mu.Lock()
		summary.Entitlements = len(items)
		mu.Unlock()
	}()

	// Fetch offerings count
	wg.Add(1)
	go func() {
		defer wg.Done()
		var resp listCount
		path := fmt.Sprintf("/projects/%s/offerings?limit=5", pid)
		if err := client.Get(ctx, path, &resp); err != nil {
			mu.Lock()
			errs = append(errs, fmt.Errorf("offerings: %w", err))
			mu.Unlock()
			return
		}
		var items []json.RawMessage
		if err := json.Unmarshal(resp.Items, &items); err != nil {
			mu.Lock()
			errs = append(errs, fmt.Errorf("offerings parse: %w", err))
			mu.Unlock()
			return
		}
		mu.Lock()
		summary.Offerings = len(items)
		mu.Unlock()
	}()

	wg.Wait()

	if len(errs) > 0 {
		for _, e := range errs {
			output.PrintWarning("%s", e)
		}
	}

	output.PrintInfo("Project: %s", summary.ProjectID)
	output.PrintInfo("  Apps:              %d", summary.Apps)
	output.PrintInfo("  Entitlements:      %d", summary.Entitlements)
	output.PrintInfo("  Offerings:         %d", summary.Offerings)
	output.PrintInfo("  Active Subscribers: %d", summary.ActiveSubscribers)
	output.PrintInfo("  Active Trials:     %d", summary.ActiveTrials)
	output.PrintInfo("  MRR:               $%.2f", summary.MRR)
	output.PrintInfo("  Revenue:           $%.2f", summary.Revenue)

	return output.Print(summary)
}
