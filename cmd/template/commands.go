package template

import (
	"fmt"
	"io/ioutil"
	"log"

	t "github.com/hckops/hckctl/pkg/template"
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

func RunTemplateRemoteCmd(name string) {
	data, err := fetchTemplate(name)
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
		return "", fmt.Errorf("file not found")
	}
	return string(data), nil
}

func fetchTemplate(name string) (string, error) {
	var data string

	req, err := t.NewTemplateReq(name)
	if err != nil {
		return "", err
	}

	// attempts remote validation and to access private templates
	data, err = req.FetchApiTemplate()
	if err != nil {

		data, err = req.FetchPublicTemplate()
		if err != nil {
			return "", fmt.Errorf("template not found")
		}
	}
	return data, nil
}
