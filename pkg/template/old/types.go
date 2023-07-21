package old

import (
	"fmt"
	"path/filepath"

	box "github.com/hckops/hckctl/pkg/box/model"
	lab "github.com/hckops/hckctl/pkg/lab/model"
	"github.com/hckops/hckctl/pkg/schema"
)

const (
	InvalidPath     = "INVALID_PATH"
	InvalidRevision = "INVALID_REVISION"
)

type BoxTemplate TemplateValue[box.BoxV1]
type LabTemplate TemplateValue[lab.LabV1]
type BoxInfo TemplateInfo[box.BoxV1]
type LabInfo TemplateInfo[lab.LabV1]

func newBoxTemplate(value *box.BoxV1) *BoxTemplate {
	return &BoxTemplate{
		Kind: schema.KindBoxV1,
		Data: *value,
	}
}

func newLabTemplate(value *lab.LabV1) *LabTemplate {
	return &LabTemplate{
		Kind: schema.KindLabV1,
		Data: *value,
	}
}

func newDefaultTemplateInfo[T TemplateType](template *TemplateValue[T], sourceType SourceType) *TemplateInfo[T] {
	return &TemplateInfo[T]{
		Value:      template,
		SourceType: sourceType,
		Cached:     false,
		Path:       InvalidPath,
		Revision:   InvalidRevision,
	}
}

func newCachedTemplateInfo[T TemplateType](template *TemplateValue[T], sourceType SourceType, path string) (*TemplateInfo[T], error) {

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve absolute path %s", path)
	}

	return &TemplateInfo[T]{
		Value:      template,
		SourceType: sourceType,
		Cached:     true,
		Path:       absolutePath,
		Revision:   InvalidRevision,
	}, nil
}

func newGitTemplateInfo[T TemplateType](template *TemplateValue[T], path string, hash string) (*TemplateInfo[T], error) {

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve absolute path %s", path)
	}

	return &TemplateInfo[T]{
		Value:      template,
		SourceType: Git,
		Cached:     false,
		Path:       absolutePath,
		Revision:   hash,
	}, nil
}

func (t *RawTemplate) toValidated(path string, isValid bool) *TemplateValidated {
	return &TemplateValidated{t, path, isValid}
}
