package magicsetup

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/AndroidPoet/revenuecat-cli/internal/api"
	"github.com/AndroidPoet/revenuecat-cli/internal/cli"
	"github.com/AndroidPoet/revenuecat-cli/internal/output"
)

// ANSI
const (
	cReset   = "\033[0m"
	cBold    = "\033[1m"
	cDim     = "\033[2m"
	cRed     = "\033[31m"
	cGreen   = "\033[32m"
	cYellow  = "\033[33m"
	cBlue    = "\033[34m"
	cMagenta = "\033[35m"
	cCyan    = "\033[36m"
)

var reader = bufio.NewReader(os.Stdin)

// ── Templates ───────────────────────────────────────────────────────────────

type packageDef struct {
	Key         string
	DisplayName string
	Type        string
	Price       string
	Entitlement string
}

type template struct {
	Name         string
	Description  string
	Entitlements []string
	Packages     []packageDef
}

var templates = []template{
	{
		Name:         "Freemium",
		Description:  "Weekly, Monthly, Annual + Lifetime",
		Entitlements: []string{"premium"},
		Packages: []packageDef{
			{"weekly", "Weekly", "subscription", "$2.99/week", "premium"},
			{"monthly", "Monthly", "subscription", "$9.99/month", "premium"},
			{"annual", "Annual", "subscription", "$49.99/year", "premium"},
			{"lifetime", "Lifetime", "one_time", "$149.99", "premium"},
		},
	},
	{
		Name:         "Simple Paywall",
		Description:  "Monthly + Annual",
		Entitlements: []string{"premium"},
		Packages: []packageDef{
			{"monthly", "Monthly", "subscription", "$9.99/month", "premium"},
			{"annual", "Annual", "subscription", "$49.99/year", "premium"},
		},
	},
	{
		Name:         "Trial-First",
		Description:  "Free trial then Monthly / Annual",
		Entitlements: []string{"premium"},
		Packages: []packageDef{
			{"monthly", "Monthly (7-day trial)", "subscription", "$9.99/month", "premium"},
			{"annual", "Annual (14-day trial)", "subscription", "$59.99/year", "premium"},
		},
	},
	{
		Name:         "Tiered",
		Description:  "Basic + Pro tiers with Monthly / Annual",
		Entitlements: []string{"basic", "pro"},
		Packages: []packageDef{
			{"basic_monthly", "Basic Monthly", "subscription", "$4.99/month", "basic"},
			{"basic_annual", "Basic Annual", "subscription", "$29.99/year", "basic"},
			{"pro_monthly", "Pro Monthly", "subscription", "$14.99/month", "basic,pro"},
			{"pro_annual", "Pro Annual", "subscription", "$99.99/year", "basic,pro"},
		},
	},
	{
		Name:         "Consumable",
		Description:  "Credit / coin packs",
		Entitlements: []string{"credits"},
		Packages: []packageDef{
			{"small_pack", "10 Credits", "one_time", "$0.99", "credits"},
			{"medium_pack", "50 Credits", "one_time", "$3.99", "credits"},
			{"large_pack", "150 Credits", "one_time", "$9.99", "credits"},
			{"mega_pack", "500 Credits", "one_time", "$24.99", "credits"},
		},
	},
}

// ── Command ─────────────────────────────────────────────────────────────────

// MagicSetupCmd is the magicsetup command
var MagicSetupCmd = &cobra.Command{
	Use:   "magicsetup",
	Short: "One-click offerings setup for iOS, Android, or both",
	Long: `Interactive wizard that creates your full RevenueCat stack in one go:
apps, products, entitlements, offerings, and packages — all wired together.

Choose your platform (iOS / Android / both), pick a template or build custom,
set prices, and everything gets created and connected automatically.`,
	RunE: runMagicSetup,
}

// ── Helpers ─────────────────────────────────────────────────────────────────

