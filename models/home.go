package  models

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

var tabLayout = lipgloss.NewStyle().
    Foreground(lipgloss.Color("#FFF7DB")).
    Background(lipgloss.Color("#888B7E")).
    Border(lipgloss.NormalBorder()).
    Padding(0, 3).
    MarginTop(1)

var activeTabLayout = tabLayout.Copy().
    Foreground(lipgloss.Color("#FFF7DB")).
    Background(lipgloss.Color("#F25D94")).
    Underline(true)


type teaModelWithName interface {
    tea.Model
    Name() string
}

type Home struct {
    Pages []teaModelWithName
    Selected int
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
    }

    return h, nil
}

func (h Home) View() string {
    var renderString string

    for i, p := range h.Pages {
        if i == 0 {
            renderString += activeTabLayout.Render(p.Name())
        } else {
            renderString += tabLayout.Render(p.Name())
        }
    }

    sideBar :=  lipgloss.JoinVertical(lipgloss.Top, renderString)

    return sideBar
}

func NewHome() Home {
    return Home{
        Pages : []teaModelWithName{NewNetwork(), NewNetwork()},
        Selected : 0,
    }
}
