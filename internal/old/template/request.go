package template

import (
	"fmt"
	"github.com/hckops/hckctl/internal/old/common"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"

	"github.com/pkg/errors"
)

type TemplateParam struct {
	TemplateKind  string
	TemplateName  string
	Revision      string
	ClientVersion string
}

// TODO better validation: should not start with dashes or contain double dashes
var isValidName = regexp.MustCompile(`^[A-Za-z-]+$`).MatchString

func ValidateTemplateParam(param *TemplateParam) error {

	if !isValidName(param.TemplateName) {
		return fmt.Errorf("invalid name")
	}
	return nil
}

// TODO e.g. https://api.hckops.com/template/box?name=official/alpine&version=main&format=json
// TODO or redirect validate https://schema.hckops.com/validate?kind=box&group=official&name=alpine
// TODO or content https://schema.hckops.com/template?kind=box&group=official&name=alpine&version=main&format=json|yaml
func (param *TemplateParam) RequestApiTemplate() (string, error) {

	templateUrl, err := url.Parse(fmt.Sprintf("%s/todo", common.ApiUrl))
	if err != nil {
		return "", errors.Wrapf(err, "invalid api url: %s", templateUrl.String())
	}

	// TODO authentication
	// TODO add header e.g. x-client=hckctl-v0.0.0
	params := url.Values{}
	params.Add("name", param.TemplateName)
	params.Add("version", param.Revision)
	params.Add("format", "yaml")
	params.Add("client", param.ClientVersion)
	templateUrl.RawQuery = params.Encode()

	template, err := httpGetString(templateUrl.String())
	if err != nil {
		return "", errors.Wrapf(err, "error requesting api template: %s", templateUrl.String())
	}

	return template, nil
}

func (param *TemplateParam) RequestPublicTemplate() (string, error) {

	path, err := buildPath(param)
	if err != nil {
		return "", errors.Wrap(err, "invalid template url")
	}

	template, err := httpGetString(path)
	if err != nil {
		return "", errors.Wrapf(err, "error requesting public template: %s", path)
	}

	return template, nil
}

// TODO replace "official" with fallback to "group/name"
func buildPath(param *TemplateParam) (string, error) {

	kind, err := templateKindToPath(param.TemplateKind)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("%s/%s/official/%s.yml", param.Revision, kind, param.TemplateName)
	fullPath := fmt.Sprintf("%s/%s", common.MegalopolisRawUrl, path)

	return fullPath, nil
}

func templateKindToPath(kind string) (string, error) {
	switch kind {
	case "box/v1":
		return "boxes", nil
	default:
		return "", fmt.Errorf("invalid template kind: %s", kind)
	}
}

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

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil || len(data) == 0 {
		return "", errors.Wrap(err, "invalid body")
	}

	return string(data), nil
}
