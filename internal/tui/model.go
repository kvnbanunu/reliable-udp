package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	clientSent int = iota
	clientRecv
	proxySent
	proxyRecv
	serverSent
	serverRecv
)

const (
	padding = 2
	maxWidth = 80
)

func Run() {
	m := model{
		lines: []progress.Model{
			progress.New(progress.WithDefaultGradient()),
			progress.New(progress.WithDefaultGradient()),
			progress.New(progress.WithDefaultGradient()),
			progress.New(progress.WithDefaultGradient()),
			progress.New(progress.WithDefaultGradient()),
			progress.New(progress.WithDefaultGradient()),
		},
	}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error", err)
	}
}

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type model struct {
	// max int
	lines []progress.Model
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		for i := range m.lines {
			m.lines[i].Width = min(msg.Width - padding*2 - 4, maxWidth)
		}
		return m, nil
	case tickMsg:
		var cmds []tea.Cmd
		for i := range m.lines {
			if m.lines[i].Percent() == 1.0 {
				return m, tea.Quit
			}

			cmd := m.lines[i].IncrPercent(float64(i) * 0.01)
			cmds = append(cmds, cmd)
		}
		cmds = append(cmds, tickCmd())
		return m, tea.Batch(cmds...)
	case progress.FrameMsg:
		var cmd tea.Cmd
		var cmds []tea.Cmd
		var pm tea.Model
		for i := range m.lines {
		pm, cmd = m.lines[i].Update(msg)
		m.lines[i] = pm.(progress.Model)
		cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	default:
		return m, nil
	}
}

func (m model) View() string {
	var view string
	pad := strings.Repeat(" ", padding)

	for i := range m.lines {
		view += "\n" + pad + m.lines[i].View()
	}
	view += "\n\n" + helpStyle("Press any key to quit")
	return view
}
