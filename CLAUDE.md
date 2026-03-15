# Bubbletea Foundation

## Build & Test
- `go build ./...` to build all packages
- `go test ./...` to run all tests
- `go test -v ./...` for verbose output
- `go test -race ./...` for race condition detection
- `go vet ./...` for linting

## Run Example
- `cd example && go run .`

## Project Structure
- `app.go` — App orchestrator (Init/Update/View, functional options)
- `view.go` — View interface, ViewContext, navigation messages
- `overlay.go` — Overlay interface and stack management
- `nav.go` — Navigation stack (push/pop/replace)
- `flash.go` — Flash message system with auto-decay
- `breadcrumb.go` — Auto breadcrumb bar from nav stack
- `statusbar.go` — Bottom status bar with key hints
- `theme/` — Theme interface and Tokyo Night default
- `components/` — Optional reusable components (table, confirm, help)
- `layout/` — Layout helpers (panes, centering)
- `example/` — Example app demonstrating all features

## Conventions
- All interfaces are minimal (3 methods or fewer)
- Components are independent — no cross-dependencies
- Theme is injected, never imported directly by components
- Views communicate through the App via tea.Msg, never directly
- Functional options pattern for App configuration
