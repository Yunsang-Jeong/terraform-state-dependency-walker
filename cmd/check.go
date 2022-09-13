package cmd

import (
	"github.com/Yunsang-Jeong/terraform-state-dependency-walker/internal"
	"github.com/spf13/cobra"
)

type CheckCmd struct{}

var checkStringFlags = map[string]stringFlag{
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
	TerraformStateObjectName: {
		shorten:     "s",
		description: "[req] The name of terraform state in AWS S3 bucket",
		requirement: true,
	},
	DynamodbTableName: {
		shorten:     "d",
		description: "[req] The name of AWS DynamoDB table to put dependecy info",
		requirement: true,
	},
	DynamodbTablePrimaryKey: {
		shorten:      "p",
		defaultValue: "StateFileLocation",
		description:  "[opt] The name of primary key in AWS DynamoDB table",
		requirement:  false,
	},
}

func (g *CheckCmd) Init() *cobra.Command {
	c := &cobra.Command{
		Use:   "check",
		Short: "Check the dependency on the backend.",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := internal.CheckConfig{}

			config.AWSClientRegion, _ = cmd.Flags().GetString(AWSClientRegion)
			config.TerraformStateBucketName, _ = cmd.Flags().GetString(TerraformStateBucketName)
			config.TerraformStateObjectName, _ = cmd.Flags().GetString(TerraformStateObjectName)
			config.DynamodbTableName, _ = cmd.Flags().GetString(DynamodbTableName)
			config.DynamodbTablePrimaryKey, _ = cmd.Flags().GetString(DynamodbTablePrimaryKey)

			if err := config.CheckStart(); err != nil {
				return err
			}

			return nil
		},
	}

	for name, flag := range checkStringFlags {
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
