package model

type LabV1 struct {
	Kind    string
	Name    string
	Tags    []string
	Boxes   []LabBox
	Network LabNetwork
	Dump    LabDump
}

type LabBox struct {
	Alias    string
	Template struct {
		Name string
		Env  []string
	}
	Size  string
	Vpn   string
	Ports []LabPort
	Dumps []string
}

type LabPort struct {
	Name   string
	Public bool
}

type LabNetwork struct {
	Vpn []LabVpn
}

type LabVpn struct {
	Name string
	Ref  string
}

type LabDump struct {
	Git []GitDump
}

type GitDump struct {
	Name     string
	Category string
	Url      string
	Branch   string
}
