package model

type DumpV1 struct {
	Kind  string
	Name  string
	Tags  []string
	Group string
	Git   GitDump
}

type GitDump struct {
	RepositoryUrl string
	Branch        string
}
