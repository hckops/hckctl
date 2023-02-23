package template

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"

	"github.com/hckops/hckctl/internal/common"
	"github.com/rs/zerolog/log"
)

// TODO remove prefix
type TemplateReq struct {
	TemplateName  string
	TemplateKind  string
	Revision      string
	ClientVersion string
}

// TODO better validation: should not start with dashes or contain double dashes
var IsValidName = regexp.MustCompile(`^[A-Za-z-]+$`).MatchString

func FetchTemplate(name, revision string) (string, error) {
	var data string

	req, err := newTemplateReq(name, revision)
	if err != nil {
		log.Err(err).Msg("fetch template")
		return "", err
	}

	// attempts remote validation and to access private templates
	data, err = req.fetchApiTemplate()
	if err != nil {
		log.Warn().Msg(err.Error())

		data, err = req.fetchPublicTemplate()
		if err != nil {
			log.Err(err).Msg("fetch template")
			return "", fmt.Errorf("unable to fetch the template")
		}
	}
	return data, nil
}

func newTemplateReq(name, revision string) (*TemplateReq, error) {

	if !IsValidName(name) {
		return nil, fmt.Errorf("invalid name")
	}

	return &TemplateReq{
		TemplateName:  name,
		TemplateKind:  "box", // TODO enum
		Revision:      revision,
		ClientVersion: "hckctl-v0.0.0", // TODO sha/tag
	}, nil
}

// TODO e.g. https://api.hckops.com/template/box?name=official/alpine&version=main&format=json
// TODO or redirect validate https://schema.hckops.com/validate?kind=box&group=official&name=alpine
// TODO or content https://schema.hckops.com/template?kind=box&group=official&name=alpine&version=main&format=json|yaml
func (req *TemplateReq) fetchApiTemplate() (string, error) {

	templateUrl, err := url.Parse(fmt.Sprintf("%s/todo", common.ApiUrl))
	if err != nil {
		return "", fmt.Errorf("invalid api url")
	}

	// TODO authentication
	// TODO add header e.g. x-client=hckctl-v0.0.0
	params := url.Values{}
	params.Add("name", req.TemplateName)
	params.Add("version", req.Revision)
	params.Add("format", "yaml")
	//params.Add("client", req.ClientVersion)
	templateUrl.RawQuery = params.Encode()

	template, err := httpGetString(templateUrl.String())
	if err != nil {
		return "", fmt.Errorf("error fetching api template: %v", err)
	}

	return template, nil
}

func (req *TemplateReq) fetchPublicTemplate() (string, error) {

	// TODO use TemplateKind i.e. box -> boxes
	path := fmt.Sprintf("%s/boxes/official/%s.yml", req.Revision, req.TemplateName)

	template, err := httpGetString(fmt.Sprintf("%s/%s", common.MegalopolisRawUrl, path))
	if err != nil {
		return "", fmt.Errorf("error fetching public template: %v", err)
	}

	return template, nil
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