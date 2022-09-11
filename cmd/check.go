package cmd

import (
	"github.com/spf13/cobra"
)

type CheckCmd struct{}

func (g *CheckCmd) Init() *cobra.Command {
	c := &cobra.Command{
		Use:   "check",
		Short: "Check the dependency on the backend.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return g.run()
		},
	}

	return c
}

func (g *CheckCmd) run() error {

	return nil
}
