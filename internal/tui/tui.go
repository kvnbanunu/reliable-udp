package tui

import (
	"reliable-udp/internal/utils"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Settings for the progress bar
func NewProgress() progress.Model {
	return progress.New(
		progress.WithDefaultGradient(),
		progress.WithoutPercentage(),
	)
}

func NewTextInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Hello World"
	ti.CharLimit = int(utils.MAX_PAYLOAD_LEN)
	ti.Width = 20

	return ti
}

// Generic Run command to start the render program
func Run(m tea.Model) error {
	p := tea.NewProgram(m)
	_, err := p.Run()
	if err != nil {
		return err
	}
	return nil
}
