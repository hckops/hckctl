package template

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/hckops/hckctl/pkg/common"
)

// TODO basic validation
func IsValidName(value string) bool {
	return true
}

func IsNotValidName(value string) bool {
	return !IsValidName(value)
}

func httpGetString(url string) (string, error) {
	// TODO context with timeout
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("network error")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("not found")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("invalid body")
	}

	return string(data), nil
}

func FetchPublicTemplate(name string) (string, error) {

	path := fmt.Sprintf("boxes/official/%s.yml", name)

	template, err := httpGetString(fmt.Sprintf("%s/%s", common.UrlMegalopolisRaw, path))
	if err != nil {
		return "", fmt.Errorf("invalid public template")
	}

	return template, nil
}

// TODO add headers e.g. client version
func FetchApiTemplate(name string) (string, error) {

	templateUrl, err := url.Parse(fmt.Sprintf("%s/template", common.UrlApi))
	if err != nil {
		return "", fmt.Errorf("invalid api url")
	}

	params := url.Values{}
	params.Add("name", name)
	templateUrl.RawQuery = params.Encode()

	template, err := httpGetString(templateUrl.String())
	if err != nil {
		return "", fmt.Errorf("invalid api template")
	}

	return template, nil
}
