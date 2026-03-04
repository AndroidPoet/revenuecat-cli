package purchases

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/AndroidPoet/revenuecat-cli/internal/api"
	"github.com/AndroidPoet/revenuecat-cli/internal/cli"
	"github.com/AndroidPoet/revenuecat-cli/internal/output"
)

var (
	purchaseID string
)

// PurchasesCmd manages purchases
var PurchasesCmd = &cobra.Command{
	Use:   "purchases",
	Short: "Manage purchases",
	Long:  `Get, list entitlements, and refund purchases in your RevenueCat project.`,
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get purchase details",
	RunE:  runGet,
}

var listEntitlementsCmd = &cobra.Command{
	Use:   "list-entitlements",
	Short: "List entitlements for a purchase",
	RunE:  runListEntitlements,
}

var refundCmd = &cobra.Command{
	Use:   "refund",
	Short: "Refund a purchase",
	RunE:  runRefund,
}

func init() {
	getCmd.Flags().StringVar(&purchaseID, "purchase-id", "", "purchase ID")
	getCmd.MarkFlagRequired("purchase-id")

	listEntitlementsCmd.Flags().StringVar(&purchaseID, "purchase-id", "", "purchase ID")
	listEntitlementsCmd.MarkFlagRequired("purchase-id")

	var confirm bool
	refundCmd.Flags().StringVar(&purchaseID, "purchase-id", "", "purchase ID")
	refundCmd.Flags().BoolVar(&confirm, "confirm", false, "confirm refund")
	refundCmd.MarkFlagRequired("purchase-id")

	PurchasesCmd.AddCommand(getCmd)
	PurchasesCmd.AddCommand(listEntitlementsCmd)
	PurchasesCmd.AddCommand(refundCmd)
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

	var purchase map[string]interface{}
	path := fmt.Sprintf("/projects/%s/purchases/%s", client.GetProjectID(), purchaseID)
	if err := client.Get(ctx, path, &purchase); err != nil {
		return err
	}

	return output.Print(purchase)
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

	path := fmt.Sprintf("/projects/%s/purchases/%s/entitlements", client.GetProjectID(), purchaseID)
	if err := client.Get(ctx, path, &resp); err != nil {
		return err
	}

	return output.Print(resp.Items)
}

func runRefund(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if err := cli.CheckConfirm(cmd); err != nil {
		return err
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would refund purchase '%s'", purchaseID)
		return nil
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}

	ctx, cancel := client.Context()
	defer cancel()

	path := fmt.Sprintf("/projects/%s/purchases/%s/actions/refund", client.GetProjectID(), purchaseID)
	if err := client.Post(ctx, path, nil, nil); err != nil {
		return err
	}

	output.PrintSuccess("Purchase '%s' refunded", purchaseID)
	return nil
}
