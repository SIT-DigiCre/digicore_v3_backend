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

	"H4sIAAAAAAAC/8xWTW/bPAz+K4be92jE3roBg2/rPooOHTY063YoclBsJlanSCold8gC/feBkvPh2l0y",
	"LAF2s0T60UPyIaUVK/XCaAXKWVasGII1WlkIi0vlABWXY8AHQNoptXKgHH1yY6QouRNaZXdWK9qzZQ0L",
	"Tl//I8xYwf7LtvBZtNrsHaJG5r1PWQW2RGEIhBVsvLQOFgmQPU3AlSPmU3ajeONqjeInVKfn8LpxNSjX",
	"ona5+LRFD9mJEMWKGdQG0ImYNAkPIOnDLQ2wglmHQs0pkAVYy+cwYPMpQ7hvBFKEty3E9oeJT9k13H/W",
	"1l3puVBvuJRTXn7vH17q6gD84LUDOhZz1ZgjotoLiEz7WA3K/VDktEWK9I4D5bhrbB/KbvZbtKnWErjq",
	"wbWOW8QbC/gR+oii2s9NVGugA0p798PtRySnHch9hf0DTNI+lA0KtxxTD0SAc+AISD1Dq2lYvde44I4V",
	"7MO3L6ztmJDSYGXp+rjaORM7UKiZDjyEk2R5K+ai1AgsZQ+ANjbm2ejZKKc20gYUNyJs5aOcpcxwVwc6",
	"mVyLbg4hsm5vX4BLgkdyc33FAhKGNr+sojVqNu1Owed5frSxs9saA8PnakPOB2uMJyt3C6jtQGRfAcVs",
	"mdCg7MXV11asL1h3rqvlEWMbmFBPTNiESk4yoa1duTlswJ+2AAeRHDdlCdYmNs4en7KXkcUQ+IZt9ujG",
	"bKtoNwPsSVlGl6d02U7AkwuzPWcoH1t+naCOoM1HQ+qk4nx01j+qzv0sjyvPzeU3KM/onESvpKwhFKmv",
	"0Yhyeo3Gc4ZKJ2USn5A2+bSmx+VfZqexgNkCftu9i2ViUM+EhKHMtE+Ek2emPecgubzIz/YnpPPy3n0B",
	"sOK2e/ffTvzERw/KnQ0OXQ4VPWo1DbHwbgvXf5Flq1pb54uV0eh8Rjc+R8GnMiap3kyUGW8kPSqkLrkM",
	"2z5l9FPX/CrPc8rAxP8KAAD//1tMoDLaDAAA",
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
