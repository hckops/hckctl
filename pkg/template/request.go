package template

import (
	"github.com/hckops/hckctl/pkg/util"
	"github.com/pkg/errors"
)

type TemplateRequest struct {
}

func RequestLocalTemplate(path string) (string, error) {
	localTemplate, err := util.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "local template not found %s", localTemplate)
	}
	// TODO validation
	return localTemplate, nil
}

func RequestRemoteTemplate() (string, error) {
	return "", nil
}
