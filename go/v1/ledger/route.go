package ledger

import (
	"net/url"

	"github.com/traust-security/traust-sdk/go/v1/ledger/internal/route"
)

// Layer operations accepted by LayerPath.
const (
	LayerOpDocument   = route.Document
	LayerOpEvents     = route.Events
	LayerOpFindings   = route.Findings
	LayerOpCumulative = route.Cumulative
	LayerOpVerify     = route.Verify
	LayerOpSubmit     = route.Submit
	LayerOpResolve    = route.Resolve
	LayerOpRestate    = route.Restate
	LayerOpSign       = route.Sign
	LayerOpStamp      = route.Stamp
)

// LayerPath returns the ledger request path for op on layerID
// (`/v1/ledger/layer<op>?layer_id=…`). The ID is query-encoded, never spliced into the
// path, so opaque IDs such as `corpus:layer:org/repo` round-trip intact. Test doubles
// use it to register responses for the exact path the client sends.
func LayerPath(op, layerID string, extra url.Values) string {
	return route.Layer(op, layerID, extra)
}
