package customers

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
	customerID    string
	entitlementID string
	offeringID    string
	targetID      string
	limit         int
	startAfter    string
	allPages      bool
	attributes    string
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

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all customers",
	RunE:  runList,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a customer",
	RunE:  runCreate,
}

var listActiveEntitlementsCmd = &cobra.Command{
	Use:   "list-active-entitlements",
	Short: "List active entitlements for a customer",
	RunE:  runListActiveEntitlements,
}

var listAliasesCmd = &cobra.Command{
	Use:   "list-aliases",
	Short: "List aliases for a customer",
	RunE:  runListAliases,
}

var listAttributesCmd = &cobra.Command{
	Use:   "list-attributes",
	Short: "List attributes for a customer",
	RunE:  runListAttributes,
}

var setAttributesCmd = &cobra.Command{
	Use:   "set-attributes",
	Short: "Set attributes for a customer",
	RunE:  runSetAttributes,
}

var listSubscriptionsCmd = &cobra.Command{
	Use:   "list-subscriptions",
	Short: "List subscriptions for a customer",
	RunE:  runListSubscriptions,
}

var listPurchasesCmd = &cobra.Command{
	Use:   "list-purchases",
	Short: "List purchases for a customer",
	RunE:  runListPurchases,
}

var listInvoicesCmd = &cobra.Command{
	Use:   "list-invoices",
	Short: "List invoices for a customer",
	RunE:  runListInvoices,
}

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer a customer to another customer",
	RunE:  runTransfer,
}

var grantEntitlementCmd = &cobra.Command{
	Use:   "grant-entitlement",
	Short: "Grant an entitlement to a customer",
	RunE:  runGrantEntitlement,
}

var revokeEntitlementCmd = &cobra.Command{
	Use:   "revoke-entitlement",
	Short: "Revoke a granted entitlement from a customer",
	RunE:  runRevokeEntitlement,
}

var assignOfferingCmd = &cobra.Command{
	Use:   "assign-offering",
	Short: "Assign an offering to a customer",
	RunE:  runAssignOffering,
}

