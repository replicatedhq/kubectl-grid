package main

import (
	"github.com/replicatedhq/kubectl-grid/pkg/cli"
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure" // required for Azure
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"   // required for GKE
)

func main() {
	cli.InitAndExecute()
}
