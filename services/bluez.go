package services

import (
	"github.com/godbus/dbus/v5"
	"log"
	"me.kryptk.marco/repository"
)

const bluezInterface = "org.bluez"

type Bluez struct {
	Conn *dbus.Conn
}

func NewBluez() Bluez {
	conn, err := dbus.ConnectSystemBus()

	if err != nil {
		log.Fatal(err)
		return Bluez{}
	}

	return Bluez{
		Conn: conn,
	}
}

type BluezAdapter struct {
	Conn *dbus.Conn
	Path dbus.ObjectPath
}

func (adapter BluezAdapter) Advertise() bool {
	//TODO implement me
	panic("implement me")
}

func (adapter BluezAdapter) StopAdvertise() bool {
	//TODO implement me
	panic("implement me")
}

type BluezDevice struct {
	Adapter BluezAdapter
	Path    dbus.ObjectPath
}

func (b Bluez) GetAdapter() (repository.BluetoothAdapter, error) {
	var objects map[dbus.ObjectPath]map[string]map[string]dbus.Variant
	err := b.Conn.Object(bluezInterface, "/").Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&objects)
	if err != nil {
		return BluezAdapter{}, err
	}

	for k, v := range objects {
		if _, ok := v["org.bluez.Adapter1"]; ok {
			return BluezAdapter{
				Conn: b.Conn,
				Path: k,
			}, nil
		}
	}

	return BluezAdapter{}, nil
}

func (adapter BluezAdapter) StartDiscovery() error {
	err := adapter.Conn.Object(bluezInterface, adapter.Path).Call("org.bluez.Adapter1.StartDiscovery", 0).Err
	if err != nil {
		return err
	}

	return nil
}

func (adapter BluezAdapter) StopDiscovery() error {
	err := adapter.Conn.Object(bluezInterface, adapter.Path).Call("org.bluez.Adapter1.StopDiscovery", 0).Err
	if err != nil {
		return err
	}

	return nil
}

func (adapter BluezAdapter) IsPowered() bool {
	powered, err := adapter.Conn.Object(bluezInterface, adapter.Path).GetProperty("org.bluez.Adapter1.Powered")
	if err != nil {
		return false
	}

	return powered.Value().(bool)
}

func (adapter BluezAdapter) SetPower(power bool) (bool, error) {
	err := adapter.Conn.Object(bluezInterface, adapter.Path).SetProperty("org.bluez.Adapter1.Powered", dbus.MakeVariant(power))
	if err != nil {
		return false, err
	}

	return power, nil
}

func (adapter BluezAdapter) GetDevices() ([]repository.BluetoothDevice, error) {
	var objects map[dbus.ObjectPath]map[string]map[string]dbus.Variant
	err := adapter.Conn.Object(bluezInterface, "/").Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&objects)
	if err != nil {
		return nil, err
	}

	devices := make([]repository.BluetoothDevice, 0)

	for k, v := range objects {
		if device, ok := v["org.bluez.Device1"]; ok {
			if device["Adapter"].Value().(dbus.ObjectPath) == adapter.Path {
				devices = append(devices, BluezDevice{
					Path:    k,
					Adapter: adapter,
				})
			}
		}
	}

	return devices, nil
}

func (device BluezDevice) Connect() error {
	return device.Adapter.Conn.Object(bluezInterface, device.Path).Call(
		"org.bluez.Device1.Connect", 0,
	).Err
}

func (device BluezDevice) Disconnect() error {
	return device.Adapter.Conn.Object(bluezInterface, device.Path).Call(
		"org.bluez.Device1.Connect", 0,
	).Err
}

func (device BluezDevice) Pair() error {
	return device.Adapter.Conn.Object(bluezInterface, device.Path).Call(
		"org.bluez.Device1.Pair", 0,
	).Err
}

func (device BluezDevice) CancelPairing() error {
	return device.Adapter.Conn.Object(bluezInterface, device.Path).Call(
		"org.bluez.Device1.CancelPairing", 0,
	).Err
}

func (device BluezDevice) Remove() error {
	return device.Adapter.Conn.Object(bluezInterface, device.Adapter.Path).Call(
		"org.bluez.Adapter1.RemoveDevice", 0, device.Path,
	).Err
}

func (device BluezDevice) IsPaired() bool {
	paired, err := device.Adapter.Conn.Object(bluezInterface, device.Path).GetProperty(
		"org.bluez.Device1.Paired",
	)

	if err != nil {
		return false
	}

	return paired.Value().(bool)
}

func (device BluezDevice) IsConnected() bool {
	connected, err := device.Adapter.Conn.Object(bluezInterface, device.Path).GetProperty(
		"org.bluez.Device1.Connected",
	)

	if err != nil {
		return false
	}

	return connected.Value().(bool)
}

func (device BluezDevice) SetAlias(alias string) error {
	return device.Adapter.Conn.Object(bluezInterface, device.Path).SetProperty(
		"org.bluez.Device1.Alias", dbus.MakeVariant(alias),
	)
}

func (device BluezDevice) GetAlias() string {
	alias, err := device.Adapter.Conn.Object(bluezInterface, device.Path).GetProperty(
		"org.bluez.Device1.Alias",
	)

	if err != nil {
		return ""
	}

	return alias.Value().(string)
}

func (device BluezDevice) GetBatteryPercentage() uint {
	percentage, err := device.Adapter.Conn.Object(bluezInterface, device.Path).GetProperty(
		"org.bluez.Battery1.Percentage",
	)

	if err != nil {
		return 0
	}

	return percentage.Value().(uint)
}
