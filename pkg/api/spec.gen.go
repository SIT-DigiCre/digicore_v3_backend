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
	"lQtPYQAA",
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
