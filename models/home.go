package  models

import tea "github.com/charmbracelet/bubbletea"

type Home struct {
    Pages []tea.Model
    Selected int
}

func NewHome() Home {
    return Home{
        Pages : []tea.Model{},
        Selected : 0,
    }
}