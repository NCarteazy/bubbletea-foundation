package foundation

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/NCarteazy/bubbletea-foundation/theme"
)

// testView is a minimal View implementation for testing.
type testView struct {
	id      string
	content string
}

func (v testView) ID() string                                    { return v.id }
func (v testView) Update(msg tea.Msg, ctx ViewContext) (View, tea.Cmd) { return v, nil }
func (v testView) Render(ctx ViewContext) string                  { return v.content }

// appTestOverlay is a minimal Overlay implementation for app tests.
type appTestOverlay struct {
	done    bool
	content string
}

func (o appTestOverlay) Update(msg tea.Msg) (Overlay, tea.Cmd) { return o, nil }
func (o appTestOverlay) View() string                          { return o.content }
func (o appTestOverlay) Done() bool                            { return o.done }

func TestNewApp(t *testing.T) {
	app := New(
		WithTitle("Test App"),
		WithTheme(theme.TokyoNight),
		WithViews(
			testView{id: "home", content: "home content"},
			testView{id: "detail", content: "detail content"},
		),
		WithInitialView("home"),
		WithStatusHints([]KeyHint{
			{Key: "q", Desc: "quit"},
		}),
		WithHelp([]KeySection{
			{Title: "General", Keys: []KeyHint{{Key: "?", Desc: "help"}}},
		}),
	)

	if app == nil {
		t.Fatal("expected non-nil App")
	}
	if app.title != "Test App" {
		t.Errorf("expected title 'Test App', got %q", app.title)
	}
	if len(app.views) != 2 {
		t.Errorf("expected 2 views, got %d", len(app.views))
	}
	if app.nav.current().viewID != "home" {
		t.Errorf("expected initial view 'home', got %q", app.nav.current().viewID)
	}
}

func TestAppView(t *testing.T) {
	app := New(
		WithViews(
			testView{id: "home", content: "hello world"},
		),
		WithInitialView("home"),
	)

	// Send a WindowSizeMsg so the app knows its dimensions.
	model, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	a := model.(*App)

	output := a.View()
	if !strings.Contains(output, "hello world") {
		t.Errorf("expected View() to contain 'hello world', got:\n%s", output)
	}
}

func TestAppNavigation(t *testing.T) {
	app := New(
		WithViews(
			testView{id: "home", content: "home content"},
			testView{id: "detail", content: "detail content"},
		),
		WithInitialView("home"),
	)

	// Navigate to detail
	model, _ := app.Update(NavigateMsg{ViewID: "detail", Data: nil})
	a := model.(*App)

	if a.nav.current().viewID != "detail" {
		t.Errorf("expected current view 'detail', got %q", a.nav.current().viewID)
	}
	if a.nav.len() != 2 {
		t.Errorf("expected nav stack len=2, got %d", a.nav.len())
	}

	// Go back
	model, _ = a.Update(BackMsg{})
	a = model.(*App)

	if a.nav.current().viewID != "home" {
		t.Errorf("expected current view 'home' after back, got %q", a.nav.current().viewID)
	}
	if a.nav.len() != 1 {
		t.Errorf("expected nav stack len=1, got %d", a.nav.len())
	}
}

func TestAppOverlay(t *testing.T) {
	app := New(
		WithViews(
			testView{id: "home", content: "home content"},
		),
		WithInitialView("home"),
	)

	// Size the app first
	model, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	a := model.(*App)

	// Push an overlay
	model, _ = a.Update(ShowOverlayMsg{Overlay: appTestOverlay{content: "overlay text"}})
	a = model.(*App)

	if !a.overlays.active() {
		t.Fatal("expected overlays to be active after push")
	}

	output := a.View()
	if !strings.Contains(output, "overlay text") {
		t.Errorf("expected View() to contain 'overlay text', got:\n%s", output)
	}
}

func TestAppFlash(t *testing.T) {
	app := New(
		WithViews(
			testView{id: "home", content: "home content"},
		),
		WithInitialView("home"),
	)

	// Size the app first
	model, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	a := model.(*App)

	// Send a flash message
	model, cmd := a.Update(FlashMsg{Message: "Saved!", Duration: 5})
	a = model.(*App)

	if !a.flash.active() {
		t.Fatal("expected flash to be active after FlashMsg")
	}
	if cmd == nil {
		t.Fatal("expected non-nil cmd from FlashMsg (tick)")
	}

	output := a.View()
	if !strings.Contains(output, "Saved!") {
		t.Errorf("expected View() to contain 'Saved!', got:\n%s", output)
	}
}

func TestAppCtrlCQuits(t *testing.T) {
	app := New(
		WithViews(
			testView{id: "home", content: "home content"},
		),
		WithInitialView("home"),
	)

	_, cmd := app.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	if cmd == nil {
		t.Fatal("expected non-nil cmd from ctrl+c (quit)")
	}

	// Execute the cmd and check it produces a QuitMsg
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Errorf("expected tea.QuitMsg from ctrl+c cmd, got %T", msg)
	}
}
