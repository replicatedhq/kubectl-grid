package grid

import (
	"github.com/pkg/errors"
	"github.com/replicatedhq/kubectl-grid/pkg/grid/types"
)

func List(configFilePath string) ([]*types.GridConfig, error) {
	c, err := loadConfig(configFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config")
	}

	return c.GridConfigs, nil
}
