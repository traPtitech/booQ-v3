package openapi

//go:generate go tool oapi-codegen -config model-cfg.yaml ../../docs/openapi.yml
//go:generate go tool oapi-codegen -config server-cfg.yaml ../../docs/openapi.yml
//go:generate go tool oapi-codegen -config spec-cfg.yaml ../../docs/openapi.yml
