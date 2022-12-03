package pages

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"me.kryptk.marco/models/widgets"
	"me.kryptk.marco/repository"
	"me.kryptk.marco/services"
	"strings"
	"time"
)

var btTitle = titleStyle.Render("Bluetooth")

type Bluetooth struct {
	Selected repository.BluetoothDevice
	service  repository.Bluetooth
	list     list.Model
	hideBar  bool
	bar      widgets.Bar
}

func (bt Bluetooth) Name() string {
	return "Bluetooth"
}

func (bt Bluetooth) Close() error {
	return nil
}

func NewBluetooth() Bluetooth {
	bt := Bluetooth{service: services.NewBluez()}

	bt.hideBar = true
	bt.bar = widgets.Bar{}
	bt.list = list.New([]list.Item{}, btDelegate{}, 0, 0)
	bt.list.SetShowTitle(false)
	bt.list.SetFilteringEnabled(false)
	bt.list.SetShowStatusBar(false)

	adapter, err := bt.service.GetAdapter()
	if err != nil {
		return bt
	}

	devices, err := adapter.GetDevices()
	if err != nil {
		return bt
	}

	bt.list.SetItems(devicesToItems(devices))
	return bt
}

func (bt Bluetooth) Init() tea.Cmd {
	return nil
}

func (bt Bluetooth) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			adapter, _ := bt.service.GetAdapter()
			_ = adapter.StartDiscovery()
			devices, _ := adapter.GetDevices()

			cmds = append(cmds, bt.list.SetItems(devicesToItems(devices)))
		case "enter":
			curItem := bt.list.SelectedItem().(btItem)
			if curItem.Connected() {
				dBar := func() tea.Msg {
					return networkMsg{
						bar: widgets.NewBar(widgets.BarOpts{
							Message:   fmt.Sprintf("Disconnect from %s", curItem.Name),
							InputType: widgets.InputTypeChoice,
						}),
					}
				}

				bt.Selected = curItem.Device
				cmds = append(cmds, dBar)
				break
			}

			connectingBar := func() tea.Msg {
				return networkMsg{
					bar: widgets.NewBar(widgets.BarOpts{
						Message: fmt.Sprintf("Connecting to %s", curItem.Name),
					}),
				}
			}

			connect := func() tea.Msg {
				err := curItem.Device.Connect()
				if err != nil {
					return networkMsg{
						hideBar: true,
						timeout: time.Second * 3,
						bar: widgets.NewBar(widgets.BarOpts{
							Message: fmt.Sprintf("Error connecting : %v", err),
						}),
					}
				}

				return networkMsg{
					bar: widgets.NewBar(widgets.BarOpts{
						Message: fmt.Sprintf("Successfully connected to %s", curItem.Name),
					}),
				}
			}

			cmds = append(cmds, connectingBar, connect)
		}
	case tea.WindowSizeMsg:
		bt.list.SetSize(msg.Width, msg.Height-lipgloss.Height(btTitle)-lipgloss.Height(bt.bar.View()))
		bt.list.SetDelegate(btDelegate{Width: msg.Width - activeTabStyle.GetHorizontalFrameSize() - 6})
	case networkMsg:
		if msg.bar.Message != "" {
			bt.bar = msg.bar
		}

		cmds = append(cmds, func() tea.Msg {
			time.Sleep(msg.timeout)
			return hideBarMsg(msg.hideBar)
		})
	case hideBarMsg:
		bt.hideBar = bool(msg)
	case widgets.BarMsg[bool]:
		if msg.Output {
			err := bt.Selected.Disconnect()
			var bar widgets.Bar
			if err != nil {
				bar = widgets.NewBar(widgets.BarOpts{
					Message: "Disconnected",
				})
			} else {
				bar = widgets.NewBar(widgets.BarOpts{
					Message: fmt.Sprintf("Error disconnecting, %v", err),
				})
			}

			cmds = append(cmds, func() tea.Msg {
				return networkMsg{
					hideBar: true,
					timeout: time.Second * 3,
					bar:     bar,
				}
			})

		} else {
			cmds = append(cmds, func() tea.Msg {
				return hideBarMsg(true)
			})
		}
	}

	var cmd tea.Cmd
	bt.list, cmd = bt.list.Update(msg)

	cmds = append(cmds, cmd)

	return bt, tea.Batch(cmds...)
}

func (bt Bluetooth) View() string {
	renderStr := []string{
		btTitle,
		bt.list.View(),
	}

	if !bt.hideBar {
		renderStr = append(renderStr, bt.bar.View())
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		renderStr...,
	)
}

type btItem struct {
	Device     repository.BluetoothDevice
	Name       string
	Connected  func() bool
	Paired     bool
	Percentage uint
}

func (b btItem) FilterValue() string {
	return b.Name
}

type btDelegate struct {
	Width int
}

func (b btDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	selected := index == m.Index()

	it, ok := item.(btItem)

	if !ok {
		return
	}

	wh := lipgloss.Width
	name := strings.Trim(it.Name, "\n")

	state := "not set up"

	connected := it.Connected()
	if connected {
		state = "connected"
	}

	if it.Paired && !connected {
		state = "disconnected"
	}

	state = lipgloss.NewStyle().Faint(true).Bold(false).Render(state)

	itemStr := lipgloss.JoinVertical(
		lipgloss.Top,
		name,
		state,
	)

	empty := lipgloss.NewStyle().Width(b.Width - wh(itemStr)).String()

	status := []string{
		state,
	}

	if connected {
		status = append(status, fmt.Sprintf("%v%%", int(it.Percentage)))
	}

	statusText := lipgloss.JoinHorizontal(
		lipgloss.Left,
		status...,
	)

	iString := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Top,
			name,
			statusText,
		),
		empty,
	)

	if selected {
		iString = activeTabStyle.Render(iString)
	} else {
		iString = tabStyle.Render(iString)
	}

	_, _ = fmt.Fprintf(w, iString)
}

func (b btDelegate) Height() int {
	return 2
}

func (b btDelegate) Spacing() int {
	return 1
}

func (b btDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func devicesToItems(devices []repository.BluetoothDevice) []list.Item {
	items := make([]list.Item, len(devices))

	for i, d := range devices {
		items[i] = btItem{
			Device:     d,
			Name:       d.GetAlias(),
			Connected:  d.IsConnected,
			Paired:     d.IsPaired(),
			Percentage: d.GetBatteryPercentage(),
		}
	}

	return items
}
