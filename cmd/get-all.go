package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Yunsang-Jeong/terraform-state-dependency-walker/internal/aws"
	"github.com/Yunsang-Jeong/terraform-state-dependency-walker/internal/terraform"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	errors "github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	AWSClientRegion   = "ap-northeast-2"
	BucketName        = ""
	DynamodbTableName = "tsdw"
	JsonFileName      = "dependency_map.json"
)

type Dependency struct {
	Data     []string `json:"data"`
	Resource []string `json:"resource"`
}

type GetAllCmd struct{}

func (g *GetAllCmd) Init() *cobra.Command {
	c := &cobra.Command{
		Use:   "get-all",
		Short: "A brief description of your command",
		RunE: func(cmd *cobra.Command, args []string) error {
			return g.run()
		},
	}

	return c
}

func (g *GetAllCmd) run() error {

	a := aws.AWS{
		AWSClientRegion: AWSClientRegion,
	}

	if err := a.SetAWSClient(); err != nil {
		return err
	}

	bufferMap, err := a.DownloadS3ObjectToBuffer(BucketName, []string{"terraform.tfstate", "terraform.state"})
	if err != nil {
		return err
	}

	dependencyMap := map[string]Dependency{}

	for stateFileName, buf := range bufferMap {
		state, err := terraform.ReadTerraformState(buf)
		if err != nil {
			return err
		}

		dataBlock, err := terraform.ParseDataBlockUsingRemoteState(state)
		if err != nil {
			return err
		}

		resourceBlock, err := terraform.ParseResourceBlockUsingRemoteStateBlock(state)
		if err != nil {
			return err
		}

		if len(*dataBlock) > 1 || len(*resourceBlock) > 1 {
			dynamodbAttributeValues := make(map[string]types.AttributeValue)

			dynamodbAttributeValues["StateFileLocation"], err = attributevalue.Marshal(fmt.Sprintf("%s/%s", BucketName, stateFileName))
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

			err := a.PutItemToAWSDynamodb(DynamodbTableName, dynamodbAttributeValues)
			if err != nil {
				return err
			}

			dependencyMap[stateFileName] = Dependency{
				Data:     *dataBlock,
				Resource: *resourceBlock,
			}

			log.Printf("[%s]", stateFileName)

			for _, d := range *dataBlock {
				log.Printf(" - [data] %s", d)
			}

			for _, r := range *resourceBlock {
				log.Printf(" - [resource] %s", r)
			}
		}
	}

	dependencyMapJson, err := json.MarshalIndent(dependencyMap, "", "  ")
	if err != nil {
		return err
	}

	jsonFile, err := os.Create(JsonFileName)
	if err != nil {
		return errors.Wrap(err, "[get-all] fail to create json file")
	}

	defer jsonFile.Close()

	jsonFile.Write(dependencyMapJson)
	jsonFile.Close()

	return nil
}
