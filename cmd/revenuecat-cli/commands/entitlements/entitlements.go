package entitlements

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
	entitlementID string
	lookupKey     string
	displayName   string
	productIDs    []string
	limit         int
	startAfter    string
	allPages      bool
)

// EntitlementsCmd manages entitlements
var EntitlementsCmd = &cobra.Command{
	Use:   "entitlements",
	Short: "Manage entitlements",
	Long: `Create, update, delete, and manage entitlements and their product associations
in your RevenueCat project.`,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all entitlements",
	RunE:  runList,
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get entitlement details",
	RunE:  runGet,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new entitlement",
	RunE:  runCreate,
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an entitlement",
	RunE:  runUpdate,
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an entitlement",
	RunE:  runDelete,
}

var listProductsCmd = &cobra.Command{
	Use:   "list-products",
	Short: "List products attached to an entitlement",
	RunE:  runListProducts,
}

var attachProductsCmd = &cobra.Command{
	Use:   "attach-products",
	Short: "Attach products to an entitlement",
	RunE:  runAttachProducts,
}

var detachProductsCmd = &cobra.Command{
	Use:   "detach-products",
	Short: "Detach products from an entitlement",
	RunE:  runDetachProducts,
}

func init() {
	listCmd.Flags().IntVar(&limit, "limit", 20, "number of results per page")
	listCmd.Flags().StringVar(&startAfter, "starting-after", "", "cursor for pagination")
	listCmd.Flags().BoolVar(&allPages, "all", false, "fetch all pages")

	getCmd.Flags().StringVar(&entitlementID, "entitlement-id", "", "entitlement ID")
	getCmd.MarkFlagRequired("entitlement-id")

	createCmd.Flags().StringVar(&lookupKey, "lookup-key", "", "entitlement lookup key")
	createCmd.Flags().StringVar(&displayName, "display-name", "", "entitlement display name")
	createCmd.MarkFlagRequired("lookup-key")
	createCmd.MarkFlagRequired("display-name")

	updateCmd.Flags().StringVar(&entitlementID, "entitlement-id", "", "entitlement ID")
	updateCmd.Flags().StringVar(&displayName, "display-name", "", "new display name")
	updateCmd.MarkFlagRequired("entitlement-id")
	updateCmd.MarkFlagRequired("display-name")

	var confirm bool
	deleteCmd.Flags().StringVar(&entitlementID, "entitlement-id", "", "entitlement ID")
	deleteCmd.Flags().BoolVar(&confirm, "confirm", false, "confirm deletion")
	deleteCmd.MarkFlagRequired("entitlement-id")

	listProductsCmd.Flags().StringVar(&entitlementID, "entitlement-id", "", "entitlement ID")
	listProductsCmd.MarkFlagRequired("entitlement-id")

	attachProductsCmd.Flags().StringVar(&entitlementID, "entitlement-id", "", "entitlement ID")
	attachProductsCmd.Flags().StringSliceVar(&productIDs, "product-ids", nil, "product IDs to attach")
	attachProductsCmd.MarkFlagRequired("entitlement-id")
	attachProductsCmd.MarkFlagRequired("product-ids")

	detachProductsCmd.Flags().StringVar(&entitlementID, "entitlement-id", "", "entitlement ID")
	detachProductsCmd.Flags().StringSliceVar(&productIDs, "product-ids", nil, "product IDs to detach")
	detachProductsCmd.MarkFlagRequired("entitlement-id")
	detachProductsCmd.MarkFlagRequired("product-ids")

	EntitlementsCmd.AddCommand(listCmd)
	EntitlementsCmd.AddCommand(getCmd)
	EntitlementsCmd.AddCommand(createCmd)
	EntitlementsCmd.AddCommand(updateCmd)
	EntitlementsCmd.AddCommand(deleteCmd)
	EntitlementsCmd.AddCommand(listProductsCmd)
	EntitlementsCmd.AddCommand(attachProductsCmd)
	EntitlementsCmd.AddCommand(detachProductsCmd)
}

type EntitlementInfo struct {
	ID          string `json:"id"`
	LookupKey   string `json:"lookup_key"`
	DisplayName string `json:"display_name"`
	CreatedAt   string `json:"created_at,omitempty"`
}

