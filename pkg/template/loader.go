package template

import (
	"github.com/hckops/hckctl/pkg/template/schema"
)

type TemplateValue struct {
	Kind   schema.SchemaKind
	Data   string
	Format Format
}

type TemplateLoader interface {
	Load() (*TemplateValue, error)
}
