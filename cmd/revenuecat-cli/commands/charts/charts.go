package charts

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/AndroidPoet/revenuecat-cli/internal/api"
	"github.com/AndroidPoet/revenuecat-cli/internal/cli"
	"github.com/AndroidPoet/revenuecat-cli/internal/output"
)

func parseTimeout() time.Duration {
	t := cli.GetTimeout()
	d, err := time.ParseDuration(t)
	if err != nil {
		d = 60 * time.Second
	}
	return d
}

// --- Chart metadata ---

type chartInfo struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

var chartCatalog = []chartInfo{
	{"revenue", "Revenue", "Total revenue generated from subscriptions and purchases", "revenue"},
	{"mrr", "Monthly Recurring Revenue", "Normalized monthly value of active paid subscriptions", "revenue"},
	{"arr", "Annual Recurring Revenue", "Annualized value of active paid subscriptions", "revenue"},
	{"mrr_movement", "MRR Movement", "Breakdown of MRR changes: new, expansion, contraction, churn", "revenue"},
	{"actives", "Active Subscriptions", "Total number of active paid subscriptions", "subscriptions"},
	{"actives_movement", "Active Subscriptions Movement", "Changes in active subscriptions: new, renewals, churned", "subscriptions"},
	{"actives_new", "New Active Subscriptions", "New paid subscriptions started in the period", "subscriptions"},
	{"trials", "Active Trials", "Total number of active trial subscriptions", "trials"},
	{"trials_movement", "Trials Movement", "Changes in trials: started, converted, expired", "trials"},
	{"trials_new", "New Trials", "New trial subscriptions started in the period", "trials"},
	{"trial_conversion_rate", "Trial Conversion Rate", "Percentage of trials that convert to paid subscriptions", "conversion"},
	{"conversion_to_paying", "Conversion to Paying", "Users converting from free to paid", "conversion"},
	{"churn", "Churn", "Rate of subscription cancellations", "retention"},
	{"refund_rate", "Refund Rate", "Percentage of transactions refunded", "retention"},
	{"subscription_retention", "Subscription Retention", "Cohort-based subscription retention over time", "retention"},
	{"subscription_status", "Subscription Status", "Active subscriptions broken down by status", "subscriptions"},
	{"customers_active", "Active Customers", "Total active customers with at least one subscription", "customers"},
	{"customers_new", "New Customers", "New customers acquired in the period", "customers"},
	{"ltv_per_customer", "LTV per Customer", "Average lifetime value across all customers", "revenue"},
	{"ltv_per_paying_customer", "LTV per Paying Customer", "Average lifetime value of paying customers only", "revenue"},
	{"cohort_explorer", "Cohort Explorer", "Explore cohort-based metrics over time", "retention"},
}

// --- API response types ---

type ChartData struct {
	Object          string           `json:"object"`
	Category        string           `json:"category"`
	DisplayType     string           `json:"display_type"`
	DisplayName     string           `json:"display_name"`
	Description     string           `json:"description"`
	DocLink         *string          `json:"documentation_link,omitempty"`
	LastComputedAt  *int64           `json:"last_computed_at,omitempty"`
	StartDate       *int64           `json:"start_date,omitempty"`
	EndDate         *int64           `json:"end_date,omitempty"`
	YAxisCurrency   string           `json:"yaxis_currency,omitempty"`
	FilterAllowed   bool             `json:"filtering_allowed"`
	SegmentAllowed  bool             `json:"segmenting_allowed"`
	Resolution      string           `json:"resolution"`
	Values          json.RawMessage  `json:"values"`
	Summary         json.RawMessage  `json:"summary,omitempty"`
	YAxis           string           `json:"yaxis"`
	Segments        json.RawMessage  `json:"segments,omitempty"`
	Measures        json.RawMessage  `json:"measures,omitempty"`
	UserSelectors   json.RawMessage  `json:"user_selectors,omitempty"`
}

type ChartOptions struct {
	Object      string             `json:"object"`
	Resolutions []ResolutionOption `json:"resolutions"`
	Segments    []SegmentOption    `json:"segments"`
	Filters     []FilterOption     `json:"filters"`
}

type ResolutionOption struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

