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
	inputTypeNone InputType = iota
	inputTypeChoice
	inputTypeText
	inputTypePassword
)

type OutputType interface {
	bool | string | types.Nil
}

type BarMsg[T OutputType] struct {
	Output T
}

type Bar struct {
	Message   string
	input     textinput.Model
	InputType InputType
}

func NewBar(message string, inputType InputType) Bar {
	bar := Bar{
		Message:   message,
		InputType: inputType,
	}

	if inputType == inputTypeText || inputType == inputTypePassword {
		bar.input = textinput.New()
	}

	if inputType == inputTypePassword {
		bar.input.EchoMode = textinput.EchoPassword
	}

	return bar
}

func (b Bar) Init() tea.Cmd {
	return nil
}

func (b Bar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch b.InputType {
		case inputTypeChoice:
			return b, func() tea.Msg {
				return BarMsg[bool]{
					Output: msg.Type == tea.KeyEnter,
				}
			}
		case inputTypeText:
			if msg.Type == tea.KeyEnter {
				return b, func() tea.Msg {
					return BarMsg[string]{
						Output: b.input.Value(),
					}
				}
			}

			b.input, _ = b.input.Update(msg)
			return b, nil
		}
	}

	return b, nil
}

func (b Bar) View() string {
	str := fmt.Sprintf("%s", b.Message)

	if b.InputType == inputTypePassword || b.InputType == inputTypeText {
		str += fmt.Sprintf(": %s", b.input.View())
	}

	return lipgloss.NewStyle().Faint(true).PaddingBottom(1).Render(str)
}
