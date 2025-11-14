package client

import tea "github.com/charmbracelet/bubbletea"

func (c Client) Init() tea.Cmd {
	return nil
}

func (c Client) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return c, tea.Quit
		}
	}
	return c, nil
}

func (c Client) View() string {
	var view string
	return view
}
