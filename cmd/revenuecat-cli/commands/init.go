package commands

import (
	"github.com/AndroidPoet/revenuecat-cli/cmd/revenuecat-cli/commands/apps"
	"github.com/AndroidPoet/revenuecat-cli/cmd/revenuecat-cli/commands/auth"
	"github.com/AndroidPoet/revenuecat-cli/cmd/revenuecat-cli/commands/completion"
	"github.com/AndroidPoet/revenuecat-cli/cmd/revenuecat-cli/commands/customers"
	"github.com/AndroidPoet/revenuecat-cli/cmd/revenuecat-cli/commands/doctor"
	"github.com/AndroidPoet/revenuecat-cli/cmd/revenuecat-cli/commands/entitlements"
	"github.com/AndroidPoet/revenuecat-cli/cmd/revenuecat-cli/commands/initcmd"
	"github.com/AndroidPoet/revenuecat-cli/cmd/revenuecat-cli/commands/offerings"
	"github.com/AndroidPoet/revenuecat-cli/cmd/revenuecat-cli/commands/packages"
	"github.com/AndroidPoet/revenuecat-cli/cmd/revenuecat-cli/commands/paywalls"
	"github.com/AndroidPoet/revenuecat-cli/cmd/revenuecat-cli/commands/products"
)

func init() {
	rootCmd.AddCommand(auth.AuthCmd)
	rootCmd.AddCommand(apps.AppsCmd)
	rootCmd.AddCommand(products.ProductsCmd)
	rootCmd.AddCommand(entitlements.EntitlementsCmd)
	rootCmd.AddCommand(offerings.OfferingsCmd)
	rootCmd.AddCommand(packages.PackagesCmd)
	rootCmd.AddCommand(customers.CustomersCmd)
	rootCmd.AddCommand(paywalls.PaywallsCmd)
	rootCmd.AddCommand(completion.CompletionCmd)
	rootCmd.AddCommand(doctor.DoctorCmd)
	rootCmd.AddCommand(initcmd.InitCmd)
}
