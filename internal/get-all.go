package internal

import (
	"fmt"
	"strings"
	"sync"

	awsApi "github.com/Yunsang-Jeong/terraform-state-dependency-walker/internal/aws-api"
	terraformApi "github.com/Yunsang-Jeong/terraform-state-dependency-walker/internal/terraform-api"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type GetAllConfig struct {
	AWSClientRegion     string
	AWSClientConfig     aws.Config
	BucketName          string
	StateFileNameFilter []string
	DynamodbTableName   string
}

type Dependency struct {
	Data     []string `json:"data"`
	Resource []string `json:"resource"`
}

func (c *GetAllConfig) analyzeTerraformState(wg *sync.WaitGroup, objectName string) error {
	defer wg.Done()

	buffer, err := awsApi.DownloadBucketObjectToBuffer(c.AWSClientConfig, c.BucketName, objectName)
	if err != nil {
		return err
	}

	state, err := terraformApi.ReadTerraformState(buffer)
	if err != nil {
		return err
	}

	dataBlock, err := terraformApi.ParseDataBlockUsingRemoteState(state)
	if err != nil {
		return err
	}

	resourceBlock, err := terraformApi.ParseResourceBlockUsingRemoteStateBlock(state)
	if err != nil {
		return err
	}

	if len(*dataBlock) > 1 || len(*resourceBlock) > 1 {
		dynamodbAttributeValues := make(map[string]types.AttributeValue)

		dynamodbAttributeValues["StateFileLocation"], err = attributevalue.Marshal(fmt.Sprintf("%s/%s", c.BucketName, objectName))
		if err != nil {
			return err
		}

		dynamodbAttributeValues["Data"], err = attributevalue.Marshal(strings.Join(*dataBlock, ","))
		if err != nil {
			return err
		}

		dynamodbAttributeValues["Resource"], err = attributevalue.Marshal(strings.Join(*resourceBlock, ","))
		if err != nil {
			return err
		}

		err := awsApi.PutItemToAWSDynamodb(c.AWSClientConfig, c.DynamodbTableName, dynamodbAttributeValues)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *GetAllConfig) GetAllStart() error {
	awsClientConfig, err := awsApi.SetAWSClient(c.AWSClientRegion)
	if err != nil {
		return err
	}

	c.AWSClientConfig = *awsClientConfig

	objectList, err := awsApi.GetFilteredBucketObjectList(c.AWSClientConfig, c.BucketName, c.StateFileNameFilter)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, object := range *objectList {
		wg.Add(1)

		go c.analyzeTerraformState(&wg, object)
	}
	wg.Wait()

	return nil
}