type SegmentOption struct {
	ID               string `json:"id"`
	DisplayName      string `json:"display_name"`
	GroupDisplayName  string `json:"group_display_name,omitempty"`
}

type FilterOption struct {
	ID               string         `json:"id"`
	DisplayName      string         `json:"display_name"`
	GroupDisplayName  string         `json:"group_display_name,omitempty"`
	Options          []FilterValue  `json:"options"`
}

type FilterValue struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

// --- Flags ---

var (
	startDate  string
	endDate    string
	resolution string
	currency   string
	segment    string
	filters    string
	selectors  string
	aggregate  string

	realtime bool

	exportFormat string
	exportFile   string
	exportCharts string
)

var resolutionMap = map[string]string{
	"day": "0", "week": "1", "month": "2",
	"d": "0", "w": "1", "m": "2",
	"0": "0", "1": "1", "2": "2",
}

// --- Commands ---

// ChartsCmd is the parent command for chart operations
var ChartsCmd = &cobra.Command{
	Use:   "charts",
	Short: "Charts & analytics",
	Long:  "View, explore, and export subscription analytics charts from your RevenueCat project.",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available chart types",
	RunE:  runList,
}

var getCmd = &cobra.Command{
	Use:   "get <chart_name>",
	Short: "Get chart data",
	Long: `Fetch time-series data for a specific chart.

Available charts: revenue, mrr, arr, actives, trials, churn, and more.
Run 'rc charts list' to see all available chart types.`,
	Args: cobra.ExactArgs(1),
	RunE: runGet,
}

var optionsCmd = &cobra.Command{
	Use:   "options <chart_name>",
	Short: "Get available options for a chart",
	Long:  "Discover available resolutions, segments, and filters for a chart.",
	Args:  cobra.ExactArgs(1),
	RunE:  runOptions,
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export charts as visual report",
	Long: `Export multiple charts as a visual analytics report.

Formats: html, pdf, json, yaml, csv

HTML/PDF reports include SVG charts with data visualization.
Default charts: revenue, mrr, actives, trials, churn, arr`,
	RunE: runExport,
}

func init() {
	ChartsCmd.AddCommand(listCmd)
	ChartsCmd.AddCommand(getCmd)
	ChartsCmd.AddCommand(optionsCmd)
	ChartsCmd.AddCommand(exportCmd)

	// get flags
	getCmd.Flags().StringVar(&startDate, "start-date", "", "start date (YYYY-MM-DD)")
	getCmd.Flags().StringVar(&endDate, "end-date", "", "end date (YYYY-MM-DD)")
	getCmd.Flags().StringVar(&resolution, "resolution", "", "time resolution: day, week, month")
	getCmd.Flags().StringVar(&currency, "currency", "", "currency code (USD, EUR, GBP, etc.)")
	getCmd.Flags().StringVar(&segment, "segment", "", "segment by dimension (country, store, product, etc.)")
	getCmd.Flags().StringVar(&filters, "filters", "", `JSON filters, e.g. '[{"name":"country","values":["US"]}]'`)
	getCmd.Flags().StringVar(&selectors, "selectors", "", `JSON selectors, e.g. '{"revenue_type":"proceeds"}'`)
	getCmd.Flags().StringVar(&aggregate, "aggregate", "", "aggregate operations: average, total")
	getCmd.Flags().BoolVar(&realtime, "realtime", false, "use real-time (v3) charts")

	// options flags
	optionsCmd.Flags().BoolVar(&realtime, "realtime", false, "use real-time (v3) charts")

	// export flags
	exportCmd.Flags().StringVar(&exportFormat, "format", "html", "output format: html, pdf, json, yaml, csv")
	exportCmd.Flags().StringVar(&exportFile, "file", "", "output file path")
	exportCmd.Flags().StringVar(&exportCharts, "charts", "revenue,mrr,actives,trials,churn,arr", "comma-separated chart names")
	exportCmd.Flags().StringVar(&startDate, "start-date", "", "start date (YYYY-MM-DD)")
	exportCmd.Flags().StringVar(&endDate, "end-date", "", "end date (YYYY-MM-DD)")
	exportCmd.Flags().StringVar(&resolution, "resolution", "", "time resolution: day, week, month")
	exportCmd.Flags().StringVar(&currency, "currency", "", "currency code (USD, EUR, GBP, etc.)")
	exportCmd.Flags().BoolVar(&realtime, "realtime", false, "use real-time (v3) charts")
}

