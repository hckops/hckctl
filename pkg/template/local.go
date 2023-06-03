package template

import (
	"fmt"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/hckops/hckctl/pkg/template/schema"
	"github.com/hckops/hckctl/pkg/util"
)

// TODO refactor

func LoadLocalTemplate(path string) (schema.SchemaKind, error) {
	localTemplate, err := util.ReadFile(path)
	if err != nil {
		return -1, errors.Wrapf(err, "local template not found %s", localTemplate)
	}
	return schema.ValidateAll(localTemplate)
}

func LoadLocalBoxTemplate(path string) (*BoxV1, error) {
	localTemplate, err := util.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "local template not found %s", localTemplate)
	}

	if err := schema.ValidateBoxV1(localTemplate); err != nil {
		return nil, err
	}

	var boxSchema BoxV1
	if err := yaml.Unmarshal([]byte(localTemplate), &boxSchema); err != nil {
		return nil, fmt.Errorf("decode error: %v", err)
	}
	return &boxSchema, nil
}

func LoadLocalLabTemplate(path string) (*LabV1, error) {
	localTemplate, err := util.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "local template not found %s", localTemplate)
	}

	if err := schema.ValidateLabV1(localTemplate); err != nil {
		return nil, err
	}

	var labSchema LabV1
	if err := yaml.Unmarshal([]byte(localTemplate), &labSchema); err != nil {
		return nil, fmt.Errorf("decode error: %v", err)
	}
	return &labSchema, nil
}
