package command

type BoxCreateBody struct {
	Name     string `json:"name"`
	Revision string `json:"revision"`
}

func (b BoxCreateBody) cmdName() commandName {
	return commandBoxCreate
}

func NewBoxCreateRequest(name, revision string) *Request[BoxCreateBody] {
	body := BoxCreateBody{
		Name:     name,
		Revision: revision,
	}
	return newRequest[BoxCreateBody](body)
}
