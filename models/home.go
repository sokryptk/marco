package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var body = lipgloss.NewStyle().Border(lipgloss.NormalBorder())

type teaModelWithName interface {
	tea.Model
	Name() string
}

type Home struct {
	width, height int
	Pages         []teaModelWithName
	Selected      int
}

func (h Home) Init() tea.Cmd {
	return nil
}

func (h Home) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "enter":
			return h, tea.Quit
		}
	case tea.WindowSizeMsg:
		h.width, h.height = msg.Width, msg.Height
	}

	return h, nil
}

func (h Home) View() string {
	borderVert, borderHor := body.GetBorderBottomSize()*2, body.GetBorderLeftSize()*2

	sidebar := Sidebar{
		items:    h.Pages,
		selected: h.Selected,
	}

	sW := lipgloss.Width(sidebar.View())

	content := Content{
		width:  h.width - sW,
		height: h.height - borderVert*2,
		model:  h.Pages[h.Selected],
	}

	layout := lipgloss.JoinHorizontal(lipgloss.Top, sidebar.View(), content.View())

	return body.Width(h.width - borderHor).Height(h.height - borderVert).Render(layout)
}

func NewHome() Home {
	return Home{
		Pages:    []teaModelWithName{NewNetwork(), NewNetwork()},
		Selected: 0,
	}
}
