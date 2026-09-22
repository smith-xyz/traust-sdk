package storage

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// Views composed by other views rather than read directly. Must stay in step
// with INTERMEDIATE_VIEWS in the Python suite: a view listed there and not
// here (or the reverse) means one language considers it a consumption surface
// and the other does not.
var intermediateViews = map[string]string{
	"binding_current":       "the latest-binding filter every other view joins",
	"current_finding":       "the spine open/hardening/distinct/census all read",
	"report_current":        "one report per subject, composed into current_finding",
	"ownership_current":     "the owner join, composed into every scoped view",
	"finding_first_seen":    "the open-clock, composed into finding_timeline",
	"sla_clock":             "the policy clock, composed into finding_sla",
	"policy_report_current": "one policy report per subject, composed into current_finding",
}

// exported name for a view, matching the Query* convention.
func queryMethodName(view string) string {
	var b strings.Builder
	b.WriteString("Query")
	for _, part := range strings.Split(view, "_") {
		switch part {
		case "pqc":
			b.WriteString("PQC")
		case "sla":
			b.WriteString("SLA")
		default:
			b.WriteString(strings.ToUpper(part[:1]) + part[1:])
		}
	}
	return b.String()
}

// TestEveryConsumptionViewHasAGoReader is the language-parity gate.
//
// Dialect parity was gated from the start; language parity never was, and it
// drifted: advisory_exposure and validation_exposure both landed in contracts
// after go/v0.4.0 and had Python readers with no Go counterpart, so a Go
// adopter -- the intended consumer of these views -- could not reach them.
//
// Reads the view directory out of the generated query constants rather than
// the contracts checkout, so it works from a released module with no
// contracts source beside it.
func TestEveryConsumptionViewHasAGoReader(t *testing.T) {
	views := viewsFromContracts(t)
	if len(views) < 10 {
		t.Fatalf("found only %d views; the discovery is wrong, not the SDK", len(views))
	}

	client := reflect.TypeOf(&Client{})
	var missing []string
	for _, view := range views {
		if _, skip := intermediateViews[view]; skip {
			continue
		}
		if _, ok := client.MethodByName(queryMethodName(view)); !ok {
			missing = append(missing, view+" -> Client."+queryMethodName(view))
		}
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		t.Fatalf(
			"consumption views with no Go reader:\n  %s\n"+
				"Python can read these and Go cannot, which is the drift this gate exists "+
				"to catch. Add the Query* method, or declare the view intermediate in BOTH "+
				"languages.",
			strings.Join(missing, "\n  "),
		)
	}

	for view := range intermediateViews {
		found := false
		for _, candidate := range views {
			if candidate == view {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("intermediateViews names %q, which contracts no longer declares", view)
		}
	}
}

// viewsFromContracts lists view names from the contracts checkout the
// generator used, when it is present. CI runs `make generate` first, so it
// always is there; a bare `go test` from a module cache skips.
func viewsFromContracts(t *testing.T) []string {
	t.Helper()
	root := os.Getenv("TRAUST_CONTRACTS_STORAGE")
	if root == "" {
		// The checkout matching the ref the generated code was built from,
		// never "the last one alphabetically" -- older checkouts linger in
		// the cache and one of them lacks views this SDK already reads.
		root = filepath.Join(os.Getenv("HOME"), ".cache", "traust-sdk", "contracts",
			"checkout-"+generatedContractsRef(t), "storage", "v1")
	}
	entries, err := os.ReadDir(filepath.Join(root, "sqlite", "views"))
	if err != nil {
		t.Skipf("no view directory under %s: %v", root, err)
	}
	var views []string
	for _, entry := range entries {
		if name := entry.Name(); strings.HasSuffix(name, ".sql") {
			views = append(views, strings.TrimSuffix(name, ".sql"))
		}
	}
	sort.Strings(views)
	return views
}

// generatedContractsRef reads the ref out of the generated header, so the
// gate always checks against the contracts the code was generated from.
func generatedContractsRef(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile("queries.go")
	if err != nil {
		t.Skipf("no generated queries.go: %v", err)
	}
	header, _, _ := strings.Cut(string(data), "\n")
	fields := strings.Fields(header)
	for i, field := range fields {
		if field == "traust-contracts" && i+1 < len(fields) {
			return fields[i+1]
		}
	}
	t.Skipf("no contracts ref in generated header: %q", header)
	return ""
}
