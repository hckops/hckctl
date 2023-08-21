package model

type LabProvider string

const (
	Cloud LabProvider = "cloud"
)

func (p LabProvider) String() string {
	return string(p)
}
