package server

//go:generate go run ../../am/overlays/mergeoverlay.go ../../am/overlays/models.yaml ../../am/overlays/operations.yaml overlay.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config cfg.yaml ../../am/openapi/openapi.yaml
