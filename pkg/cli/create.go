package cli

import (
	"errors"
	"io/ioutil"

	"github.com/replicatedhq/kubectl-grid/pkg/grid"
	"github.com/replicatedhq/kubectl-grid/pkg/grid/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sigs.k8s.io/yaml"
)

func CreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "create",
		Short:         "Create a new test grid",
		SilenceErrors: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			v := viper.GetViper()

			if v.GetString("like") != "" {
				return errors.New("like is not yet supported")
			}

			data, err := ioutil.ReadFile(v.GetString("from-yaml"))
			if err != nil {
				return err
			}

			gridSpec := types.Grid{}
			if err := yaml.Unmarshal(data, &gridSpec); err != nil {
				return err
			}

			if v.GetString("name") != "" {
				gridSpec.Name = v.GetString("name")
			}

			if err := grid.Create(&gridSpec); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "Name of the grid, overriding the name in the yaml metadata.name field")
	cmd.Flags().String("from-yaml", "", "Path to YAML manifest describing the grid to create")
	cmd.Flags().String("like", "", "Name of an existing grid to clone, into a new grid")

	return cmd
}
