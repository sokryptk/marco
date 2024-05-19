package pages

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"io"
	"me.kryptk.marco/models/components"
	"me.kryptk.marco/repository"
	"me.kryptk.marco/services"
	"me.kryptk.marco/utils"
	"time"
)

var title = titleStyle.Render("Wireless")

const disconnectID string = "disconnect"

type Network struct {
	width, height int
	selected      repository.AccessPoint
	state         int
	list          list.Model
	overlay       *components.Overlay
	Service       repository.Network
}

func NewNetwork() Network {
	network := Network{Service: services.NewNM()}
	_ = network.Service.GetDevices()

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

	if w.overlay != nil && w.overlay.Visible() {
		var cmd tea.Cmd
		w.overlay, cmd = w.overlay.Update(msg)

		switch msg := msg.(type) {
		case components.ShowOverlayMsg:
			var cmd tea.Cmd
			w.overlay = w.overlay.Close()
			w.overlay, cmd = msg.Overlay.Update(msg)
			cmds = append(cmds, w.overlay.Init(), cmd)
			return w, tea.Batch(cmds...)
		}

		form, ok := w.overlay.GetChild().(*huh.Form)
		if !ok {
			return w, cmd
		}

		if form.State == huh.StateCompleted {
			w.overlay = w.overlay.Close()

			w.overlay, cmd = components.NewOverlay(
				ProgressDialog("Connecting", false),
				0,
				components.OverlayViewProps{
					X: 0.95,
					Y: 0.95,
				},
			)

			cmds = append(cmds, cmd)

			bgProcess := func() tea.Msg {
				status := w.selected.Connect(repository.ConnectOptions{Password: utils.Ptr[string](form.Get("pass").(string))})

				if status == repository.ConnectionStatusActivated {
					_, cmd = components.NewOverlay(
						ProgressDialog(fmt.Sprintf("Connected to %s", w.selected.GetSSID()), true),
						time.Second*5,
						components.OverlayViewProps{
							X: 0.5,
							Y: 0.5,
						},
					)

					return cmd()
				}

				_, cmd = components.NewOverlay(
					ProgressDialog(fmt.Sprintf("Failed to connect to %s", w.selected.GetSSID()), true),
					time.Second*5,
					components.OverlayViewProps{
						X: 0.5,
						Y: 0.5,
					},
				)

				return cmd()
			}

			cmds = append(cmds, bgProcess)
			return w, tea.Batch(cmds...)
		}

		cmds = append(cmds, cmd)
		return w, tea.Batch(cmds...)
	}

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
		case "enter":
			w.selected = w.list.SelectedItem().(item).AccessPoint

			var cmd tea.Cmd

			w.overlay, cmd = components.NewOverlay(
				ConnectDialog(w.selected),
				0,
				components.OverlayViewProps{
					X: 0.5,
					Y: 0.5,
				},
			)

			return w, cmd
		}

		var cmd tea.Cmd
		w.list, cmd = w.list.Update(msg)

		cmds = append(cmds, cmd)
	case tea.WindowSizeMsg:
		w.list.SetSize(msg.Width-2, msg.Height-lipgloss.Height(title))
		w.list.SetDelegate(itemDelegate{Width: msg.Width - activeTabStyle.GetHorizontalFrameSize() - 8})
	case components.ShowOverlayMsg:
		var cmd tea.Cmd
		w.overlay = w.overlay.Close()
		w.overlay, cmd = msg.Overlay.Update(msg)
		cmds = append(cmds, w.overlay.Init(), cmd)
	}

	return w, tea.Batch(cmds...)
}

func (w Network) View() string {
	renderList := []string{
		title,
		w.list.View(),
	}

	main := lipgloss.JoinVertical(
		lipgloss.Left,
		renderList...,
	)

	if w.overlay != nil && w.overlay.Visible() {
		return w.overlay.SetBg(main).View()
	}
	return main
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

func ConnectDialog(ap repository.AccessPoint) *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title(fmt.Sprintf("Connecting to %s", ap.GetSSID())),
			huh.NewInput().Key("pass").Title("Password").Password(true),
			huh.NewConfirm().
				Key("connect").
				Title("All done?").
				Validate(func(v bool) error {
					if !v {
						return fmt.Errorf("Welp, finish up then")
					}
					return nil
				}).
				Affirmative("Connect").
				Negative("Cancel"),
		),
	).WithWidth(60).WithHeight(20)

	return form
}

func ProgressDialog(progress string, static bool) *spinner.Spinner {
	return spinner.New().Accessible(static).Type(spinner.Globe).Title(progress)
}
