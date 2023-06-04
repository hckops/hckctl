package loader

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/template"
	"github.com/hckops/hckctl/pkg/template/schema"
	"github.com/hckops/hckctl/pkg/util"
)

type LocalTemplateOpts struct {
	Path   string
	Format string
}

type LocalTemplateLoader struct {
	opts *LocalTemplateOpts
}

func NewLocalTemplateLoader(opts *LocalTemplateOpts) *LocalTemplateLoader {
	return &LocalTemplateLoader{
		opts: opts,
	}
}

func NewDefaultLocalTemplateLoader(path string) *LocalTemplateLoader {
	return &LocalTemplateLoader{
		opts: &LocalTemplateOpts{
			Path:   path,
			Format: template.YamlFormat.String(),
		},
	}
}

func (l *LocalTemplateLoader) Load() (*TemplateValue, error) {
	data, err := util.ReadFile(l.opts.Path)
	if err != nil {
		return nil, errors.Wrapf(err, "local template not found %s", l.opts.Path)
	}

	kind, err := schema.ValidateAll(data)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid schema %s", data)
	}

	switch l.opts.Format {
	case template.YamlFormat.String():
		return &TemplateValue{kind, data, template.YamlFormat}, nil
	case template.JsonFormat.String():
		if jsonValue, err := template.ConvertFromYamlToJson(kind, data); err != nil {
			return nil, errors.Wrap(err, "conversion from yaml to json failed")
		} else {
			// adds newline only for json
			return &TemplateValue{kind, fmt.Sprintf("%s\n", jsonValue), template.JsonFormat}, nil
		}
	default:
		return nil, fmt.Errorf("invalid Format %s", l.opts.Format)
	}
}
