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
	"6rj9Of+fAAAA//8pdBI5kigAAA==",
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
