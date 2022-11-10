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
	"C2FUAAA=",
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
