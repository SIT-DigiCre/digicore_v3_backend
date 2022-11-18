// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.2 DO NOT EDIT.
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
	"AAD//5X9IYDChgAA",
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
