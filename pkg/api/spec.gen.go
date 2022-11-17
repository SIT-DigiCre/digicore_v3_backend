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
	"AP//eiGndcKGAAA=",
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
