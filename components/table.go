package components

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/NCarteazy/bubbletea-foundation/theme"
)

const columnGap = 2

// Column describes a single table column.
type Column struct {
	Name  string
	Width int
	Align lipgloss.Position
}

// Table is a sortable, scrollable table with configurable columns.
type Table struct {
	columns  []Column
	rows     [][]string
	cursor   int
	offset   int
	height   int
	width    int
	sortCol  int
	sortDesc bool
	th       theme.Theme
}

// NewTable creates a new Table with the given columns, rows, and theme.
func NewTable(columns []Column, rows [][]string, th theme.Theme) Table {
	return Table{
		columns:  columns,
		rows:     rows,
		cursor:   0,
		offset:   0,
		sortCol:  0,
		sortDesc: true,
		th:       th,
	}
}

// Update handles key messages for table navigation and sorting.
func (t Table) Update(msg tea.Msg) (Table, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.Type == tea.KeyUp || (msg.Type == tea.KeyRunes && string(msg.Runes) == "k"):
			if t.cursor > 0 {
				t.cursor--
			}
		case msg.Type == tea.KeyDown || (msg.Type == tea.KeyRunes && string(msg.Runes) == "j"):
			if len(t.rows) > 0 && t.cursor < len(t.rows)-1 {
				t.cursor++
			}
		case msg.Type == tea.KeyHome || (msg.Type == tea.KeyRunes && string(msg.Runes) == "g"):
			t.cursor = 0
		case msg.Type == tea.KeyEnd || (msg.Type == tea.KeyRunes && string(msg.Runes) == "G"):
			if len(t.rows) > 0 {
				t.cursor = len(t.rows) - 1
			}
		case msg.Type == tea.KeyRunes && string(msg.Runes) == "s":
			t.cycleSort()
		}
		t.clampCursor()
		t.scrollToView()
	}
	return t, nil
}

// View renders the table as a string.
func (t Table) View() string {
	var b strings.Builder

	// Render header.
	b.WriteString(t.renderHeader())
	b.WriteByte('\n')

	// Determine visible rows based on height (subtract 1 for the header).
	visibleRows := t.height - 1
	if visibleRows < 0 {
		visibleRows = 0
	}

	for i := t.offset; i < len(t.rows) && i < t.offset+visibleRows; i++ {
		row := t.rows[i]
		b.WriteString(t.renderRow(row, i))
		if i < t.offset+visibleRows-1 && i < len(t.rows)-1 {
			b.WriteByte('\n')
		}
	}

	return b.String()
}

// SelectedRow returns the current cursor position.
func (t Table) SelectedRow() int {
	return t.cursor
}

// RowCount returns the number of rows in the table.
func (t Table) RowCount() int {
	return len(t.rows)
}

// SetRows replaces the table's row data.
func (t *Table) SetRows(rows [][]string) {
	t.rows = rows
	t.clampCursor()
}

// SetSize sets the visible width and height of the table.
func (t *Table) SetSize(width, height int) {
	t.width = width
	t.height = height
}

// renderHeader builds the styled header line with sort indicators.
func (t Table) renderHeader() string {
	headerStyle := t.th.TableHeader()
	parts := make([]string, len(t.columns))
	for i, col := range t.columns {
		name := col.Name
		if i == t.sortCol {
			if t.sortDesc {
				name += " v"
			} else {
				name += " ^"
			}
		}
		parts[i] = alignText(truncate(name, col.Width), col.Width, col.Align)
	}
	line := strings.Join(parts, strings.Repeat(" ", columnGap))
	return headerStyle.Render(line)
}

// renderRow builds a single styled row line.
func (t Table) renderRow(row []string, idx int) string {
	parts := make([]string, len(t.columns))
	for i, col := range t.columns {
		cell := ""
		if i < len(row) {
			cell = row[i]
		}
		parts[i] = alignText(truncate(cell, col.Width), col.Width, col.Align)
	}
	line := strings.Join(parts, strings.Repeat(" ", columnGap))

	if idx == t.cursor {
		return t.th.SelectedRow().Render(line)
	}
	// Alternating row shading: even-indexed rows get surface background.
	if idx%2 == 0 {
		style := lipgloss.NewStyle().
			Background(lipgloss.Color(t.th.Surface())).
			Foreground(lipgloss.Color(t.th.Text()))
		return style.Render(line)
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.th.Text())).
		Render(line)
}

