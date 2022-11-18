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
	"AAD///VCtrMDhwAA",
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
