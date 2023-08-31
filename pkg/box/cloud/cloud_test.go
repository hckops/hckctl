package cloud

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	v1 "github.com/hckops/hckctl/pkg/api/v1"
	boxModel "github.com/hckops/hckctl/pkg/box/model"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
)

func TestToBoxDetails(t *testing.T) {
	created := "2042-12-08T10:30:05.265113665Z"
	createdTime, _ := time.Parse(time.RFC3339, created)

	message := v1.NewBoxDescribeResponse("hckadm-0.0.0-info", v1.BoxDescribeResponseBody{
		Id:      "myId",
		Name:    "myName",
		Created: created,
		Healthy: true,
		Size:    "M",
		Template: &v1.BoxDescribeTemplateInfo{
			Public:   true,
			Url:      "infoUrl",
			Revision: "infoRevision",
			Commit:   "infoCommit",
			Name:     "infoName",
		},
		Env:   []string{"KEY_1=VALUE_1", "KEY_2=VALUE_2", "INVALID", "=INVALID="},
		Ports: []string{"alias-1/123", "alias-2/456", "INVALID", "/INVALID/"},
	})
	expected := &boxModel.BoxDetails{
		Info: boxModel.BoxInfo{
			Id:      "myId",
			Name:    "myName",
			Healthy: true,
		},
		TemplateInfo: &boxModel.BoxTemplateInfo{
			GitTemplate: &commonModel.GitTemplateInfo{
				Url:      "infoUrl",
				Revision: "infoRevision",
				Commit:   "infoCommit",
				Name:     "infoName",
			},
		},
		ProviderInfo: &boxModel.BoxProviderInfo{
			Provider: boxModel.Cloud,
		},
		Size: boxModel.Medium,
		Env: []boxModel.BoxEnv{
			{Key: "KEY_1", Value: "VALUE_1"},
			{Key: "KEY_2", Value: "VALUE_2"},
		},
		Ports: []boxModel.BoxPort{
			{Alias: "alias-1", Local: "none", Remote: "123", Public: false},
			{Alias: "alias-2", Local: "none", Remote: "456", Public: false},
		},
		Created: createdTime,
	}
	result, err := toBoxDetails(message)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
