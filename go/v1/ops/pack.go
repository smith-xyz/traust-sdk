package ops

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/traust-security/traust-sdk/go/v1/api"
	"gopkg.in/yaml.v3"
)

// The schema core validates packs with, generated from core at CORE_REF (make api-generate).
//
//go:embed pack.schema.json
var packSchemaJSON []byte

var packSchema = sync.OnceValues(func() (*jsonschema.Schema, error) {
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(packSchemaJSON))
	if err != nil {
		return nil, err
	}
	compiler := jsonschema.NewCompiler()
	compiler.AssertFormat()
	if err := compiler.AddResource("pack.schema.json", doc); err != nil {
		return nil, err
	}
	return compiler.Compile("pack.schema.json")
})

// Pack is a workload pack (assets.yaml), validated and pinned by the digest of its bytes.
type Pack struct {
	name   string
	digest string
	assets map[string]packAsset
}

type packAsset struct {
	Run        []string                   `json:"run"`
	Params     map[string]json.RawMessage `json:"params"`
	Partitions string                     `json:"partitions"`
	Model      json.RawMessage            `json:"model"`
}

// command is an asset's pinned argv. "{run_id}" is the one parameter a fleet-wide command takes.
type command []string

// LoadPack validates a pack against the schema core loads packs with.
func LoadPack(raw []byte) (*Pack, error) {
	var doc any
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("ops: parse pack: %w", err)
	}
	encoded, err := json.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("ops: parse pack: %w", err)
	}
	instance, err := jsonschema.UnmarshalJSON(bytes.NewReader(encoded))
	if err != nil {
		return nil, fmt.Errorf("ops: parse pack: %w", err)
	}
	schema, err := packSchema()
	if err != nil {
		return nil, fmt.Errorf("ops: pack schema: %w", err)
	}
	if err := schema.Validate(instance); err != nil {
		return nil, fmt.Errorf("ops: invalid pack: %w", err)
	}
	var pack struct {
		Name   string               `json:"pack"`
		Assets map[string]packAsset `json:"assets"`
	}
	if err := json.Unmarshal(encoded, &pack); err != nil {
		return nil, fmt.Errorf("ops: decode pack: %w", err)
	}
	digest := sha256.Sum256(raw)
	return &Pack{name: pack.Name, digest: hex.EncodeToString(digest[:]), assets: pack.Assets}, nil
}

// ref is the name core records the pack's digest under in an order's manifest.
func (p *Pack) ref() string { return "pack:" + p.name }

// Pin is the pack version an order must run under: pass it in Submit's config so runners with
// this exact pack, and only those, run the order.
func (p *Pack) Pin() api.ConfigDigest { return api.ConfigDigest{Name: p.ref(), Digest: p.digest} }

// Assets are the assets the pack declares, sorted.
func (p *Pack) Assets() []string { return slices.Sorted(maps.Keys(p.assets)) }

// command returns the asset's argv if a fleet-wide runner can run it: no partitions, model or
// params, and "{run_id}" as a whole argument is its only placeholder.
func (p *Pack) command(name string) (command, error) {
	asset, ok := p.assets[name]
	switch {
	case !ok || len(asset.Run) == 0:
		return nil, fmt.Errorf("ops: pack %s has no command for asset %q", p.name, name)
	case asset.Partitions != "" || len(asset.Model) != 0 || len(asset.Params) != 0:
		return nil, fmt.Errorf("ops: asset %q needs inputs a fleet-wide runner does not resolve", name)
	}
	for _, arg := range asset.Run {
		if arg != "{run_id}" && strings.ContainsAny(arg, "{}") {
			return nil, fmt.Errorf("ops: asset %q has an unresolvable argument %q", name, arg)
		}
	}
	return asset.Run, nil
}

func (c command) resolve(runID uuid.UUID) []string {
	argv := make([]string, len(c))
	for i, arg := range c {
		if arg == "{run_id}" {
			arg = runID.String()
		}
		argv[i] = arg
	}
	return argv
}
