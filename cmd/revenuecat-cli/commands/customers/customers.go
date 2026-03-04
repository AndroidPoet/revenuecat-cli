package customers

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/AndroidPoet/revenuecat-cli/internal/api"
	"github.com/AndroidPoet/revenuecat-cli/internal/cli"
	"github.com/AndroidPoet/revenuecat-cli/internal/output"
)

var (
	customerID string
)

// CustomersCmd manages customers
var CustomersCmd = &cobra.Command{
	Use:   "customers",
	Short: "Manage customers",
	Long:  `Get and manage customer information in your RevenueCat project.`,
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get customer details",
	RunE:  runGet,
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a customer",
	RunE:  runDelete,
}

func init() {
	getCmd.Flags().StringVar(&customerID, "customer-id", "", "customer ID")
	getCmd.MarkFlagRequired("customer-id")

	var confirm bool
	deleteCmd.Flags().StringVar(&customerID, "customer-id", "", "customer ID")
	deleteCmd.Flags().BoolVar(&confirm, "confirm", false, "confirm deletion")
	deleteCmd.MarkFlagRequired("customer-id")

	CustomersCmd.AddCommand(getCmd)
	CustomersCmd.AddCommand(deleteCmd)
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

	var customer map[string]interface{}
	path := fmt.Sprintf("/projects/%s/customers/%s", client.GetProjectID(), customerID)
	if err := client.Get(ctx, path, &customer); err != nil {
		return err
	}

	return output.Print(customer)
}

func runDelete(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if err := cli.CheckConfirm(cmd); err != nil {
		return err
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would delete customer '%s'", customerID)
		return nil
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}

	ctx, cancel := client.Context()
	defer cancel()

	path := fmt.Sprintf("/projects/%s/customers/%s", client.GetProjectID(), customerID)
	if err := client.Delete(ctx, path); err != nil {
		return err
	}

	output.PrintSuccess("Customer '%s' deleted", customerID)
	return nil
}