// --- List ---

func runList(_ *cobra.Command, _ []string) error {
	output.Print(chartCatalog)
	return nil
}

// --- Get ---

func runGet(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	chartName := args[0]
	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}
	ctx, cancel := client.Context()
	defer cancel()

	path := buildChartPath(cli.GetProjectID(), chartName)

	var result ChartData
	if err := client.Get(ctx, path, &result); err != nil {
		return fmt.Errorf("fetching chart %s: %w", chartName, err)
	}

	output.Print(result)
	return nil
}

// --- Options ---

func runOptions(cmd *cobra.Command, args []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	chartName := args[0]
	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}
	ctx, cancel := client.Context()
	defer cancel()

	optionsPath := fmt.Sprintf("/projects/%s/charts/%s/options", cli.GetProjectID(), chartName)
	if !realtime {
		optionsPath += "?realtime=false"
	}
	path := optionsPath

	var result ChartOptions
	if err := client.Get(ctx, path, &result); err != nil {
		return fmt.Errorf("fetching options for %s: %w", chartName, err)
	}

	output.Print(result)
	return nil
}

// --- Export ---

type exportedChart struct {
	Name        string          `json:"name" yaml:"name"`
	DisplayName string          `json:"display_name" yaml:"display_name"`
	Description string          `json:"description" yaml:"description"`
	Resolution  string          `json:"resolution" yaml:"resolution"`
	YAxis       string          `json:"yaxis" yaml:"yaxis"`
	Currency    string          `json:"currency,omitempty" yaml:"currency,omitempty"`
	StartDate   string          `json:"start_date,omitempty" yaml:"start_date,omitempty"`
	EndDate     string          `json:"end_date,omitempty" yaml:"end_date,omitempty"`
	Values      json.RawMessage `json:"values" yaml:"-"`
	ValuesYAML  interface{}     `json:"-" yaml:"values,omitempty"`
	Summary     json.RawMessage `json:"summary,omitempty" yaml:"-"`
	SummaryYAML interface{}     `json:"-" yaml:"summary,omitempty"`
}

type chartExportReport struct {
	GeneratedAt string          `json:"generated_at" yaml:"generated_at"`
	ProjectID   string          `json:"project_id" yaml:"project_id"`
	Charts      []exportedChart `json:"charts" yaml:"charts"`
}

