package command

type BoxDeleteBody struct {
	Name string `json:"name"`
}

func (b BoxDeleteBody) cmdName() commandName {
	return commandBoxDelete
}

func NewBoxDeleteRequest(name string) *Request[BoxDeleteBody] {
	return newRequest[BoxDeleteBody](BoxDeleteBody{name})
}
