// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
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
	"n/7F/n8CAAD//7GfOA5mqgAA",
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
