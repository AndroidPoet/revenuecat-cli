package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/AndroidPoet/revenuecat-cli/cmd/revenuecat-cli/commands/initcmd"
	"github.com/AndroidPoet/revenuecat-cli/internal/cli"
	"github.com/AndroidPoet/revenuecat-cli/internal/config"
	"github.com/AndroidPoet/revenuecat-cli/internal/output"
)

var (
	cfgFile   string
	projectID string
	profile   string
	outputFmt string
	pretty    bool
	quiet     bool
	debug     bool
	timeout   string
	dryRun    bool

	versionStr string
	commitStr  string
	dateStr    string
)

var rootCmd = &cobra.Command{
	Use:   "revenuecat-cli",
	Short: "RevenueCat CLI",
	Long: `revenuecat-cli is a fast, lightweight, and scriptable CLI for RevenueCat.

It provides comprehensive automation for in-app subscription management,
designed for CI/CD pipelines and developer productivity.

Design Philosophy:
  • JSON-first output for automation
  • Explicit flags over cryptic shortcuts
  • No interactive prompts
  • Clean exit codes (0=success, 1=error, 2=validation)`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip for completion and help commands
		if cmd.Name() == "completion" || cmd.Name() == "help" || cmd.Name() == "__complete" {
			return nil
		}

		// Load .rc.yaml project config if it exists
		if cwd, err := os.Getwd(); err == nil {
			if projectCfg := initcmd.FindProjectConfig(cwd); projectCfg != "" {
				viper.SetConfigFile(projectCfg)
				viper.SetConfigType("yaml")
				_ = viper.MergeInConfig()
			}
		}

		// Sync flags to cli package
		cli.SetProjectID(projectID)
		cli.SetProfile(profile)
		cli.SetTimeout(timeout)
		cli.SetDryRun(dryRun)

		// Initialize config
		if err := config.Init(cfgFile, profile); err != nil {
			return err
		}

		// Setup output formatter
		output.Setup(outputFmt, pretty, quiet)

		// Set debug mode
		if debug {
			config.SetDebug(true)
		}

		return nil
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

// SetVersionInfo sets the version information
func SetVersionInfo(version, commit, date string) {
	versionStr = version
	commitStr = commit
	dateStr = date
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default $HOME/.revenuecat-cli/config.json)")
	rootCmd.PersistentFlags().StringVarP(&projectID, "project", "p", "", "RevenueCat project ID (or RC_PROJECT env)")
	rootCmd.PersistentFlags().StringVar(&profile, "profile", "", "auth profile name (or RC_PROFILE env)")
	rootCmd.PersistentFlags().StringVarP(&outputFmt, "output", "o", "json", "output format: json, table, minimal, tsv, csv, yaml")
	rootCmd.PersistentFlags().BoolVar(&pretty, "pretty", false, "pretty-print JSON output")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress non-essential output")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "show API requests/responses")
	rootCmd.PersistentFlags().StringVar(&timeout, "timeout", "60s", "request timeout")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "preview changes without applying")

	// Bind to viper
	viper.BindPFlag("project", rootCmd.PersistentFlags().Lookup("project"))
	viper.BindPFlag("profile", rootCmd.PersistentFlags().Lookup("profile"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))

	// Environment variable bindings
	viper.BindEnv("project", "RC_PROJECT")
	viper.BindEnv("profile", "RC_PROFILE")
	viper.BindEnv("output", "RC_OUTPUT")
	viper.BindEnv("debug", "RC_DEBUG")
	viper.BindEnv("timeout", "RC_TIMEOUT")
	viper.BindEnv("api_key", "RC_API_KEY")

	// Add version command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("revenuecat-cli %s\n", versionStr)
			fmt.Printf("  commit: %s\n", commitStr)
			fmt.Printf("  built:  %s\n", dateStr)
		},
	})
}

// GetRootCmd returns the root command for adding subcommands
func GetRootCmd() *cobra.Command {
	return rootCmd
}
