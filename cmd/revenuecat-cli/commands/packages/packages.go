package packages

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
	offeringID  string
	packageID   string
	lookupKey   string
	displayName string
	productIDs  []string
	limit       int
	startAfter  string
	allPages    bool
)

// PackagesCmd manages packages
var PackagesCmd = &cobra.Command{
	Use:   "packages",
	Short: "Manage packages",
	Long:  `List, create, and manage packages within offerings in your RevenueCat project.`,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List packages in an offering",
	RunE:  runList,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new package",
	RunE:  runCreate,
}

var attachProductsCmd = &cobra.Command{
	Use:   "attach-products",
	Short: "Attach products to a package",
	RunE:  runAttachProducts,
}

var detachProductsCmd = &cobra.Command{
	Use:   "detach-products",
	Short: "Detach products from a package",
	RunE:  runDetachProducts,
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get package details",
	RunE:  runGet,
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a package",
	RunE:  runUpdate,
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a package",
	RunE:  runDelete,
}

var listProductsCmd = &cobra.Command{
	Use:   "list-products",
	Short: "List products in a package",
	RunE:  runListProducts,
}

func init() {
	listCmd.Flags().StringVar(&offeringID, "offering-id", "", "offering ID")
	listCmd.Flags().IntVar(&limit, "limit", 20, "number of results per page")
	listCmd.Flags().StringVar(&startAfter, "starting-after", "", "cursor for pagination")
	listCmd.Flags().BoolVar(&allPages, "all", false, "fetch all pages")
	listCmd.MarkFlagRequired("offering-id")

	createCmd.Flags().StringVar(&offeringID, "offering-id", "", "offering ID")
	createCmd.Flags().StringVar(&lookupKey, "lookup-key", "", "package lookup key")
	createCmd.Flags().StringVar(&displayName, "display-name", "", "package display name")
	createCmd.MarkFlagRequired("offering-id")
	createCmd.MarkFlagRequired("lookup-key")
	createCmd.MarkFlagRequired("display-name")

	attachProductsCmd.Flags().StringVar(&packageID, "package-id", "", "package ID")
	attachProductsCmd.Flags().StringSliceVar(&productIDs, "product-ids", nil, "product IDs to attach")
	attachProductsCmd.MarkFlagRequired("package-id")
	attachProductsCmd.MarkFlagRequired("product-ids")

	detachProductsCmd.Flags().StringVar(&packageID, "package-id", "", "package ID")
	detachProductsCmd.Flags().StringSliceVar(&productIDs, "product-ids", nil, "product IDs to detach")
	detachProductsCmd.MarkFlagRequired("package-id")
	detachProductsCmd.MarkFlagRequired("product-ids")

	getCmd.Flags().StringVar(&packageID, "package-id", "", "package ID")
	getCmd.MarkFlagRequired("package-id")

	updateCmd.Flags().StringVar(&packageID, "package-id", "", "package ID")
	updateCmd.Flags().StringVar(&displayName, "display-name", "", "new display name")
	updateCmd.MarkFlagRequired("package-id")

	var confirm bool
	deleteCmd.Flags().StringVar(&packageID, "package-id", "", "package ID")
	deleteCmd.Flags().BoolVar(&confirm, "confirm", false, "confirm deletion")
	deleteCmd.MarkFlagRequired("package-id")

	listProductsCmd.Flags().StringVar(&packageID, "package-id", "", "package ID")
	listProductsCmd.Flags().IntVar(&limit, "limit", 20, "number of results per page")
	listProductsCmd.Flags().StringVar(&startAfter, "starting-after", "", "cursor for pagination")
	listProductsCmd.Flags().BoolVar(&allPages, "all", false, "fetch all pages")
	listProductsCmd.MarkFlagRequired("package-id")

	PackagesCmd.AddCommand(listCmd)
	PackagesCmd.AddCommand(createCmd)
	PackagesCmd.AddCommand(attachProductsCmd)
	PackagesCmd.AddCommand(detachProductsCmd)
	PackagesCmd.AddCommand(getCmd)
	PackagesCmd.AddCommand(updateCmd)
	PackagesCmd.AddCommand(deleteCmd)
	PackagesCmd.AddCommand(listProductsCmd)
}

type PackageInfo struct {
	ID          string `json:"id"`
	LookupKey   string `json:"lookup_key"`
	DisplayName string `json:"display_name"`
	OfferingID  string `json:"offering_id,omitempty"`
	CreatedAt   interface{} `json:"created_at,omitempty"`
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

	path := fmt.Sprintf("/projects/%s/offerings/%s/packages", client.GetProjectID(), offeringID)

	if allPages {
		var items []PackageInfo
		err := client.ListAll(ctx, path, limit, func(raw json.RawMessage) error {
			var page []PackageInfo
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
		Items    []PackageInfo `json:"items"`
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

func runCreate(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would create package '%s' in offering '%s'", lookupKey, offeringID)
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

	var pkg PackageInfo
	path := fmt.Sprintf("/projects/%s/offerings/%s/packages", client.GetProjectID(), offeringID)
	if err := client.Post(ctx, path, body, &pkg); err != nil {
		return err
	}

	output.PrintSuccess("Package '%s' created successfully", lookupKey)
	return output.Print(pkg)
}

func runAttachProducts(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would attach %d products to package '%s'", len(productIDs), packageID)
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

	path := fmt.Sprintf("/projects/%s/packages/%s/products/attach", client.GetProjectID(), packageID)
	if err := client.Post(ctx, path, body, nil); err != nil {
		return err
	}

	output.PrintSuccess("Attached %d products to package '%s'", len(productIDs), packageID)
	return nil
}

func runDetachProducts(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would detach %d products from package '%s'", len(productIDs), packageID)
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

	path := fmt.Sprintf("/projects/%s/packages/%s/products/detach", client.GetProjectID(), packageID)
	if err := client.Post(ctx, path, body, nil); err != nil {
		return err
	}

	output.PrintSuccess("Detached %d products from package '%s'", len(productIDs), packageID)
	return nil
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

	var pkg PackageInfo
	path := fmt.Sprintf("/projects/%s/packages/%s", client.GetProjectID(), packageID)
	if err := client.Get(ctx, path, &pkg); err != nil {
		return err
	}
	return output.Print(pkg)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would update package '%s'", packageID)
		return nil
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}

	ctx, cancel := client.Context()
	defer cancel()

	body := map[string]interface{}{}
	if displayName != "" {
		body["display_name"] = displayName
	}

	var pkg PackageInfo
	path := fmt.Sprintf("/projects/%s/packages/%s", client.GetProjectID(), packageID)
	if err := client.Post(ctx, path, body, &pkg); err != nil {
		return err
	}

	output.PrintSuccess("Package '%s' updated successfully", packageID)
	return output.Print(pkg)
}

func runDelete(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}
	if err := cli.CheckConfirm(cmd); err != nil {
		return err
	}
	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would delete package '%s'", packageID)
		return nil
	}
	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}
	ctx, cancel := client.Context()
	defer cancel()

	path := fmt.Sprintf("/projects/%s/packages/%s", client.GetProjectID(), packageID)
	if err := client.Delete(ctx, path); err != nil {
		return err
	}
	output.PrintSuccess("Package '%s' deleted", packageID)
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

	path := fmt.Sprintf("/projects/%s/packages/%s/products", client.GetProjectID(), packageID)

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
