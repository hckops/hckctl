package model

type BoxV1 struct {
	Kind  string
	Name  string
	Tags  []string
	Image struct {
		Repository string
		Version    string
	}
	Shell   string
	Network struct {
		Ports []string
	}
}

type BoxPort struct {
	Alias  string
	Local  string
	Remote string
}
