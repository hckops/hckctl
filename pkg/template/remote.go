package template

import (
	"github.com/pkg/errors"
)

type RemoteSource[T TemplateType] struct {
	cacheOpts *CacheSourceOpts
	url       string
}

func (src *RemoteSource[T]) Parse() (*RawTemplate, error) {
	return nil, errors.New("not implemented")
}

func (src *RemoteSource[T]) Validate() ([]*TemplateValidated, error) {
	return nil, errors.New("not implemented")
}

func (src *RemoteSource[T]) Read() (*TemplateInfo[T], error) {
	return nil, errors.New("not implemented")
}
