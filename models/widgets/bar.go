package widgets

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go/types"
)

type InputType int

const (
	InputTypeNone InputType = iota
	InputTypeChoice
	InputTypeText
	InputTypePassword
)

type OutputType interface {
	bool | string | types.Nil
}

type BarMsg[T OutputType] struct {
	ID     string
	Output T
}

type BarOpts Bar

type Bar struct {
	ID        string
	Message   string
	input     textinput.Model
	InputType InputType
}

func NewBar(opts BarOpts) Bar {
	bar := opts

	if opts.InputType == InputTypeText || opts.InputType == InputTypePassword {
		bar.input = textinput.New()
	}

	if opts.InputType == InputTypePassword {
		bar.input.EchoMode = textinput.EchoPassword
		bar.input.Prompt = ""
		bar.input.Focus()
	}

	return Bar(bar)
}

func (b Bar) Init() tea.Cmd {
	return nil
}

func (b Bar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch b.InputType {
		case InputTypeChoice:
			return b, func() tea.Msg {
				return BarMsg[bool]{
					ID:     b.ID,
					Output: msg.Type == tea.KeyEnter,
				}
			}
		case InputTypeText, InputTypePassword:
			if msg.Type == tea.KeyEnter {
				return b, func() tea.Msg {
					return BarMsg[string]{
						ID:     b.ID,
						Output: b.input.Value(),
					}
				}
			}

			var cmd tea.Cmd
			b.input, cmd = b.input.Update(msg)
			return b, cmd
		}
	}

	return b, nil
}

func (b Bar) View() string {
	str := fmt.Sprintf("%s", b.Message)

	if b.InputType == InputTypePassword || b.InputType == InputTypeText {
		str += fmt.Sprintf(": %s", b.input.View())
	}

	return lipgloss.NewStyle().Faint(true).PaddingBottom(1).Render(str)
}
