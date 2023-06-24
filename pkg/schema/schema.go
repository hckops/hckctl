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
	KindCommandV1
	KindBoxV1
	KindLabV1
	KindTaskV1
	KindFlowV1
)

var kinds = map[SchemaKind]string{
	KindConfigV1:  "config/v1",
	KindCommandV1: "command/v1",
	KindBoxV1:     "box/v1",
	KindLabV1:     "lab/v1",
	KindTaskV1:    "task/v1",
	KindFlowV1:    "flow/v1",
}

func (s SchemaKind) String() string {
	return kinds[s]
}
