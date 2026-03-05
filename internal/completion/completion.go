package completion

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/AndroidPoet/revenuecat-cli/internal/api"
	"github.com/AndroidPoet/revenuecat-cli/internal/cli"
)

// AppIDs returns a completion function for app IDs
func AppIDs() func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		items := fetchIDs("/projects/%s/apps", "id", "name")
		return items, cobra.ShellCompDirectiveNoFileComp
	}
}

// ProductIDs returns a completion function for product IDs
func ProductIDs() func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		items := fetchIDs("/projects/%s/products", "id", "store_identifier")
		return items, cobra.ShellCompDirectiveNoFileComp
	}
}

// EntitlementIDs returns a completion function for entitlement IDs
func EntitlementIDs() func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		items := fetchIDs("/projects/%s/entitlements", "id", "lookup_key")
		return items, cobra.ShellCompDirectiveNoFileComp
	}
}

// OfferingIDs returns a completion function for offering IDs
func OfferingIDs() func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		items := fetchIDs("/projects/%s/offerings", "id", "lookup_key")
		return items, cobra.ShellCompDirectiveNoFileComp
	}
}

// PackageIDs returns a completion function for package IDs
func PackageIDs() func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Packages require an offering ID, so just disable file comp
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
}

// fetchIDs is a helper that queries an API list endpoint and returns ID suggestions.
// Each suggestion is formatted as "id\tdescription" for shell completion.
func fetchIDs(pathTemplate, idField, descField string) []string {
	projectID := cli.GetProjectID()
	if projectID == "" {
		return nil
	}

	client, err := api.NewClient(projectID, 5*time.Second)
	if err != nil {
		return nil
	}

	ctx, cancel := client.Context()
	defer cancel()

	path := fmt.Sprintf(pathTemplate+"?limit=20", projectID)

	var resp struct {
		Items []map[string]interface{} `json:"items"`
	}
	if err := client.Get(ctx, path, &resp); err != nil {
		return nil
	}

	var results []string
	for _, item := range resp.Items {
		id, ok := item[idField].(string)
		if !ok {
			continue
		}
		desc, _ := item[descField].(string)
		if desc != "" {
			results = append(results, fmt.Sprintf("%s\t%s", id, desc))
		} else {
			results = append(results, id)
		}
	}
	return results
}