func prompt(text, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s? %s%s [%s]: ", cCyan, cReset, text, defaultVal)
	} else {
		fmt.Printf("%s? %s%s: ", cCyan, cReset, text)
	}
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return defaultVal
	}
	return line
}

func confirm(text string) bool {
	fmt.Printf("%s? %s%s (y/N): ", cYellow, cReset, text)
	line, _ := reader.ReadString('\n')
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(line)), "y")
}

func banner() {
	fmt.Println()
	fmt.Printf("%s%s", cBold, cCyan)
	fmt.Println("  ┌─────────────────────────────────────────────┐")
	fmt.Println("  │            Magic Setup                      │")
	fmt.Println("  │                                             │")
	fmt.Println("  │  Apps  Products  Entitlements  Offerings    │")
	fmt.Println("  │  All wired together, ready to go.           │")
	fmt.Println("  └─────────────────────────────────────────────┘")
	fmt.Printf("%s\n", cReset)
}

func divider() {
	fmt.Printf("%s────────────────────────────────────────────────────────%s\n", cDim, cReset)
}

func stepHeader(n, total int, text string) {
	fmt.Printf("\n%s%s[%d/%d]%s %s%s%s\n", cBold, cCyan, n, total, cReset, cBold, text, cReset)
}

func parseTimeout() time.Duration {
	d, err := time.ParseDuration(cli.GetTimeout())
	if err != nil {
		return 60 * time.Second
	}
	return d
}

// ── Setup State ─────────────────────────────────────────────────────────────

type setupState struct {
	platform    string // "ios", "android", "both"
	appName     string
	iosBundleID string
	androidPkg  string
	baseID      string

	entitlementKeys []string
	offeringKey     string
	offeringName    string
	packages        []packageDef

	// Created IDs
	iosAppID     string
	androidAppID string
	entIDs       map[string]string
	iosProdIDs   map[string]string
	andProdIDs   map[string]string
	offeringID   string
	packageIDs   []string

	client *api.Client
	dryRun bool
}

// ── Main Flow ───────────────────────────────────────────────────────────────

func runMagicSetup(cmd *cobra.Command, args []string) error {
	banner()

	s := &setupState{
		entIDs:     make(map[string]string),
		iosProdIDs: make(map[string]string),
		andProdIDs: make(map[string]string),
		dryRun:     cli.IsDryRun(),
	}

	if s.dryRun {
		output.PrintWarning("DRY RUN — no resources will be created")
		fmt.Println()
	}

	// Check auth
	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return fmt.Errorf("not authenticated — run: rc auth login --api-key YOUR_KEY")
	}
	s.client = client
	output.PrintSuccess("Authenticated (project: %s)", client.GetProjectID())

	// ── Platform ──
	s.selectPlatform()

	// ── App details ──
	divider()
	fmt.Printf("\n%sApp Details%s\n\n", cBold, cReset)
	s.appName = prompt("App name", "My App")

	if s.platform == "ios" || s.platform == "both" {
		s.iosBundleID = prompt("iOS Bundle ID", "com.example.myapp")
	}
	if s.platform == "android" || s.platform == "both" {
		s.androidPkg = prompt("Android Package Name", "com.example.myapp")
	}

	if s.iosBundleID != "" {
		s.baseID = s.iosBundleID
	} else {
		s.baseID = s.androidPkg
	}

	// ── Template ──
	divider()
	s.selectTemplate()

	// ── Review & customize ──
	s.reviewAndCustomize()

	// ── Summary ──
	s.showSummary()

	fmt.Println()
	if !confirm("Create everything?") {
		output.PrintInfo("Aborted.")
		return nil
	}

	// ── Execute ──
	fmt.Println()
	return s.execute()
}

func (s *setupState) selectPlatform() {
	fmt.Printf("\n%sPlatform%s\n\n", cBold, cReset)
	fmt.Printf("  %s1)%s iOS only           (App Store)\n", cCyan, cReset)
	fmt.Printf("  %s2)%s Android only       (Play Store)\n", cCyan, cReset)
	fmt.Printf("  %s3)%s Both iOS + Android  (recommended)\n", cCyan, cReset)
	fmt.Println()

	choice := prompt("Select platform", "3")
	switch choice {
	case "1":
		s.platform = "ios"
	case "2":
		s.platform = "android"
	default:
		s.platform = "both"
	}
	output.PrintSuccess("Platform: %s", s.platform)
}

