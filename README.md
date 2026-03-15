# Bubbletea Foundation

A Go framework that wraps [Bubble Tea](https://github.com/charmbracelet/bubbletea) to eliminate TUI boilerplate. Define views as a simple interface — the framework handles routing, overlays, resize, breadcrumbs, status bar, flash messages, and theming.

## Install

```bash
go get github.com/NCarteazy/bubbletea-foundation@latest
```

## Quick Start

```go
package main

import (
    "fmt"
    "os"

    tea "github.com/charmbracelet/bubbletea"
    foundation "github.com/NCarteazy/bubbletea-foundation"
    "github.com/NCarteazy/bubbletea-foundation/theme"
)

type homeView struct{}

func (v homeView) ID() string { return "Home" }

func (v homeView) Update(msg tea.Msg, ctx foundation.ViewContext) (foundation.View, tea.Cmd) {
    if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "q" {
        return v, tea.Quit
    }
    return v, nil
}

func (v homeView) Render(ctx foundation.ViewContext) string {
    return ctx.Theme.Title().Render("Hello, World!")
}

func main() {
    app := foundation.New(
        foundation.WithTheme(theme.TokyoNight),
        foundation.WithInitialView("Home"),
        foundation.WithViews(homeView{}),
        foundation.WithStatusHints([]foundation.KeyHint{
            {Key: "q", Desc: "quit"},
        }),
    )
    if _, err := tea.NewProgram(app, tea.WithAltScreen()).Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

## What It Handles

| Feature | How |
|---------|-----|
| View routing | Navigation stack with push/pop/replace |
| Resize | Auto-propagates dimensions to all views via `ViewContext` |
| Overlays | Stack-based modals that block input to views underneath |
| Breadcrumbs | Auto-generated from navigation stack |
| Status bar | Configurable key hints, shows flash messages |
| Flash messages | Timed auto-decay feedback messages |
| Help | Built-in `?` toggle with grouped key binding sections |
| Theming | Injected via `ViewContext`, Tokyo Night default |
| Quit | `ctrl+c` always works |

## The View Interface

Views are fully flexible — render whatever you want:

```go
type View interface {
    ID() string
    Update(msg tea.Msg, ctx ViewContext) (View, tea.Cmd)
    Render(ctx ViewContext) string
}
```

`ViewContext` gives you everything you need:

```go
type ViewContext struct {
    Width  int          // content area width (chrome excluded)
    Height int          // content area height (chrome excluded)
    Theme  theme.Theme  // active color palette and styles
    Data   any          // data passed via Navigate()
}
```

## Navigation

Views navigate by returning commands:

```go
// Push a new view
return v, foundation.Navigate("detail", itemData)

// Go back
return v, foundation.Back()

// Replace current view
return v, foundation.Replace("settings", nil)
```

## Components

Optional building blocks — use them or don't:

### Table

Sortable, scrollable table with keyboard navigation (`j/k`, `g/G`, `s` to sort).

```go
import "github.com/NCarteazy/bubbletea-foundation/components"

tbl := components.NewTable(columns, rows, ctx.Theme)
tbl.SetSize(ctx.Width, ctx.Height)
```

### Confirm Dialog

Generic confirmation overlay with `y/n/esc` handling.

```go
c := components.NewConfirm("Delete this file?", ctx.Theme)
c.OnConfirm = func() tea.Cmd {
    return foundation.Flash("Deleted!", 30)
}
return v, foundation.ShowOverlay(foundation.ConfirmOverlay(c))
```

### Help Overlay

Key binding reference panel, auto-wired with `?`.

```go
foundation.WithHelp([]foundation.KeySection{
    {Title: "Navigation", Keys: []foundation.KeyHint{
        {Key: "j/k", Desc: "move up/down"},
        {Key: "enter", Desc: "select"},
    }},
})
```

## Theme

The `Theme` interface provides colors, pre-built styles, and utilities. Tokyo Night ships as default.

```go
// Use in views via ctx.Theme
title := ctx.Theme.Title().Render("My Title")
bar := ctx.Theme.SizeBar(0.65, 40)
size := ctx.Theme.FormatSize(1073741824) // "1.1 GB"
```

Create custom themes by implementing the `theme.Theme` interface.

## Layout Helpers

```go
import "github.com/NCarteazy/bubbletea-foundation/layout"

// Two-pane split
layout.TwoPane(leftContent, rightContent, width, height)
layout.TwoPaneWithRatio(left, right, 0.3, width, height)

// Center an overlay on a background
layout.OverlayCenter(width, height, background, overlay)
```

## Example

See [example/main.go](example/main.go) for a complete working app with three views, a table, confirm dialog, help overlay, and flash messages.

```bash
cd example && go run .
```

## License

MIT
