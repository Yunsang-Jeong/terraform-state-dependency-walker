package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	errors "github.com/pkg/errors"
)

type AWS struct {
	AWSClientRegion string
	AWSClientConfig aws.Config
}

func (a *AWS) SetAWSClient() error {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(a.AWSClientRegion),
	)
	if err != nil {
		return errors.Wrap(err, "fail to create a new aws client")
	}
	a.AWSClientConfig = cfg

	return nil
}
