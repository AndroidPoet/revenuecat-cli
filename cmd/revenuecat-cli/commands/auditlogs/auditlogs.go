package auditlogs

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
	limit      int
	startAfter string
	allPages   bool
)

// AuditLogsCmd manages audit logs
var AuditLogsCmd = &cobra.Command{
	Use:   "audit-logs",
	Short: "View audit logs",
	Long:  `List audit log entries for your RevenueCat project.`,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List audit log entries",
	RunE:  runList,
}

func init() {
	listCmd.Flags().IntVar(&limit, "limit", 20, "number of results per page")
	listCmd.Flags().StringVar(&startAfter, "starting-after", "", "cursor for pagination")
	listCmd.Flags().BoolVar(&allPages, "all", false, "fetch all pages")

	AuditLogsCmd.AddCommand(listCmd)
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

	path := fmt.Sprintf("/projects/%s/audit_logs", client.GetProjectID())

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
