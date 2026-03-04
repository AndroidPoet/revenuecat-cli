package apps

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
	appID       string
	appName     string
	appType     string
	bundleID    string
	packageName string
	limit       int
	startAfter  string
	allPages    bool
)

// AppsCmd manages applications
var AppsCmd = &cobra.Command{
	Use:   "apps",
	Short: "Manage apps",
	Long:  `List, create, update, and delete apps in your RevenueCat project.`,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all apps",
	RunE:  runList,
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get app details",
	RunE:  runGet,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new app",
	RunE:  runCreate,
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an app",
	RunE:  runUpdate,
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an app",
	RunE:  runDelete,
}

var apiKeysCmd = &cobra.Command{
	Use:   "api-keys",
	Short: "List public API keys for an app",
	RunE:  runAPIKeys,
}

func init() {
	// List flags
	listCmd.Flags().IntVar(&limit, "limit", 20, "number of results per page")
	listCmd.Flags().StringVar(&startAfter, "starting-after", "", "cursor for pagination")
	listCmd.Flags().BoolVar(&allPages, "all", false, "fetch all pages")

	// Get flags
	getCmd.Flags().StringVar(&appID, "app-id", "", "app ID")
	getCmd.MarkFlagRequired("app-id")

	// Create flags
	createCmd.Flags().StringVar(&appName, "name", "", "app name")
	createCmd.Flags().StringVar(&appType, "type", "", "app type (app_store, play_store, stripe, amazon, mac_app_store, roku, web)")
	createCmd.Flags().StringVar(&bundleID, "bundle-id", "", "iOS bundle ID")
	createCmd.Flags().StringVar(&packageName, "package-name", "", "Android package name")
	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("type")

	// Update flags
	updateCmd.Flags().StringVar(&appID, "app-id", "", "app ID")
	updateCmd.Flags().StringVar(&appName, "name", "", "new app name")
	updateCmd.MarkFlagRequired("app-id")

	// Delete flags
	var confirm bool
	deleteCmd.Flags().StringVar(&appID, "app-id", "", "app ID")
	deleteCmd.Flags().BoolVar(&confirm, "confirm", false, "confirm deletion")
	deleteCmd.MarkFlagRequired("app-id")

	// API keys flags
	apiKeysCmd.Flags().StringVar(&appID, "app-id", "", "app ID")
	apiKeysCmd.MarkFlagRequired("app-id")

	AppsCmd.AddCommand(listCmd)
	AppsCmd.AddCommand(getCmd)
	AppsCmd.AddCommand(createCmd)
	AppsCmd.AddCommand(updateCmd)
	AppsCmd.AddCommand(deleteCmd)
	AppsCmd.AddCommand(apiKeysCmd)
}

type AppInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	BundleID    string `json:"bundle_id,omitempty"`
	PackageName string `json:"package_name,omitempty"`
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

	path := fmt.Sprintf("/projects/%s/apps", client.GetProjectID())

	if allPages {
		var items []AppInfo
		err := client.ListAll(ctx, path, limit, func(raw json.RawMessage) error {
			var page []AppInfo
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

	// Single page
	query := fmt.Sprintf("?limit=%d", limit)
	if startAfter != "" {
		query += "&starting_after=" + startAfter
	}

	var resp struct {
		Items      []AppInfo `json:"items"`
		NextPage   string    `json:"next_page,omitempty"`
		TotalCount int       `json:"total_count,omitempty"`
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

	var app AppInfo
	path := fmt.Sprintf("/projects/%s/apps/%s", client.GetProjectID(), appID)
	if err := client.Get(ctx, path, &app); err != nil {
		return err
	}

	return output.Print(app)
}

func runCreate(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would create app '%s' of type '%s'", appName, appType)
		return nil
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}

	ctx, cancel := client.Context()
	defer cancel()

	body := map[string]interface{}{
		"name": appName,
		"type": appType,
	}
	if bundleID != "" {
		body["bundle_id"] = bundleID
	}
	if packageName != "" {
		body["package_name"] = packageName
	}

	var app AppInfo
	path := fmt.Sprintf("/projects/%s/apps", client.GetProjectID())
	if err := client.Post(ctx, path, body, &app); err != nil {
		return err
	}

	output.PrintSuccess("App '%s' created successfully", appName)
	return output.Print(app)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would update app '%s'", appID)
		return nil
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}

	ctx, cancel := client.Context()
	defer cancel()

	body := map[string]interface{}{}
	if appName != "" {
		body["name"] = appName
	}
	if bundleID != "" {
		body["bundle_id"] = bundleID
	}
	if packageName != "" {
		body["package_name"] = packageName
	}

	var app AppInfo
	path := fmt.Sprintf("/projects/%s/apps/%s", client.GetProjectID(), appID)
	if err := client.Patch(ctx, path, body, &app); err != nil {
		return err
	}

	output.PrintSuccess("App '%s' updated successfully", appID)
	return output.Print(app)
}

func runDelete(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if err := cli.CheckConfirm(cmd); err != nil {
		return err
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would delete app '%s'", appID)
		return nil
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}

	ctx, cancel := client.Context()
	defer cancel()

	path := fmt.Sprintf("/projects/%s/apps/%s", client.GetProjectID(), appID)
	if err := client.Delete(ctx, path); err != nil {
		return err
	}

	output.PrintSuccess("App '%s' deleted", appID)
	return nil
}

func runAPIKeys(cmd *cobra.Command, args []string) error {
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
			Key  string `json:"key"`
			Name string `json:"name,omitempty"`
		} `json:"items"`
	}

	path := fmt.Sprintf("/projects/%s/apps/%s/api_keys", client.GetProjectID(), appID)
	if err := client.Get(ctx, path, &resp); err != nil {
		return err
	}

	return output.Print(resp.Items)
}
