package template

import (
	"fmt"
	"github.com/hckops/hckctl/pkg/template/schema"
	"github.com/hckops/hckctl/pkg/util"
	"github.com/pkg/errors"
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

func (l *LocalTemplateLoader) Load() (*TemplateValue, error) {
	return LoadLocalTemplate(l.opts)
}

func LoadLocalTemplate(opts *LocalTemplateOpts) (*TemplateValue, error) {
	data, err := util.ReadFile(opts.Path)
	if err != nil {
		return nil, errors.Wrapf(err, "local template not found %s", opts.Path)
	}

	kind, err := schema.ValidateAll(data)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid schema %s", data)
	}

	switch opts.Format {
	case YamlFormat.String():
		return &TemplateValue{kind, data, YamlFormat}, nil
	case JsonFormat.String():
		if jsonValue, err := ConvertFromYamlToJson(kind, data); err != nil {
			return nil, errors.Wrap(err, "conversion from yaml to json failed")
		} else {
			// adds newline only for json
			return &TemplateValue{kind, fmt.Sprintf("%s\n", jsonValue), JsonFormat}, nil
		}
	default:
		return nil, fmt.Errorf("invalid Format %s", opts.Format)
	}
}
