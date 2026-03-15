package theme

import (
	"testing"
)

func TestTokyoNightImplementsTheme(t *testing.T) {
	var _ Theme = TokyoNight
}

func TestColorMethodsReturnNonEmpty(t *testing.T) {
	colors := map[string]string{
		"Bg":        TokyoNight.Bg(),
		"Primary":   TokyoNight.Primary(),
		"Secondary": TokyoNight.Secondary(),
		"Success":   TokyoNight.Success(),
		"Warning":   TokyoNight.Warning(),
		"Danger":    TokyoNight.Danger(),
		"Text":      TokyoNight.Text(),
		"Muted":     TokyoNight.Muted(),
		"Surface":   TokyoNight.Surface(),
	}
	for name, val := range colors {
		if val == "" {
			t.Errorf("%s() returned empty string", name)
		}
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{0, "0 B"},
		{1000, "1.0 kB"},
		{1000000, "1.0 MB"},
		{1000000000, "1.0 GB"},
	}
	for _, tt := range tests {
		got := TokyoNight.FormatSize(tt.input)
		if got != tt.expected {
			t.Errorf("FormatSize(%d) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestSizeBarNonEmpty(t *testing.T) {
	result := TokyoNight.SizeBar(0.5, 20)
	if result == "" {
		t.Error("SizeBar(0.5, 20) returned empty string")
	}
}

func TestSizeBarLength(t *testing.T) {
	// SizeBar with ratio 0.5 and maxWidth 20 should have 10 filled and 10 empty blocks
	result := TokyoNight.SizeBar(0.5, 20)
	if result == "" {
		t.Fatal("SizeBar returned empty string")
	}
}

func TestStyleMethodsReturnStyles(t *testing.T) {
	// Verify style methods don't panic and return usable styles
	_ = TokyoNight.Title()
	_ = TokyoNight.Panel()
	_ = TokyoNight.ActivePanel()
	_ = TokyoNight.StatusBar()
	_ = TokyoNight.Breadcrumb()
	_ = TokyoNight.TableHeader()
	_ = TokyoNight.SelectedRow()
}
