package grid

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/replicatedhq/kubectl-grid/pkg/grid/types"
)

// Create will create the grid defined in the gridSpec
// the name of the grid will be the name in the metadata.name field
// This function is synchronous and will not return until all clusters are ready
func Create(configFilePath string, g *types.Grid) error {
	completed := map[int]bool{}
	completedChans := make([]chan string, len(g.Spec.Clusters))
	for i := range g.Spec.Clusters {
		completedChans[i] = make(chan string)
		completed[i] = false
	}

	if err := addGridToConfig(configFilePath, g.Name); err != nil {
		return errors.Wrap(err, "failed to add grid to config file")
	}

	// start listening for completed events
	finished := make(chan bool)
	go func() {
		cases := make([]reflect.SelectCase, len(completedChans))
		for i, ch := range completedChans {
			cases[i] = reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ch),
			}
		}

		for {
			i, completedErr, ok := reflect.Select(cases)
			if ok {
				if completedErr.String() != "" {
					fmt.Printf("cluster %#v failed with error: %s\n", g.Spec.Clusters[i], completedErr.String())
				}

				completed[i] = true
			}

			allCompleted := true
			for _, v := range completed {
				if !v {
					allCompleted = false
				}
			}

			if allCompleted {
				finished <- true
				return
			}
		}
	}()

	// start each
	for i, cluster := range g.Spec.Clusters {
		go createCluster(g.Name, cluster, completedChans[i], configFilePath)
	}

	// wait for all channels to be closed
	<-finished

	return nil
}

func addGridToConfig(configFilePath string, name string) error {
	lockConfig()
	defer unlockConfig()
	c, err := loadConfig(configFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to load config")
	}

	if c.GridConfigs == nil {
		c.GridConfigs = []*types.GridConfig{}
	}

	// if the grid already exists, err, this is an add function
	for _, gc := range c.GridConfigs {
		if gc.Name == name {
			return fmt.Errorf("grid with name %s already exists. if you want to delete it, run kubectl grid delete %s", name, name)
		}
	}

	gridConfig := types.GridConfig{
		Name:           name,
		ClusterConfigs: []*types.ClusterConfig{},
	}
	c.GridConfigs = append(c.GridConfigs, &gridConfig)

	if err := saveConfig(c, configFilePath); err != nil {
		return errors.Wrap(err, "failed to save config")
	}

	return nil
}

// createCluster will create the cluster synchronously
// when it's completed, it will return the error or "" as a string on the channel
func createCluster(gridName string, cluster *types.ClusterSpec, completedCh chan string, configFilePath string) {
	if cluster.EKS != nil {
		createEKSCluster(gridName, cluster.EKS, completedCh, configFilePath)
		return
	}

	completedCh <- "unknown cluster"
}

func createEKSCluster(gridName string, eksCluster *types.EKSSpec, completedCh chan string, configFilePath string) {
	if eksCluster.ExistingCluster != nil {
		connectExistingEKSCluster(gridName, eksCluster.ExistingCluster, completedCh, configFilePath)
		return
	} else if eksCluster.NewCluster != nil {
		createNewEKSCluter(gridName, eksCluster.NewCluster, completedCh, configFilePath)
		return
	}

	completedCh <- "eks cluster must have new or existing"
}

func connectExistingEKSCluster(gridName string, existingEKSCluster *types.EKSExistingClusterSpec, completedCh chan string, configFilePath string) {
	accessKeyID, err := existingEKSCluster.AccessKeyID.String()
	if err != nil {
		completedCh <- fmt.Sprintf("failed to read access key id: %s", err.Error())
	}
	secretAccessKey, err := existingEKSCluster.SecretAccessKey.String()
	if err != nil {
		completedCh <- fmt.Sprintf("failed to read secret access key: %s", err.Error())
	}

	kubeConfig, err := GetEKSClusterKubeConfig(existingEKSCluster.Region, accessKeyID, secretAccessKey, existingEKSCluster.ClusterName)
	if err != nil {
		completedCh <- fmt.Sprintf("failed to get kubeconfig from eks cluster: %s", err.Error())
	}

	lockConfig()
	defer unlockConfig()
	c, err := loadConfig(configFilePath)
	if err != nil {
		completedCh <- fmt.Sprintf("failed to load config: %s", err.Error())
		return
	}

	clusterConfig := types.ClusterConfig{
		Provider:   "aws",
		IsExisting: true,
		Region:     existingEKSCluster.Region,
		Kubeconfig: kubeConfig,
	}

	for _, gridConfig := range c.GridConfigs {
		if gridConfig.Name == gridName {
			gridConfig.ClusterConfigs = append(gridConfig.ClusterConfigs, &clusterConfig)
		}
	}
	if err := saveConfig(c, configFilePath); err != nil {
		completedCh <- fmt.Sprintf("error saving config: %s", err.Error())
	}
	completedCh <- "not implemented"
}

func createNewEKSCluter(gridName string, newEKSCluster *types.EKSNewClusterSpec, completedCh chan string, configFilePath string) {
	completedCh <- "not implemented"
}
