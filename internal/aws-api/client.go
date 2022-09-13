package awsApi

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	errors "github.com/pkg/errors"
)

func SetAWSClient(awsClientRegion string) (*aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsClientRegion),
	)
	if err != nil {
		return nil, errors.Wrap(err, "[awsApi:SetAWSClient] fail to create a new aws client")
	}

	return &cfg, nil
}