func (s *setupState) selectTemplate() {
	fmt.Printf("\n%sSetup Template%s\n\n", cBold, cReset)
	for i, t := range templates {
		fmt.Printf("  %s%d)%s %s%-16s%s %s\n", cCyan, i+1, cReset, cBold, t.Name, cReset, t.Description)
	}
	fmt.Printf("  %s6)%s %sCustom%s            Build your own from scratch\n", cCyan, cReset, cBold, cReset)
	fmt.Println()

	choice := prompt("Choose template (1-6)", "1")
	idx, err := strconv.Atoi(choice)

	if err == nil && idx >= 1 && idx <= len(templates) {
		tmpl := templates[idx-1]
		s.entitlementKeys = tmpl.Entitlements
		s.offeringKey = "default"
		s.offeringName = "Default Offering"

		// Apply template packages with correct product IDs
		for _, p := range tmpl.Packages {
			pkg := packageDef{
				Key:         p.Key,
				DisplayName: p.DisplayName,
				Type:        p.Type,
				Price:       p.Price,
				Entitlement: p.Entitlement,
			}
			s.packages = append(s.packages, pkg)
		}

		fmt.Println()
		output.PrintSuccess("Template loaded: %s", tmpl.Name)
	} else {
		// Custom mode
		s.customSetup()
	}
}

func (s *setupState) customSetup() {
	fmt.Println()
	ents := prompt("Entitlements (comma-separated)", "premium")
	s.entitlementKeys = splitTrim(ents)

	divider()
	fmt.Printf("\n%sOffering%s\n\n", cBold, cReset)
	s.offeringKey = prompt("Offering lookup key", "default")
	s.offeringName = prompt("Offering display name", "Default Offering")

	divider()
	fmt.Printf("\n%sPackages%s\n", cBold, cReset)
	fmt.Printf("%sDefine each subscription/purchase tier.%s\n", cDim, cReset)

	s.addPackages()
}

func (s *setupState) addPackages() {
	for {
		fmt.Println()
		key := prompt("Package lookup key (e.g. monthly, annual, lifetime)", "")
		if key == "" {
			break
		}
		name := prompt("Display name", capitalize(key))
		typ := prompt("Type (subscription / one_time)", "subscription")
		price := prompt("Price (e.g. $9.99/month)", "")
		ent := prompt("Attach to entitlements (comma-separated)", s.entitlementKeys[0])

		s.packages = append(s.packages, packageDef{
			Key:         key,
			DisplayName: name,
			Type:        typ,
			Price:       price,
			Entitlement: ent,
		})

		output.PrintSuccess("Added: %s %s%s%s", name, cGreen, price, cReset)
		fmt.Println()
		if !confirm("Add another package?") {
			break
		}
	}
}

func (s *setupState) reviewAndCustomize() {
	divider()
	fmt.Printf("\n%sReview & Customize%s\n\n", cBold, cReset)

	fmt.Printf("  %sEntitlements:%s %s\n", cBold, cReset, strings.Join(s.entitlementKeys, ", "))
	if confirm("Change entitlements?") {
		ents := prompt("Entitlements (comma-separated)", strings.Join(s.entitlementKeys, ","))
		s.entitlementKeys = splitTrim(ents)
	}

	fmt.Println()
	fmt.Printf("  %sOffering:%s %s (%s)\n", cBold, cReset, s.offeringKey, s.offeringName)
	if confirm("Change offering?") {
		s.offeringKey = prompt("Offering lookup key", s.offeringKey)
		s.offeringName = prompt("Offering display name", s.offeringName)
	}

	fmt.Println()
	fmt.Printf("  %sPackages:%s\n", cBold, cReset)
	for i, p := range s.packages {
		iosID := s.baseID + "." + p.Key
		androidID := s.baseID + "." + p.Key
		fmt.Printf("    %s%d.%s %-20s %s%-14s%s %s%s%s\n",
			cCyan, i+1, cReset, p.DisplayName, cDim, p.Type, cReset, cGreen, p.Price, cReset)
		if s.platform == "ios" || s.platform == "both" {
			fmt.Printf("       iOS: %s\n", iosID)
		}
		if s.platform == "android" || s.platform == "both" {
			fmt.Printf("       Android: %s\n", androidID)
		}
	}

	fmt.Println()
	if confirm("Edit any package?") {
		s.editPackages()
	}

	if confirm("Add more packages?") {
		s.addPackages()
	}
}

