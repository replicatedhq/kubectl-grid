package grid

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/replicatedhq/kubectl-grid/pkg/grid/types"
)

func Delete(g *types.GridConfig, configFilePath string) error {
	// delete any clusters that we created
	for _, c := range g.ClusterConfigs {
		if c.IsExisting {
			continue
		}

		if err := deleteCluster(c); err != nil {
			return errors.Wrap(err, "failed to delete cluster")
		}
	}

	if err := removeGridFromConfig(g.Name, configFilePath); err != nil {
		return errors.Wrap(err, "failed to remove grid from config")
	}

	return nil
}

func deleteCluster(c *types.ClusterConfig) error {
	if c.Provider == "aws" {
		return deleteEKSCluster(c)
	}

	return nil
}

func deleteEKSCluster(c *types.ClusterConfig) error {
	clusterName := c.GetDeterministicClusterName()

	fmt.Printf("%s\n", clusterName)
	return nil
}
