package cmd

import (
	"github.com/Yunsang-Jeong/terraform-state-dependency-walker/internal"
	"github.com/spf13/cobra"
)

type GetAllCmd struct{}

var stateFileNameFilter = []string{"terraform.tfstate", "terraform.state"}

var getAllStringFlags = map[string]stringFlag{
	AWSClientRegion: {
		shorten:      "r",
		defaultValue: "ap-northeast-2",
		description:  "[opt] The name of AWS client region",
		requirement:  false,
	},
	TerraformStateBucketName: {
		shorten:     "b",
		description: "[req] The name of AWS S3 bucket to search terraform state",
		requirement: true,
	},
	DynamodbTableName: {
		shorten:     "d",
		description: "[req] The name of AWS DynamoDB table to put dependecy info",
		requirement: true,
	},
}

func (g *GetAllCmd) Init() *cobra.Command {
	c := &cobra.Command{
		Use:   "get-all",
		Short: "Get terraform state from AWS S3 bucket, make dependency-map by analyzing it, and put dependency-map to AWS Dynamodb",
		Long:  "Get terraform state from AWS S3 bucket, make dependency-map by analyzing it, and put dependency-map to AWS Dynamodb",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := internal.GetAllConfig{
				StateFileNameFilter: stateFileNameFilter,
			}

			config.AWSClientRegion, _ = cmd.Flags().GetString(AWSClientRegion)
			config.TerraformStateBucketName, _ = cmd.Flags().GetString(TerraformStateBucketName)
			config.DynamodbTableName, _ = cmd.Flags().GetString(DynamodbTableName)

			if err := config.GetAllStart(); err != nil {
				return err
			}

			return nil
		},
	}

	for name, flag := range getAllStringFlags {
		c.Flags().StringP(
			name,
			flag.shorten,
			flag.defaultValue,
			flag.description,
		)

		if flag.requirement {
			c.MarkFlagRequired(name)
		}
	}

	return c
}
