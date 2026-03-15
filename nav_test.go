package foundation

import "testing"

func TestPushTwoItems(t *testing.T) {
	s := &navStack{}
	s.push("a", nil)
	s.push("b", 42)

	if s.len() != 2 {
		t.Fatalf("expected len=2, got %d", s.len())
	}
	cur := s.current()
	if cur == nil || cur.viewID != "b" {
		t.Fatalf("expected current viewID=b, got %v", cur)
	}
	if cur.data != 42 {
		t.Fatalf("expected current data=42, got %v", cur.data)
	}
}

func TestPushTwoPop(t *testing.T) {
	s := &navStack{}
	s.push("a", "dataA")
	s.push("b", "dataB")
	s.pop()

	if s.len() != 1 {
		t.Fatalf("expected len=1, got %d", s.len())
	}
	cur := s.current()
	if cur == nil || cur.viewID != "a" {
		t.Fatalf("expected current viewID=a, got %v", cur)
	}
	if cur.data != "dataA" {
		t.Fatalf("expected current data=dataA, got %v", cur.data)
	}
}

func TestPopOnEmpty(t *testing.T) {
	s := &navStack{}
	s.pop() // should not panic

	if s.len() != 0 {
		t.Fatalf("expected len=0, got %d", s.len())
	}
}

func TestReplace(t *testing.T) {
	s := &navStack{}
	s.push("a", "d1")
	s.push("b", "d2")
	s.replace("c", "d3")

	if s.len() != 2 {
		t.Fatalf("expected len=2, got %d", s.len())
	}
	cur := s.current()
	if cur == nil || cur.viewID != "c" {
		t.Fatalf("expected current viewID=c, got %v", cur)
	}
	if cur.data != "d3" {
		t.Fatalf("expected current data=d3, got %v", cur.data)
	}
}

func TestBreadcrumbs(t *testing.T) {
	s := &navStack{}
	s.push("a", nil)
	s.push("b", nil)
	s.push("c", nil)

	bc := s.breadcrumbs()
	if len(bc) != 3 {
		t.Fatalf("expected 3 breadcrumbs, got %d", len(bc))
	}
	expected := []string{"a", "b", "c"}
	for i, v := range expected {
		if bc[i] != v {
			t.Fatalf("breadcrumbs[%d]: expected %q, got %q", i, v, bc[i])
		}
	}
}