// cycleSort cycles through columns and toggles sort direction.
func (t *Table) cycleSort() {
	if t.sortDesc {
		// Toggle to ascending on the same column.
		t.sortDesc = false
	} else {
		// Move to next column, reset to descending.
		t.sortCol = (t.sortCol + 1) % len(t.columns)
		t.sortDesc = true
	}
	t.sortRows()
}

// sortRows sorts the rows by the current sort column.
func (t *Table) sortRows() {
	col := t.sortCol
	desc := t.sortDesc
	sort.SliceStable(t.rows, func(i, j int) bool {
		a := cellVal(t.rows[i], col)
		b := cellVal(t.rows[j], col)

		aNum, aOk := parseNumeric(a)
		bNum, bOk := parseNumeric(b)

		var less bool
		if aOk && bOk {
			less = aNum < bNum
		} else {
			less = strings.ToLower(a) < strings.ToLower(b)
		}

		if desc {
			return !less
		}
		return less
	})
}

// clampCursor ensures cursor is within valid range.
func (t *Table) clampCursor() {
	if len(t.rows) == 0 {
		t.cursor = 0
		return
	}
	if t.cursor < 0 {
		t.cursor = 0
	}
	if t.cursor >= len(t.rows) {
		t.cursor = len(t.rows) - 1
	}
}

// scrollToView adjusts offset so cursor is visible.
func (t *Table) scrollToView() {
	visibleRows := t.height - 1
	if visibleRows <= 0 {
		return
	}
	if t.cursor < t.offset {
		t.offset = t.cursor
	}
	if t.cursor >= t.offset+visibleRows {
		t.offset = t.cursor - visibleRows + 1
	}
}

// truncate cuts a string to fit within maxWidth.
func truncate(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= maxWidth {
		return s
	}
	if maxWidth <= 1 {
		return string(r[:maxWidth])
	}
	return string(r[:maxWidth-1]) + "~"
}

// alignText pads or aligns a string within the given width.
func alignText(s string, width int, align lipgloss.Position) string {
	r := []rune(s)
	if len(r) >= width {
		return s
	}
	pad := width - len(r)
	switch align {
	case lipgloss.Right:
		return strings.Repeat(" ", pad) + s
	case lipgloss.Center:
		left := pad / 2
		right := pad - left
		return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
	default: // Left
		return s + strings.Repeat(" ", pad)
	}
}

// cellVal safely retrieves a cell value from a row.
func cellVal(row []string, col int) string {
	if col < len(row) {
		return row[col]
	}
	return ""
}

// sizeMultipliers maps suffixes to byte values for numeric-aware sorting.
var sizeMultipliers = map[string]float64{
	"B":  1,
	"KB": 1e3,
	"MB": 1e6,
	"GB": 1e9,
	"TB": 1e12,
	"PB": 1e15,
}

// sizePattern matches values like "1.5 GB" or "300 MB".
var sizePattern = regexp.MustCompile(`^([\d.]+)\s*(B|KB|MB|GB|TB|PB)$`)

// percentPattern matches values like "45.2%".
var percentPattern = regexp.MustCompile(`^([\d.]+)\s*%$`)

// parseNumeric attempts to interpret a string as a numeric value.
// It handles plain numbers, sizes (e.g., "1.5 GB"), percentages, and
// comma-separated numbers (e.g., "1,234").
func parseNumeric(s string) (float64, bool) {
	s = strings.TrimSpace(s)

	// Try size format: "1.5 GB"
	if m := sizePattern.FindStringSubmatch(s); m != nil {
		val, err := strconv.ParseFloat(m[1], 64)
		if err == nil {
			if mult, ok := sizeMultipliers[m[2]]; ok {
				return val * mult, true
			}
		}
	}

	// Try percentage: "45.2%"
	if m := percentPattern.FindStringSubmatch(s); m != nil {
		val, err := strconv.ParseFloat(m[1], 64)
		if err == nil {
			return val, true
		}
	}

	// Try comma-separated number: "1,234"
	cleaned := strings.ReplaceAll(s, ",", "")
	if val, err := strconv.ParseFloat(cleaned, 64); err == nil {
		return val, true
	}

	return 0, false
}