func parseTimeout() time.Duration {
	d, err := time.ParseDuration(cli.GetTimeout())
	if err != nil {
		return 60 * time.Second
	}
	return d
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

	path := fmt.Sprintf("/projects/%s/entitlements", client.GetProjectID())

	if allPages {
		var items []EntitlementInfo
		err := client.ListAll(ctx, path, limit, func(raw json.RawMessage) error {
			var page []EntitlementInfo
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
		Items    []EntitlementInfo `json:"items"`
		NextPage string           `json:"next_page,omitempty"`
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

	var ent EntitlementInfo
	path := fmt.Sprintf("/projects/%s/entitlements/%s", client.GetProjectID(), entitlementID)
	if err := client.Get(ctx, path, &ent); err != nil {
		return err
	}

	return output.Print(ent)
}

func runCreate(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would create entitlement '%s'", lookupKey)
		return nil
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}

	ctx, cancel := client.Context()
	defer cancel()

	body := map[string]interface{}{
		"lookup_key":   lookupKey,
		"display_name": displayName,
	}

	var ent EntitlementInfo
	path := fmt.Sprintf("/projects/%s/entitlements", client.GetProjectID())
	if err := client.Post(ctx, path, body, &ent); err != nil {
		return err
	}

	output.PrintSuccess("Entitlement '%s' created successfully", lookupKey)
	return output.Print(ent)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would update entitlement '%s'", entitlementID)
		return nil
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}

	ctx, cancel := client.Context()
	defer cancel()

	body := map[string]interface{}{
		"display_name": displayName,
	}

	var ent EntitlementInfo
	path := fmt.Sprintf("/projects/%s/entitlements/%s", client.GetProjectID(), entitlementID)
	if err := client.Patch(ctx, path, body, &ent); err != nil {
		return err
	}

	output.PrintSuccess("Entitlement '%s' updated successfully", entitlementID)
	return output.Print(ent)
}

func runDelete(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if err := cli.CheckConfirm(cmd); err != nil {
		return err
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would delete entitlement '%s'", entitlementID)
		return nil
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}

	ctx, cancel := client.Context()
	defer cancel()

	path := fmt.Sprintf("/projects/%s/entitlements/%s", client.GetProjectID(), entitlementID)
	if err := client.Delete(ctx, path); err != nil {
		return err
	}

	output.PrintSuccess("Entitlement '%s' deleted", entitlementID)
	return nil
}

func runListProducts(cmd *cobra.Command, args []string) error {
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
		Items []struct {
			ID              string `json:"id"`
			StoreIdentifier string `json:"store_identifier"`
			Type            string `json:"type"`
		} `json:"items"`
	}

	path := fmt.Sprintf("/projects/%s/entitlements/%s/products", client.GetProjectID(), entitlementID)
	if err := client.Get(ctx, path, &resp); err != nil {
		return err
	}

	return output.Print(resp.Items)
}

func runAttachProducts(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would attach %d products to entitlement '%s'", len(productIDs), entitlementID)
		return nil
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}

	ctx, cancel := client.Context()
	defer cancel()

	body := map[string]interface{}{
		"product_ids": productIDs,
	}

	path := fmt.Sprintf("/projects/%s/entitlements/%s/products/attach", client.GetProjectID(), entitlementID)
	if err := client.Post(ctx, path, body, nil); err != nil {
		return err
	}

	output.PrintSuccess("Attached %d products to entitlement '%s'", len(productIDs), entitlementID)
	return nil
}

func runDetachProducts(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would detach %d products from entitlement '%s'", len(productIDs), entitlementID)
		return nil
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}

	ctx, cancel := client.Context()
	defer cancel()

	body := map[string]interface{}{
		"product_ids": productIDs,
	}

	path := fmt.Sprintf("/projects/%s/entitlements/%s/products/detach", client.GetProjectID(), entitlementID)
	if err := client.Post(ctx, path, body, nil); err != nil {
		return err
	}

	output.PrintSuccess("Detached %d products from entitlement '%s'", len(productIDs), entitlementID)
	return nil
}