func runExport(cmd *cobra.Command, _ []string) error {
	if err := cli.RequireProject(cmd); err != nil {
		return err
	}

	chartNames := strings.Split(exportCharts, ",")
	for i := range chartNames {
		chartNames[i] = strings.TrimSpace(chartNames[i])
	}

	client, err := api.NewClient(cli.GetProjectID(), parseTimeout())
	if err != nil {
		return err
	}
	ctx, cancel := client.Context()
	defer cancel()

	var charts []exportedChart
	for _, name := range chartNames {
		output.PrintInfo("Fetching %s...", name)
		path := buildChartPath(cli.GetProjectID(), name)

		var data ChartData
		if err := client.Get(ctx, path, &data); err != nil {
			output.PrintWarning("Skipping %s: %v", name, err)
			continue
		}

		ec := exportedChart{
			Name:        name,
			DisplayName: data.DisplayName,
			Description: data.Description,
			Resolution:  data.Resolution,
			YAxis:       data.YAxis,
			Currency:    data.YAxisCurrency,
			Values:      data.Values,
			Summary:     data.Summary,
		}

		if data.StartDate != nil {
			ec.StartDate = time.Unix(*data.StartDate, 0).UTC().Format("2006-01-02")
		}
		if data.EndDate != nil {
			ec.EndDate = time.Unix(*data.EndDate, 0).UTC().Format("2006-01-02")
		}

		// Parse values for YAML (so it's not raw JSON)
		var vals interface{}
		json.Unmarshal(data.Values, &vals)
		ec.ValuesYAML = vals

		var summ interface{}
		json.Unmarshal(data.Summary, &summ)
		ec.SummaryYAML = summ

		charts = append(charts, ec)
	}

	if len(charts) == 0 {
		return fmt.Errorf("no chart data retrieved")
	}

	report := chartExportReport{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		ProjectID:   cli.GetProjectID(),
		Charts:      charts,
	}

	// Determine output file
	outFile := exportFile
	if outFile == "" {
		switch exportFormat {
		case "html":
			outFile = "rc-charts.html"
		case "pdf":
			outFile = "rc-charts.pdf"
		case "json":
			outFile = "rc-charts.json"
		case "yaml":
			outFile = "rc-charts.yaml"
		case "csv":
			outFile = "rc-charts.csv"
		default:
			outFile = "rc-charts." + exportFormat
		}
	}

	switch exportFormat {
	case "json":
		data, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling JSON: %w", err)
		}
		if err := os.WriteFile(outFile, data, 0644); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}

	case "yaml":
		data, err := yaml.Marshal(report)
		if err != nil {
			return fmt.Errorf("marshaling YAML: %w", err)
		}
		if err := os.WriteFile(outFile, data, 0644); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}

	case "csv":
		if err := writeCSV(outFile, charts); err != nil {
			return err
		}

	case "html":
		htmlContent, err := renderChartsHTML(report)
		if err != nil {
			return err
		}
		if err := os.WriteFile(outFile, []byte(htmlContent), 0644); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}

	case "pdf":
		htmlContent, err := renderChartsHTML(report)
		if err != nil {
			return err
		}
		if err := htmlToPDF(htmlContent, outFile); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unsupported format: %s (use html, pdf, json, yaml, csv)", exportFormat)
	}

	output.PrintSuccess("Charts exported to %s", outFile)
	return nil
}

// --- Path builder ---

func buildChartPath(projectID, chartName string) string {
	path := fmt.Sprintf("/projects/%s/charts/%s", projectID, chartName)

	var params []string
	if !realtime {
		params = append(params, "realtime=false")
	}
	if startDate != "" {
		params = append(params, "start_date="+startDate)
	}
	if endDate != "" {
		params = append(params, "end_date="+endDate)
	}
	if resolution != "" {
		r := resolution
		if mapped, ok := resolutionMap[strings.ToLower(r)]; ok {
			r = mapped
		}
		params = append(params, "resolution="+r)
	}
	if currency != "" {
		params = append(params, "currency="+strings.ToUpper(currency))
	}
	if segment != "" {
		params = append(params, "segment="+segment)
	}
	if filters != "" {
		params = append(params, "filters="+filters)
	}
	if selectors != "" {
		params = append(params, "selectors="+selectors)
	}
	if aggregate != "" {
		params = append(params, "aggregate="+aggregate)
	}

	if len(params) > 0 {
		path += "?" + strings.Join(params, "&")
	}
	return path
}

// --- CSV export ---

