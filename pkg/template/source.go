package template

type SourceValidator interface {
	Parse() (*RawTemplate, error)
	Validate() ([]*TemplateValidated, error)
}

func NewLocalValidator(path string) SourceValidator {
	return &LocalSource[string]{path: path}
}

func NewRemoteValidator(url string) SourceValidator {
	return &RemoteSource[string]{url: url}
}

func NewGitValidator(opts *GitSourceOptions, name string) SourceValidator {
	return &GitSource[string]{opts, name}
}

type SourceLoader[T TemplateType] interface {
	Read() (*TemplateInfo[T], error)
}

type CacheSourceOpts struct {
	cacheDir  string
	cacheName string
}

func NewLocalLoader[T TemplateType](path string) SourceLoader[T] {
	return &LocalSource[T]{path: path}
}

func NewLocalCachedLoader[T TemplateType](path string, cacheDir string) SourceLoader[T] {
	return &LocalSource[T]{
		path: path,
		cacheOpts: &CacheSourceOpts{
			cacheDir:  cacheDir,
			cacheName: Local.String(),
		},
	}
}

func NewRemoteLoader[T TemplateType](url string, cacheDir string) SourceLoader[T] {
	return &RemoteSource[T]{
		url: url,
		cacheOpts: &CacheSourceOpts{
			cacheDir:  cacheDir,
			cacheName: Remote.String(),
		},
	}
}

func NewGitLoader[T TemplateType](opts *GitSourceOptions, name string) SourceLoader[T] {
	return &GitSource[T]{opts, name}
}
