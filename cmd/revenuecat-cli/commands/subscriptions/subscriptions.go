package subscriptions

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/AndroidPoet/revenuecat-cli/internal/api"
	"github.com/AndroidPoet/revenuecat-cli/internal/cli"
	"github.com/AndroidPoet/revenuecat-cli/internal/output"
)

var (
	subscriptionID string
)

// SubscriptionsCmd manages subscriptions
var SubscriptionsCmd = &cobra.Command{
	Use:   "subscriptions",
	Short: "Manage subscriptions",
	Long:  `Get, list entitlements, cancel, and refund subscriptions in your RevenueCat project.`,
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get subscription details",
	RunE:  runGet,
}

var listEntitlementsCmd = &cobra.Command{
	Use:   "list-entitlements",
	Short: "List entitlements for a subscription",
	RunE:  runListEntitlements,
}

var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel a subscription",
	RunE:  runCancel,
}

var refundCmd = &cobra.Command{
	Use:   "refund",
	Short: "Refund a subscription",
	RunE:  runRefund,
}

func init() {
	getCmd.Flags().StringVar(&subscriptionID, "subscription-id", "", "subscription ID")
	getCmd.MarkFlagRequired("subscription-id")

	listEntitlementsCmd.Flags().StringVar(&subscriptionID, "subscription-id", "", "subscription ID")
	listEntitlementsCmd.MarkFlagRequired("subscription-id")

	var confirm bool
	cancelCmd.Flags().StringVar(&subscriptionID, "subscription-id", "", "subscription ID")
	cancelCmd.Flags().BoolVar(&confirm, "confirm", false, "confirm cancellation")
	cancelCmd.MarkFlagRequired("subscription-id")

	refundCmd.Flags().StringVar(&subscriptionID, "subscription-id", "", "subscription ID")
	refundCmd.Flags().BoolVar(&confirm, "confirm", false, "confirm refund")
	refundCmd.MarkFlagRequired("subscription-id")

	SubscriptionsCmd.AddCommand(getCmd)
	SubscriptionsCmd.AddCommand(listEntitlementsCmd)
	SubscriptionsCmd.AddCommand(cancelCmd)
	SubscriptionsCmd.AddCommand(refundCmd)
}

func parseTimeout() time.Duration {
	d, err := time.ParseDuration(cli.GetTimeout())
	if err != nil {
		return 60 * time.Second
	}
	return d
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

	var subscription map[string]interface{}
	path := fmt.Sprintf("/projects/%s/subscriptions/%s", client.GetProjectID(), subscriptionID)
	if err := client.Get(ctx, path, &subscription); err != nil {
		return err
	}

	return output.Print(subscription)
}

func runListEntitlements(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}

	ctx, cancel := client.Context()
	defer cancel()

	var resp struct {
		Items []map[string]interface{} `json:"items"`
	}

	path := fmt.Sprintf("/projects/%s/subscriptions/%s/entitlements", client.GetProjectID(), subscriptionID)
	if err := client.Get(ctx, path, &resp); err != nil {
		return err
	}

	return output.Print(resp.Items)
}

func runCancel(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if err := cli.CheckConfirm(cmd); err != nil {
		return err
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would cancel subscription '%s'", subscriptionID)
		return nil
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}

	ctx, cancel := client.Context()
	defer cancel()

	path := fmt.Sprintf("/projects/%s/subscriptions/%s/actions/cancel", client.GetProjectID(), subscriptionID)
	if err := client.Post(ctx, path, nil, nil); err != nil {
		return err
	}

	output.PrintSuccess("Subscription '%s' cancelled", subscriptionID)
	return nil
}

func runRefund(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if err := cli.CheckConfirm(cmd); err != nil {
		return err
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would refund subscription '%s'", subscriptionID)
		return nil
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}

	ctx, cancel := client.Context()
	defer cancel()

	path := fmt.Sprintf("/projects/%s/subscriptions/%s/actions/refund", client.GetProjectID(), subscriptionID)
	if err := client.Post(ctx, path, nil, nil); err != nil {
		return err
	}

	output.PrintSuccess("Subscription '%s' refunded", subscriptionID)
	return nil
}
