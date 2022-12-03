package pages

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"me.kryptk.marco/models/widgets"
	"me.kryptk.marco/repository"
	"me.kryptk.marco/services"
	"me.kryptk.marco/utils"
	"time"
)

var title = titleStyle.Render("Wireless")

const disconnectID string = "disconnect"

type networkMsg struct {
	timeout time.Duration
	hideBar bool
	bar     widgets.Bar
}

type hideBarMsg bool

type Network struct {
	width, height int
	selected      repository.AccessPoint
	state         int
	list          list.Model
	bar           tea.Model
	barState      int
	Service       repository.Network
}

func NewNetwork() Network {
	network := Network{Service: services.NewNM()}
	_ = network.Service.GetDevices()

	network.bar = widgets.Bar{}
	network.list = list.New([]list.Item{}, itemDelegate{10}, 0, 0)
	network.list.SetFilteringEnabled(false)
	network.list.SetShowTitle(false)
	network.list.Styles.HelpStyle = network.list.Styles.HelpStyle.PaddingBottom(1)
	network.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "scan the network")),
		}
	}
	network.list.SetShowStatusBar(false)
	return network
}
func (w Network) Init() tea.Cmd {
	return nil
}

func (w Network) Close() error {
	return w.Service.Close()
}

func (w Network) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 2)

	cmds = append(cmds, w.list.SetItems(w.getItems()))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch w.state {
		case 0:
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
			case "enter":
				currentItem := w.list.SelectedItem().(item).AccessPoint
				var netMsg networkMsg

				if !currentItem.IsConnected() {
					netMsg.bar = widgets.NewBar(widgets.BarOpts{Message: "Connecting..."})
					cmds = append(cmds, func() tea.Msg {
						return connectWithBar(w, repository.ConnectOptions{})
					})
				} else {
					netMsg.bar = widgets.NewBar(widgets.BarOpts{
						Message:   fmt.Sprintf("Disconnect from %s?", currentItem.GetSSID()),
						InputType: widgets.InputTypeChoice,
						ID:        disconnectID,
					})
				}

				cmds = append(cmds, func() tea.Msg {
					return netMsg
				})
			}

			var cmd tea.Cmd
			w.list, cmd = w.list.Update(msg)

			cmds = append(cmds, cmd)
		case 1:
			var cmd tea.Cmd
			w.bar, cmd = w.bar.Update(msg)
			cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		w.list.SetSize(msg.Width-2, msg.Height-lipgloss.Height(title)-lipgloss.Height(w.bar.View())-2)
		w.list.SetDelegate(itemDelegate{Width: msg.Width - activeTabStyle.GetHorizontalFrameSize() - 8})
	case widgets.BarMsg[bool]:
		currentItem := w.list.SelectedItem().(item).AccessPoint
		if msg.ID == disconnectID {
			var timeout time.Duration
			var bar widgets.Bar

			if msg.Output {
				timeout = time.Second * 3
			}

			netMsg := networkMsg{
				hideBar: msg.Output,
				timeout: timeout,
				bar:     bar,
			}

			if !msg.Output {
				cmds = append(cmds, func() tea.Msg {
					return netMsg
				})

				break
			}

			err := currentItem.Disconnect()
			if err != nil {
				netMsg.bar = widgets.NewBar(widgets.BarOpts{Message: fmt.Sprintf("Error while disconnecting : %v", err)})
			} else {
				netMsg.bar = widgets.NewBar(widgets.BarOpts{Message: fmt.Sprintf("Disconnected from %s", currentItem.GetSSID())})
			}

			cmds = append(cmds, func() tea.Msg {
				return netMsg
			})
		}

	case widgets.BarMsg[string]:
		cmds = append(cmds, func() tea.Msg {
			return connectWithBar(w, repository.ConnectOptions{Password: utils.Ptr(msg.Output)})
		})
	case networkMsg:
		if msg.bar.Message != "" {
			w.bar = msg.bar
		}

		cmds = append(cmds, func() tea.Msg {
			time.Sleep(msg.timeout)
			return hideBarMsg(msg.hideBar)
		})

	case hideBarMsg:
		if msg {
			w.state = 0
		} else {
			w.state = 1
		}
	}

	return w, tea.Batch(cmds...)
}

func (w Network) View() string {
	renderList := []string{
		title,
		w.list.View(),
	}

	if w.state == 1 {
		renderList = append(renderList, w.bar.View())
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		renderList...,
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
					AccessPoint: p,
					Title:       p.GetSSID(),
					Strength:    renderWifi(p.GetStrength()),
					Frequency:   fmt.Sprintf("%0.1fGHz", float64(p.GetFrequency())/float64(1000)),
					Connected:   p.IsConnected(),
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
	AccessPoint repository.AccessPoint
	Title       string
	Strength    string
	Frequency   string
	Connected   bool
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

	var connectedState string
	if it.Connected {
		connectedState = "connected"
	}

	connected := lipgloss.NewStyle().Faint(true).PaddingLeft(1).Render(connectedState)
	empty := lipgloss.NewStyle().Width(i.Width - wh(ssid) - wh(strength) - wh(connectedState)).String()

	iString := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			ssid,
			connected,
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

func connectWithBar(w Network, options repository.ConnectOptions) networkMsg {
	status := w.list.SelectedItem().(item).AccessPoint.Connect(options)

	switch status {
	case repository.ConnectionStatusNeedAuth:
		return networkMsg{
			bar: widgets.NewBar(widgets.BarOpts{
				Message:   fmt.Sprintf("Password for %s", w.list.SelectedItem().(item).Title),
				InputType: widgets.InputTypePassword,
			}),
		}
	case repository.ConnectionStatusActivated:
		return networkMsg{
			hideBar: true,
			timeout: time.Second * 3,
			bar:     widgets.NewBar(widgets.BarOpts{Message: "Activated connection!"}),
		}
	default:
		return networkMsg{
			hideBar: true,
			timeout: time.Second * 3,
			bar: widgets.NewBar(widgets.BarOpts{
				Message: fmt.Sprintf("Connection failed for %s, err: %d", w.list.SelectedItem().(item).Title, status),
			}),
		}
	}
}