func writeCSV(outFile string, charts []exportedChart) error {
	f, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("creating CSV file: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	w.Write([]string{"chart", "index", "date", "value"})

	for _, ch := range charts {
		points := parseDataPoints(ch.Values)
		for i, pt := range points {
			w.Write([]string{ch.Name, fmt.Sprintf("%d", i), pt.date, fmt.Sprintf("%.2f", pt.value)})
		}
	}
	return nil
}

// --- Data point extraction ---

type dataPoint struct {
	date  string
	value float64
}

// parseDataPoints extracts numeric data from chart values.
// RevenueCat charts return values as arrays: [timestamp_seconds, value, ...] per point.
func parseDataPoints(raw json.RawMessage) []dataPoint {
	var points []dataPoint

	// Try array of arrays: [[timestamp, value, ...], ...]
	var arrays [][]json.Number
	if err := json.Unmarshal(raw, &arrays); err == nil {
		for _, arr := range arrays {
			if len(arr) >= 2 {
				ts, _ := arr[0].Float64()
				val, _ := arr[1].Float64()
				// RevenueCat uses epoch seconds (not ms)
				date := time.Unix(int64(ts), 0).UTC().Format("2006-01-02")
				points = append(points, dataPoint{date: date, value: val})
			}
		}
		return points
	}

	// Try array of objects with date/value fields
	var objects []map[string]interface{}
	if err := json.Unmarshal(raw, &objects); err == nil {
		for _, obj := range objects {
			var date string
			var val float64

			if d, ok := obj["date"]; ok {
				date = fmt.Sprintf("%v", d)
			} else if ts, ok := obj["timestamp"]; ok {
				if f, ok := ts.(float64); ok {
					date = time.Unix(int64(f), 0).UTC().Format("2006-01-02")
				}
			}

			if v, ok := obj["value"]; ok {
				if f, ok := v.(float64); ok {
					val = f
				}
			}

			points = append(points, dataPoint{date: date, value: val})
		}
	}

	return points
}

// --- SVG chart generation ---

type svgChart struct {
	Title  string
	Desc   string
	YAxis  string
	Points []dataPoint
}

func renderSVG(ch svgChart) string {
	points := ch.Points
	if len(points) == 0 {
		return fmt.Sprintf(`<div class="chart-empty">No data available for %s</div>`, template.HTMLEscapeString(ch.Title))
	}

	width := 700.0
	height := 300.0
	padL := 70.0  // left padding for y-axis labels
	padR := 20.0
	padT := 20.0
	padB := 60.0  // bottom padding for x-axis labels
	plotW := width - padL - padR
	plotH := height - padT - padB

	// Find min/max
	minVal, maxVal := points[0].value, points[0].value
	for _, p := range points {
		if p.value < minVal {
			minVal = p.value
		}
		if p.value > maxVal {
			maxVal = p.value
		}
	}

	// Add 10% padding to range
	valRange := maxVal - minVal
	if valRange == 0 {
		valRange = 1
		minVal -= 0.5
		maxVal += 0.5
	} else {
		minVal -= valRange * 0.05
		maxVal += valRange * 0.05
		if minVal < 0 && points[0].value >= 0 {
			minVal = 0
		}
		valRange = maxVal - minVal
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf(`<svg viewBox="0 0 %.0f %.0f" xmlns="http://www.w3.org/2000/svg" style="width:100%%;max-width:%.0fpx;height:auto">`, width, height, width))

	// Background
	b.WriteString(fmt.Sprintf(`<rect width="%.0f" height="%.0f" fill="#fafafa" rx="4"/>`, width, height))

	// Grid lines (5 horizontal)
	for i := 0; i <= 4; i++ {
		y := padT + plotH - (float64(i)/4.0)*plotH
		val := minVal + (float64(i)/4.0)*valRange
		b.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#e0e0e0" stroke-width="1"/>`, padL, y, width-padR, y))
		label := formatAxisLabel(val, ch.YAxis)
		b.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" text-anchor="end" font-size="11" fill="#666" font-family="system-ui,sans-serif">%s</text>`, padL-8, y+4, template.HTMLEscapeString(label)))
	}

	// Data polyline + area fill
	n := len(points)
	var polyPoints, areaPoints strings.Builder
	for i, p := range points {
		x := padL + (float64(i)/float64(n-max(1, 1)))*plotW
		if n > 1 {
			x = padL + (float64(i)/float64(n-1))*plotW
		}
		y := padT + plotH - ((p.value-minVal)/valRange)*plotH
		if i == 0 {
			areaPoints.WriteString(fmt.Sprintf("%.1f,%.1f ", x, padT+plotH))
		}
		polyPoints.WriteString(fmt.Sprintf("%.1f,%.1f ", x, y))
		areaPoints.WriteString(fmt.Sprintf("%.1f,%.1f ", x, y))
	}
	// Close area
	lastX := padL + (float64(n-1)/float64(max(n-1, 1)))*plotW
	areaPoints.WriteString(fmt.Sprintf("%.1f,%.1f", lastX, padT+plotH))

	// Area fill
	b.WriteString(fmt.Sprintf(`<polygon points="%s" fill="rgba(75,72,242,0.1)" />`, strings.TrimSpace(areaPoints.String())))
	// Line
	b.WriteString(fmt.Sprintf(`<polyline points="%s" fill="none" stroke="#4B48F2" stroke-width="2.5" stroke-linejoin="round" stroke-linecap="round"/>`, strings.TrimSpace(polyPoints.String())))

	// Data points (circles) — show all if <=30 points, else every Nth
	step := 1
	if n > 30 {
		step = n / 15
	}
	for i := 0; i < n; i += step {
		p := points[i]
		x := padL + (float64(i)/float64(max(n-1, 1)))*plotW
		y := padT + plotH - ((p.value-minVal)/valRange)*plotH
		label := formatAxisLabel(p.value, ch.YAxis)
		b.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="3" fill="#4B48F2"><title>%s: %s</title></circle>`, x, y, template.HTMLEscapeString(p.date), template.HTMLEscapeString(label)))
	}

	// X-axis labels
	labelStep := 1
	if n > 12 {
		labelStep = n / 6
	}
	for i := 0; i < n; i += labelStep {
		x := padL + (float64(i)/float64(max(n-1, 1)))*plotW
		date := points[i].date
		if len(date) > 10 {
			date = date[:10]
		}
		b.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" text-anchor="middle" font-size="10" fill="#999" font-family="system-ui,sans-serif" transform="rotate(-45 %.1f %.1f)">%s</text>`, x, padT+plotH+15, x, padT+plotH+15, template.HTMLEscapeString(date)))
	}

	b.WriteString("</svg>")
	return b.String()
}

