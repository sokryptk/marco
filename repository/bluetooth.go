package repository

type Bluetooth interface {
	GetAdapter() (BluetoothAdapter, error)
}

type BluetoothAdapter interface {
	StartDiscovery() error
	StopDiscovery() error

	Advertise() bool
	StopAdvertise() bool
	GetDevices() ([]BluetoothDevice, error)

	SetPower(on bool) (bool, error)
	IsPowered() bool
}

type BluetoothDevice interface {
	Connect() error
	Disconnect() error
	Remove() error

	Pair() error
	CancelPairing() error
	IsPaired() bool
	IsConnected() bool

	SetAlias(alias string) error
	GetAlias() string

	GetBatteryPercentage() uint
}
