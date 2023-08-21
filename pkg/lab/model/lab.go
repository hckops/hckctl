package model

type LabV1 struct {
	Kind    string
	Name    string
	Tags    []string
	Boxes   []LabBox
	Infra   []LabInfra
	Network LabNetwork
	Dump    LabDump
}

type LabBox struct {
	Alias    string
	Template string
	Env      []string
	Size     string
	Vpn      string
	Ports    []string
	Dumps    []string
}

type LabNetwork struct {
	Vpn []LabVpn
}

type LabVpn struct {
	Name   string
	Config string
}

type LabDump struct {
	Git []GitDump
}

type GitDump struct {
	Name   string
	Group  string
	Url    string
	Branch string
}

type LabInfra struct {
	Alias          string
	Source         string // helm|compose
	RepositoryUrl  string
	TargetRevision string
	Path           string
}
