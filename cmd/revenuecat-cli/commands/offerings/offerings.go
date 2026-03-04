package offerings

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
	lookupKey   string
	displayName string
	isCurrent   bool
	metadata    string
	limit       int
	startAfter  string
	allPages    bool
)

// OfferingsCmd manages offerings
var OfferingsCmd = &cobra.Command{
	Use:   "offerings",
	Short: "Manage offerings",
	Long:  `List, create, and update offerings in your RevenueCat project.`,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all offerings",
	RunE:  runList,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new offering",
	RunE:  runCreate,
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an offering",
	RunE:  runUpdate,
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get offering details",
	RunE:  runGet,
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an offering",
	RunE:  runDelete,
}

func init() {
	listCmd.Flags().IntVar(&limit, "limit", 20, "number of results per page")
	listCmd.Flags().StringVar(&startAfter, "starting-after", "", "cursor for pagination")
	listCmd.Flags().BoolVar(&allPages, "all", false, "fetch all pages")

	createCmd.Flags().StringVar(&lookupKey, "lookup-key", "", "offering lookup key")
	createCmd.Flags().StringVar(&displayName, "display-name", "", "offering display name")
	createCmd.Flags().StringVar(&metadata, "metadata", "", "JSON metadata string")
	createCmd.MarkFlagRequired("lookup-key")
	createCmd.MarkFlagRequired("display-name")

	updateCmd.Flags().StringVar(&offeringID, "offering-id", "", "offering ID")
	updateCmd.Flags().StringVar(&displayName, "display-name", "", "new display name")
	updateCmd.Flags().BoolVar(&isCurrent, "is-current", false, "set as current offering")
	updateCmd.Flags().StringVar(&metadata, "metadata", "", "JSON metadata string")
	updateCmd.MarkFlagRequired("offering-id")

	getCmd.Flags().StringVar(&offeringID, "offering-id", "", "offering ID")
	getCmd.MarkFlagRequired("offering-id")

	var confirm bool
	deleteCmd.Flags().StringVar(&offeringID, "offering-id", "", "offering ID")
	deleteCmd.Flags().BoolVar(&confirm, "confirm", false, "confirm deletion")
	deleteCmd.MarkFlagRequired("offering-id")

	OfferingsCmd.AddCommand(listCmd)
	OfferingsCmd.AddCommand(createCmd)
	OfferingsCmd.AddCommand(updateCmd)
	OfferingsCmd.AddCommand(getCmd)
	OfferingsCmd.AddCommand(deleteCmd)
}

type OfferingInfo struct {
	ID          string `json:"id"`
	LookupKey   string `json:"lookup_key"`
	DisplayName string `json:"display_name"`
	IsCurrent   bool   `json:"is_current"`
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

	path := fmt.Sprintf("/projects/%s/offerings", client.GetProjectID())

	if allPages {
		var items []OfferingInfo
		err := client.ListAll(ctx, path, limit, func(raw json.RawMessage) error {
			var page []OfferingInfo
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
		Items    []OfferingInfo `json:"items"`
		NextPage string        `json:"next_page,omitempty"`
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
		output.PrintInfo("Dry run: would create offering '%s'", lookupKey)
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
	if metadata != "" {
		var meta map[string]interface{}
		if err := json.Unmarshal([]byte(metadata), &meta); err != nil {
			return fmt.Errorf("invalid metadata JSON: %w", err)
		}
		body["metadata"] = meta
	}

	var offering OfferingInfo
	path := fmt.Sprintf("/projects/%s/offerings", client.GetProjectID())
	if err := client.Post(ctx, path, body, &offering); err != nil {
		return err
	}

	output.PrintSuccess("Offering '%s' created successfully", lookupKey)
	return output.Print(offering)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would update offering '%s'", offeringID)
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
	if cmd.Flags().Changed("is-current") {
		body["is_current"] = isCurrent
	}
	if metadata != "" {
		var meta map[string]interface{}
		if err := json.Unmarshal([]byte(metadata), &meta); err != nil {
			return fmt.Errorf("invalid metadata JSON: %w", err)
		}
		body["metadata"] = meta
	}

	var offering OfferingInfo
	path := fmt.Sprintf("/projects/%s/offerings/%s", client.GetProjectID(), offeringID)
	if err := client.Patch(ctx, path, body, &offering); err != nil {
		return err
	}

	output.PrintSuccess("Offering '%s' updated successfully", offeringID)
	return output.Print(offering)
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

	var offering OfferingInfo
	path := fmt.Sprintf("/projects/%s/offerings/%s", client.GetProjectID(), offeringID)
	if err := client.Get(ctx, path, &offering); err != nil {
		return err
	}
	return output.Print(offering)
}

func runDelete(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}
	if err := cli.CheckConfirm(cmd); err != nil {
		return err
	}
	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would delete offering '%s'", offeringID)
		return nil
	}
	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}
	ctx, cancel := client.Context()
	defer cancel()

	path := fmt.Sprintf("/projects/%s/offerings/%s", client.GetProjectID(), offeringID)
	if err := client.Delete(ctx, path); err != nil {
		return err
	}
	output.PrintSuccess("Offering '%s' deleted", offeringID)
	return nil
}
