package template

import (
	"github.com/pkg/errors"
)

type GitSource[T TemplateType] struct {
	opts *GitSourceOptions
	name string
}

func (src *GitSource[T]) Parse() (*RawTemplate, error) {
	return nil, errors.New("not implemented")
}

func (src *GitSource[T]) Validate() ([]*TemplateValidated, error) {
	// TODO wildcard := fmt.Sprintf("%s/**/*.{yml,yaml}", src.opts.CacheBaseDir)
	return nil, errors.New("not implemented")
}

func (src *GitSource[T]) Read() (*TemplateInfo[T], error) {
	return nil, errors.New("not implemented")
}
