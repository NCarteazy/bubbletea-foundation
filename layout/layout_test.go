package layout

import "testing"

func TestTwoPane(t *testing.T) {
	result := TwoPane("LEFT", "RIGHT", 80, 24)
	if result == "" {
		t.Fatal("TwoPane returned empty string")
	}
}

func TestTwoPaneWithRatio(t *testing.T) {
	result := TwoPaneWithRatio("LEFT", "RIGHT", 0.3, 80, 24)
	if result == "" {
		t.Fatal("TwoPaneWithRatio returned empty string")
	}
}

func TestOverlayCenter(t *testing.T) {
	result := OverlayCenter(80, 24, "background", "overlay")
	if result == "" {
		t.Fatal("OverlayCenter returned empty string")
	}
}

func TestOverlayCenterZeroDimensions(t *testing.T) {
	result := OverlayCenter(0, 0, "bg", "overlay")
	if result != "overlay" {
		t.Fatalf("expected %q, got %q", "overlay", result)
	}
}
