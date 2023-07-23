package cloud

import (
	"testing"

	"github.com/stretchr/testify/assert"

	v1 "github.com/hckops/hckctl/pkg/api/v1"
	"github.com/hckops/hckctl/pkg/box/model"
)

func TestToBoxDetails(t *testing.T) {
	message := v1.NewBoxDescribeResponse("hckadm-0.0.0-info", "myName")
	expected := &model.BoxDetails{
		Info: model.BoxInfo{
			Name: "myName",
		},
	}

	assert.Equal(t, expected, toBoxDetails(message))
}
