package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	errors "github.com/pkg/errors"
)

func (a *AWS) PutItemToAWSDynamodb(tableName string, item map[string]types.AttributeValue) error {
	cli := dynamodb.NewFromConfig(a.AWSClientConfig)

	// reponse: https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/dynamodb#PutItemOutput
	// AWS Docs: https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_PutItem.html
	_, err := cli.PutItem(
		context.TODO(),
		&dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      item,
		},
	)
	if err != nil {
		return errors.Wrap(err, "fail to put item to AWS DynamoDB")
	}

	return nil
}

func (a *AWS) GetItemToAWSDynamodb(tableName string, item map[string]types.AttributeValue) (*map[string]types.AttributeValue, error) {
	cli := dynamodb.NewFromConfig(a.AWSClientConfig)

	resp, err := cli.GetItem(
		context.TODO(),
		&dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key:       item,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "fail to get item to AWS DynamoDB")
	}

	return &resp.Item, nil
}