func init() {
	// Get flags
	getCmd.Flags().StringVar(&customerID, "customer-id", "", "customer ID")
	getCmd.MarkFlagRequired("customer-id")

	// Delete flags
	var confirm bool
	deleteCmd.Flags().StringVar(&customerID, "customer-id", "", "customer ID")
	deleteCmd.Flags().BoolVar(&confirm, "confirm", false, "confirm deletion")
	deleteCmd.MarkFlagRequired("customer-id")

	// List flags
	listCmd.Flags().IntVar(&limit, "limit", 20, "number of results per page")
	listCmd.Flags().StringVar(&startAfter, "starting-after", "", "cursor for pagination")
	listCmd.Flags().BoolVar(&allPages, "all", false, "fetch all pages")

	// Create flags
	createCmd.Flags().StringVar(&customerID, "customer-id", "", "customer app_user_id")
	createCmd.MarkFlagRequired("customer-id")

	// List active entitlements flags
	listActiveEntitlementsCmd.Flags().StringVar(&customerID, "customer-id", "", "customer ID")
	listActiveEntitlementsCmd.Flags().IntVar(&limit, "limit", 20, "number of results per page")
	listActiveEntitlementsCmd.Flags().StringVar(&startAfter, "starting-after", "", "cursor for pagination")
	listActiveEntitlementsCmd.Flags().BoolVar(&allPages, "all", false, "fetch all pages")
	listActiveEntitlementsCmd.MarkFlagRequired("customer-id")

	// List aliases flags
	listAliasesCmd.Flags().StringVar(&customerID, "customer-id", "", "customer ID")
	listAliasesCmd.MarkFlagRequired("customer-id")

	// List attributes flags
	listAttributesCmd.Flags().StringVar(&customerID, "customer-id", "", "customer ID")
	listAttributesCmd.MarkFlagRequired("customer-id")

	// Set attributes flags
	setAttributesCmd.Flags().StringVar(&customerID, "customer-id", "", "customer ID")
	setAttributesCmd.Flags().StringVar(&attributes, "attributes", "", "JSON attributes to set")
	setAttributesCmd.MarkFlagRequired("customer-id")
	setAttributesCmd.MarkFlagRequired("attributes")

	// List subscriptions flags
	listSubscriptionsCmd.Flags().StringVar(&customerID, "customer-id", "", "customer ID")
	listSubscriptionsCmd.Flags().IntVar(&limit, "limit", 20, "number of results per page")
	listSubscriptionsCmd.Flags().StringVar(&startAfter, "starting-after", "", "cursor for pagination")
	listSubscriptionsCmd.Flags().BoolVar(&allPages, "all", false, "fetch all pages")
	listSubscriptionsCmd.MarkFlagRequired("customer-id")

	// List purchases flags
	listPurchasesCmd.Flags().StringVar(&customerID, "customer-id", "", "customer ID")
	listPurchasesCmd.Flags().IntVar(&limit, "limit", 20, "number of results per page")
	listPurchasesCmd.Flags().StringVar(&startAfter, "starting-after", "", "cursor for pagination")
	listPurchasesCmd.Flags().BoolVar(&allPages, "all", false, "fetch all pages")
	listPurchasesCmd.MarkFlagRequired("customer-id")

	// List invoices flags
	listInvoicesCmd.Flags().StringVar(&customerID, "customer-id", "", "customer ID")
	listInvoicesCmd.Flags().IntVar(&limit, "limit", 20, "number of results per page")
	listInvoicesCmd.Flags().StringVar(&startAfter, "starting-after", "", "cursor for pagination")
	listInvoicesCmd.Flags().BoolVar(&allPages, "all", false, "fetch all pages")
	listInvoicesCmd.MarkFlagRequired("customer-id")

	// Transfer flags
	transferCmd.Flags().StringVar(&customerID, "customer-id", "", "customer ID")
	transferCmd.Flags().StringVar(&targetID, "target-id", "", "target customer ID")
	transferCmd.MarkFlagRequired("customer-id")
	transferCmd.MarkFlagRequired("target-id")

	// Grant entitlement flags
	grantEntitlementCmd.Flags().StringVar(&customerID, "customer-id", "", "customer ID")
	grantEntitlementCmd.Flags().StringVar(&entitlementID, "entitlement-id", "", "entitlement ID")
	grantEntitlementCmd.MarkFlagRequired("customer-id")
	grantEntitlementCmd.MarkFlagRequired("entitlement-id")

	// Revoke entitlement flags
	revokeEntitlementCmd.Flags().StringVar(&customerID, "customer-id", "", "customer ID")
	revokeEntitlementCmd.Flags().StringVar(&entitlementID, "entitlement-id", "", "entitlement ID")
	revokeEntitlementCmd.MarkFlagRequired("customer-id")
	revokeEntitlementCmd.MarkFlagRequired("entitlement-id")

	// Assign offering flags
	assignOfferingCmd.Flags().StringVar(&customerID, "customer-id", "", "customer ID")
	assignOfferingCmd.Flags().StringVar(&offeringID, "offering-id", "", "offering ID")
	assignOfferingCmd.MarkFlagRequired("customer-id")
	assignOfferingCmd.MarkFlagRequired("offering-id")

	CustomersCmd.AddCommand(getCmd)
	CustomersCmd.AddCommand(deleteCmd)
	CustomersCmd.AddCommand(listCmd)
	CustomersCmd.AddCommand(createCmd)
	CustomersCmd.AddCommand(listActiveEntitlementsCmd)
	CustomersCmd.AddCommand(listAliasesCmd)
	CustomersCmd.AddCommand(listAttributesCmd)
	CustomersCmd.AddCommand(setAttributesCmd)
	CustomersCmd.AddCommand(listSubscriptionsCmd)
	CustomersCmd.AddCommand(listPurchasesCmd)
	CustomersCmd.AddCommand(listInvoicesCmd)
	CustomersCmd.AddCommand(transferCmd)
	CustomersCmd.AddCommand(grantEntitlementCmd)
	CustomersCmd.AddCommand(revokeEntitlementCmd)
	CustomersCmd.AddCommand(assignOfferingCmd)
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

	path := fmt.Sprintf("/projects/%s/customers", client.GetProjectID())

	if allPages {
		var items []map[string]interface{}
		err := client.ListAll(ctx, path, limit, func(raw json.RawMessage) error {
			var page []map[string]interface{}
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
		Items    []map[string]interface{} `json:"items"`
		NextPage string                   `json:"next_page,omitempty"`
	}
	if err := client.Get(ctx, path+query, &resp); err != nil {
		return err
	}
	if resp.NextPage != "" {
		output.PrintInfo("More results available. Use --starting-after=%s for next page, or --all for everything.", resp.NextPage)
	}
	return output.Print(resp.Items)
}

func runCreate(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would create customer '%s'", customerID)
		return nil
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}
	ctx, cancel := client.Context()
	defer cancel()

	body := map[string]interface{}{"app_user_id": customerID}
	path := fmt.Sprintf("/projects/%s/customers", client.GetProjectID())
	var customer map[string]interface{}
	if err := client.Post(ctx, path, body, &customer); err != nil {
		return err
	}

	output.PrintSuccess("Customer '%s' created successfully", customerID)
	return output.Print(customer)
}

func runListActiveEntitlements(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}
	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}
	ctx, cancel := client.Context()
	defer cancel()

	path := fmt.Sprintf("/projects/%s/customers/%s/active_entitlements", client.GetProjectID(), customerID)

	if allPages {
		var items []map[string]interface{}
		err := client.ListAll(ctx, path, limit, func(raw json.RawMessage) error {
			var page []map[string]interface{}
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
		Items    []map[string]interface{} `json:"items"`
		NextPage string                   `json:"next_page,omitempty"`
	}
	if err := client.Get(ctx, path+query, &resp); err != nil {
		return err
	}
	if resp.NextPage != "" {
		output.PrintInfo("More results available. Use --starting-after=%s for next page, or --all for everything.", resp.NextPage)
	}
	return output.Print(resp.Items)
}

func runListAliases(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}
	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}
	ctx, cancel := client.Context()
	defer cancel()

	var resp map[string]interface{}
	path := fmt.Sprintf("/projects/%s/customers/%s/aliases", client.GetProjectID(), customerID)
	if err := client.Get(ctx, path, &resp); err != nil {
		return err
	}
	return output.Print(resp)
}

