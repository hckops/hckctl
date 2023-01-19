package template

import (
	"fmt"
	"io/ioutil"
	"log"

	t "github.com/hckops/hckctl/internal/template"
)

func RunTemplateLocalCmd(path string) {
	data, err := loadTemplate(path)
	if err != nil {
		log.Fatalln(err)
	}

	err = t.ValidateAllSchema(data)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print(data)
}

func RunTemplateRemoteCmd(name, revision string) {
	data, err := t.FetchTemplate(name, revision)
	if err != nil {
		log.Fatalln(err)
	}

	err = t.ValidateAllSchema(data)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print(data)
}

func loadTemplate(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("unable to load the template")
	}
	return string(data), nil
}
