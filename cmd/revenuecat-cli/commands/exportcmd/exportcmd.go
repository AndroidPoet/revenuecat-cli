package exportcmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/AndroidPoet/revenuecat-cli/internal/api"
	"github.com/AndroidPoet/revenuecat-cli/internal/cli"
	"github.com/AndroidPoet/revenuecat-cli/internal/output"
)

var (
	exportFile string
	importFile string
	confirm    bool
)

// ExportCmd exports project configuration to YAML
var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export project configuration",
	Long:  `Export entitlements, offerings, and packages to a YAML file for backup or migration.`,
	RunE:  runExport,
}

// ImportCmd imports project configuration from YAML
var ImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import project configuration",
	Long:  `Import entitlements, offerings, and packages from a YAML file.`,
	RunE:  runImport,
}

func init() {
	ExportCmd.Flags().StringVar(&exportFile, "file", "rc-export.yaml", "output file path")

	ImportCmd.Flags().StringVar(&importFile, "file", "", "input YAML file path")
	ImportCmd.Flags().BoolVar(&confirm, "confirm", false, "confirm import (required for safety)")
	ImportCmd.MarkFlagRequired("file")
}

type ProjectConfig struct {
	Version      string             `yaml:"version" json:"version"`
	ProjectID    string             `yaml:"project_id" json:"project_id"`
	ExportedAt   string             `yaml:"exported_at" json:"exported_at"`
	Entitlements []EntitlementExport `yaml:"entitlements" json:"entitlements"`
	Offerings    []OfferingExport    `yaml:"offerings" json:"offerings"`
}

type EntitlementExport struct {
	LookupKey   string `yaml:"lookup_key" json:"lookup_key"`
	DisplayName string `yaml:"display_name" json:"display_name"`
}

type OfferingExport struct {
	LookupKey   string          `yaml:"lookup_key" json:"lookup_key"`
	DisplayName string          `yaml:"display_name" json:"display_name"`
	IsCurrent   bool            `yaml:"is_current" json:"is_current"`
	Packages    []PackageExport `yaml:"packages" json:"packages"`
}

type PackageExport struct {
	LookupKey   string `yaml:"lookup_key" json:"lookup_key"`
	DisplayName string `yaml:"display_name" json:"display_name"`
}

func parseTimeout() time.Duration {
	d, err := time.ParseDuration(cli.GetTimeout())
	if err != nil {
		return 60 * time.Second
	}
	return d
}

