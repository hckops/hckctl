package model

// TODO out-of-scope for initial release

type LabV1 struct {
	Kind  string
	Name  string
	Tags  []string
	Boxes []LabBox
}

type LabBox struct {
	Name     string
	Template struct {
		Name string
	}
	Size    string
	Vpn     bool
	Envs    []LabEnv
	Secrets []LabSecret
}

type LabEnv struct {
	Name  string
	Value string
}

type LabSecret struct {
	Name      string
	LocalRef  string
	RemoteRef string
	Alias     string
}