func formatAxisLabel(val float64, yaxis string) string {
	switch yaxis {
	case "$":
		if math.Abs(val) >= 1_000_000 {
			return fmt.Sprintf("$%.1fM", val/1_000_000)
		}
		if math.Abs(val) >= 1_000 {
			return fmt.Sprintf("$%.1fK", val/1_000)
		}
		return fmt.Sprintf("$%.0f", val)
	case "%":
		return fmt.Sprintf("%.1f%%", val)
	default:
		if math.Abs(val) >= 1_000_000 {
			return fmt.Sprintf("%.1fM", val/1_000_000)
		}
		if math.Abs(val) >= 1_000 {
			return fmt.Sprintf("%.1fK", val/1_000)
		}
		return fmt.Sprintf("%.0f", val)
	}
}

// --- HTML rendering ---

type htmlData struct {
	GeneratedAt string
	ProjectID   string
	Charts      []htmlChartData
}

type htmlChartData struct {
	DisplayName string
	Description string
	Resolution  string
	DateRange   string
	Currency    string
	SVG         template.HTML
	SummaryHTML template.HTML
}

func renderChartsHTML(report chartExportReport) (string, error) {
	var htmlCharts []htmlChartData

	for _, ch := range report.Charts {
		points := parseDataPoints(ch.Values)
		svg := renderSVG(svgChart{
			Title:  ch.DisplayName,
			Desc:   ch.Description,
			YAxis:  ch.YAxis,
			Points: points,
		})

		dateRange := ""
		if ch.StartDate != "" && ch.EndDate != "" {
			dateRange = ch.StartDate + " to " + ch.EndDate
		}

		summaryHTML := renderSummary(ch.Summary)

		htmlCharts = append(htmlCharts, htmlChartData{
			DisplayName: ch.DisplayName,
			Description: ch.Description,
			Resolution:  ch.Resolution,
			DateRange:   dateRange,
			Currency:    ch.Currency,
			SVG:         template.HTML(svg),
			SummaryHTML: summaryHTML,
		})
	}

	data := htmlData{
		GeneratedAt: report.GeneratedAt,
		ProjectID:   report.ProjectID,
		Charts:      htmlCharts,
	}

	tmpl, err := template.New("charts").Parse(chartsHTMLTemplate)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}

func renderSummary(raw json.RawMessage) template.HTML {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}

	var summary map[string]interface{}
	if err := json.Unmarshal(raw, &summary); err != nil {
		return ""
	}

	if len(summary) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(`<div class="summary-row">`)
	for key, val := range summary {
		label := strings.ReplaceAll(strings.Title(strings.ReplaceAll(key, "_", " ")), " ", " ")
		b.WriteString(fmt.Sprintf(`<div class="summary-item"><span class="summary-label">%s</span><span class="summary-value">%v</span></div>`, template.HTMLEscapeString(label), val))
	}
	b.WriteString(`</div>`)
	return template.HTML(b.String())
}

