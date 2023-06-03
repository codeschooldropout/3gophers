package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/codeschooldropout/3gophers/internal/signals"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyOrder     = "order"
	columnKeyContracts = "contracts"
	columnKeyPrice     = "price"
	columnKeyTicker    = "ticker"
	columnKeyInterval  = "interval"
	columnKeyPosition  = "position"
	columnKeyTimeNow   = "timenow"
	columnKeyAsset     = "asset"
)

var (
	styleSubtle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888"))

	styleBase = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#fff")).
			BorderForeground(lipgloss.Color("#888")).
			Align(lipgloss.Right)
)

type Model struct {
	alertTable table.Model
}

func makeRow(alert signals.Alert) table.Row {
	return table.NewRow(table.RowData{
		columnKeyOrder:     alert.Order,
		columnKeyContracts: alert.Contracts,
		columnKeyPrice:     alert.Price,
		columnKeyTicker:    alert.Ticker,
		columnKeyInterval:  alert.Interval,
		columnKeyPosition:  alert.Position,
		columnKeyTimeNow:   alert.TimeNow,
		columnKeyAsset:     alert.Asset,
	})
}

func NewModel() Model {
	return Model{
		alertTable: table.New([]table.Column{
			table.NewColumn(columnKeyOrder, "Order", 13),
			table.NewColumn(columnKeyContracts, "Contracts", 13),
			table.NewColumn(columnKeyPrice, "Price", 13),
			table.NewColumn(columnKeyTicker, "Ticker", 13),
			table.NewColumn(columnKeyInterval, "Interval", 13),
			table.NewColumn(columnKeyPosition, "Position", 13),
			table.NewColumn(columnKeyTimeNow, "TimeNow", 13),
		}).WithRows([]table.Row{
			makeRow(*signals.NewAlert("", 0, 0, "", 0, 0, "", *signals.NewAsset("", "", "", ""))),
		}).
			BorderRounded().
			WithBaseStyle(styleBase).
			WithPageSize(6).
			Focused(true),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.alertTable, cmd = m.alertTable.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)
		}
	}

	return m, tea.Batch(cmds...)

}

func (m Model) View() string {
	selected := m.alertTable.HighlightedRow().Data[columnKeyOrder].(string)
	view := lipgloss.JoinVertical(
		lipgloss.Left,
		styleSubtle.Render("press q/esc/ctrl+c to quit"),
		styleSubtle.Render("Hilighted: "+selected),
		m.alertTable.View(),
	) + "\n"
	return lipgloss.NewStyle().UnsetMarginLeft().Render(view)
}
