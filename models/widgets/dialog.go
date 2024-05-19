package widgets

type DialogConfigurationType int

const (
	DialogConfigurationTypeString = iota
	DialogConfigurationTypeInput  = iota
)

type Dialog struct {
	Config []DialogConfiguration
}

type DialogConfiguration struct {
	Type interface{}
}
