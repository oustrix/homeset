//go:build tools
// +build tools

package tools

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=api/openapi/config.yaml api/openapi/api.yaml

import (
	_ "github.com/gojuno/minimock/v3/cmd/minimock"
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
)
