package paywalls

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/AndroidPoet/revenuecat-cli/internal/api"
	"github.com/AndroidPoet/revenuecat-cli/internal/cli"
	"github.com/AndroidPoet/revenuecat-cli/internal/output"
)

var (
	offeringID string
	paywallID  string
	limit      int
	startAfter string
	allPages   bool
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

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all paywalls",
	RunE:  runList,
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get paywall details",
	RunE:  runGet,
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a paywall",
	RunE:  runDelete,
}

func init() {
	createCmd.Flags().StringVar(&offeringID, "offering-id", "", "offering ID to create paywall for")
	createCmd.MarkFlagRequired("offering-id")

	listCmd.Flags().IntVar(&limit, "limit", 20, "number of results per page")
	listCmd.Flags().StringVar(&startAfter, "starting-after", "", "cursor for pagination")
	listCmd.Flags().BoolVar(&allPages, "all", false, "fetch all pages")

	getCmd.Flags().StringVar(&paywallID, "paywall-id", "", "paywall ID")
	getCmd.MarkFlagRequired("paywall-id")

	var confirm bool
	deleteCmd.Flags().StringVar(&paywallID, "paywall-id", "", "paywall ID")
	deleteCmd.Flags().BoolVar(&confirm, "confirm", false, "confirm deletion")
	deleteCmd.MarkFlagRequired("paywall-id")

	PaywallsCmd.AddCommand(createCmd)
	PaywallsCmd.AddCommand(listCmd)
	PaywallsCmd.AddCommand(getCmd)
	PaywallsCmd.AddCommand(deleteCmd)
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

type PaywallInfo struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at,omitempty"`
}

func runList(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}

	ctx, cancel := client.Context()
	defer cancel()

	path := fmt.Sprintf("/projects/%s/paywalls", client.GetProjectID())

	if allPages {
		var items []PaywallInfo
		err := client.ListAll(ctx, path, limit, func(raw json.RawMessage) error {
			var page []PaywallInfo
			if err := json.Unmarshal(raw, &page); err != nil {
				return err
			}
			items = append(items, page...)
			return nil
		})
		if err != nil {
			return err
		}
		return output.Print(items)
	}

	query := fmt.Sprintf("?limit=%d", limit)
	if startAfter != "" {
		query += "&starting_after=" + startAfter
	}

	var resp struct {
		Items    []PaywallInfo `json:"items"`
		NextPage string       `json:"next_page,omitempty"`
	}
	if err := client.Get(ctx, path+query, &resp); err != nil {
		return err
	}

	if resp.NextPage != "" {
		output.PrintInfo("More results available. Use --starting-after=%s for next page, or --all for everything.", resp.NextPage)
	}

	return output.Print(resp.Items)
}

func runGet(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}
	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}
	ctx, cancel := client.Context()
	defer cancel()

	var paywall PaywallInfo
	path := fmt.Sprintf("/projects/%s/paywalls/%s", client.GetProjectID(), paywallID)
	if err := client.Get(ctx, path, &paywall); err != nil {
		return err
	}
	return output.Print(paywall)
}

func runDelete(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}
	if err := cli.CheckConfirm(cmd); err != nil {
		return err
	}
	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would delete paywall '%s'", paywallID)
		return nil
	}
	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}
	ctx, cancel := client.Context()
	defer cancel()

	path := fmt.Sprintf("/projects/%s/paywalls/%s", client.GetProjectID(), paywallID)
	if err := client.Delete(ctx, path); err != nil {
		return err
	}
	output.PrintSuccess("Paywall '%s' deleted", paywallID)
	return nil
}
