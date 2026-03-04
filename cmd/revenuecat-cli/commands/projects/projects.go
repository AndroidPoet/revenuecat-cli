package projects

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
	projectName string
	limit       int
	startAfter  string
	allPages    bool
)

// ProjectsCmd manages projects
var ProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Manage projects",
	Long:  `List and create projects in your RevenueCat account.`,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	RunE:  runList,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	RunE:  runCreate,
}

func init() {
	listCmd.Flags().IntVar(&limit, "limit", 20, "number of results per page")
	listCmd.Flags().StringVar(&startAfter, "starting-after", "", "cursor for pagination")
	listCmd.Flags().BoolVar(&allPages, "all", false, "fetch all pages")

	createCmd.Flags().StringVar(&projectName, "name", "", "project name")
	createCmd.MarkFlagRequired("name")

	ProjectsCmd.AddCommand(listCmd)
	ProjectsCmd.AddCommand(createCmd)
}

type ProjectInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at,omitempty"`
}

func parseTimeout() time.Duration {
	d, err := time.ParseDuration(cli.GetTimeout())
	if err != nil {
		return 60 * time.Second
	}
	return d
}

func runList(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient("", parseTimeout())
	if err != nil {
		return err
	}

	ctx, cancel := client.Context()
	defer cancel()

	path := "/projects"

	if allPages {
		var items []ProjectInfo
		err := client.ListAll(ctx, path, limit, func(raw json.RawMessage) error {
			var page []ProjectInfo
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
		Items    []ProjectInfo `json:"items"`
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
	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would create project '%s'", projectName)
		return nil
	}

	client, err := api.NewClient("", parseTimeout())
	if err != nil {
		return err
	}

	ctx, cancel := client.Context()
	defer cancel()

	body := map[string]interface{}{
		"name": projectName,
	}

	var project ProjectInfo
	path := "/projects"
	if err := client.Post(ctx, path, body, &project); err != nil {
		return err
	}

	output.PrintSuccess("Project '%s' created successfully", projectName)
	return output.Print(project)
}
