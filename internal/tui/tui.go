package tui

import (
	"reliable-udp/internal/packet"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func NewHelpModel() help.Model {
	h := help.New()
	h.ShowAll = false
	return h
}

// Settings for the progress bar
func NewProgressModel() progress.Model {
	return progress.New(
		progress.WithDefaultGradient(),
		progress.WithoutPercentage(),
	)
}

func NewTextInputModel() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Hello World"
	ti.CharLimit = int(packet.MAX_PAYLOAD_LEN)
	ti.Width = 20
	ti.Focus()

	return ti
}

// Generic Run command to start the render program
func Run(m tea.Model) error {
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		return err
	}
	return nil
}
