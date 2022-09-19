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

	"H4sIAAAAAAAC/8xY3W4cNRR+lZHhcpLZ/kSClXpBgVZFKaCElIsoWjkzZ2cdPPbEPrNkWc0FCRJCqBJv",
	"gIRQhSr1ArUCVeRtlqh5DGR7sj/zs7tdslHvZsb2Od853+fj4xmSUCapFCBQk/aQKNCpFBrsy30a7cBx",
	"BhrNWygFgrCPNE05CykyKYIjLYX5psMeJNQ8va+gS9rkvWBiOnCjOvhUKalInuc+iUCHiqXGCGmTJ5Sz",
	"yFr0wMzxPcBwk+Q+eSQQlKB8F1Qf1PqR7A40QlJG8bnEBzIT0fr974CWmQrBiyRoT0j04IRpnEDZEzTD",
	"nlTsO7gBOB9l2AOBhdXZtOR+Yd3KxZloD0mqZAoKmVMRhz5w84CDFEibaFRMxCaQBLSmMdSM5T5RcJwx",
	"ZSLcL0xMFhzkPtmB4y+lxm0ZM/Ex5fyQht9UnYcyqrHvk5MNSVO2YYZjEBtwgopuII3toiNK2uTy+dPL",
	"P85Hpy9HZ+ejs5+IT/pOosbOGFxehmodTuHbZbHI0ncVYIZ7GtRjqOJioRSdTPHVsI1Ofxud/m6xvdzb",
	"2a7H5idM3LvlJ/Tk3u2trSs1Sck7saIzWWECIQa1lOuLF88uXr9a6PFD568nFXY08G6HCVQyykIn+1WC",
	"fvPri9H3P1z++Pzi7z/fvHr97z8/Lx13pk2JS1ZUwujsmdHA6V+js/OLX54u6bWkizEEf8J9iZDmfDk9",
	"6Yfg9mNVTrVKqkBQfMqS2znXYwopZrpqSo+/F9YOpeRARcVcMXFisWnb0BBZHzqcJQxra17EdChV1DHp",
	"ZlHtlOatt3iHvI2mzVzMIhDYEVly6M7WyhQDtNOAtFm2deoyVvxpnZW8zxNeKW3NUfqzFBSMLXFSHH2L",
	"i8Mwk6ZMLirub2HTkAFhphgOds2RWnRgQBUocwSbt0P79kCqhCJpk8++/ooUB7DVrh0l/pW7HmLqtjkT",
	"XWlxMORm5BMWs1Aqk9c+KO3O+TubtzZbhlSZgqAps59amy3ik5Riz8IJ+NXujsFGNtsqPAT07AzP1XyT",
	"Cts1PIrcqCsO/myXebvVurYuZroG1fQy22NwZszV0H1iuinHgIsvCKcJlbom0iegWHfg2ZXlOKtac3yD",
	"xvsyGlxjrDUNUEMD5xkJGNm4LTKRH6oM8vUSshTI3SwMQWtPu6Kf+2TLoagzPkYblO4GDazq8UnSKFs3",
	"pUm3xVG0duEWfuryM8E3N8hr0G6pqK1VvCVf76h6F6Ncr3zHXUqtfN1iz83ywh5Y0qoadlbWr2Hnp45K",
	"zj13udbeF1fwKF89WzY7piMIXAvSuLuTgZcq2WUc6jJT9HJrz0zhZ4584lmsuU/utu4sTszMDwG76O7i",
	"ReMfGv9DrFcdC2nvz/Yq+wf5wUTLhiLbNqVZDUF7qbmpzONock1dXzHK5tDzeBrajVafpVSTVRJoNbAE",
	"pVM/+FbV2s1Ix60wi7VdMJuICPrApTmf7a3Jdr7tIBj2pMa8PUylwjywt2LF6CF3TPXGh2OXZtz001yG",
	"lNvPRqpSlYY/aLVahoYJuuH8X2Xa1jriE3dFciU99yv616Z6l+bauPOD/L8AAAD//0Ppa36gFQAA",
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
