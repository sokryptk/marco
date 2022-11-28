package layout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
	"math"
	"me.kryptk.marco/models/pages"
)

var body = lipgloss.NewStyle().Border(lipgloss.ThickBorder()).Padding(1, 3).Margin(1, 3)

type PageSwitchMsg struct {
	new tea.Model
}

type teaModelWithName interface {
	tea.Model
	Name() string
	Close() error
}

type Home struct {
	width, height int
	paneState     int
	Pages         []teaModelWithName
	Sidebar       Sidebar
	Content       Content
	Selected      int
}

func (h Home) Close() {
	for _, page := range h.Pages {
		if err := page.Close(); err != nil {
			log.Printf("err in %s: %e", page.Name(), err)
		}
	}
}

func (h Home) Init() tea.Cmd {
	return nil
}

func (h Home) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	var cmd tea.Cmd
	switch h.paneState {
	case 0:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up":
				selected := int(math.Abs(float64((h.Selected - 1) % len(h.Pages))))
				h.Selected = selected
				h.Sidebar.selected = selected

				h.Content, _ = h.Content.Update(PageSwitchMsg{h.Pages[selected]})
			case "down":
				selected := (h.Selected + 1) % len(h.Pages)
				h.Selected = selected
				h.Sidebar.selected = selected

				h.Content, _ = h.Content.Update(PageSwitchMsg{h.Pages[selected]})
			}
		}

		h.Sidebar, cmd = h.Sidebar.Update(msg)
	case 1:
		h.Content, cmd = h.Content.Update(msg)
	}

	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return h, tea.Quit
		case "tab":
			h.paneState = (h.paneState + 1) % 2
		}
	case tea.WindowSizeMsg:
		h.width, h.height = msg.Width, msg.Height
		sW := lipgloss.Width(h.Sidebar.View())

		sizeMsg := tea.WindowSizeMsg{
			Width:  h.width - sW - body.GetHorizontalFrameSize(),
			Height: h.height - body.GetVerticalFrameSize() - 2,
		}

		var comd tea.Cmd
		h.Content, comd = h.Content.Update(sizeMsg)

		for i, p := range h.Pages {
			page, _ := p.Update(sizeMsg)

			h.Pages[i] = page.(teaModelWithName)
		}

		cmds = append(cmds, comd)
	}

	return h, tea.Batch(cmds...)
}

func (h Home) View() string {
	// The height and the width of the body doesn't take into account the borders and the margins
	// The Height and Width are set of the inner connect of the body.
	// Therefore, subtract the margins, borders before setting the size of the body
	w, he := body.GetHorizontalMargins()+body.GetHorizontalBorderSize(), body.GetVerticalMargins()+body.GetVerticalBorderSize()
	altered := body.Copy().Width(h.width - w).Height(h.height - he)

	selectedTab := lipgloss.NewStyle().Faint(true)
	sideView := h.Sidebar.View()

	contentView := h.Content.View()

	if h.paneState != 0 {
		sideView = selectedTab.Render(sideView)
	} else {
		contentView = selectedTab.Render(contentView)
	}

	layout := lipgloss.JoinHorizontal(lipgloss.Top, sideView, contentView)

	return altered.Render(layout)
}

func NewHome() Home {
	withNames := []teaModelWithName{
		pages.NewNetwork(),
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
