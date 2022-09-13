package internal

import (
	"fmt"
	"strings"

	awsApi "github.com/Yunsang-Jeong/terraform-state-dependency-walker/internal/aws-api"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
)

type CheckConfig struct {
	AWSClientRegion          string
	AWSClientConfig          aws.Config
	TerraformStateBucketName string
	TerraformStateObjectName string
	DynamodbTableName        string
}

type ScanResult struct {
	TerraformStateLocation         string   `dynamodbav:"TerraformStateLocation,string"`
	ResourceBlocksUsingRemoteState []string `dynamodbav:"ResourceBlocksUsingRemoteState,stringset"`
	DataBlocksWithTypeRemoteState  []string `dynamodbav:"DataBlocksWithTypeRemoteState,stringset"`
}

func (c *CheckConfig) CheckStart() error {
	awsClientConfig, err := awsApi.SetAWSClient(c.AWSClientRegion)
	if err != nil {
		return errors.Wrap(err, "[internal:CheckStart]")
	}

	c.AWSClientConfig = *awsClientConfig

	ddbQuery := make(map[string]types.AttributeValue)
	ddbQuery[":key"], err = attributevalue.Marshal(fmt.Sprintf("%s/%s", c.TerraformStateBucketName, c.TerraformStateObjectName))
	if err != nil {
		return errors.Wrap(err, "[internal:CheckStart]")
	}

	ddbResp, err := awsApi.QueryToAWSDynamodb(c.AWSClientConfig, c.DynamodbTableName, ddbQuery)
	if err != nil {
		return errors.Wrap(err, "[internal:CheckStart]")
	}

	result := []ScanResult{}
	if err := attributevalue.UnmarshalListOfMaps(*ddbResp, &result); err != nil {
		return errors.Wrap(err, "[internal:CheckStart]")
	}

	fmt.Printf("✨ Find %d dependencies!\n", len(*ddbResp))
	for index, r := range result {
		fmt.Printf(" ✔ [%d] %s", index, r.TerraformStateLocation)
		if len(r.ResourceBlocksUsingRemoteState) > 0 {
			fmt.Printf(" (maybe used in '%s')\n", strings.Join(r.ResourceBlocksUsingRemoteState, "', '"))
		} else {
			fmt.Print("\n")
		}
	}

	return nil
}
