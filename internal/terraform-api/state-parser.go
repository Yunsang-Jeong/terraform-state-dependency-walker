package terraformApi

import (
	"encoding/json"
	"strings"

	errors "github.com/pkg/errors"
)

type TerraformState struct {
	Version          int                      `json:"version"`
	TerraformVersion string                   `json:"terraform_version"`
	Resources        []TerraformStateResource `json:"resources"`
	TheOthers        map[string]interface{}   `json:"-"`
}

type TerraformStateResource struct {
	Module    string                           `json:"module"`
	Mode      string                           `json:"mode"`
	Type      string                           `json:"type"`
	Name      string                           `json:"name"`
	Provider  string                           `json:"provider"`
	Instances []TerraformStateResourceInstance `json:"instances"`
}

type TerraformStateResourceInstance struct {
	IndexKey            interface{}       `json:"index_key,omitempty"`
	Status              string            `json:"status,omitempty"`
	Deposed             string            `json:"deposed,omitempty"`
	SchemaVersion       uint64            `json:"schema_version"`
	Attributes          json.RawMessage   `json:"attributes,omitempty"`
	AttributesFlat      map[string]string `json:"attributes_flat,omitempty"`
	SensitiveAttributes json.RawMessage   `json:"sensitive_attributes,omitempty"`
	Private             []byte            `json:"private,omitempty"`
	Dependencies        []string          `json:"dependencies,omitempty"`
	CreateBeforeDestroy bool              `json:"create_before_destroy,omitempty"`
}

type TerraformStateResourceInstanceAttribute struct {
	Backend   string                                         `json:"backend"`
	Config    TerraformStateResourceInstanceAttributeConfig  `json:"config"`
	Outputs   TerraformStateResourceInstanceAttributeOutputs `json:"outputs"`
	TheOthers map[string]interface{}                         `json:"-"`
}

type TerraformStateResourceInstanceAttributeConfig struct {
	Value     map[string]string      `json:"value"`
	TheOthers map[string]interface{} `json:"-"`
}

type TerraformStateResourceInstanceAttributeOutputs struct {
	Value     map[string]string      `json:"value"`
	TheOthers map[string]interface{} `json:"-"`
}

type ParsedTerraformBlock struct {
	TerraformStateLocation              string   `dynamodbav:"TerraformStateLocation,string"`
	ResourceBlocksUsingRemoteState      []string `dynamodbav:"ResourceBlocksUsingRemoteState,stringset"`
	ResourceBlocksUsingRemoteStateCount int      `dynamodbav:"-"`
	DataBlocksWithTypeRemoteState       []string `dynamodbav:"DataBlocksWithTypeRemoteState,stringset"`
	DataBlocksWithTypeRemoteStateCount  int      `dynamodbav:"-"`
}

func removeDuplicate[T string | int](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}

	return list
}

func MakeTerraformStateObjectFromData(terraformStateData []byte) (*TerraformState, error) {
	terraformStateObject := &TerraformState{}

	if err := json.Unmarshal(terraformStateData, terraformStateObject); err != nil {
		return nil, errors.Wrap(err, "[terraformApi:ReadTerraformState] fail to unmarshal terraform state")
	}

	return terraformStateObject, nil
}

func (p *ParsedTerraformBlock) ParseTerraformBlocksAssociatedWithRemoteState(terraformState *TerraformState) error {
	for _, resource := range terraformState.Resources {
		for _, instance := range resource.Instances {

			if resource.Mode == "data" && resource.Type == "terraform_remote_state" {
				/*
					Parse the data block using remote state
				*/
				attribute := &TerraformStateResourceInstanceAttribute{}

				if err := json.Unmarshal(instance.Attributes, attribute); err != nil {
					return errors.Wrap(err, "[terraformApi:ParseTerraformBlocksUsingRemoteState] fail to unmarshal terraform state")
				}

				dataBlockInfo := strings.Join([]string{attribute.Config.Value["bucket"], attribute.Config.Value["key"]}, "/")
				p.DataBlocksWithTypeRemoteState = append(p.DataBlocksWithTypeRemoteState, dataBlockInfo)
			} else {
				/*
					Parse the resource block using remote state
				*/
				for _, depndency := range instance.Dependencies {
					if strings.HasPrefix(depndency, "data.terraform_remote_state") {
						resourceBlockName := strings.Join([]string{resource.Module, resource.Type, resource.Name}, ".")
						p.ResourceBlocksUsingRemoteState = append(p.ResourceBlocksUsingRemoteState, resourceBlockName)
					}
				}
			}

		}

	}
	p.DataBlocksWithTypeRemoteState = removeDuplicate(p.DataBlocksWithTypeRemoteState)
	p.DataBlocksWithTypeRemoteStateCount = len(p.DataBlocksWithTypeRemoteState)

	p.ResourceBlocksUsingRemoteState = removeDuplicate(p.ResourceBlocksUsingRemoteState)
	p.ResourceBlocksUsingRemoteStateCount = len(p.ResourceBlocksUsingRemoteState)

	return nil
}
