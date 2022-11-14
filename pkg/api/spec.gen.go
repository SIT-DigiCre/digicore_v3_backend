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
	"AAD//+TZjsT+hgAA",
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
