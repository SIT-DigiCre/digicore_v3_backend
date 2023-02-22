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
