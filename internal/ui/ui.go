package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/codeschooldropout/3gophers/internal/signals"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyCall            = "call"
	columnKeyPosition        = "position"
	columnKeyPrice           = "price"
	columnKeyPNL             = "pnl"
	columnKeyBars            = "bars"
	columnKeyStopLoss        = "sl"
	columnKeyStopLossPercent = "slp"
	columnKeyExchange        = "exchange"
	columnKeyBase            = "base"
	columnKeyQuote           = "quote"
	columnKeyTF              = "tf"
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
		columnKeyCall:            alert.Call,
		columnKeyPosition:        alert.Position,
		columnKeyPrice:           alert.Price,
		columnKeyPNL:             alert.PNL,
		columnKeyBars:            alert.Bars,
		columnKeyStopLoss:        alert.StopLoss,
		columnKeyStopLossPercent: alert.StopLossPercent,
		columnKeyExchange:        alert.Asset.Exchange,
		columnKeyBase:            alert.Asset.Base,
		columnKeyQuote:           alert.Asset.Quote,
		columnKeyTF:              alert.Asset.Timeframe,
	})
}

func NewModel() Model {
	return Model{
		alertTable: table.New([]table.Column{
			table.NewColumn(columnKeyCall, "Call", 13),
			table.NewColumn(columnKeyPosition, "Position", 10),
			table.NewColumn(columnKeyPrice, "Price", 5),
			table.NewColumn(columnKeyPNL, "PNL", 5),
			table.NewColumn(columnKeyBars, "Bars", 5),
			table.NewColumn(columnKeyStopLoss, "SL", 5),
			table.NewColumn(columnKeyStopLossPercent, "SLP", 5),
			table.NewColumn(columnKeyExchange, "Exchange", 10),
			table.NewColumn(columnKeyBase, "Base", 10),
			table.NewColumn(columnKeyQuote, "Quote", 10),
			table.NewColumn(columnKeyTF, "TF", 10),
		}).WithRows([]table.Row{
			makeRow(*signals.NewAlert("exit", "short", 63.44, 0.78, 49, 0, 0, *signals.NewAsset("COINBASE", "ATOM", "USD", "5m"))),
			makeRow(*signals.NewAlert("exit", "short", 63.44, 0.78, 49, 0, 0, *signals.NewAsset("COINBASE", "ATOM", "USD", "5m"))),
			makeRow(*signals.NewAlert("exit", "short", 63.44, 0.78, 49, 0, 0, *signals.NewAsset("COINBASE", "ATOM", "USD", "5m"))),
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
	selected := m.alertTable.HighlightedRow().Data[columnKeyCall].(string)
	view := lipgloss.JoinVertical(
		lipgloss.Left,
		styleSubtle.Render("press q/esc/ctrl+c to quit"),
		styleSubtle.Render("Hilighted: "+selected),
		m.alertTable.View(),
	) + "\n"
	return lipgloss.NewStyle().UnsetMarginLeft().Render(view)
}