func (s *setupState) editPackages() {
	for {
		idx := prompt(fmt.Sprintf("Package number to edit (1-%d, or 'done')", len(s.packages)), "done")
		if idx == "done" {
			break
		}
		i, err := strconv.Atoi(idx)
		if err != nil || i < 1 || i > len(s.packages) {
			output.PrintWarning("Invalid number")
			continue
		}
		i-- // zero-based

		p := &s.packages[i]
		fmt.Printf("\n  %sEditing: %s%s\n\n", cBold, p.DisplayName, cReset)
		p.Key = prompt("  Lookup key", p.Key)
		p.DisplayName = prompt("  Display name", p.DisplayName)
		p.Type = prompt("  Type (subscription/one_time)", p.Type)
		p.Price = prompt("  Price", p.Price)
		p.Entitlement = prompt("  Entitlements (comma-separated)", p.Entitlement)

		output.PrintSuccess("Updated %s", p.DisplayName)
		fmt.Println()
	}
}

func (s *setupState) showSummary() {
	divider()
	fmt.Printf("\n%s%sSetup Summary%s\n\n", cBold, cMagenta, cReset)

	fmt.Printf("  %sPlatform:%s      %s\n", cBold, cReset, s.platform)
	fmt.Printf("  %sApp:%s           %s\n", cBold, cReset, s.appName)
	if s.iosBundleID != "" {
		fmt.Printf("  %siOS:%s           %s\n", cBold, cReset, s.iosBundleID)
	}
	if s.androidPkg != "" {
		fmt.Printf("  %sAndroid:%s       %s\n", cBold, cReset, s.androidPkg)
	}
	fmt.Printf("  %sEntitlements:%s  %s\n", cBold, cReset, strings.Join(s.entitlementKeys, ", "))
	fmt.Printf("  %sOffering:%s      %s (%s)\n", cBold, cReset, s.offeringKey, s.offeringName)
	fmt.Println()

	// Table header
	fmt.Printf("  %s%-22s %-14s %-14s %s%s\n", cBold, "Package", "Type", "Price", "Product IDs", cReset)
	fmt.Printf("  %s%-22s %-14s %-14s %s%s\n", cDim, "──────────────────────", "──────────────", "──────────────", "─────────────────────────────────────", cReset)

	totalProducts := 0
	for _, p := range s.packages {
		ids := ""
		if s.platform == "ios" || s.platform == "both" {
			ids += "iOS:" + s.baseID + "." + p.Key
			totalProducts++
		}
		if s.platform == "android" || s.platform == "both" {
			if ids != "" {
				ids += " | "
			}
			ids += "Android:" + s.baseID + "." + p.Key
			totalProducts++
		}
		fmt.Printf("  %-22s %-14s %s%-14s%s %s%s%s\n",
			p.DisplayName, p.Type, cGreen, p.Price, cReset, cDim, ids, cReset)
	}

	apps := 0
	if s.iosBundleID != "" {
		apps++
	}
	if s.androidPkg != "" {
		apps++
	}

	fmt.Printf("\n  %sWill create: %d app(s), %d product(s), %d entitlement(s), 1 offering, %d package(s)%s\n",
		cDim, apps, totalProducts, len(s.entitlementKeys), len(s.packages), cReset)
	divider()
}

