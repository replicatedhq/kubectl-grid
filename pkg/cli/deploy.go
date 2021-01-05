package cli

import (
	"errors"
	"io/ioutil"

	"github.com/replicatedhq/kubectl-grid/pkg/app"
	"github.com/replicatedhq/kubectl-grid/pkg/grid"
	"github.com/replicatedhq/kubectl-grid/pkg/grid/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sigs.k8s.io/yaml"
)

func DeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "deploy",
		Short:         "Deploy an application to a grid",
		SilenceErrors: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			v := viper.GetViper()

			data, err := ioutil.ReadFile(v.GetString("application"))
			if err != nil {
				return err
			}

			application := types.Application{}
			if err := yaml.Unmarshal(data, &application); err != nil {
				return err
			}

			grids, err := grid.List(v.GetString("config-file"))
			if err != nil {
				return err
			}

			for _, g := range grids {
				if g.Name == v.GetString("grid") {
					if err := app.Deploy(g, &application); err != nil {
						return err
					}

					return nil
				}
			}

			return errors.New("unable to find grid")
		},
	}

	cmd.Flags().StringP("grid", "g", "", "Name of the grid")
	cmd.Flags().String("application", "", "Path to YAML manifest describing the application to deploy")

	return cmd
}
