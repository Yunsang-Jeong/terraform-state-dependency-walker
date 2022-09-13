package internal

import (
	"fmt"
	"log"
	"sync"

	awsApi "github.com/Yunsang-Jeong/terraform-state-dependency-walker/internal/aws-api"
	terraformApi "github.com/Yunsang-Jeong/terraform-state-dependency-walker/internal/terraform-api"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/pkg/errors"
)

type GetAllConfig struct {
	AWSClientRegion          string
	AWSClientConfig          aws.Config
	TerraformStateBucketName string
	StateFileNameFilter      []string
	DynamodbTableName        string
}

func (c *GetAllConfig) analyzeTerraformStateAndInputResult(wg *sync.WaitGroup, errChannel chan<- error, logChannel chan<- terraformApi.ParsedTerraformBlock, terraformStateObjectName string) {
	defer wg.Done()

	terraformStateLocation := fmt.Sprintf("%s/%s", c.TerraformStateBucketName, terraformStateObjectName)

	buffer, err := awsApi.DownloadBucketObjectToBuffer(c.AWSClientConfig, c.TerraformStateBucketName, terraformStateObjectName)
	if err != nil {
		errChannel <- errors.Wrap(err, "[internal:analyzeTerraformState]")
		return
	}

	terraformStateObject, err := terraformApi.MakeTerraformStateObjectFromData(buffer)
	if err != nil {
		errChannel <- errors.Wrap(err, "[internal:analyzeTerraformState]")
		return
	}

	result := terraformApi.ParsedTerraformBlock{
		TerraformStateLocation: terraformStateLocation,
	}

	if err := result.ParseTerraformBlocksAssociatedWithRemoteState(terraformStateObject); err != nil {
		errChannel <- errors.Wrap(err, "[internal:analyzeTerraformState]")
		return
	}

	if result.ResourceBlocksUsingRemoteStateCount+result.DataBlocksWithTypeRemoteStateCount > 0 {
		dynamodbAttributeValues, err := attributevalue.MarshalMap(result)
		if err != nil {
			errChannel <- errors.Wrap(err, "[internal:analyzeTerraformState] fail to create dynamodb attribute(StateFileLocation)")
			return
		}

		err = awsApi.PutItemToAWSDynamodb(c.AWSClientConfig, c.DynamodbTableName, dynamodbAttributeValues)
		if err != nil {
			errChannel <- errors.Wrap(err, "[internal:analyzeTerraformState]")
			return
		}
	}

	logChannel <- result
	errChannel <- nil
}

func (c *GetAllConfig) GetAllStart() error {
	awsClientConfig, err := awsApi.SetAWSClient(c.AWSClientRegion)
	if err != nil {
		return errors.Wrap(err, "[internal:GetAllStart]")
	}

	c.AWSClientConfig = *awsClientConfig

	objectList, err := awsApi.GetFilteredBucketObjectList(c.AWSClientConfig, c.TerraformStateBucketName, c.StateFileNameFilter)
	if err != nil {
		return errors.Wrap(err, "[internal:GetAllStart]")
	}

	var wg sync.WaitGroup

	logChannel := make(chan terraformApi.ParsedTerraformBlock, len(*objectList))
	errChannel := make(chan error, len(*objectList))

	for _, object := range *objectList {
		wg.Add(1)
		go c.analyzeTerraformStateAndInputResult(&wg, errChannel, logChannel, object)
	}

	wg.Wait()

	close(logChannel)
	close(errChannel)

	for log := range logChannel {
		if log.ResourceBlocksUsingRemoteStateCount+log.DataBlocksWithTypeRemoteStateCount > 0 {
			fmt.Printf("✨ Find dependency in %s!\n", log.TerraformStateLocation)
			fmt.Printf(" ✔ Data Blocks With Type Remote State : %d!\n", log.DataBlocksWithTypeRemoteStateCount)
			fmt.Printf(" ✔ Resource Blocks Using Remote State : %d!\n", log.ResourceBlocksUsingRemoteStateCount)
		}
	}

	errFlag := false
	for err := range errChannel {
		if err != nil {
			log.Println(err)
			errFlag = true
		}
	}

	if errFlag {
		return errors.New("Fail to run get-all command. Review printed log.")
	}

	return nil
}
