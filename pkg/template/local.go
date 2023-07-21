package template

import (
	"github.com/pkg/errors"
)

type LocalSource[T TemplateType] struct {
	cacheOpts *CacheSourceOpts
	path      string
}

func (src *LocalSource[T]) Parse() (*RawTemplate, error) {
	return nil, errors.New("not implemented")
}

func (src *LocalSource[T]) Validate() ([]*TemplateValidated, error) {
	return nil, errors.New("not implemented")
}

func (src *LocalSource[T]) Read() (*TemplateInfo[T], error) {
	return nil, errors.New("not implemented")
}
