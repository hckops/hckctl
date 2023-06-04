package loader

import (
	"github.com/hckops/hckctl/pkg/template"
	"github.com/hckops/hckctl/pkg/template/schema"
)

type TemplateValue struct {
	Kind   schema.SchemaKind
	Data   string
	Format template.Format
}

type TemplateLoader interface {
	Load() (*TemplateValue, error)
}
