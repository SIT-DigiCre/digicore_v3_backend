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

	"H4sIAAAAAAAC/8xYW2/jNhP9KwK/71GJvJcArYF96LbdxRbZtoib7UNgGDQ1lrmlSGU4cuMa+u8FSfmm",
	"i+2mdtA3Sxwezsw5nBl5xYTJC6NBk2XDFUOwhdEW/MN7nt7BYwmW3JMwmkD7n7wolBScpNHJV2u0e2fF",
	"HHLufv0fYcaG7H/JFjoJqzb5EdEgq6oqZilYgbJwIGzIvnAlU48YgbOJIyBxzaqYfdIEqLkaAS4AL+/J",
	"aGkJ8qYXPxv6YEqdXv78O7CmRAFRasBG2lAET9LS1pV7zUuaG5R/wQu4811Jc9BUo+6npYprdC+XADFc",
	"sQJNAUgyqEjBApT7QcsC2JBZQqkzF0gO1vIMOtaqmCE8lhJdhA81xHbDuIrZHTz+aizdmkzq77lSUy7+",
	"aB8uTHoCvrfaAR3JTJfFuVHvLeBnaKNJYfSkxI4cxezpyvBCXjmoDPQVPBHyK+KZ37gId8ZtyKV+9yrO",
	"+dO71zc3a2KMUZMM+Z6vUhNkgM+C/jYAzw3SxIKaTaQmNGkpglTO6X1p3Z3P4ZyoDXo2R8RbBhpp6w82",
	"0Go/QhBgm9ROPlsuoNpBCqo7DxRxKm0bym7e12hTYxRw3YKrDbeIfeLlguQCJkrmkjoveSqtMJhOXLpl",
	"2mnSfwGO6/ifCNLZUpmCpoku82loJi0T5+ikx9N+WXapy6HEuzprnH5IeI209UcZ71NQM3ZCafz6Jx0P",
	"wxntQB4rjM/B7FOWLYUAe5JYa8uxL08gSpS0HLnGVM8xwBHQNTL3NPVPHwzmnNiQ/fT7b6xuY/4Mv8ri",
	"9ZlzoiLUDqlnxjsjSbmVH2QmhUFH1gLQhm755vrV9cApxRSgeSH9q8H1gMWs4DT37iRqXTIy8Onab7gf",
	"gSJvEd3f3TKPhL73fkrDaqg48f6s9nowONsssFvYOiaC241zbi0U3gfmZpLAQIgvEbsqMbYj0i+AcraM",
	"/M5mnG0BB87B0nuTLs8Ya8cY0TMGRU4CTjbh3m0lSFhCdVlCTnJyFC5CZEMnqWJ2E7zoAt94mzQm7B5W",
	"7aY99co2mPTptu5vFxdufU5Xfrb+HQzyDNptVMqLirdx1n9Uvce9vKx8N6NPp3zD5ihYRWIOnrS2hgPK",
	"5TUczumiUqkofKLa6Je1e1w9P1s+O27MSMJc03u782VUoJlJBV2Zqdv4xTNTn3NAPtm+r1XM3g7eHE/M",
	"3me13/T2+KbN3wL/QqzriYUNH/ZnlYdxNd5q2VHk56bugnRfuM+eQyTtDFsXLUf9DH3ede7FC9Bx5ZSt",
	"HHodnEDrzl9lz9Xby8gn7HCbrd+wn4gUFqCM69H+c8xPv8MkWc2NpWq4KgxSlbiBl6PkUxWomm/0OOOl",
	"cjO1MoIr/9rLFRvL3wwGA0fD1rvV4T+drK93LGbh2yuU9SpuXQHrKnjD1sddjau/AwAA//+UQ6+n6hQA",
	"AA==",
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