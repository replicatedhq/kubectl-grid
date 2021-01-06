package cli

import (
	"errors"

	"github.com/replicatedhq/kubectl-grid/pkg/grid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func DeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "delete",
		Short:         "Delete a grid",
		SilenceErrors: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			v := viper.GetViper()

			grids, err := grid.List(v.GetString("config-file"))
			if err != nil {
				return err
			}

			for _, g := range grids {
				if g.Name == args[0] {
					if err := grid.Delete(g, v.GetString("config-file")); err != nil {
						return err
					}

					return nil
				}
			}

			return errors.New("grid not found")
		},
	}

	cmd.Flags().StringP("name", "n", "", "Name of the grid, overriding the name in the yaml metadata.name field")
	cmd.Flags().String("from-yaml", "", "Path to YAML manifest describing the grid to create")
	cmd.Flags().String("like", "", "Name of an existing grid to clone, into a new grid")

	return cmd
}
