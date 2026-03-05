package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/AndroidPoet/revenuecat-cli/internal/api"
	"github.com/AndroidPoet/revenuecat-cli/internal/cli"
	"github.com/AndroidPoet/revenuecat-cli/internal/output"
)

var (
	reportFile   string
	reportFormat string
)

// ReportCmd exports a full project report
var ReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a full project report",
	Long:  `Export a comprehensive project report with apps, products, entitlements, offerings, packages, and metrics into HTML, JSON, or YAML.`,
	RunE:  runReport,
}

func init() {
	ReportCmd.Flags().StringVar(&reportFile, "file", "", "output file path (default: rc-report.{format})")
	ReportCmd.Flags().StringVar(&reportFormat, "format", "html", "output format: html, json, yaml, pdf")
}

// --- Data types ---

type ProjectReport struct {
	GeneratedAt  string          `json:"generated_at" yaml:"generated_at"`
	ProjectID    string          `json:"project_id" yaml:"project_id"`
	Summary      ReportSummary   `json:"summary" yaml:"summary"`
	Apps         []AppData       `json:"apps" yaml:"apps"`
	Products     []ProductData   `json:"products" yaml:"products"`
	Entitlements []EntitlementData `json:"entitlements" yaml:"entitlements"`
	Offerings    []OfferingData  `json:"offerings" yaml:"offerings"`
}

type ReportSummary struct {
	TotalApps         int     `json:"total_apps" yaml:"total_apps"`
	TotalProducts     int     `json:"total_products" yaml:"total_products"`
	TotalEntitlements int     `json:"total_entitlements" yaml:"total_entitlements"`
	TotalOfferings    int     `json:"total_offerings" yaml:"total_offerings"`
	TotalPackages     int     `json:"total_packages" yaml:"total_packages"`
	ActiveSubscribers int     `json:"active_subscribers" yaml:"active_subscribers"`
	ActiveTrials      int     `json:"active_trials" yaml:"active_trials"`
	MRR               float64 `json:"mrr" yaml:"mrr"`
	Revenue           float64 `json:"revenue" yaml:"revenue"`
}

type AppData struct {
	ID   string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
	Type string `json:"type" yaml:"type"`
}

type ProductData struct {
	ID              string `json:"id" yaml:"id"`
	StoreIdentifier string `json:"store_identifier" yaml:"store_identifier"`
	Type            string `json:"type" yaml:"type"`
	AppID           string `json:"app_id" yaml:"app_id"`
}

type EntitlementData struct {
	ID          string `json:"id" yaml:"id"`
	LookupKey   string `json:"lookup_key" yaml:"lookup_key"`
	DisplayName string `json:"display_name" yaml:"display_name"`
}

type OfferingData struct {
	ID          string        `json:"id" yaml:"id"`
	LookupKey   string        `json:"lookup_key" yaml:"lookup_key"`
	DisplayName string        `json:"display_name" yaml:"display_name"`
	IsCurrent   bool          `json:"is_current" yaml:"is_current"`
	Packages    []PackageData `json:"packages" yaml:"packages"`
}

type PackageData struct {
	ID          string `json:"id" yaml:"id"`
	LookupKey   string `json:"lookup_key" yaml:"lookup_key"`
	DisplayName string `json:"display_name" yaml:"display_name"`
}

func parseTimeout() time.Duration {
	d, err := time.ParseDuration(cli.GetTimeout())
	if err != nil {
		return 60 * time.Second
	}
	return d
}

