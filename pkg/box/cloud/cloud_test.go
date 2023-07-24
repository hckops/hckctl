package cloud

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	v1 "github.com/hckops/hckctl/pkg/api/v1"
	"github.com/hckops/hckctl/pkg/box/model"
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
	expected := &model.BoxDetails{
		Info: model.BoxInfo{
			Id:      "myId",
			Name:    "myName",
			Healthy: true,
		},
		TemplateInfo: &model.BoxTemplateInfo{
			GitTemplate: &model.GitTemplateInfo{
				Url:      "infoUrl",
				Revision: "infoRevision",
				Commit:   "infoCommit",
				Name:     "infoName",
			},
		},
		ProviderInfo: &model.BoxProviderInfo{
			Provider: model.Cloud,
		},
		Size: model.Medium,
		Env: []model.BoxEnv{
			{Key: "KEY_1", Value: "VALUE_1"},
			{Key: "KEY_2", Value: "VALUE_2"},
		},
		Ports: []model.BoxPort{
			{Alias: "alias-1", Local: "TODO", Remote: "123", Public: false},
			{Alias: "alias-2", Local: "TODO", Remote: "456", Public: false},
		},
		Created: createdTime,
	}
	result, err := toBoxDetails(message)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
