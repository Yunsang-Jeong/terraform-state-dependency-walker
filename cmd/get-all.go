package cmd

import (
	"github.com/Yunsang-Jeong/terraform-state-dependency-walker/internal"
	"github.com/spf13/cobra"
)

type GetAllCmd struct{}

func (g *GetAllCmd) Init() *cobra.Command {
	c := &cobra.Command{
		Use:   "get-all",
		Short: "A brief description of your command",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := g.run(); err != nil {
				return err
			}
			return nil
		},
	}

	return c
}

func (g *GetAllCmd) run() error {
	config := internal.GetAllConfig{
		AWSClientRegion:     "ap-northeast-2",
		BucketName:          "tgdf1lk345adsgf0g45n2kl3",
		StateFileNameFilter: []string{"terraform.tfstate", "terraform.state"},
		DynamodbTableName:   "tsdw",
	}

	return config.GetAllStart()
}
