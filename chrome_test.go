package foundation

import (
	"strings"
	"testing"

	"github.com/NCarteazy/bubbletea-foundation/theme"
)

func TestRenderBreadcrumbsWithCrumbs(t *testing.T) {
	th := theme.TokyoNight
	crumbs := []string{"Dashboard", "Reports", "Detail"}
	result := renderBreadcrumbs(crumbs, th)

	if result == "" {
		t.Fatal("expected non-empty breadcrumb string")
	}
	for _, crumb := range crumbs {
		if !strings.Contains(result, crumb) {
			t.Errorf("breadcrumb output missing %q", crumb)
		}
	}
}

func TestRenderBreadcrumbsNil(t *testing.T) {
	th := theme.TokyoNight
	result := renderBreadcrumbs(nil, th)

	if result != "" {
		t.Errorf("expected empty string for nil crumbs, got %q", result)
	}
}

func TestRenderBreadcrumbsEmpty(t *testing.T) {
	th := theme.TokyoNight
	result := renderBreadcrumbs([]string{}, th)

	if result != "" {
		t.Errorf("expected empty string for empty crumbs, got %q", result)
	}
}

func TestRenderStatusBarWithHints(t *testing.T) {
	th := theme.TokyoNight
	hints := []KeyHint{
		{Key: "q", Desc: "quit"},
		{Key: "enter", Desc: "select"},
	}
	result := renderStatusBar(hints, nil, 80, th)

	if result == "" {
		t.Fatal("expected non-empty status bar string")
	}
	for _, h := range hints {
		if !strings.Contains(result, h.Desc) {
			t.Errorf("status bar output missing hint desc %q", h.Desc)
		}
	}
}

func TestRenderStatusBarWithFlash(t *testing.T) {
	th := theme.TokyoNight
	hints := []KeyHint{
		{Key: "q", Desc: "quit"},
	}
	flash := &flashState{}
	flash.set("File saved!", 5)

	result := renderStatusBar(hints, flash, 80, th)

	if !strings.Contains(result, "File saved!") {
		t.Errorf("expected flash message in status bar, got %q", result)
	}
}
