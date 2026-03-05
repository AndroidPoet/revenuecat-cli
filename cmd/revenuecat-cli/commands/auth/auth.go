package auth

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/AndroidPoet/revenuecat-cli/internal/config"
	"github.com/AndroidPoet/revenuecat-cli/internal/output"
)

var (
	profileName    string
	apiKey         string
	defaultProject string
)

// AuthCmd manages authentication profiles
var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication profiles",
	Long:  `Manage authentication profiles for RevenueCat API access.`,
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Configure authentication credentials",
	Long: `Configure API key credentials for a profile.

Get your API key from the RevenueCat dashboard:
  https://app.revenuecat.com/settings/api-keys

Examples:
  rc auth login --api-key sk_xxxx
  rc auth login --api-key sk_xxxx --name production --default-project proj_xxxx`,
	RunE: runLogin,
}

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch active profile",
	RunE:  runSwitch,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all profiles",
	RunE:  runList,
}

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show active profile",
	RunE:  runCurrent,
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a profile",
	RunE:  runDelete,
}

func init() {
	// Login flags
	loginCmd.Flags().StringVar(&profileName, "name", "default", "profile name")
	loginCmd.Flags().StringVar(&apiKey, "api-key", "", "RevenueCat API v2 secret key")
	loginCmd.Flags().StringVar(&defaultProject, "default-project", "", "default project ID for this profile")
	loginCmd.MarkFlagRequired("api-key")

	// Switch flags
	switchCmd.Flags().StringVar(&profileName, "name", "", "profile name to switch to")
	switchCmd.MarkFlagRequired("name")

	// Delete flags
	var confirm bool
	deleteCmd.Flags().StringVar(&profileName, "name", "", "profile name to delete")
	deleteCmd.Flags().BoolVar(&confirm, "confirm", false, "confirm deletion")
	deleteCmd.MarkFlagRequired("name")

	AuthCmd.AddCommand(loginCmd)
	AuthCmd.AddCommand(switchCmd)
	AuthCmd.AddCommand(listCmd)
	AuthCmd.AddCommand(currentCmd)
	AuthCmd.AddCommand(deleteCmd)
}

func runLogin(cmd *cobra.Command, args []string) error {
	p := config.Profile{
		Name:           profileName,
		APIKey:         apiKey,
		DefaultProject: defaultProject,
	}

	config.SetProfile(p)
	config.SetDefaultProfile(profileName)

	if err := config.Save(); err != nil {
		return err
	}

	output.PrintSuccess("Profile '%s' configured successfully", profileName)
	return output.Print(map[string]interface{}{
		"profile":         profileName,
		"default_project": defaultProject,
	})
}

func runSwitch(cmd *cobra.Command, args []string) error {
	// Verify profile exists
	profiles := config.ListProfiles()
	found := false
	for _, p := range profiles {
		if p == profileName {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("profile '%s' not found. Run 'rc auth list' to see available profiles", profileName)
	}

	config.SetDefaultProfile(profileName)
	if err := config.Save(); err != nil {
		return err
	}

	output.PrintSuccess("Switched to profile '%s'", profileName)
	return nil
}

type profileInfo struct {
	Name           string `json:"name"`
	DefaultProject string `json:"default_project,omitempty"`
	Active         bool   `json:"active"`
}

func runList(cmd *cobra.Command, args []string) error {
	cfg := config.GetConfig()
	if cfg == nil || len(cfg.Profiles) == 0 {
		output.PrintInfo("No profiles configured. Run 'rc auth login' to get started.")
		return output.Print([]profileInfo{})
	}

	result := make([]profileInfo, 0, len(cfg.Profiles))
	for _, p := range cfg.Profiles {
		result = append(result, profileInfo{
			Name:           p.Name,
			DefaultProject: p.DefaultProject,
			Active:         p.Name == cfg.DefaultProfile,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return output.Print(result)
}

func runCurrent(cmd *cobra.Command, args []string) error {
	p := config.GetProfile()
	if p == nil {
		return fmt.Errorf("no active profile. Run 'rc auth login' first")
	}

	return output.Print(profileInfo{
		Name:           p.Name,
		DefaultProject: p.DefaultProject,
		Active:         true,
	})
}

func runDelete(cmd *cobra.Command, args []string) error {
	confirm, _ := cmd.Flags().GetBool("confirm")
	if !confirm {
		return fmt.Errorf("use --confirm to delete profile '%s'", profileName)
	}

	// Verify profile exists
	profiles := config.ListProfiles()
	found := false
	for _, p := range profiles {
		if p == profileName {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("profile '%s' not found", profileName)
	}

	config.DeleteProfile(profileName)

	// Reset default if we deleted the default profile
	cfg := config.GetConfig()
	if cfg.DefaultProfile == profileName {
		remaining := config.ListProfiles()
		if len(remaining) > 0 {
			config.SetDefaultProfile(remaining[0])
		} else {
			config.SetDefaultProfile("")
		}
	}

	if err := config.Save(); err != nil {
		return err
	}

	output.PrintSuccess("Profile '%s' deleted", profileName)
	return nil
}
