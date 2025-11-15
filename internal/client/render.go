package client

import (
	"reliable-udp/internal/tui"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (c Client) Init() tea.Cmd {
	return nil
}

func (c Client) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, tui.Keys.Quit):
			return c, tea.Quit
		}
	}
	return c, nil
}

func (c Client) View() string {
	var view string
	return view
}