func runExport(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}

	pid := client.GetProjectID()

	// Fetch all entitlements
	ctx, cancel := client.Context()
	defer cancel()

	var entitlements []EntitlementExport
	entPath := fmt.Sprintf("/projects/%s/entitlements", pid)
	err = client.ListAll(ctx, entPath, 100, func(raw json.RawMessage) error {
		var items []struct {
			LookupKey   string `json:"lookup_key"`
			DisplayName string `json:"display_name"`
		}
		if err := json.Unmarshal(raw, &items); err != nil {
			return err
		}
		for _, item := range items {
			entitlements = append(entitlements, EntitlementExport{
				LookupKey:   item.LookupKey,
				DisplayName: item.DisplayName,
			})
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("fetching entitlements: %w", err)
	}

	// Fetch all offerings
	type offeringRaw struct {
		ID          string `json:"id"`
		LookupKey   string `json:"lookup_key"`
		DisplayName string `json:"display_name"`
		IsCurrent   bool   `json:"is_current"`
	}

	var rawOfferings []offeringRaw
	offPath := fmt.Sprintf("/projects/%s/offerings", pid)
	err = client.ListAll(ctx, offPath, 100, func(raw json.RawMessage) error {
		var items []offeringRaw
		if err := json.Unmarshal(raw, &items); err != nil {
			return err
		}
		rawOfferings = append(rawOfferings, items...)
		return nil
	})
	if err != nil {
		return fmt.Errorf("fetching offerings: %w", err)
	}

	// For each offering, fetch its packages
	var offerings []OfferingExport
	for _, off := range rawOfferings {
		pkgPath := fmt.Sprintf("/projects/%s/offerings/%s/packages", pid, off.ID)
		var packages []PackageExport
		err = client.ListAll(ctx, pkgPath, 100, func(raw json.RawMessage) error {
			var items []struct {
				LookupKey   string `json:"lookup_key"`
				DisplayName string `json:"display_name"`
			}
			if err := json.Unmarshal(raw, &items); err != nil {
				return err
			}
			for _, item := range items {
				packages = append(packages, PackageExport{
					LookupKey:   item.LookupKey,
					DisplayName: item.DisplayName,
				})
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("fetching packages for offering %s: %w", off.LookupKey, err)
		}

		offerings = append(offerings, OfferingExport{
			LookupKey:   off.LookupKey,
			DisplayName: off.DisplayName,
			IsCurrent:   off.IsCurrent,
			Packages:    packages,
		})
	}

	config := ProjectConfig{
		Version:      "1",
		ProjectID:    pid,
		ExportedAt:   time.Now().UTC().Format(time.RFC3339),
		Entitlements: entitlements,
		Offerings:    offerings,
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("marshaling YAML: %w", err)
	}

	if err := os.WriteFile(exportFile, data, 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	output.PrintSuccess("Exported to %s", exportFile)
	return nil
}

func runImport(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	if !confirm {
		return fmt.Errorf("this is a destructive operation. Use --confirm to proceed")
	}

	data, err := os.ReadFile(importFile)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	var config ProjectConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("parsing YAML: %w", err)
	}

	if cli.IsDryRun() {
		output.PrintInfo("Dry run: would import %d entitlements, %d offerings from %s",
			len(config.Entitlements), len(config.Offerings), importFile)
		return nil
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}

	pid := client.GetProjectID()
	var createdEnts, createdOffs, createdPkgs int

	// Import entitlements
	for _, ent := range config.Entitlements {
		ctx, cancel := client.Context()
		body := map[string]interface{}{
			"lookup_key":   ent.LookupKey,
			"display_name": ent.DisplayName,
		}
		path := fmt.Sprintf("/projects/%s/entitlements", pid)
		if err := client.Post(ctx, path, body, nil); err != nil {
			cancel()
			output.PrintWarning("Entitlement '%s': %s (skipping)", ent.LookupKey, err)
			continue
		}
		cancel()
		createdEnts++
	}

	// Import offerings and their packages
	for _, off := range config.Offerings {
		ctx, cancel := client.Context()
		body := map[string]interface{}{
			"lookup_key":   off.LookupKey,
			"display_name": off.DisplayName,
		}
		var created struct {
			ID string `json:"id"`
		}
		path := fmt.Sprintf("/projects/%s/offerings", pid)
		if err := client.Post(ctx, path, body, &created); err != nil {
			cancel()
			output.PrintWarning("Offering '%s': %s (skipping)", off.LookupKey, err)
			continue
		}
		cancel()
		createdOffs++

		// Create packages for this offering
		for _, pkg := range off.Packages {
			ctx2, cancel2 := client.Context()
			pkgBody := map[string]interface{}{
				"lookup_key":   pkg.LookupKey,
				"display_name": pkg.DisplayName,
			}
			pkgPath := fmt.Sprintf("/projects/%s/offerings/%s/packages", pid, created.ID)
			if err := client.Post(ctx2, pkgPath, pkgBody, nil); err != nil {
				cancel2()
				output.PrintWarning("Package '%s' in offering '%s': %s (skipping)", pkg.LookupKey, off.LookupKey, err)
				continue
			}
			cancel2()
			createdPkgs++
		}
	}

	output.PrintSuccess("Import complete: %d entitlements, %d offerings, %d packages created", createdEnts, createdOffs, createdPkgs)
	return nil
}
