// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"j61v6XboioPlDt3o+uj/AQAA///nZ0gzoVAAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
