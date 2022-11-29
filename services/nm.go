package services

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	"log"
	"me.kryptk.marco/repository"
	"time"
)

const nmInterface = "org.freedesktop.NetworkManager"

type NM struct {
	Conn *dbus.Conn
}

func NewNM() NM {
	conn, err := dbus.ConnectSystemBus()

	if err != nil {
		log.Fatal(err)
		return NM{}
	}

	return NM{
		Conn: conn,
	}
}

type NMDevice struct {
	conn *dbus.Conn
	Path dbus.ObjectPath
}

type NMWifiDevice struct {
	NMDevice
}

type NMEthernetDevice struct {
	NMDevice
}

type NMAccessPoint struct {
	dev  *NMWifiDevice
	Path dbus.ObjectPath
}

func (ap NMAccessPoint) GetSSID() (s string) {
	v, err := ap.dev.conn.Object(nmInterface, ap.Path).GetProperty("org.freedesktop.NetworkManager.AccessPoint.Ssid")
	if err != nil {
		log.Println(err)
	}

	return string(v.Value().([]byte))
}

func (dev NMWifiDevice) RequestScan() error {
	err := dev.conn.Object(nmInterface, dev.Path).Call("org.freedesktop.NetworkManager.Device.Wireless.RequestScan", 0, map[string]interface{}{})
	if err != nil {
		return err.Err
	}

	return nil
}

func (net NM) Close() error {
	return net.Conn.Close()
}

func (dev NMWifiDevice) GetActiveConnection() (repository.AccessPoint, error) {
	objectPath, err := dev.conn.Object(nmInterface, dev.Path).GetProperty("org.freedesktop.NetworkManager.Device.ActiveConnection")
	if err != nil {
		return NMAccessPoint{}, err
	}

	d, err := dev.conn.Object(nmInterface, objectPath.Value().(dbus.ObjectPath)).GetProperty("org.freedesktop.NetworkManager.Connection.Active.Id")
	if err != nil {
		return NMAccessPoint{}, err
	}

	points := dev.GetAccessPoints()
	for _, c := range points {
		if c.GetSSID() == d.Value().(string) {
			return c, nil
		}
	}

	return NMAccessPoint{}, fmt.Errorf("Unknown error occurred")
}

func (dev NMDevice) GetDeviceType() repository.DeviceType {
	rawDeviceType, err := dev.conn.Object(nmInterface, dev.Path).GetProperty("org.freedesktop.NetworkManager.Device.DeviceType")
	if err != nil {
		return repository.DeviceTypeUnknown
	}

	var deviceType uint
	err = rawDeviceType.Store(&deviceType)
	if err != nil {
		return repository.DeviceTypeUnknown
	}

	return repository.DeviceType(deviceType)
}

func (dev NMDevice) GetHwAddress() string {
	udi, err := dev.conn.Object(nmInterface, dev.Path).GetProperty("org.freedesktop.NetworkManager.Device.Interface")
	if err != nil {
		log.Println(err)
		return ""
	}

	return udi.Value().(string)
}

func (net NM) GetDevices() []repository.Device {
	var devicePaths []dbus.ObjectPath
	err := net.Conn.Object(nmInterface, "/org/freedesktop/NetworkManager").Call("org.freedesktop.NetworkManager.GetAllDevices", 0).Store(&devicePaths)
	if err != nil {
		return nil
	}

	devices := make([]repository.Device, len(devicePaths))

	for i, d := range devicePaths {
		device := NMDevice{conn: net.Conn, Path: d}

		switch device.GetDeviceType() {
		case repository.DeviceTypeWifi:
			devices[i] = NMWifiDevice{NMDevice: device}
		case repository.DeviceTypeEthernet:
			devices[i] = NMEthernetDevice{NMDevice: device}
		default:
			devices[i] = device
		}
	}

	return devices
}

func (ap NMAccessPoint) GetFrequency() uint {
	v, err := ap.dev.conn.Object(nmInterface, ap.Path).GetProperty("org.freedesktop.NetworkManager.AccessPoint.Frequency")
	if err != nil {
		log.Println(err)
	}

	var frequency uint
	err = v.Store(&frequency)
	if err != nil {
		log.Println(err)
		return 0
	}

	return frequency
}

func (ap NMAccessPoint) GetStrength() uint {
	v, err := ap.dev.conn.Object(nmInterface, ap.Path).GetProperty("org.freedesktop.NetworkManager.AccessPoint.Strength")
	if err != nil {
		log.Println(err)
	}

	var strength uint
	err = v.Store(&strength)
	if err != nil {
		log.Println(err)
		return 0
	}

	return strength
}

func (dev NMWifiDevice) GetAccessPoints() []repository.AccessPoint {
	var accessPaths []dbus.ObjectPath
	err := dev.conn.Object(nmInterface, dev.Path).Call("org.freedesktop.NetworkManager.Device.Wireless.GetAllAccessPoints", 0).Store(&accessPaths)
	if err != nil {
		return nil
	}

	points := make([]repository.AccessPoint, len(accessPaths))

	for i, p := range accessPaths {
		if p != "" {
			points[i] = NMAccessPoint{
				Path: p,
				dev:  &dev,
			}
		}
	}

	return points
}

