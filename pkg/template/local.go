package template

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/template/schema"
	"github.com/hckops/hckctl/pkg/util"
)

type LocalTemplateOpts struct {
	Path   string
	Format string
}

type RemoteTemplateOpts struct {
	SourceDir string
	SourceUrl string
	Revision  string
	Name      string
	Format    string
}

type TemplateValue struct {
	Kind   schema.SchemaKind
	Data   string
	Format Format
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

func LoadRemoteTemplate() {
	// check if git source exists
	// > if not download --> (if fail exit)
	// > otherwise update --> (if fail WARN offline but continue)
	// load local template
	// list all templates
}
