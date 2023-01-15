package box

import "fmt"

func RunBoxListCmd() {
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
