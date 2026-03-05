package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func setupTest(format string, pretty, quiet bool) *bytes.Buffer {
	buf := &bytes.Buffer{}
	SetWriter(buf)
	Setup(format, pretty, quiet)
	return buf
}

func TestPrintJSON(t *testing.T) {
	buf := setupTest("json", false, false)

	data := []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}{
		{ID: "1", Name: "test"},
	}

	if err := Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	var result []map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Invalid JSON output: %v", err)
	}

	if result[0]["id"] != "1" {
		t.Errorf("expected id=1, got %s", result[0]["id"])
	}
}

func TestPrintPrettyJSON(t *testing.T) {
	buf := setupTest("json", true, false)

	data := map[string]string{"key": "value"}
	if err := Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "  ") {
		t.Error("expected indented JSON for pretty mode")
	}
}

func TestPrintCSV(t *testing.T) {
	buf := setupTest("csv", false, false)

	type Item struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	data := []Item{
		{ID: "1", Name: "alpha"},
		{ID: "2", Name: "beta"},
	}

	if err := Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 { // header + 2 rows
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "ID") {
		t.Error("expected header with ID")
	}
}

func TestPrintYAML(t *testing.T) {
	buf := setupTest("yaml", false, false)

	data := map[string]string{"key": "value"}
	if err := Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	if !strings.Contains(buf.String(), "key: value") {
		t.Error("expected YAML output")
	}
}

func TestPrintMinimal(t *testing.T) {
	buf := setupTest("minimal", false, false)

	type Item struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	data := []Item{
		{ID: "abc", Name: "test"},
	}

	if err := Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	output := strings.TrimSpace(buf.String())
	if output != "abc" {
		t.Errorf("expected 'abc', got '%s'", output)
	}
}

func TestPrintTable(t *testing.T) {
	buf := setupTest("table", false, false)

	type Item struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	data := []Item{
		{ID: "1", Name: "alpha"},
	}

	if err := Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "ID") || !strings.Contains(output, "NAME") {
		t.Error("expected table headers")
	}
	if !strings.Contains(output, "alpha") {
		t.Error("expected table row data")
	}
}

func TestPrintEmptySlice(t *testing.T) {
	buf := setupTest("table", false, false)

	type Item struct {
		ID string `json:"id"`
	}

	if err := Print([]Item{}); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	if !strings.Contains(buf.String(), "no results") {
		t.Error("expected 'no results' for empty slice")
	}
}

func TestQuietMode(t *testing.T) {
	buf := setupTest("json", false, true)

	PrintSuccess("should not appear")
	PrintInfo("should not appear")

	if buf.Len() != 0 {
		t.Error("expected no output in quiet mode")
	}
}

func TestPrintTSV(t *testing.T) {
	buf := setupTest("tsv", false, false)

	type Item struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	data := []Item{{ID: "1", Name: "test"}}

	if err := Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "\t") {
		t.Error("expected tab-separated output")
	}
}
