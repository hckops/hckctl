package template

import (
	"fmt"
	"log"

	t "github.com/hckops/hckctl/pkg/template"
)

// TODO client side schema validation
func fetchBox(name string) *t.BoxV1 {

	req, err := t.NewTemplateReq(name)
	if err != nil {
		log.Fatalln(err)
	}

	var data string
	// attempt remote validation and to access private templates
	data, err = req.FetchApiTemplate()
	if err != nil {
		log.Fatalln(err)
	}
	data, err = req.FetchPublicTemplate()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print(data)

	box, err := t.ParseBoxV1(data)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print(box.Name)

	return box
}
