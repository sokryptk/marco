package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"me.kryptk.marco/repository"
	"me.kryptk.marco/services"
	"strconv"
)

type Network struct {
	Service repository.Network
}

func NewNetwork() Network {
	return Network{
		Service: services.NewNM(),
	}
}
func (w Network) Init() tea.Cmd {
	return nil
}

func (w Network) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return w, nil
}

func (w Network) View() string {
	var render string

	devices := w.Service.GetDevices()
	for _, device := range devices {
		render += "\n"

		switch d := device.(type) {
		case repository.WiFiDevice:
			render += "Wireless"
			points := d.GetAccessPoints()

			d.RequestScan()
			for _, p := range points {
				render += "\n" + p.GetSSID() + " " + strconv.Itoa(int(p.GetStrength()))
			}
		}
	}

	return render
}

func (w Network) Name() string {
	return "Network"
}
