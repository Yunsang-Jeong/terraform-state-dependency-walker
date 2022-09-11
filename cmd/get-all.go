package cmd

import (
	"github.com/Yunsang-Jeong/terraform-state-dependency-walker/internal"
	"github.com/spf13/cobra"
)

type GetAllCmd struct{}

var (
	AWSClientRegion   string
	BucketName        string
	DynamodbTableName string
)

func (g *GetAllCmd) Init() *cobra.Command {
	c := &cobra.Command{
		Use:          "get-all",
		Short:        "Get terraform state from AWS S3 bucket, analyze it, and put result(dependency map) to AWS Dynamodb",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := g.run(); err != nil {
				return err
			}
			return nil
		},
	}

	c.Flags().StringVar(&AWSClientRegion, "aws-region", "ap-northeast-2", "The name of AWS region")

	c.Flags().StringVar(&BucketName, "bucket-name", "", "[req] The name of AWS S3 bucket to search terraform state")
	c.MarkFlagRequired("bucket-name")

	c.Flags().StringVar(&DynamodbTableName, "ddb-name", "", "[req] The name of AWS DynamoDB table to put dependecy info")
	c.MarkFlagRequired("ddb-name")

	return c
}

func (g *GetAllCmd) run() error {
	config := internal.GetAllConfig{
		AWSClientRegion:     AWSClientRegion,
		BucketName:          BucketName,
		DynamodbTableName:   DynamodbTableName,
		StateFileNameFilter: []string{"terraform.tfstate", "terraform.state"},
	}

	return config.GetAllStart()
}
