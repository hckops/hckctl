package template

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
	return readCachedTemplateInfo[T](src.cacheOpts, src.path, Local)
}
