package schema

import (
	_ "embed"
)

// see http://json-schema.org

//go:embed box-v1.json
var boxV1Schema string

//go:embed lab-v1.json
var labV1Schema string

type SchemaKind int

const (
	KindConfigV1 SchemaKind = iota
	KindBoxV1
	KindLabV1
)

var kinds = map[SchemaKind]string{
	KindConfigV1: "config/v1",
	KindBoxV1:    "box/v1",
	KindLabV1:    "lab/v1",
}

func (s SchemaKind) String() string {
	return kinds[s]
}
