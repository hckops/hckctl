package template

import (
	"fmt"
	"io"
	"net/http"

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

// TODO not used e.g. https://raw.githubusercontent.com/hckops/megalopolis/main/box/base/alpine.yml
func httpGetString(url string) (string, error) {
	// TODO context with timeout
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.Wrap(err, "network error")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("not found")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil || len(data) == 0 {
		return "", errors.Wrap(err, "invalid body")
	}

	return string(data), nil
}