func (ap NMAccessPoint) Connect(options repository.ConnectOptions) repository.ConnectionStatus {
	var activeConn dbus.ObjectPath

	err := ap.dev.conn.Object(nmInterface, "/org/freedesktop/NetworkManager").Call(
		"org.freedesktop.NetworkManager.ActivateConnection", 0, dbus.ObjectPath("/"), ap.dev.Path, ap.Path,
	).Store(&activeConn)

	if err != nil {
		state, reason, err := ap.addConnectionAndConnect(options)
		if err != nil {
			log.Println(state, reason, err)

			return repository.ConnectionStatus(state)
		}
	}

	state, reason, err := ap.dev.listenForState()
	if err != nil {
		log.Println(state, reason, err)
		return repository.ConnectionStatus(state)
	}

	if repository.ConnectionStatusErrNoSecrets.Equal(reason) || repository.ConnectionStatusNeedAuth.Equal(state) {
		// Need auth

		settings, err := ap.dev.conn.Object(nmInterface, activeConn).GetProperty("org.freedesktop.NetworkManager.Connection.Active.Connection")
		if err != nil {
			log.Println(err)
			return repository.ConnectionStatusErrFailed
		}

		err = ap.updateConnection(settings.Value().(dbus.ObjectPath), options)
		if err != nil {
			log.Println(err)
			return repository.ConnectionStatusErrFailed
		}

		err = ap.dev.conn.Object(nmInterface, "/org/freedesktop/NetworkManager").Call(
			"org.freedesktop.NetworkManager.ActivateConnection", 0, settings.Value().(dbus.ObjectPath), ap.dev.Path, ap.Path,
		).Err

		if err != nil {
			return repository.ConnectionStatusErrFailed
		}

		state, reason, err := ap.dev.listenForState()
		if err != nil {
			log.Println(state, reason, err)
			return repository.ConnectionStatus(state)
		}

		return repository.ConnectionStatus(state)
	}

	return repository.ConnectionStatus(state)
}

func (ap NMAccessPoint) updateConnection(connectionPath dbus.ObjectPath, options repository.ConnectOptions) error {
	properties := map[string]map[string]interface{}{
		"connection": {
			"id": ap.GetSSID(),
		},
		"802-11-wireless": {
			"mode": "infrastructure",
			"ssid": []byte(ap.GetSSID()),
		},
	}

	if options.Password != nil {
		properties["802-11-wireless-security"] = map[string]interface{}{
			"key-mgmt": "wpa-psk",
			"psk":      options.Password,
		}
	}

	err := ap.dev.conn.Object(nmInterface, connectionPath).Call("org.freedesktop.NetworkManager.Settings.Connection.Update", 0, properties).Err
	if err != nil {
		return err
	}

	return nil
}

func (ap NMAccessPoint) addConnectionAndConnect(options repository.ConnectOptions) (uint32, uint32, error) {
	properties := map[string]map[string]interface{}{
		"connection": {
			"id": ap.GetSSID(),
		},
		"802-11-wireless": {
			"mode": "infrastructure",
			"ssid": []byte(ap.GetSSID()),
		},
	}

	if options.Password != nil {
		properties["802-11-wireless-security"] = map[string]interface{}{
			"key-mgmt": "wpa-psk",
			"psk":      *options.Password,
		}
	}

	err := ap.dev.conn.Object(nmInterface, "/org/freedesktop/NetworkManager").Call("org.freedesktop.NetworkManager.AddAndActivateConnection", 0, properties, ap.dev.Path, ap.Path).Err
	if err != nil {
		return 120, 0, err
	}

	return ap.dev.listenForState()
}

func (dev NMDevice) listenForState() (state uint32, reason uint32, err error) {
	opts := []dbus.MatchOption{
		dbus.WithMatchInterface("org.freedesktop.NetworkManager.Device"),
		dbus.WithMatchObjectPath(dev.Path),
		dbus.WithMatchMember("StateChanged"),
	}

	err = dev.conn.AddMatchSignal(opts...)
	if err != nil {
		return 0, 0, err
	}

	signals := make(chan *dbus.Signal)
	dev.conn.Signal(signals)

	defer func() {
		dev.conn.RemoveMatchSignal(opts...)
		dev.conn.RemoveSignal(signals)
	}()

	for {
		select {
		case signal := <-signals:
			state = signal.Body[0].(uint32)
			reason = signal.Body[2].(uint32)

			if repository.ConnectionStatusErrFailed.Equal(state) || repository.ConnectionStatusActivated.Equal(state) {
				return state, reason, nil
			}
		case <-time.After(time.Second * 10):
			return state, reason, nil
		}
	}
}

func (ap NMAccessPoint) IsConnected() bool {
	path, err := ap.dev.conn.Object(nmInterface, ap.dev.Path).GetProperty("org.freedesktop.NetworkManager.Device.Wireless.ActiveAccessPoint")
	if err != nil {
		return false // handle it later
	}

	if path.Value().(dbus.ObjectPath) == ap.Path {
		return true
	}

	return false
}
