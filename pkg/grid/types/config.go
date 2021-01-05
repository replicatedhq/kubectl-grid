package types

type GridsConfig struct {
	GridConfigs []*GridConfig `json:"grids,omitempty"`
}

type GridConfig struct {
	Name           string           `json:"name"`
	ClusterConfigs []*ClusterConfig `json:"clusters,omitempty"`
}

type ClusterConfig struct {
	Provider   string `json:"provider"`
	IsExisting bool   `json:"isExisting"`
	Region     string `json:"region"`
	Kubeconfig string `json:"kubeconfig,omitempty"`
}
