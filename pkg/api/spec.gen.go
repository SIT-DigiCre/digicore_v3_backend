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

	"H4sIAAAAAAAC/+xaW28UyRX+K6NKHhKp7RluUjISD5AERAQJsmPyYFmjcveZmYLuqqaq2vHEasmxEwVI",
	"LHhZIS0Pq73ALqy8LIJdlsvyZ4axzb9YVdW4Z/o203OzELtvnu7qc75zzlenvj7tDWQzz2cUqBSouoE4",
	"CJ9RAfrHWewswPUAhFS/bEYlUP0n9n2X2FgSRstXBaPqmrCb4GH116851FEV/arcM102d0X5T5wzjsIw",
	"tJADwubEV0ZQFV3BLnG0xRKoNVYJpD2PQgtdoBI4xe4i8DXgs0ey2BISvCSKvzB5jgXUmb3/BRAs4DaU",
	"HAaiRJkswToRsgdlieJANhkn/4QjgHMmkE2gsms1npbQ6lrXdDEmqhvI58wHLolhkQtr4Ko/ZMsHVEVC",
	"ckIbKhAPhMANyLgXWojD9YBwFeFy10TvgZXQQgtw/TIT8iJrEPoH7Lqr2L6Wdm4zJ8O+hdbnGPbJnLrd",
	"ADoH65LjOYkb+qGrGFXRwaOdg69et7eetrdft7dvIgutGYoqOxG4MAlVO+zDt0gaNPDfV4CBXBLAL0Ea",
	"F7EZrQXcHQ9be+uz9tYXGtvTpYWL2dgsj9DTxywPr58+furUIZsYc2sNjmNZIVRCA3gh153dB50Xz4Z6",
	"/L3x12Rc1gS49RqhkjMnsA3txwl6/5Pd9r/+ffDfR53nT/afvXj76n+F4w6EanHemExobz9QHNj6vr39",
	"unNnp6DXBC8iCFav9omC5OcryafLuOV1W1KcVpJjKurAa+NHu/f/xwdvVJz7bx4WTnALMB/PXefFs87L",
	"BzmOaOABJ3Y6m9qflQg3lSVO1rTBZJaw43AQYjzAb3+8vXdzs3Bm6oQLOUE5Ond2Ojd3xnBXu4YpnsTn",
	"b9pbX7e3b/22sG8iah52++NcZcwFTIvxbvPLzo37MWc0cF1l2MUTp3D3buEwIm8TZnD37sgZ9DEHKmsT",
	"8XP/+a29zfvvNj/f/+7Tzn9ujEjXLgIbXNdvMgo1GnirwKcCZe/Oy84Pj9/de3Xw8Mn+R486t58P3vZ9",
	"8I7F0DWZB1NH17n3svPNx1NANz5Rk/n69nbx48ZCk2dk/OCTgjLasKntFOuI6X7V6yKJgKyoa8cTnU/Z",
	"fLqkNpo5OcR5MHo3fWJkKrXUEc/dPktGmU7HlMQyEGlTIrqe6LdJc92FPYt5shTbkqxBzSUekZnvFA4R",
	"NuNOTckZ4mQuyZe2wxXoKJpRrZWBo5mQQ/qu9KvlIM2XhVnqTVmx+nVcwvsgYZdIW36UVrwEyYrlCj+7",
	"CfY1cLK4YA1ThSoTvtroTg1nVz0h76JqFVFlVoQt5iYV2chibYi4KiKGBoqWwcKjgFAocKCPcuKOcv4N",
	"O42GHhcfSEcvMMK4+g85PHy1qM/ksKnDCDZVFwM74ES2Fu0meN3RIGAO/Ewgm+rXqv51jnFP7U/057//",
	"DXUnQ5qv+q7ad8ZdU0rfHMmE1pnGQaRiOPojaRCbcZXiNeDCDKBOzB+bryhCMB8o9om+VJk/oXMqmxpO",
	"2T08FhugI4vPsM6DLOkVJTOMUKnQ46wLjrlrTlUrPv48XqlMbbzWf3hnDNkuRuDUPaN5lhFW6dUVMPGV",
	"7f6CMpER6RXgpN4q6SeTcaa5ZuoNQp5lTmuKsWZM5nImiyVFAUUbc7b06Cd5AOFsC1II5GJg2yBESRi1",
	"FFrolEGRZTxCW04MrXOqKiIJlktbsySPt10NN3Pidv1k5aeHb2CQU+BuoqnNlLwJX+8pe4ejnC19I3mf",
	"SV/zcMmsKmmRlclhY2X2HDZ+skrpuiXz1UeU/noID7vjZ0tnR0npstE2ubvba5V8zupEi49UZrovQTPP",
	"TNfPAPo04lhDC52snBiemNiXKv3QyeEPRV/aJiDroWJB1eW4VlleCVd6XFYl0rLJDzIKtKRfBgbVqPf9",
	"ZHbNKBhQnkv90I60+xRiTZBKoOZAgZL2fXkel2tHQ52+jV72e+++Azd8d1nuhr8cLZiogkT1s1FKeeg3",
	"jJQ65hy3ivSF3pM/r76QU8rUd7DZt4eoeNldIoZjel3iKDgWpNL9YXeR3pxpsGzQy4bLh8O51RGdB4fu",
	"ioiJeAS/iIpk8+ir3IybR37R3gOJUYBS/Uojg1UfVK/QT6iHhX4gnhAH1sBlPrLMZxw9a6uWyxtNJmRY",
	"3fAZl2FZfzrjBK+6pnDN6HW8jgNXoipymY1dfVnxmPHE7d9VKhVVjh66jcH/NSb02xWykJn2mpfI0Ept",
	"DqHeFxNrddzhSvhTAAAA//9/AfzbqygAAA==",
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
