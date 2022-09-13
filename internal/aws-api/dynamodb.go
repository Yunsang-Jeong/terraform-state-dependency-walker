package awsApi

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	errors "github.com/pkg/errors"
)

func PutItemToAWSDynamodb(awsConfig aws.Config, tableName string, item map[string]types.AttributeValue) error {
	cli := dynamodb.NewFromConfig(awsConfig)

	// reponse: https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/dynamodb#PutItemOutput
	// AWS Docs: https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_PutItem.html
	_, err := cli.PutItem(context.TODO(),
		&dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      item,
		},
	)
	if err != nil {
		return errors.Wrap(err, "[awsAPi:PutItemToAWSDynamodb] fail to put item to AWS DynamoDB")
	}

	return nil
}

// func GetItemToAWSDynamodb(awsConfig aws.Config, tableName string, item map[string]types.AttributeValue) (*map[string]types.AttributeValue, error) {
// 	cli := dynamodb.NewFromConfig(awsConfig)

// 	resp, err := cli.GetItem(context.TODO(),
// 		&dynamodb.GetItemInput{
// 			TableName: aws.String(tableName),
// 			Key:       item,
// 		},
// 	)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "[awsAPi:GetItemToAWSDynamodb] fail to get item to AWS DynamoDB")
// 	}

// 	return &resp.Item, nil
// }

func QueryToAWSDynamodb(awsConfig aws.Config, tableName string, item map[string]types.AttributeValue) (*[]map[string]types.AttributeValue, error) {
	cli := dynamodb.NewFromConfig(awsConfig)

	resp, err := cli.Scan(context.TODO(),
		&dynamodb.ScanInput{
			TableName:                 aws.String(tableName),
			FilterExpression:          aws.String("contains(DataBlocksWithTypeRemoteState, :key)"),
			ExpressionAttributeValues: item,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "[awsAPi:QueryToAWSDynamodb] fail to query item to AWS DynamoDB")
	}

	return &resp.Items, nil
}
