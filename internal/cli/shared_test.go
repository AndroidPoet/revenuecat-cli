package cli

import "testing"

func TestGetProjectID_Flag(t *testing.T) {
	SetProjectID("proj_test")
	defer SetProjectID("")

	if got := GetProjectID(); got != "proj_test" {
		t.Errorf("expected proj_test, got %s", got)
	}
}

func TestGetProfile_Default(t *testing.T) {
	SetProfile("")

	if got := GetProfile(); got != "default" {
		t.Errorf("expected 'default', got %s", got)
	}
}

func TestGetTimeout_Default(t *testing.T) {
	SetTimeout("")

	if got := GetTimeout(); got != "60s" {
		t.Errorf("expected '60s', got %s", got)
	}
}

func TestGetTimeout_Custom(t *testing.T) {
	SetTimeout("30s")
	defer SetTimeout("")

	if got := GetTimeout(); got != "30s" {
		t.Errorf("expected '30s', got %s", got)
	}
}

func TestIsDryRun(t *testing.T) {
	SetDryRun(true)
	defer SetDryRun(false)

	if !IsDryRun() {
		t.Error("expected dry run to be true")
	}
}
