package common

type BoxV1 struct {
	Version string   `json:"version"`
	Name    string   `json:"name"`
	Tags    []string `json:"tags"`
	Image   struct {
		Repository string `json:"repository"`
	} `json:"image"`
}