func runReport(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}

	ctx, cancel := client.Context()
	defer cancel()

	pid := client.GetProjectID()

	report := ProjectReport{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		ProjectID:   pid,
	}

	// 1. Fetch apps
	var appMaps []map[string]interface{}
	appPath := fmt.Sprintf("/projects/%s/apps", pid)
	if err := client.ListAll(ctx, appPath, 100, func(raw json.RawMessage) error {
		var page []map[string]interface{}
		if err := json.Unmarshal(raw, &page); err != nil {
			return err
		}
		appMaps = append(appMaps, page...)
		return nil
	}); err != nil {
		output.PrintWarning("Failed to fetch apps: %s", err)
	} else {
		for _, m := range appMaps {
			report.Apps = append(report.Apps, AppData{
				ID:   str(m, "id"),
				Name: str(m, "name"),
				Type: str(m, "type"),
			})
		}
	}

	// 2. Fetch products
	var prodMaps []map[string]interface{}
	prodPath := fmt.Sprintf("/projects/%s/products", pid)
	if err := client.ListAll(ctx, prodPath, 100, func(raw json.RawMessage) error {
		var page []map[string]interface{}
		if err := json.Unmarshal(raw, &page); err != nil {
			return err
		}
		prodMaps = append(prodMaps, page...)
		return nil
	}); err != nil {
		output.PrintWarning("Failed to fetch products: %s", err)
	} else {
		for _, m := range prodMaps {
			report.Products = append(report.Products, ProductData{
				ID:              str(m, "id"),
				StoreIdentifier: str(m, "store_identifier"),
				Type:            str(m, "type"),
				AppID:           str(m, "app_id"),
			})
		}
	}

	// 3. Fetch entitlements
	var entMaps []map[string]interface{}
	entPath := fmt.Sprintf("/projects/%s/entitlements", pid)
	if err := client.ListAll(ctx, entPath, 100, func(raw json.RawMessage) error {
		var page []map[string]interface{}
		if err := json.Unmarshal(raw, &page); err != nil {
			return err
		}
		entMaps = append(entMaps, page...)
		return nil
	}); err != nil {
		output.PrintWarning("Failed to fetch entitlements: %s", err)
	} else {
		for _, m := range entMaps {
			report.Entitlements = append(report.Entitlements, EntitlementData{
				ID:          str(m, "id"),
				LookupKey:   str(m, "lookup_key"),
				DisplayName: str(m, "display_name"),
			})
		}
	}

	// 4. Fetch offerings
	var offMaps []map[string]interface{}
	offPath := fmt.Sprintf("/projects/%s/offerings", pid)
	if err := client.ListAll(ctx, offPath, 100, func(raw json.RawMessage) error {
		var page []map[string]interface{}
		if err := json.Unmarshal(raw, &page); err != nil {
			return err
		}
		offMaps = append(offMaps, page...)
		return nil
	}); err != nil {
		output.PrintWarning("Failed to fetch offerings: %s", err)
	} else {
		for _, m := range offMaps {
			off := OfferingData{
				ID:          str(m, "id"),
				LookupKey:   str(m, "lookup_key"),
				DisplayName: str(m, "display_name"),
				IsCurrent:   boolVal(m, "is_current"),
			}

			// 4a. Fetch packages for this offering
			var pkgMaps []map[string]interface{}
			pkgPath := fmt.Sprintf("/projects/%s/offerings/%s/packages", pid, off.ID)
			if err := client.ListAll(ctx, pkgPath, 100, func(raw json.RawMessage) error {
				var page []map[string]interface{}
				if err := json.Unmarshal(raw, &page); err != nil {
					return err
				}
				pkgMaps = append(pkgMaps, page...)
				return nil
			}); err != nil {
				output.PrintWarning("Failed to fetch packages for offering %s: %s", off.LookupKey, err)
			} else {
				for _, pm := range pkgMaps {
					off.Packages = append(off.Packages, PackageData{
						ID:          str(pm, "id"),
						LookupKey:   str(pm, "lookup_key"),
						DisplayName: str(pm, "display_name"),
					})
				}
			}

			report.Offerings = append(report.Offerings, off)
		}
	}

	// 5. Fetch metrics
	var metrics struct {
		ActiveSubscribers int     `json:"active_subscribers"`
		ActiveTrials      int     `json:"active_trials"`
		MRR               float64 `json:"mrr"`
		Revenue           float64 `json:"revenue"`
	}
	metricsPath := fmt.Sprintf("/projects/%s/metrics/overview", pid)
	if err := client.Get(ctx, metricsPath, &metrics); err != nil {
		output.PrintWarning("Failed to fetch metrics: %s (using zeros)", err)
	} else {
		report.Summary.ActiveSubscribers = metrics.ActiveSubscribers
		report.Summary.ActiveTrials = metrics.ActiveTrials
		report.Summary.MRR = metrics.MRR
		report.Summary.Revenue = metrics.Revenue
	}

	// Build summary counts
	totalPkgs := 0
	for _, off := range report.Offerings {
		totalPkgs += len(off.Packages)
	}
	report.Summary.TotalApps = len(report.Apps)
	report.Summary.TotalProducts = len(report.Products)
	report.Summary.TotalEntitlements = len(report.Entitlements)
	report.Summary.TotalOfferings = len(report.Offerings)
	report.Summary.TotalPackages = totalPkgs

	// Set default filename based on format
	if reportFile == "" {
		switch reportFormat {
		case "pdf":
			reportFile = "rc-report.pdf"
		case "json":
			reportFile = "rc-report.json"
		case "yaml":
			reportFile = "rc-report.yaml"
		default:
			reportFile = "rc-report.html"
		}
	}

	// Write output
	switch reportFormat {
	case "json":
		data, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling JSON: %w", err)
		}
		if err := os.WriteFile(reportFile, data, 0644); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}

	case "yaml":
		data, err := yaml.Marshal(report)
		if err != nil {
			return fmt.Errorf("marshaling YAML: %w", err)
		}
		if err := os.WriteFile(reportFile, data, 0644); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}

	case "html":
		html, err := renderHTML(report)
		if err != nil {
			return fmt.Errorf("rendering HTML: %w", err)
		}
		if err := os.WriteFile(reportFile, []byte(html), 0644); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}

	case "pdf":
		html, err := renderHTML(report)
		if err != nil {
			return fmt.Errorf("rendering HTML: %w", err)
		}
		if err := htmlToPDF(html, reportFile); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unsupported format: %s (use html, json, yaml, or pdf)", reportFormat)
	}

	output.PrintSuccess("Report saved to %s", reportFile)
	if reportFormat == "html" {
		output.PrintInfo("Open in browser to view or print to PDF")
	}

	return nil
}

