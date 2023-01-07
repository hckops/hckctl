package template

import (
	"fmt"

	"github.com/hckops/hckctl/pkg/template"
)

func fetchBox(name string) {

	if template.IsNotValidName(name) {
		fmt.Println("invalid name")
		return
	}

	// attempt remote validation and access to private templates
	apiTemplate, err := template.FetchApiTemplate(name)
	if err != nil {
		publicTemplate, err := template.FetchPublicTemplate(name)
		if err != nil {
			return
		}
		fmt.Print(publicTemplate)
		return
	}
	fmt.Print(apiTemplate)
}
