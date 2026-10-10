// Package route builds ledger service request paths.
//
// Layer operations are addressed by query (`/v1/ledger/layer<op>?layer_id=`) so any
// opaque layer ID — including migrated corpus IDs with ':' and '/' — travels as an
// encoded value instead of being spliced into a path segment.
package route

import "net/url"

// Layer operation suffixes appended to LayerBase.
const (
	LayerBase  = "/v1/ledger/layer"
	Document   = ""
	Events     = "/events"
	Findings   = "/findings"
	Cumulative = "/cumulative"
	Verify     = "/verify"
	Submit     = "/submit"
	Resolve    = "/resolve"
	Restate    = "/restate"
	Sign       = "/sign"
	Stamp      = "/stamp"

	layerIDParam = "layer_id"
)

// Layer returns the request path for op on layerID, with any extra query values.
func Layer(op, layerID string, extra url.Values) string {
	values := url.Values{}
	for key, vals := range extra {
		values[key] = append([]string(nil), vals...)
	}
	values.Set(layerIDParam, layerID)
	return LayerBase + op + "?" + values.Encode()
}
