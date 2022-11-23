package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"me.kryptk.marco/models/pages"
)

var body = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 3).Margin(1, 3)

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
	// The height and the width of the body doesn't take into account the borders and the margins
	// The Height and Width are set of the inner connect of the body.
	// Therefore, subtract the margins, borders before setting the size of the body
	w, he := body.GetHorizontalMargins()+body.GetHorizontalBorderSize(), body.GetVerticalMargins()+body.GetVerticalBorderSize()
	body = body.Width(h.width - w).Height(h.height - he)

	sidebar := Sidebar{
		items:    h.Pages,
		selected: h.Selected,
	}

	sW := lipgloss.Width(sidebar.View())

	content := Content{
		width:  h.width - sW - body.GetHorizontalFrameSize(),
		height: h.height - body.GetVerticalFrameSize(),
		model:  h.Pages[h.Selected],
	}

	layout := lipgloss.JoinHorizontal(lipgloss.Top, sidebar.View(), content.View())

	return body.Render(layout)
}

func NewHome() Home {
	return Home{
		Pages:    []teaModelWithName{pages.NewNetwork(), pages.NewNetwork()},
		Selected: 0,
	}
}
