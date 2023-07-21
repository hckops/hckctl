package template

import (
	"path/filepath"

	"github.com/google/uuid"

	"github.com/hckops/hckctl/pkg/util"
)

type LocalSource[T TemplateType] struct {
	cacheOpts *CacheSourceOpts
	path      string
}

func (src *LocalSource[T]) Parse() (*RawTemplate, error) {
	return readRawTemplate(src.path)
}

func (src *LocalSource[T]) Validate() ([]*TemplateValidated, error) {
	return readTemplates(src.path)
}

func (src *LocalSource[T]) Read() (*TemplateInfo[T], error) {
	if src.cacheOpts != nil {
		return readLocalCachedTemplateInfo[T](src.cacheOpts, src.path, Local)
	}
	return readTemplateInfo[T](Local, src.path, Local.String())
}

func readLocalCachedTemplateInfo[T TemplateType](cacheOpts *CacheSourceOpts, path string, sourceType SourceType) (*TemplateInfo[T], error) {

	value, err := readTemplate[T](path)
	if err != nil {
		return nil, err
	}

	cachedPath := filepath.Join(cacheOpts.cacheDir, cacheOpts.cacheName, uuid.New().String())
	if _, err := util.CopyFile(path, cachedPath); err != nil {
		return nil, err
	}

	return &TemplateInfo[T]{
		Value:      value,
		SourceType: sourceType,
		Cached:     true,
		Path:       cachedPath,
		Revision:   sourceType.String(),
	}, nil
}
