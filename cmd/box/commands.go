package box

import (
	"fmt"
	"log"

	// TODO
	l "github.com/rs/zerolog/log"

	b "github.com/hckops/hckctl/internal/box"
	t "github.com/hckops/hckctl/internal/template"
)

func runCloudBoxCmd(name, revision string) {
	l.Info().Msg("CLOUD")
	log.Println("CLOUD")
}

func runKubeBoxCmd(name, revision string) {
	log.Println("KUBE")
}

func runDockerBoxCmd(name, revision string) {
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

func runBoxListCmd() {
	for _, box := range getBoxes() {
		fmt.Println(box)
	}
}

// TODO api endpoint or fetch from https://github.com/hckops/megalopolis/tree/main/boxes
// TODO see github api
func getBoxes() []string {
	// TODO struct: name, alias e.g. alpine -> official/alpine
	// TODO revision param

	return []string{
		"alpine",
		"parrot",
	}
}
