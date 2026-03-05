package diff

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
	source string
	target string
)

// DiffCmd compares two RevenueCat projects
var DiffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compare two projects",
	Long:  `Compare entitlements and offerings between two RevenueCat projects.`,
	RunE:  runDiff,
}

func init() {
	DiffCmd.Flags().StringVar(&source, "source", "", "source project ID")
	DiffCmd.Flags().StringVar(&target, "target", "", "target project ID")
	DiffCmd.MarkFlagRequired("source")
	DiffCmd.MarkFlagRequired("target")
}

type DiffResult struct {
	Entitlements DiffSection `json:"entitlements"`
	Offerings    DiffSection `json:"offerings"`
}

type DiffSection struct {
	OnlyInSource []string `json:"only_in_source"`
	OnlyInTarget []string `json:"only_in_target"`
	InBoth       []string `json:"in_both"`
}

func parseTimeout() time.Duration {
	d, err := time.ParseDuration(cli.GetTimeout())
	if err != nil {
		return 60 * time.Second
	}
	return d
}

type namedItem struct {
	LookupKey string `json:"lookup_key"`
}

func collectKeys(client *api.Client, path string) (map[string]bool, error) {
	ctx, cancel := client.Context()
	defer cancel()

	keys := make(map[string]bool)
	err := client.ListAll(ctx, path, 100, func(raw json.RawMessage) error {
		var items []namedItem
		if err := json.Unmarshal(raw, &items); err != nil {
			return err
		}
		for _, item := range items {
			keys[item.LookupKey] = true
		}
		return nil
	})
	return keys, err
}

func diffKeys(srcKeys, tgtKeys map[string]bool) DiffSection {
	section := DiffSection{
		OnlyInSource: []string{},
		OnlyInTarget: []string{},
		InBoth:       []string{},
	}
	for k := range srcKeys {
		if tgtKeys[k] {
			section.InBoth = append(section.InBoth, k)
		} else {
			section.OnlyInSource = append(section.OnlyInSource, k)
		}
	}
	for k := range tgtKeys {
		if !srcKeys[k] {
			section.OnlyInTarget = append(section.OnlyInTarget, k)
		}
	}
	return section
}

func runDiff(cmd *cobra.Command, args []string) error {
	timeout := parseTimeout()

	srcClient, err := api.NewClient(source, timeout)
	if err != nil {
		return fmt.Errorf("source client: %w", err)
	}

	tgtClient, err := api.NewClient(target, timeout)
	if err != nil {
		return fmt.Errorf("target client: %w", err)
	}

	// Fetch entitlements from both
	srcEntKeys, err := collectKeys(srcClient, fmt.Sprintf("/projects/%s/entitlements", source))
	if err != nil {
		return fmt.Errorf("source entitlements: %w", err)
	}

	tgtEntKeys, err := collectKeys(tgtClient, fmt.Sprintf("/projects/%s/entitlements", target))
	if err != nil {
		return fmt.Errorf("target entitlements: %w", err)
	}

	// Fetch offerings from both
	srcOffKeys, err := collectKeys(srcClient, fmt.Sprintf("/projects/%s/offerings", source))
	if err != nil {
		return fmt.Errorf("source offerings: %w", err)
	}

	tgtOffKeys, err := collectKeys(tgtClient, fmt.Sprintf("/projects/%s/offerings", target))
	if err != nil {
		return fmt.Errorf("target offerings: %w", err)
	}

	result := DiffResult{
		Entitlements: diffKeys(srcEntKeys, tgtEntKeys),
		Offerings:    diffKeys(srcOffKeys, tgtOffKeys),
	}

	output.PrintInfo("Comparing project %s (source) vs %s (target)", source, target)
	output.PrintInfo("")
	output.PrintInfo("Entitlements:")
	output.PrintInfo("  Only in source: %d", len(result.Entitlements.OnlyInSource))
	output.PrintInfo("  Only in target: %d", len(result.Entitlements.OnlyInTarget))
	output.PrintInfo("  In both:        %d", len(result.Entitlements.InBoth))
	output.PrintInfo("")
	output.PrintInfo("Offerings:")
	output.PrintInfo("  Only in source: %d", len(result.Offerings.OnlyInSource))
	output.PrintInfo("  Only in target: %d", len(result.Offerings.OnlyInTarget))
	output.PrintInfo("  In both:        %d", len(result.Offerings.InBoth))

	return output.Print(result)
}
