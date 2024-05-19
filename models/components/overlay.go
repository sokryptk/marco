package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"me.kryptk.marco/utils"
	"time"
)

// Overlay is a simple over the screen modal
type Overlay struct {
	child     tea.Model
	visible   bool
	viewProps OverlayViewProps
}

type ShowOverlayMsg struct {
	Overlay Overlay
	Timeout time.Duration
}

type CloseOverlayMsg struct{}

type OverlayViewProps struct {
	Background string
	X          lipgloss.Position
	Y          lipgloss.Position
}

func NewOverlay(child tea.Model, duration time.Duration, view OverlayViewProps) (*Overlay, tea.Cmd) {
	var cmds []tea.Cmd

	overlay := Overlay{
		visible:   true,
		child:     child,
		viewProps: view,
	}

	cmd := func() tea.Msg {
		return ShowOverlayMsg{Overlay: overlay}
	}

	cmds = append(cmds, cmd, overlay.child.Init())

	if duration > 0 {
		cd := func() tea.Msg {
			time.Sleep(duration)
			return CloseOverlayMsg{}
		}

		cmds = append(cmds, cd)
	}

	return &overlay, tea.Batch(cmds...)
}

func (o *Overlay) Init() tea.Cmd {
	return o.child.Init()
}

func (o *Overlay) Update(msg tea.Msg) (*Overlay, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.(type) {
	case CloseOverlayMsg:
		o.Close()
	}

	o.child, cmd = o.child.Update(msg)

	return o, cmd
}

func (o *Overlay) GetChild() tea.Model {
	return o.child
}

func (o *Overlay) SetView(props OverlayViewProps) *Overlay {
	o.viewProps = props
	return o
}

func (o *Overlay) Close() *Overlay {
	o.visible = false
	return o
}

func (o *Overlay) SetBg(bg string) *Overlay {
	o.viewProps.Background = bg
	return o
}

func (o *Overlay) Visible() bool {
	return o.visible
}

func (o *Overlay) View() string {
	width, height := lipgloss.Width(o.viewProps.Background), lipgloss.Height(o.viewProps.Background)
	cWidth, cHeight := lipgloss.Width(o.child.View()), lipgloss.Height(o.child.View())

	overlay := lipgloss.NewStyle().Border(lipgloss.ThickBorder()).BorderForeground(lipgloss.Color("43")).Render(o.child.View())

	insertAtX := int(float64(width)*float64(o.viewProps.X) - float64(cWidth/2))
	insertAtY := int(float64(height)*float64(o.viewProps.Y) - float64(cHeight)/2)
	return utils.UI.PlaceOverlay(insertAtX, insertAtY, overlay, o.viewProps.Background)
}
