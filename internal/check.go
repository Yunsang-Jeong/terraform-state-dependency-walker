package internal

import (
	"fmt"

	awsApi "github.com/Yunsang-Jeong/terraform-state-dependency-walker/internal/aws-api"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type CheckConfig struct {
	AWSClientRegion          string
	AWSClientConfig          aws.Config
	TerraformStateBucketName string
	TerraformStateObjectName string
	DynamodbTableName        string
	DynamodbTablePrimaryKey  string
}

func (c *CheckConfig) CheckStart() error {
	awsClientConfig, err := awsApi.SetAWSClient(c.AWSClientRegion)
	if err != nil {
		return err
	}

	c.AWSClientConfig = *awsClientConfig

	ddbQuery := make(map[string]types.AttributeValue)
	ddbQuery["Data"], err = attributevalue.Marshal(fmt.Sprintf("%s/%s", c.TerraformStateBucketName, c.TerraformStateObjectName))
	if err != nil {
		return err
	}

	ddbResp, err := awsApi.GetItemToAWSDynamodb(c.AWSClientConfig, c.DynamodbTableName, ddbQuery)
	if err != nil {
		return err
	}

	fmt.Print(*ddbResp)
	// for k, _ := range *ddbResp {
	// 	fmt.Printf("[%s]", k)
	// }

	return nil
}
