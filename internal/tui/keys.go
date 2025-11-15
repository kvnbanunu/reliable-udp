package tui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit key.Binding
	Enter key.Binding
}

var Keys = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc/ctrl+c", "quit"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "send message"),
	),
}

// returns mini help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Enter}
}

// returns expanded help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Enter},
	}
}
