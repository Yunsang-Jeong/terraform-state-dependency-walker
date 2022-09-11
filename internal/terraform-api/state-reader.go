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

func ReadTerraformState(terraformStateData []byte) (*TerraformState, error) {
	terraformState := &TerraformState{}

	if err := json.Unmarshal(terraformStateData, terraformState); err != nil {
		return nil, errors.Wrap(err, "fail to unmarshal terraform state")
	}

	return terraformState, nil
}

func ParseDataBlockUsingRemoteState(terraformState *TerraformState) (*[]string, error) {
	dataBlockUsingRemoteState := []string{}

	for _, resource := range terraformState.Resources {
		if resource.Mode == "data" && resource.Type == "terraform_remote_state" {
			for _, instance := range resource.Instances {
				attribute := &TerraformStateResourceInstanceAttribute{}

				if err := json.Unmarshal(instance.Attributes, attribute); err != nil {
					return nil, errors.Wrap(err, "fail to unmarshal terraform state")
				}

				dataBlockInfo := strings.Join([]string{attribute.Config.Value["bucket"], attribute.Config.Value["key"]}, "/")
				dataBlockUsingRemoteState = append(dataBlockUsingRemoteState, dataBlockInfo)
			}
		}
	}

	return &dataBlockUsingRemoteState, nil
}

func ParseResourceBlockUsingRemoteStateBlock(terraformState *TerraformState) (*[]string, error) {
	resourceBlockUsingRemoteStateBlock := []string{}

	for _, resource := range terraformState.Resources {
		for _, instance := range resource.Instances {
			for _, depndency := range instance.Dependencies {
				if strings.HasPrefix(depndency, "data.terraform_remote_state") {
					resourceBlockName := strings.Join([]string{resource.Module, resource.Type, resource.Name}, ".")
					resourceBlockUsingRemoteStateBlock = append(resourceBlockUsingRemoteStateBlock, resourceBlockName)
				}
			}
		}
	}

	return &resourceBlockUsingRemoteStateBlock, nil
}
