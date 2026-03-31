package swagger

import _ "embed"

//go:embed swagger.json
var spec []byte

// Spec returns the OpenAPI specification JSON.
func Spec() []byte {
	return spec
}
