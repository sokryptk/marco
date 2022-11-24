package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"math"
	"me.kryptk.marco/models/pages"
)

var body = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 3).Margin(1, 3)

type teaModelWithName interface {
	tea.Model
	Name() string
}

type Home struct {
	width, height int
	paneState     int
	Pages         []teaModelWithName
	Sidebar       Sidebar
	Content       Content
	Selected      int
}

func (h Home) Init() tea.Cmd {
	return nil
}

func (h Home) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmd := make([]tea.Cmd, 2)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "enter":
			return h, tea.Quit
		case "tab":
			h.paneState = (h.paneState + 1) % 2
		}
	case tea.WindowSizeMsg:
		h.width, h.height = msg.Width, msg.Height
		sW := lipgloss.Width(h.Sidebar.View())
		h.Content.width = h.width - sW - body.GetHorizontalFrameSize()
		h.Content.height = h.height - body.GetVerticalFrameSize()
	}

	switch h.paneState {
	case 0:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up":
				selected := int(math.Abs(float64((h.Selected - 1) % len(h.Pages))))
				h.Selected = selected
				h.Sidebar.selected = selected
				h.Content.model = h.Pages[selected]
			case "down":
				selected := (h.Selected + 1) % len(h.Pages)
				h.Selected = selected
				h.Sidebar.selected = selected
				h.Content.model = h.Pages[selected]
			}
		}

		h.Sidebar, cmd[0] = h.Sidebar.Update(msg)
	case 1:
		h.Content, cmd[1] = h.Content.Update(msg)
	}

	return h, tea.Batch(cmd...)
}

func (h Home) View() string {
	// The height and the width of the body doesn't take into account the borders and the margins
	// The Height and Width are set of the inner connect of the body.
	// Therefore, subtract the margins, borders before setting the size of the body
	w, he := body.GetHorizontalMargins()+body.GetHorizontalBorderSize(), body.GetVerticalMargins()+body.GetVerticalBorderSize()
	body = body.Width(h.width - w).Height(h.height - he)

	selectedTab := lipgloss.NewStyle().Faint(true)
	sideView := h.Sidebar.View()
	contentView := h.Content.View()

	if h.paneState != 0 {
		sideView = selectedTab.Render(sideView)
	} else {
		contentView = selectedTab.Render(contentView)
	}

	layout := lipgloss.JoinHorizontal(lipgloss.Top, sideView, contentView)

	return body.Render(layout)
}

func NewHome() Home {
	withNames := []teaModelWithName{
		pages.NewNetwork(),
	}

	return Home{
		Pages: withNames,
		Sidebar: Sidebar{
			items: withNames,
		},
		Content: Content{
			model: withNames[0],
		},
	}
}
