package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetNamespacesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "namespaces",
		Aliases: []string{
			"namespace",
			"ns",
		},
		Short:         "List the namespaces in a cluster on the grid",
		SilenceErrors: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("not implemented\n")
			return nil
		},
	}

	return cmd
}
