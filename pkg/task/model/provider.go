package model

type TaskProvider string

const (
	Docker TaskProvider = "docker"
)

func (p TaskProvider) String() string {
	return string(p)
}
