package box

import (
	"log"

	b "github.com/hckops/hckctl/pkg/box"
	t "github.com/hckops/hckctl/pkg/template"
)

func RunBoxDockerCmd(name, revision string) {
	data, err := t.FetchTemplate(name, revision)
	if err != nil {
		log.Fatalln(err)
	}

	box, err := t.ParseValidBoxV1(data)
	if err != nil {
		log.Fatalln(err)
	}

	b.NewDockerBox(box).InitBox()
}

func RunBoxCloudCmd(name string) {
	log.Println("TODO")
}
