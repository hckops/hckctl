package command

type BoxListBody struct{}

func (b BoxListBody) cmdName() commandName {
	return commandBoxList
}

func NewBoxListRequest() *Request[BoxListBody] {
	return newRequest[BoxListBody](BoxListBody{})
}
