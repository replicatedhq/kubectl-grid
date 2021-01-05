package grid

import (
	"fmt"
	"reflect"
	"time"

	"github.com/replicatedhq/kubectl-grid/pkg/grid/types"
)

// Create will create the grid defined in the gridSpec
// the name of the grid will be the name in the metadata.name field
// This function is synchronous and will not return until all clusters are ready
func Create(g *types.Grid) error {
	completed := map[int]bool{}
	completedChans := make([]chan string, len(g.Spec.Clusters))
	for i := range g.Spec.Clusters {
		completedChans[i] = make(chan string)
		completed[i] = false
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
		go createCluster(cluster, completedChans[i])
	}

	// wait for all channels to be closed
	<-finished

	return nil
}

// createCluster will create the cluster synchronously
// when it's completed, it will return the error or "" as a string on the channel
func createCluster(cluster *types.ClusterSpec, completedCh chan string) {
	time.Sleep(time.Second)
	completedCh <- "test"
}
