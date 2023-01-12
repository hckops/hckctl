package box

import (
	"log"

	b "github.com/hckops/hckctl/pkg/box"
	t "github.com/hckops/hckctl/pkg/template"
)

func RunBoxDockerCmd(name string) {
	data, err := t.FetchTemplate(name)
	if err != nil {
		log.Fatalln(err)
	}

	box, err := t.ParseValidBoxV1(data)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(box.Name)
	b.InitDockerBox()
}

func RunBoxCloudCmd(name string) {
	log.Println("TODO")
}
