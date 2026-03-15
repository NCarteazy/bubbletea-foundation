package foundation

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/NCarteazy/bubbletea-foundation/components"
	"github.com/NCarteazy/bubbletea-foundation/layout"
	"github.com/NCarteazy/bubbletea-foundation/theme"
)

// KeySection groups key hints under a title for the help overlay.
type KeySection struct {
	Title string
	Keys  []KeyHint
}

// Option configures an App during construction.
type Option func(*App)

// App is the top-level Bubble Tea model. It manages view routing,
// overlays, flash messages, breadcrumbs, status bar, and help.
type App struct {
	title       string
	th          theme.Theme
	views       map[string]View
	initialView string
	statusHints []KeyHint
	helpKeys    []KeySection

	nav      navStack
	overlays overlayStack
	flash    flashState

	width, height int
}

// --- Functional options ---

func WithTitle(title string) Option {
	return func(a *App) { a.title = title }
}

func WithTheme(th theme.Theme) Option {
	return func(a *App) { a.th = th }
}

func WithInitialView(id string) Option {
	return func(a *App) { a.initialView = id }
}

func WithViews(views ...View) Option {
	return func(a *App) {
		for _, v := range views {
			a.views[v.ID()] = v
		}
	}
}

func WithStatusHints(hints []KeyHint) Option {
	return func(a *App) { a.statusHints = hints }
}

func WithHelp(sections []KeySection) Option {
	return func(a *App) { a.helpKeys = sections }
}

// New creates a new App with the given options.
func New(opts ...Option) *App {
	a := &App{
		th:    theme.TokyoNight,
		views: make(map[string]View),
	}
	for _, opt := range opts {
		opt(a)
	}
	if a.initialView != "" {
		a.nav.push(a.initialView, nil)
	}
	return a
}

// Init implements tea.Model.
func (a *App) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		return a, nil

	case tea.KeyMsg:
		// ctrl+c always quits
		if msg.Type == tea.KeyCtrlC {
			return a, tea.Quit
		}

		// If overlays are active, route input to the top overlay
		if a.overlays.active() {
			top := a.overlays.top()
			updated, cmd := top.Update(msg)
			if updated.Done() {
				a.overlays.pop()
			} else {
				a.overlays.entries[len(a.overlays.entries)-1] = updated
			}
			return a, cmd
		}

		// "?" toggles the help overlay
		if msg.String() == "?" && len(a.helpKeys) > 0 {
			a.overlays.push(newHelpOverlay(a.helpKeys, a.th))
			return a, nil
		}

		// Route to current view
		return a.updateCurrentView(msg)

	case NavigateMsg:
		a.nav.push(msg.ViewID, msg.Data)
		return a, nil

	case BackMsg:
		if a.nav.len() > 1 {
			a.nav.pop()
		}
		return a, nil

	case ReplaceMsg:
		a.nav.replace(msg.ViewID, msg.Data)
		return a, nil

	case ShowOverlayMsg:
		a.overlays.push(msg.Overlay)
		return a, nil

	case FlashMsg:
		a.flash.set(msg.Message, msg.Duration)
		return a, flashTickCmd()

	case flashTickMsg:
		a.flash.tick()
		if a.flash.active() {
			return a, flashTickCmd()
		}
		return a, nil

	default:
		return a.updateCurrentView(msg)
	}
}

// updateCurrentView forwards a message to the current view.
func (a *App) updateCurrentView(msg tea.Msg) (tea.Model, tea.Cmd) {
	cur := a.nav.current()
	if cur == nil {
		return a, nil
	}
	v, ok := a.views[cur.viewID]
	if !ok {
		return a, nil
	}
	updated, cmd := v.Update(msg, a.viewContext())
	a.views[cur.viewID] = updated
	return a, cmd
}

// viewContext builds a ViewContext for the current view.
func (a *App) viewContext() ViewContext {
	// Chrome = breadcrumb line (1) + status bar line (1)
	const chromeLines = 2
	contentHeight := a.height - chromeLines
	if contentHeight < 0 {
		contentHeight = 0
	}

	var data any
	if cur := a.nav.current(); cur != nil {
		data = cur.data
	}

	return ViewContext{
		Width:  a.width,
		Height: contentHeight,
		Theme:  a.th,
		Data:   data,
	}
}

// View implements tea.Model.
func (a *App) View() string {
	// Breadcrumbs
	crumbs := renderBreadcrumbs(a.nav.breadcrumbs(), a.th)

	// Current view content
	var content string
	if cur := a.nav.current(); cur != nil {
		if v, ok := a.views[cur.viewID]; ok {
			content = v.Render(a.viewContext())
		}
	}

	// Status bar
	status := renderStatusBar(a.statusHints, &a.flash, a.width, a.th)

	// Compose vertically
	result := lipgloss.JoinVertical(lipgloss.Left, crumbs, content, status)

	// If an overlay is active, render it centered on top
	if a.overlays.active() {
		overlayContent := a.overlays.top().View()
		result = layout.OverlayCenter(a.width, a.height, result, overlayContent)
	}

	return result
}

// --- Built-in help overlay ---

type helpOverlay struct {
	sections []KeySection
	th       theme.Theme
	done     bool
}

func newHelpOverlay(sections []KeySection, th theme.Theme) *helpOverlay {
	return &helpOverlay{
		sections: sections,
		th:       th,
	}
}

func (h *helpOverlay) Update(msg tea.Msg) (Overlay, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "?", "esc", "q":
			h.done = true
		}
	}
	return h, nil
}

func (h *helpOverlay) View() string {
	var sb strings.Builder

	titleStyle := h.th.Title()
	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(h.th.Primary())).
		Bold(true)
	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(h.th.Text()))
	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(h.th.Secondary())).
		Bold(true)

	sb.WriteString(titleStyle.Render("Help"))
	sb.WriteString("\n\n")

	for i, section := range h.sections {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(sectionStyle.Render(section.Title))
		sb.WriteString("\n")
		for _, key := range section.Keys {
			sb.WriteString(fmt.Sprintf("  %s  %s\n", keyStyle.Render(key.Key), descStyle.Render(key.Desc)))
		}
	}

	return h.th.Panel().Render(sb.String())
}

func (h *helpOverlay) Done() bool {
	return h.done
}

// --- Confirm overlay adapter ---

// confirmOverlay wraps a components.Confirm as an Overlay.
type confirmOverlay struct {
	confirm *components.Confirm
}

// ConfirmOverlay creates an Overlay from a Confirm component.
func ConfirmOverlay(c *components.Confirm) Overlay {
	return &confirmOverlay{confirm: c}
}

func (o *confirmOverlay) Update(msg tea.Msg) (Overlay, tea.Cmd) {
	updated, cmd := o.confirm.Update(msg)
	o.confirm = updated
	return o, cmd
}

func (o *confirmOverlay) View() string {
	return o.confirm.View()
}

func (o *confirmOverlay) Done() bool {
	return o.confirm.Done()
}
