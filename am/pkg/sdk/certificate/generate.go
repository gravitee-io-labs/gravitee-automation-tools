package certificate

//go:generate go run ../../../overlays/mergeoverlay.go ../../../overlays/models.yaml ../../../overlays/operations.yaml ../overlay-paths.yaml ../overlay.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config cfg.yaml ../../../openapi/openapi.yaml
