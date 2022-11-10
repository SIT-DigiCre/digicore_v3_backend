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
