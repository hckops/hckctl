package model

type DumpV1 struct {
	Kind  string
	Name  string
	Tags  []string
	Group string
	Git   GitDump
}

type GitDump struct {
	Repository string
	Branch     string
}
