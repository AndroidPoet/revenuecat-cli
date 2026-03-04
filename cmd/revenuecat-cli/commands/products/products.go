package products

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
	storeIdentifier string
	productType     string
	appID           string
	limit           int
	startAfter      string
	allPages        bool
)

// ProductsCmd manages products
var ProductsCmd = &cobra.Command{
	Use:   "products",
	Short: "Manage products",
	Long:  `List and create products in your RevenueCat project.`,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all products",
	RunE:  runList,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new product",
	RunE:  runCreate,
}

func init() {
	listCmd.Flags().IntVar(&limit, "limit", 20, "number of results per page")
	listCmd.Flags().StringVar(&startAfter, "starting-after", "", "cursor for pagination")
	listCmd.Flags().BoolVar(&allPages, "all", false, "fetch all pages")

	createCmd.Flags().StringVar(&storeIdentifier, "store-identifier", "", "store product identifier")
	createCmd.Flags().StringVar(&productType, "type", "", "product type (subscription, one_time)")
	createCmd.Flags().StringVar(&appID, "app-id", "", "app ID this product belongs to")
	createCmd.MarkFlagRequired("store-identifier")
	createCmd.MarkFlagRequired("type")
	createCmd.MarkFlagRequired("app-id")

	ProductsCmd.AddCommand(listCmd)
	ProductsCmd.AddCommand(createCmd)
}

type ProductInfo struct {
	ID              string `json:"id"`
	StoreIdentifier string `json:"store_identifier"`
	Type            string `json:"type"`
	AppID           string `json:"app_id,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
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

	path := fmt.Sprintf("/projects/%s/products", client.GetProjectID())

	if allPages {
		var items []ProductInfo
		err := client.ListAll(ctx, path, limit, func(raw json.RawMessage) error {
			var page []ProductInfo
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
		Items    []ProductInfo `json:"items"`
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
		output.PrintInfo("Dry run: would create product '%s' of type '%s'", storeIdentifier, productType)
		return nil
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}

	ctx, cancel := client.Context()
	defer cancel()

	body := map[string]interface{}{
		"store_identifier": storeIdentifier,
		"type":             productType,
		"app_id":           appID,
	}

	var product ProductInfo
	path := fmt.Sprintf("/projects/%s/products", client.GetProjectID())
	if err := client.Post(ctx, path, body, &product); err != nil {
		return err
	}

	output.PrintSuccess("Product '%s' created successfully", storeIdentifier)
	return output.Print(product)
}