func runListAttributes(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}
	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}
	ctx, cancel := client.Context()
	defer cancel()

	var resp map[string]interface{}
	path := fmt.Sprintf("/projects/%s/customers/%s/attributes", client.GetProjectID(), customerID)
	if err := client.Get(ctx, path, &resp); err != nil {
		return err
	}
	return output.Print(resp)
}

func runSetAttributes(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would set attributes on customer '%s'", customerID)
		return nil
	}

	var body map[string]interface{}
	if err := json.Unmarshal([]byte(attributes), &body); err != nil {
		return fmt.Errorf("invalid JSON for --attributes: %w", err)
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}
	ctx, cancel := client.Context()
	defer cancel()

	path := fmt.Sprintf("/projects/%s/customers/%s/attributes", client.GetProjectID(), customerID)
	if err := client.Post(ctx, path, body, nil); err != nil {
		return err
	}

	output.PrintSuccess("Attributes set on customer '%s'", customerID)
	return nil
}

func runListSubscriptions(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}
	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}
	ctx, cancel := client.Context()
	defer cancel()

	path := fmt.Sprintf("/projects/%s/customers/%s/subscriptions", client.GetProjectID(), customerID)

	if allPages {
		var items []map[string]interface{}
		err := client.ListAll(ctx, path, limit, func(raw json.RawMessage) error {
			var page []map[string]interface{}
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
		Items    []map[string]interface{} `json:"items"`
		NextPage string                   `json:"next_page,omitempty"`
	}
	if err := client.Get(ctx, path+query, &resp); err != nil {
		return err
	}
	if resp.NextPage != "" {
		output.PrintInfo("More results available. Use --starting-after=%s for next page, or --all for everything.", resp.NextPage)
	}
	return output.Print(resp.Items)
}

func runListPurchases(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}
	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}
	ctx, cancel := client.Context()
	defer cancel()

	path := fmt.Sprintf("/projects/%s/customers/%s/purchases", client.GetProjectID(), customerID)

	if allPages {
		var items []map[string]interface{}
		err := client.ListAll(ctx, path, limit, func(raw json.RawMessage) error {
			var page []map[string]interface{}
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
		Items    []map[string]interface{} `json:"items"`
		NextPage string                   `json:"next_page,omitempty"`
	}
	if err := client.Get(ctx, path+query, &resp); err != nil {
		return err
	}
	if resp.NextPage != "" {
		output.PrintInfo("More results available. Use --starting-after=%s for next page, or --all for everything.", resp.NextPage)
	}
	return output.Print(resp.Items)
}

func runListInvoices(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}
	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}
	ctx, cancel := client.Context()
	defer cancel()

	path := fmt.Sprintf("/projects/%s/customers/%s/invoices", client.GetProjectID(), customerID)

	if allPages {
		var items []map[string]interface{}
		err := client.ListAll(ctx, path, limit, func(raw json.RawMessage) error {
			var page []map[string]interface{}
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
		Items    []map[string]interface{} `json:"items"`
		NextPage string                   `json:"next_page,omitempty"`
	}
	if err := client.Get(ctx, path+query, &resp); err != nil {
		return err
	}
	if resp.NextPage != "" {
		output.PrintInfo("More results available. Use --starting-after=%s for next page, or --all for everything.", resp.NextPage)
	}
	return output.Print(resp.Items)
}

func runTransfer(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would transfer customer '%s' to '%s'", customerID, targetID)
		return nil
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}
	ctx, cancel := client.Context()
	defer cancel()

	body := map[string]interface{}{"target_customer_id": targetID}
	path := fmt.Sprintf("/projects/%s/customers/%s/actions/transfer", client.GetProjectID(), customerID)
	if err := client.Post(ctx, path, body, nil); err != nil {
		return err
	}

	output.PrintSuccess("Customer '%s' transferred to '%s'", customerID, targetID)
	return nil
}

