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

	"H4sIAAAAAAAC/+xc3Y8UxxH/V1aTSEmkgT0bIyUr+QFzgEk4jA6OPCAe+mZqbxvPTg/dPWc2p5UIlyjg",
	"BMFLZCl+iJzETkxEHAsnxDbxP7Ne7vxfRN09Xz3T83H7xWbhxb7d6an6VdWvP6q6lj3LIf2A+OBzZnX2",
	"LAosID4D+eEt5G7CzRAYF58c4nPw5Z8oCDzsII6J377BiC++Y04P+kj89X0KXatjfa+dim6rp6x9hlJC",
	"reFwaFsuMIfiQAixOtZV5GFXSmyBGGO3gDvHraFtnfc5UB95l4HuAp0/kssDxqGfR3GR8LMk9N35698E",
	"RkLqQMslwFo+4S24hRlPoWz5KOQ9QvEvYAFwToW8Bz6PpOpuGdqRdEkXJaKzZwWUBEA5VizyYBc88Qcf",
	"BGB1LMYp9neEIX1gDO2A4dnQtijcDDEVFl6LRKQvXB/a1ibcPAf8zC746j/n1zeBAd2VMM+vb0ARydbm",
	"haIu27p1jKAAH3OICzvgH4NbnKJjHO3Id24gqyNftK1dRVHxbkg9u49uvfn6yZMSrEP6/SgER5c+uvNk",
	"tP/RaP/JaP+uriZVkXNIrM+W0CJ3XCKMXyA72D+NPG8bOe8WPSBQTAby8NH9w789k1Cfjfbv6TgTaAag",
	"LmTxXcY7fhgsMcAtBnQD1jFzCHWXFWcYwSziwg7xt6g3KRH/NLrzFwntSYHyMRi7j/03X9PYz5weId45",
	"ijSfYJ/DDtBGmsePPxl/+UWtwp8odT1C+XmfU+KGjlqlJrH24I+PR7/81eFvHo2ffn7wxZfffv3bxgaH",
	"TOxI/QkZMNr/RMT+zr9H+8/GD+831JrjQwLBTmKuB8Lkpzx/8k7McWlqFzd1btE8XAX7EhrEq62OmFPk",
	"sy7QixPH5vnvPjv8RkTl4JtPJwuMhqGAnOJdKS6PHLkuBcYmA/3tfx88v3e7MXu7mDI+uYvGD++P790/",
	"urafIR9No/GHozt/H+2//6PGmjHbQF7WyG1CPEB+MyLc/uv47seaLj/0PCHXQ9O67/EHjY2IlU3nvccf",
	"HNl7AaLg81PT8PLg6fvPb3/83e0/H/zro/Gv7x6RpgrAafC8oEd8uBj2t9Xxf2ogzx9+Nf7PZ999+PXh",
	"p58f/P7R+MHTElB+2AeKnQy417LY3iZ9mDW28Ydfjf/xh+mxTU7QvLP++aD5JmVbU/tjctPzSUM8TXOT",
	"KLsA5penZNHQDbGT9VlzbxlHy/iRn1Rqd2BxGlPcFSD+GnPoM8MpFHngu4hexfCeYZ0b2krC+XVj+mU+",
	"wUgvykRq27h4ps+5Sj/zz3NhiBFE+mwdtKZMkyycE4lGlKKBWXDeh2dSew2uXKwj0ufRASqJYlUiXjTm",
	"ne0b4PBMdms1ck3W52VO1iGafVlUb+BhgBzMB4bTf76eYPByF/uY9dajc1HRi10KcLpSwyICeNaEsqBJ",
	"FR/qRlzmiPJSUczwtGx26WqTeGddnhWoObshLc4a3tCNsFMC5KJlZtRm3ldLMFVFPjXFHNVMUjNGnPiP",
	"MFNLIlkeIom4gYMzaJpVxKrqWVUhEICMYTOkrtqGYChmCYtkLauIOTSVNwo6qJeRpKpOsxHFEQ8NWzFL",
	"vq+ZsdHAVOIVQryiPFdVoLaaIMyMTcWaAz4BzYWgIzA6R8uyQhVyON6FC7iPzSyLbSrjlF1e6qqtSDWr",
	"IcmVOHTFYa/kOFvF+aoiUemEyBR1dNXlRR7dTybDbM3X+chElc6ZzI3pi0oNykCpotIyUA8zTuhgAp4n",
	"QhXh344E1XE+VliCUBdWPD31wHm3bF+qLmkJmgUiPXJPmafRABA1zYCcAXJYTpmdAMsqKZh45IpWdQWq",
	"vmJUVdmprM7UVlPqqx7NyxLNiwQ1KXtdTv3/nPbmtpbmtxmNV/C5rc+Z9bikzB4bmG5hr3bA5dsB1Wvz",
	"3rMaXM3eeI/XyxaDMiLrblOPIFPEGpyQYj64LLbFqAMEEAV6KuQ98WlbfjpLaB9xq2P99OdXrKgBQK7B",
	"8qmVbJU9zgNVlcN+l0gcmItV21rHO9ghVARvFyhTfQYnjr92fE2whgTgowDLr9aOn5DLB+9JOO2kILYD",
	"8n96q8I54C05ouVhJvIJ4QuVBbnqsaqzyQUJ9YEDZVbnmgiw1bFuhkAHcVbSsUi3y4DHBiLjJnrd1htm",
	"Xl9bm1lDRrYwaOpSCR0HGGvt6DYPbeuNtRNlshOwba2HRL70Rv1LSQ/M0LZOKlOrX8i17mRJJv2epde1",
	"68KdqhScFveGdhT09l6UKg/rw18a+TNJtm0igKBZGv9sZh5PFk5DMBAimVgL4kNsRhNavAyMaO9pxZMG",
	"DGlly541bNnMFWbmQx3bKClfFFpWKuo+arRcUb2e/bJxtK2ORS54wKHI1HX5fas/aMBXNbaqE+4Vaaci",
	"rVsRjFWjrm0FoWHlPE0BNeXjpZCvDBll9/NbxB3MkId1fasGJm6Y/a4bM1zO2eNUUGcVF34vvrYoPYLI",
	"ES3VYVk4eqhbj7mHUqkxhO1CAk48iw0U7s/a13aySSdhBkuvAsXdQUu+WVgiCvnw3CZbUVVJk3tLpKki",
	"tX0Bc6sRyHhCMXWbNQWrDVFlyRVZKW3VkDLeRndscydupMfknxRfpZEz4G6u8DJX8uZ0LSl761HOl77J",
	"9auRvurllhrVkvcaRg4rKfPnsNJjCqXntdQPkFjrnRge8ib3lvQOj66YS6e2GFA2seX99NxdIrUYHHIl",
	"RpYNvLRHBT6+0C41TQworQfKa4eZlQNt86sM5B3aC60jqqv76rw89dSqHMuiXoSYJ1HOXUqV/qAVUNLF",
	"8iLMyJUNsBYSKnMWkA1WBuuKRaskA92SV9BVMUp/FjW//TisCM9GFtqiU7N61oQFB0oONAhp5nfAk3Jt",
	"4RO97aaNLaUTPhpTnZDpfTILCmSszhDP9QLomThKPxWH5Yfi2Gnmw3FY8gvKuR6RzSqX8KTcaKJOf0Ce",
	"nAz5q+6qvRLrd+wlEyd3E78Q72o66zdSzZCXbTetjKL5l6Lz317rArhRwP1ipnFTooVmd6/0Bhyk7ZmV",
	"J+5oWOkCcikZsJCQxuoanL/jkS/b+bskYoXfZ89/nagI1kYW6GxXh0mbig19xE1O61marexikbYRV6fn",
	"clh9mh63JS9q0YjUNUnadQteJe/5xSMTuTkvHuVBW4JUvgGlshm9gVUruVbsqVbhYX3NF7s/YHULRdI0",
	"XN+SkLQoL1M7S8aGJsVdzSWrWuSNGTLzXNbQKr4CxHmVK6cUkm+Il1UwdUe4sAseCSxb/TRO9rd32u29",
	"HmF82NkLCOXDtvwXKyhG2170U7TkermLQo9bHcsjDvLk12KzJDT3+Mdra2siDCm6vep/kI/J28KUbrIs",
	"N7QLO7CY//mxofpBZ36s7PApDI7bifdMN3T5wfKGbnh9+L8AAAD//zrQA+JgUgAA",
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
