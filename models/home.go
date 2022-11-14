package  models

import tea "github.com/charmbracelet/bubbletea"

type Home struct {
    Pages []tea.Model
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
    return "Home"
}

func NewHome() Home {
    return Home{
        Pages : []tea.Model{},
        Selected : 0,
    }
}