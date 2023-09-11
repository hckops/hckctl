package common

import (
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/template"
)

func AddTemplateLabels[T template.TemplateType](info *template.TemplateInfo[T], labels commonModel.Labels) commonModel.Labels {
	var allLabels commonModel.Labels
	switch info.SourceType {
	case template.Local:
		allLabels = labels.AddLocal(info.Path)
	case template.Git:
		allLabels = labels.AddGit(info.Path, info.Revision)
	}
	return allLabels
}
