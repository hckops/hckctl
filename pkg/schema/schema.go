package schema

import (
	_ "embed"
)

// see http://json-schema.org

//go:embed box-v1.json
var boxV1Schema string

//go:embed lab-v1.json
var labV1Schema string

//go:embed task-v1.json
var taskV1Schema string

//go:embed dump-v1.json
var dumpV1Schema string

type SchemaKind int

const (
	KindConfigV1 SchemaKind = iota
	KindApiV1
	KindSidecarV1
	KindBoxV1
	KindLabV1
	KindTaskV1
	KindFlowV1
	KindDumpV1
)

var kinds = map[SchemaKind]string{
	KindConfigV1:  "config/v1",
	KindApiV1:     "api/v1",
	KindSidecarV1: "sidecar/v1",
	KindBoxV1:     "box/v1",
	KindLabV1:     "lab/v1",
	KindTaskV1:    "task/v1",
	KindFlowV1:    "flow/v1",
	KindDumpV1:    "dump/v1",
}

func (s SchemaKind) String() string {
	return kinds[s]
}
