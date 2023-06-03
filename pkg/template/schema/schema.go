package schema

import (
	_ "embed"
)

// see http://json-schema.org

//go:embed box-v1.json
var BoxV1Schema string

//go:embed lab-v1.json
var LabV1Schema string
