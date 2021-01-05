package types

type ValueOrValueFrom struct {
	Value     string     `json:"value,omitempty"`
	ValueFrom *ValueFrom `json:"valueFrom,omitempty"`
}

type ValueFrom struct {
	OSEnv string `json:"osEnv,omitempty"`
}