// --- helpers ---

func str(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func boolVal(m map[string]interface{}, key string) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// --- HTML template ---

type htmlData struct {
	ProjectReport
	MRRFormatted     string
	RevenueFormatted string
}

func renderHTML(report ProjectReport) (string, error) {
	data := htmlData{
		ProjectReport:    report,
		MRRFormatted:     fmt.Sprintf("$%.2f", report.Summary.MRR),
		RevenueFormatted: fmt.Sprintf("$%.2f", report.Summary.Revenue),
	}

	tmpl, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}

// htmlToPDF converts HTML to PDF using Chrome/Chromium headless
func htmlToPDF(html, pdfPath string) error {
	chromePath := findChrome()
	if chromePath == "" {
		return fmt.Errorf("PDF export requires Chrome or Chromium.\nInstall Chrome, or use --format html and print to PDF from your browser")
	}

	// Write HTML to temp file
	tmpFile, err := os.CreateTemp("", "rc-report-*.html")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(html); err != nil {
		tmpFile.Close()
		return fmt.Errorf("writing temp file: %w", err)
	}
	tmpFile.Close()

	// Convert absolute path for the PDF output
	absPDF, err := filepath.Abs(pdfPath)
	if err != nil {
		absPDF = pdfPath
	}

	// Run Chrome headless
	cmd := exec.Command(chromePath,
		"--headless",
		"--disable-gpu",
		"--no-sandbox",
		"--print-to-pdf="+absPDF,
		"--print-to-pdf-no-header",
		tmpFile.Name(),
	)
	cmd.Stderr = nil
	cmd.Stdout = nil

	output.PrintInfo("Generating PDF...")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Chrome PDF conversion failed: %w\nTry: rc report --format html, then print to PDF from browser", err)
	}

	return nil
}