// ── Execution ───────────────────────────────────────────────────────────────

func (s *setupState) execute() error {
	totalSteps := 6
	ctx, cancel := s.client.Context()
	defer cancel()
	projectID := s.client.GetProjectID()

	// ── Step 1: Create Apps ──
	stepHeader(1, totalSteps, "Creating app(s)")

	if s.platform == "ios" || s.platform == "both" {
		body := map[string]interface{}{
			"name":      s.appName + " iOS",
			"type":      "app_store",
			"bundle_id": s.iosBundleID,
		}
		var result map[string]interface{}
		path := fmt.Sprintf("/projects/%s/apps", projectID)

		if s.dryRun {
			fmt.Printf("  %s(dry-run)%s POST %s\n", cDim, cReset, path)
			s.iosAppID = "app_ios_dry_run"
		} else {
			err := s.client.Post(ctx, path, body, &result)
			if err != nil {
				output.PrintWarning("iOS app: %v", err)
				s.iosAppID = prompt("Enter existing iOS app ID (or skip)", "")
			} else {
				s.iosAppID = getID(result)
			}
		}
		if s.iosAppID != "" {
			output.PrintSuccess("iOS app: %s", s.iosAppID)
		}
	}

	if s.platform == "android" || s.platform == "both" {
		body := map[string]interface{}{
			"name":         s.appName + " Android",
			"type":         "play_store",
			"package_name": s.androidPkg,
		}
		var result map[string]interface{}
		path := fmt.Sprintf("/projects/%s/apps", projectID)

		if s.dryRun {
			fmt.Printf("  %s(dry-run)%s POST %s\n", cDim, cReset, path)
			s.androidAppID = "app_android_dry_run"
		} else {
			err := s.client.Post(ctx, path, body, &result)
			if err != nil {
				output.PrintWarning("Android app: %v", err)
				s.androidAppID = prompt("Enter existing Android app ID (or skip)", "")
			} else {
				s.androidAppID = getID(result)
			}
		}
		if s.androidAppID != "" {
			output.PrintSuccess("Android app: %s", s.androidAppID)
		}
	}

	// ── Step 2: Create Products ──
	stepHeader(2, totalSteps, "Creating products")

	for _, p := range s.packages {
		iosProductID := s.baseID + "." + p.Key
		androidProductID := s.baseID + "." + p.Key

		if (s.platform == "ios" || s.platform == "both") && s.iosAppID != "" {
			body := map[string]interface{}{
				"store_identifier": iosProductID,
				"type":             p.Type,
				"app_id":           s.iosAppID,
			}
			var result map[string]interface{}
			path := fmt.Sprintf("/projects/%s/products", projectID)

			if s.dryRun {
				fmt.Printf("  %s(dry-run)%s iOS product [%s] %s%s%s\n", cDim, cReset, p.Key, cGreen, p.Price, cReset)
				s.iosProdIDs[p.Key] = "prod_ios_" + p.Key + "_dry"
			} else {
				err := s.client.Post(ctx, path, body, &result)
				if err != nil {
					output.PrintWarning("  iOS [%s]: %v", p.Key, err)
				} else {
					s.iosProdIDs[p.Key] = getID(result)
					output.PrintSuccess("  iOS [%s]: %s  %s%s%s", p.Key, s.iosProdIDs[p.Key], cGreen, p.Price, cReset)
				}
			}
		}

		if (s.platform == "android" || s.platform == "both") && s.androidAppID != "" {
			body := map[string]interface{}{
				"store_identifier": androidProductID,
				"type":             p.Type,
				"app_id":           s.androidAppID,
			}
			var result map[string]interface{}
			path := fmt.Sprintf("/projects/%s/products", projectID)

			if s.dryRun {
				fmt.Printf("  %s(dry-run)%s Android product [%s] %s%s%s\n", cDim, cReset, p.Key, cGreen, p.Price, cReset)
				s.andProdIDs[p.Key] = "prod_and_" + p.Key + "_dry"
			} else {
				err := s.client.Post(ctx, path, body, &result)
				if err != nil {
					output.PrintWarning("  Android [%s]: %v", p.Key, err)
				} else {
					s.andProdIDs[p.Key] = getID(result)
					output.PrintSuccess("  Android [%s]: %s  %s%s%s", p.Key, s.andProdIDs[p.Key], cGreen, p.Price, cReset)
				}
			}
		}
	}

	// ── Step 3: Create Entitlements ──
	stepHeader(3, totalSteps, "Creating entitlements")

	for _, ent := range s.entitlementKeys {
		body := map[string]interface{}{
			"lookup_key":   ent,
			"display_name": capitalize(ent),
		}
		var result map[string]interface{}
		path := fmt.Sprintf("/projects/%s/entitlements", projectID)

		if s.dryRun {
			fmt.Printf("  %s(dry-run)%s Entitlement [%s]\n", cDim, cReset, ent)
			s.entIDs[ent] = "entl_" + ent + "_dry"
		} else {
			err := s.client.Post(ctx, path, body, &result)
			if err != nil {
				output.PrintWarning("  [%s]: %v", ent, err)
			} else {
				s.entIDs[ent] = getID(result)
				output.PrintSuccess("  [%s]: %s", ent, s.entIDs[ent])
			}
		}
	}

	// ── Step 4: Attach Products to Entitlements ──
	stepHeader(4, totalSteps, "Attaching products to entitlements")

	for _, p := range s.packages {
		ents := splitTrim(p.Entitlement)
		for _, ent := range ents {
			entID, ok := s.entIDs[ent]
			if !ok {
				output.PrintWarning("  Skipping: entitlement '%s' not found", ent)
				continue
			}

			var prodIDs []string
			if id, ok := s.iosProdIDs[p.Key]; ok && id != "" {
				prodIDs = append(prodIDs, id)
			}
			if id, ok := s.andProdIDs[p.Key]; ok && id != "" {
				prodIDs = append(prodIDs, id)
			}

			if len(prodIDs) == 0 {
				continue
			}

			body := map[string]interface{}{
				"product_ids": prodIDs,
			}
			path := fmt.Sprintf("/projects/%s/entitlements/%s/actions/attach_products", projectID, entID)

			if s.dryRun {
				fmt.Printf("  %s(dry-run)%s [%s] -> [%s]\n", cDim, cReset, p.Key, ent)
			} else {
				err := s.client.Post(ctx, path, body, nil)
				if err != nil {
					output.PrintWarning("  [%s] -> [%s]: %v", p.Key, ent, err)
				} else {
					output.PrintSuccess("  [%s] -> [%s]", p.Key, ent)
				}
			}
		}
	}

	// ── Step 5: Create Offering + Packages ──
	stepHeader(5, totalSteps, "Creating offering and packages")

	// Create offering
	offeringBody := map[string]interface{}{
		"lookup_key":   s.offeringKey,
		"display_name": s.offeringName,
	}
	var offeringResult map[string]interface{}
	offeringPath := fmt.Sprintf("/projects/%s/offerings", projectID)

	if s.dryRun {
		fmt.Printf("  %s(dry-run)%s Offering [%s]\n", cDim, cReset, s.offeringKey)
		s.offeringID = "ofrngs_dry_run"
	} else {
		err := s.client.Post(ctx, offeringPath, offeringBody, &offeringResult)
		if err != nil {
			output.PrintWarning("Offering: %v", err)
			s.offeringID = prompt("Enter existing offering ID", "")
		} else {
			s.offeringID = getID(offeringResult)
		}
	}

	if s.offeringID != "" {
		output.PrintSuccess("Offering: %s (%s)", s.offeringID, s.offeringKey)

		// Set as current
		updateBody := map[string]interface{}{"is_current": true}
		updatePath := fmt.Sprintf("/projects/%s/offerings/%s", projectID, s.offeringID)
		if !s.dryRun {
			_ = s.client.Patch(ctx, updatePath, updateBody, nil)
		}
		output.PrintSuccess("Set as current offering")
	}

	// Create packages
	for _, p := range s.packages {
		pkgBody := map[string]interface{}{
			"lookup_key":   p.Key,
			"display_name": p.DisplayName,
		}
		var pkgResult map[string]interface{}
		pkgPath := fmt.Sprintf("/projects/%s/offerings/%s/packages", projectID, s.offeringID)

		var pkgID string
		if s.dryRun {
			fmt.Printf("  %s(dry-run)%s Package [%s]\n", cDim, cReset, p.Key)
			pkgID = "pkg_" + p.Key + "_dry"
		} else {
			err := s.client.Post(ctx, pkgPath, pkgBody, &pkgResult)
			if err != nil {
				output.PrintWarning("  Package [%s]: %v", p.Key, err)
				continue
			}
			pkgID = getID(pkgResult)
			output.PrintSuccess("  Package [%s]: %s", p.Key, pkgID)
		}

		s.packageIDs = append(s.packageIDs, pkgID)

		// Attach products to package
		var attachIDs []string
		if id, ok := s.iosProdIDs[p.Key]; ok && id != "" {
			attachIDs = append(attachIDs, id)
		}
		if id, ok := s.andProdIDs[p.Key]; ok && id != "" {
			attachIDs = append(attachIDs, id)
		}

		if len(attachIDs) > 0 && pkgID != "" {
			attachBody := map[string]interface{}{
				"product_ids": attachIDs,
			}
			attachPath := fmt.Sprintf("/projects/%s/offerings/%s/packages/%s/actions/attach_products",
				projectID, s.offeringID, pkgID)

			if !s.dryRun {
				err := s.client.Post(ctx, attachPath, attachBody, nil)
				if err != nil {
					output.PrintWarning("  Attach [%s]: %v", p.Key, err)
				} else {
					output.PrintSuccess("  Attached products to [%s]", p.Key)
				}
			}
		}
	}

	// ── Step 6: Done ──
	stepHeader(6, totalSteps, "Complete")
	fmt.Println()
	divider()
	fmt.Printf("\n%s%sSetup Complete!%s\n\n", cBold, cGreen, cReset)

	totalProducts := len(s.iosProdIDs) + len(s.andProdIDs)
	totalApps := 0
	if s.iosAppID != "" {
		totalApps++
	}
	if s.androidAppID != "" {
		totalApps++
	}

	fmt.Printf("  %sCreated:%s\n", cBold, cReset)
	fmt.Printf("    Apps:          %d\n", totalApps)
	if s.iosAppID != "" {
		fmt.Printf("      iOS:         %s\n", s.iosAppID)
	}
	if s.androidAppID != "" {
		fmt.Printf("      Android:     %s\n", s.androidAppID)
	}
	fmt.Printf("    Products:      %d\n", totalProducts)
	fmt.Printf("    Entitlements:  %d\n", len(s.entIDs))
	fmt.Printf("    Offering:      %s\n", s.offeringID)
	fmt.Printf("    Packages:      %d\n", len(s.packageIDs))

	fmt.Printf("\n  %sPrice Reference:%s\n", cBold, cReset)
	for _, p := range s.packages {
		fmt.Printf("    %-22s %s\n", p.DisplayName, p.Price)
	}

	fmt.Printf("\n  %sVerify:%s  rc status\n", cDim, cReset)
	fmt.Printf("  %sView:%s    rc offerings get --offering-id %s\n", cDim, cReset, s.offeringID)
	fmt.Println()
	divider()

	return nil
}

// ── Utilities ───────────────────────────────────────────────────────────────

func splitTrim(s string) []string {
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func getID(m map[string]interface{}) string {
	if id, ok := m["id"]; ok {
		return fmt.Sprintf("%v", id)
	}
	return ""
}
