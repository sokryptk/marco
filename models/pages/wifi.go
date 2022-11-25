package pages

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"me.kryptk.marco/repository"
	"me.kryptk.marco/services"
	"time"
)

var titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Padding(1, 0)
var tabStyle = lipgloss.NewStyle().PaddingLeft(1).Bold(true)
var activeTabStyle = tabStyle.Copy().Foreground(lipgloss.Color("5")).Border(bo, false, false, false, true).BorderForeground(lipgloss.Color("5"))
var bo = lipgloss.Border{
	Left: "▐",
}

type Network struct {
	width, height int
	selected      repository.AccessPoint
	list          list.Model
	Service       repository.Network
}

func NewNetwork() Network {
	network := Network{Service: services.NewNM()}
	_ = network.Service.GetDevices()

	network.list = list.New([]list.Item{}, itemDelegate{10}, 0, 0)
	network.list.SetFilteringEnabled(false)
	network.list.SetShowTitle(false)
	network.list.SetShowStatusBar(false)
	return network
}
func (w Network) Init() tea.Cmd {
	return nil
}

func (w Network) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 2)

	cmds = append(cmds, w.list.SetItems(w.getItems()))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			for _, device := range w.Service.GetDevices() {
				switch device := device.(type) {
				case repository.WiFiDevice:
					_ = device.RequestScan()
					time.Sleep(time.Second * 5)
					cmds = append(cmds, w.list.SetItems(w.getItems()))
				}
			}
		}
	case tea.WindowSizeMsg:
		w.list.SetSize(msg.Width, msg.Height-titleStyle.GetVerticalFrameSize())
		w.list.SetDelegate(itemDelegate{Width: msg.Width - activeTabStyle.GetHorizontalFrameSize() - 4})
	}

	var cmd tea.Cmd
	w.list, cmd = w.list.Update(msg)

	cmds = append(cmds, cmd)

	return w, tea.Batch(cmds...)
}

func (w Network) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		titleStyle.Render("Wireless"),
		w.list.View(),
	)
}

func (w Network) Name() string {
	return "Network"
}

func (w Network) getItems() []list.Item {
	items := make([]list.Item, 0)

	devices := w.Service.GetDevices()
	for _, device := range devices {
		switch d := device.(type) {
		case repository.WiFiDevice:
			points := d.GetAccessPoints()

			for _, p := range points {
				items = append(items, item{
					Title:     p.GetSSID(),
					Strength:  renderWifi(p.GetStrength()),
					Frequency: fmt.Sprintf("%0.1fGHz", float64(p.GetFrequency())/float64(1000)),
				})
			}
		}
	}

	return items
}

func renderWifi(strength uint) string {
	switch {
	case strength >= 66:
		return "\uF1EB"
	case strength >= 33:
		return "\uF6AB"
	default:
		return "\uF6AA"
	}
}

type item struct {
	Title     string
	Strength  string
	Frequency string
}

func (i item) FilterValue() string {
	return i.Title
}

type itemDelegate struct {
	Width int
}

func (i itemDelegate) Render(w io.Writer, m list.Model, index int, curItem list.Item) {
	selected := index == m.Index()

	it, ok := curItem.(item)
	if !ok {
		return
	}

	wh := lipgloss.Width
	ssid := it.Title
	strength := it.Strength
	empty := lipgloss.NewStyle().Width(i.Width - wh(ssid) - wh(strength)).String()

	iString := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			ssid,
			empty,
			strength,
		),
		it.Frequency,
	)

	if selected {
		iString = activeTabStyle.Render(iString)
	} else {
		iString = tabStyle.Render(iString)
	}

	_, _ = fmt.Fprintf(w, iString)
}

func (i itemDelegate) Height() int {
	return 2
}

func (i itemDelegate) Spacing() int {
	return 1
}

func (i itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
