package template

import (
	"fmt"

	t "github.com/hckops/hckctl/pkg/template"
)

// TODO client side schema validation
func fetchBox(name string) {

	req, err := t.NewTemplateReq(name)
	if err != nil {
		fmt.Println(err)
		return
	}

	// attempt remote validation and access to private templates
	apiTemplate, err := req.FetchApiTemplate()
	if err != nil {
		publicTemplate, err := req.FetchPublicTemplate()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Print(publicTemplate)
		return
	}
	fmt.Print(apiTemplate)
}