// --- PDF generation (same pattern as report.go) ---

func htmlToPDF(html, pdfPath string) error {
	chromePath := findChrome()
	if chromePath == "" {
		return fmt.Errorf("PDF export requires Chrome or Chromium.\nInstall Chrome, or use --format html and print to PDF from your browser")
	}

	tmpFile, err := os.CreateTemp("", "rc-charts-*.html")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(html); err != nil {
		tmpFile.Close()
		return fmt.Errorf("writing temp file: %w", err)
	}
	tmpFile.Close()

	absPDF, err := filepath.Abs(pdfPath)
	if err != nil {
		absPDF = pdfPath
	}

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
		return fmt.Errorf("Chrome PDF conversion failed: %w\nTry: rc charts export --format html, then print to PDF from browser", err)
	}

	return nil
}

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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// --- HTML Template ---

const chartsHTMLTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>RevenueCat Charts Report</title>
<style>
  *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    background: #fff;
    color: #1a1a2e;
    padding: 40px;
    max-width: 900px;
    margin: 0 auto;
  }
  .header {
    text-align: center;
    margin-bottom: 40px;
    padding-bottom: 24px;
    border-bottom: 3px solid #4B48F2;
  }
  .header h1 {
    font-size: 28px;
    color: #4B48F2;
    margin-bottom: 8px;
  }
  .header .meta {
    color: #666;
    font-size: 14px;
  }
  .chart-section {
    margin-bottom: 48px;
    page-break-inside: avoid;
  }
  .chart-title {
    font-size: 20px;
    font-weight: 700;
    color: #1a1a2e;
    margin-bottom: 4px;
  }
  .chart-desc {
    font-size: 13px;
    color: #666;
    margin-bottom: 4px;
  }
  .chart-meta {
    font-size: 12px;
    color: #999;
    margin-bottom: 16px;
  }
  .chart-meta span {
    margin-right: 16px;
  }
  .chart-svg-container {
    background: #fafafa;
    border: 1px solid #eee;
    border-radius: 8px;
    padding: 16px;
    margin-bottom: 12px;
  }
  .chart-empty {
    text-align: center;
    color: #999;
    padding: 40px;
    font-size: 14px;
  }
  .summary-row {
    display: flex;
    gap: 24px;
    flex-wrap: wrap;
    margin-top: 8px;
  }
  .summary-item {
    display: flex;
    flex-direction: column;
    background: #f5f4ff;
    padding: 12px 16px;
    border-radius: 6px;
    min-width: 120px;
  }
  .summary-label {
    font-size: 11px;
    color: #666;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    margin-bottom: 4px;
  }
  .summary-value {
    font-size: 18px;
    font-weight: 700;
    color: #4B48F2;
  }
  .footer {
    text-align: center;
    color: #999;
    font-size: 12px;
    margin-top: 40px;
    padding-top: 20px;
    border-top: 1px solid #eee;
  }
  @media print {
    body { padding: 20px; }
    .chart-section { page-break-inside: avoid; }
  }
</style>
</head>
<body>
  <div class="header">
    <h1>Charts & Analytics Report</h1>
    <div class="meta">
      Project: {{.ProjectID}} &middot; Generated: {{.GeneratedAt}}
    </div>
  </div>

  {{range .Charts}}
  <div class="chart-section">
    <div class="chart-title">{{.DisplayName}}</div>
    <div class="chart-desc">{{.Description}}</div>
    <div class="chart-meta">
      {{if .Resolution}}<span>Resolution: {{.Resolution}}</span>{{end}}
      {{if .DateRange}}<span>{{.DateRange}}</span>{{end}}
      {{if .Currency}}<span>Currency: {{.Currency}}</span>{{end}}
    </div>
    <div class="chart-svg-container">
      {{.SVG}}
    </div>
    {{.SummaryHTML}}
  </div>
  {{end}}

  <div class="footer">
    Generated by revenuecat-cli &middot; https://github.com/AndroidPoet/revenuecat-cli
  </div>
</body>
</html>`