// findChrome locates Chrome/Chromium binary on the system
func findChrome() string {
	switch runtime.GOOS {
	case "darwin":
		paths := []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
			"/Applications/Brave Browser.app/Contents/MacOS/Brave Browser",
			"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
		}
		for _, p := range paths {
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
	case "linux":
		names := []string{"google-chrome", "google-chrome-stable", "chromium", "chromium-browser"}
		for _, name := range names {
			if p, err := exec.LookPath(name); err == nil {
				return p
			}
		}
	case "windows":
		paths := []string{
			filepath.Join(os.Getenv("PROGRAMFILES"), "Google", "Chrome", "Application", "chrome.exe"),
			filepath.Join(os.Getenv("PROGRAMFILES(X86)"), "Google", "Chrome", "Application", "chrome.exe"),
			filepath.Join(os.Getenv("LOCALAPPDATA"), "Google", "Chrome", "Application", "chrome.exe"),
		}
		for _, p := range paths {
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
	}
	return ""
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>RevenueCat Project Report</title>
<style>
  *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
    color: #1a1a2e;
    background: #fff;
    line-height: 1.6;
    padding: 2rem 1rem;
  }
  .container { max-width: 900px; margin: 0 auto; }
  header {
    border-bottom: 3px solid #4B48F2;
    padding-bottom: 1rem;
    margin-bottom: 2rem;
  }
  header h1 { font-size: 1.8em; color: #4B48F2; }
  header .meta { color: #666; font-size: 0.9em; margin-top: 0.25rem; }
  h2 {
    font-size: 1.3em;
    border-bottom: 2px solid #e0e0e0;
    padding-bottom: 0.4rem;
    margin: 2rem 0 1rem;
  }
  .metrics-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 1rem;
    margin-bottom: 1rem;
  }
  .metric-card {
    background: #f5f5f7;
    border-radius: 8px;
    padding: 1rem;
    text-align: center;
  }
  .metric-card .value {
    font-size: 1.6em;
    font-weight: 700;
    color: #4B48F2;
  }
  .metric-card .label {
    font-size: 0.8em;
    color: #666;
    margin-top: 0.25rem;
  }
  table {
    width: 100%;
    border-collapse: collapse;
    margin-bottom: 1.5rem;
    font-size: 0.9em;
  }
  th {
    background: #4B48F2;
    color: #fff;
    text-align: left;
    padding: 0.6rem 0.8rem;
    font-weight: 600;
  }
  td { padding: 0.5rem 0.8rem; border-bottom: 1px solid #eee; }
  tr:nth-child(even) td { background: #fafafa; }
  .offering-block { margin-bottom: 1.5rem; }
  .offering-block h3 { font-size: 1.05em; margin-bottom: 0.5rem; }
  .offering-block .badge {
    display: inline-block;
    background: #E8514A;
    color: #fff;
    font-size: 0.7em;
    padding: 0.15rem 0.5rem;
    border-radius: 4px;
    margin-left: 0.5rem;
    vertical-align: middle;
  }
  .empty { color: #999; font-style: italic; margin-bottom: 1.5rem; }
  footer {
    margin-top: 3rem;
    padding-top: 1rem;
    border-top: 1px solid #e0e0e0;
    font-size: 0.8em;
    color: #999;
    text-align: center;
  }
  footer a { color: #4B48F2; text-decoration: none; }
  footer a:hover { text-decoration: underline; }
  @media print {
    body { padding: 0; }
    .metric-card { break-inside: avoid; }
    table { break-inside: auto; }
    tr { break-inside: avoid; }
    footer { break-before: avoid; }
  }
  @media (max-width: 600px) {
    .metrics-grid { grid-template-columns: repeat(2, 1fr); }
  }
</style>
</head>
<body>
<div class="container">
  <header>
    <h1>RevenueCat Project Report</h1>
    <div class="meta">Generated: {{.GeneratedAt}} &middot; Project: {{.ProjectID}}</div>
  </header>

  <h2>Summary</h2>
  <div class="metrics-grid">
    <div class="metric-card"><div class="value">{{.Summary.ActiveSubscribers}}</div><div class="label">Active Subscribers</div></div>
    <div class="metric-card"><div class="value">{{.Summary.ActiveTrials}}</div><div class="label">Active Trials</div></div>
    <div class="metric-card"><div class="value">{{.MRRFormatted}}</div><div class="label">MRR</div></div>
    <div class="metric-card"><div class="value">{{.RevenueFormatted}}</div><div class="label">Revenue</div></div>
  </div>
  <div class="metrics-grid">
    <div class="metric-card"><div class="value">{{.Summary.TotalApps}}</div><div class="label">Apps</div></div>
    <div class="metric-card"><div class="value">{{.Summary.TotalProducts}}</div><div class="label">Products</div></div>
    <div class="metric-card"><div class="value">{{.Summary.TotalEntitlements}}</div><div class="label">Entitlements</div></div>
    <div class="metric-card"><div class="value">{{.Summary.TotalOfferings}}</div><div class="label">Offerings</div></div>
  </div>

  <h2>Apps</h2>
  {{if .Apps}}
  <table>
    <tr><th>ID</th><th>Name</th><th>Type</th></tr>
    {{range .Apps}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Type}}</td></tr>
    {{end}}
  </table>
  {{else}}<p class="empty">No apps found.</p>{{end}}

  <h2>Products</h2>
  {{if .Products}}
  <table>
    <tr><th>ID</th><th>Store Identifier</th><th>Type</th><th>App ID</th></tr>
    {{range .Products}}<tr><td>{{.ID}}</td><td>{{.StoreIdentifier}}</td><td>{{.Type}}</td><td>{{.AppID}}</td></tr>
    {{end}}
  </table>
  {{else}}<p class="empty">No products found.</p>{{end}}

  <h2>Entitlements</h2>
  {{if .Entitlements}}
  <table>
    <tr><th>ID</th><th>Lookup Key</th><th>Display Name</th></tr>
    {{range .Entitlements}}<tr><td>{{.ID}}</td><td>{{.LookupKey}}</td><td>{{.DisplayName}}</td></tr>
    {{end}}
  </table>
  {{else}}<p class="empty">No entitlements found.</p>{{end}}

  <h2>Offerings</h2>
  {{if .Offerings}}
  {{range .Offerings}}
  <div class="offering-block">
    <h3>{{.DisplayName}} <code>({{.LookupKey}})</code>{{if .IsCurrent}}<span class="badge">CURRENT</span>{{end}}</h3>
    {{if .Packages}}
    <table>
      <tr><th>ID</th><th>Lookup Key</th><th>Display Name</th></tr>
      {{range .Packages}}<tr><td>{{.ID}}</td><td>{{.LookupKey}}</td><td>{{.DisplayName}}</td></tr>
      {{end}}
    </table>
    {{else}}<p class="empty">No packages in this offering.</p>{{end}}
  </div>
  {{end}}
  {{else}}<p class="empty">No offerings found.</p>{{end}}

  <footer>
    Generated by <a href="https://github.com/AndroidPoet/revenuecat-cli">revenuecat-cli</a>
  </footer>
</div>
</body>
</html>`
