package foundation

import "testing"

func TestFlashSetAndActive(t *testing.T) {
	f := &flashState{}
	f.set("Hello!", 5)

	if !f.active() {
		t.Fatal("expected flash to be active after set")
	}
	if f.message != "Hello!" {
		t.Fatalf("expected message=%q, got %q", "Hello!", f.message)
	}
}

func TestFlashTickDecay(t *testing.T) {
	f := &flashState{}
	f.set("Hello!", 5)

	for i := 0; i < 5; i++ {
		if !f.active() {
			t.Fatalf("expected flash active at tick %d", i)
		}
		f.tick()
	}

	if f.active() {
		t.Fatal("expected flash inactive after 5 ticks")
	}
	if f.message != "" {
		t.Fatalf("expected message=%q, got %q", "", f.message)
	}
}

func TestFlashDefaultInactive(t *testing.T) {
	f := &flashState{}

	if f.active() {
		t.Fatal("expected default flash state to be inactive")
	}
}

func TestFlashTickOnDefault(t *testing.T) {
	f := &flashState{}
	f.tick() // should not panic
}