func runGrantEntitlement(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would grant entitlement '%s' to customer '%s'", entitlementID, customerID)
		return nil
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}
	ctx, cancel := client.Context()
	defer cancel()

	body := map[string]interface{}{"entitlement_id": entitlementID}
	path := fmt.Sprintf("/projects/%s/customers/%s/actions/grant_entitlement", client.GetProjectID(), customerID)
	if err := client.Post(ctx, path, body, nil); err != nil {
		return err
	}

	output.PrintSuccess("Entitlement '%s' granted to customer '%s'", entitlementID, customerID)
	return nil
}

func runRevokeEntitlement(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would revoke entitlement '%s' from customer '%s'", entitlementID, customerID)
		return nil
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}
	ctx, cancel := client.Context()
	defer cancel()

	body := map[string]interface{}{"entitlement_id": entitlementID}
	path := fmt.Sprintf("/projects/%s/customers/%s/actions/revoke_granted_entitlement", client.GetProjectID(), customerID)
	if err := client.Post(ctx, path, body, nil); err != nil {
		return err
	}

	output.PrintSuccess("Entitlement '%s' revoked from customer '%s'", entitlementID, customerID)
	return nil
}

func runAssignOffering(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would assign offering '%s' to customer '%s'", offeringID, customerID)
		return nil
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}
	ctx, cancel := client.Context()
	defer cancel()

	body := map[string]interface{}{"offering_id": offeringID}
	path := fmt.Sprintf("/projects/%s/customers/%s/actions/assign_offering", client.GetProjectID(), customerID)
	if err := client.Post(ctx, path, body, nil); err != nil {
		return err
	}

	output.PrintSuccess("Offering '%s' assigned to customer '%s'", offeringID, customerID)
	return nil
}
