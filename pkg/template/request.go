package template

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/util"
)

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
