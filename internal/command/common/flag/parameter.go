package flag

import (
	"fmt"
	"strings"

	"github.com/hckops/hckctl/pkg/common/model"
)

func ValidateParametersFlag(inputs []string) (model.Parameters, error) {
	parameters := model.Parameters{}
	for _, input := range inputs {
		keyValue := strings.Split(input, "=")
		if len(keyValue) != 2 {
			return nil, fmt.Errorf("invalid parameter format [%s], expected KEY=VALUE", input)
		}
		key := strings.TrimSpace(keyValue[0])
		if len(key) == 0 {
			return nil, fmt.Errorf("invalid parameter key format [%s], expected KEY=VALUE", input)
		}
		value := strings.TrimSpace(keyValue[1])
		if len(value) == 0 {
			return nil, fmt.Errorf("invalid parameter value format [%s], expected KEY=VALUE", input)
		}
		parameters[key] = value
	}
	return parameters, nil
}
