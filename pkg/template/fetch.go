package template

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"

	"github.com/hckops/hckctl/pkg/common"
)

type TemplateReq struct {
	TemplateName  string
	TemplateKind  string
	SourceVersion string
	ClientVersion string
}

// TODO better validation: should not start with dashes or contain double dashes
var isValidName = regexp.MustCompile(`^[A-Za-z-]+$`).MatchString

func NewTemplateReq(name string) (*TemplateReq, error) {

	if !isValidName(name) {
		return nil, fmt.Errorf("invalid name")
	}

	return &TemplateReq{
		TemplateName:  name,
		TemplateKind:  "box", // TODO enum
		SourceVersion: "main",
		ClientVersion: "hckctl-v0.0.0", // TODO sha/tag
	}, nil
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

func (req *TemplateReq) FetchPublicTemplate() (string, error) {

	// TODO use TemplateKind
	path := fmt.Sprintf("%s/boxes/official/%s.yml", req.SourceVersion, req.TemplateName)

	template, err := httpGetString(fmt.Sprintf("%s/%s", common.UrlMegalopolisRaw, path))
	if err != nil {
		return "", fmt.Errorf("error fetching public template: %v", err)
	}

	return template, nil
}

// TODO e.g. https://api.hckops.com/template/box?name=official/alpine&version=main&format=json
func (req *TemplateReq) FetchApiTemplate() (string, error) {

	templateUrl, err := url.Parse(fmt.Sprintf("%s/template", common.UrlApi))
	if err != nil {
		return "", fmt.Errorf("invalid api url")
	}

	// TODO authentication
	// TODO add header e.g. x-client=hckctl-v0.0.0
	params := url.Values{}
	params.Add("name", req.TemplateName)
	params.Add("version", req.SourceVersion)
	params.Add("format", "yaml")
	//params.Add("client", req.ClientVersion)
	templateUrl.RawQuery = params.Encode()

	template, err := httpGetString(templateUrl.String())
	if err != nil {
		return "", fmt.Errorf("error fetching api template: %v", err)
	}

	return template, nil
}
