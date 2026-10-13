package ledger

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// LayerListFilter is one pluggable parameter for a layer listing. Typed
// constructors model the keys the ledger's server-side filter registry
// allowlists; RawFilter reaches registry keys this SDK does not yet model.
// Filters compose freely: ListLayers(ctx, ReviewStatus("pending"),
// LayerIDs("a", "b"), WithReviews()).
type LayerListFilter interface {
	apply(values url.Values) error
}

// eqFilter renders one key:eq:value filter.
type eqFilter struct {
	key   string
	value string
}

func (f eqFilter) apply(values url.Values) error {
	if f.key == "" || f.value == "" {
		return fmt.Errorf("filter %q: key and value must be non-empty", f.key)
	}
	values.Add("filter", f.key+":eq:"+f.value)
	return nil
}

// inFilter renders one key:in:v1,v2 filter.
type inFilter struct {
	key    string
	values []string
}

func (f inFilter) apply(values url.Values) error {
	if f.key == "" {
		return fmt.Errorf("filter %q: key must be non-empty", f.key)
	}
	if len(f.values) == 0 {
		return fmt.Errorf("filter %q: in-list requires at least one value", f.key)
	}
	for _, v := range f.values {
		if v == "" {
			return fmt.Errorf("filter %q: in-list must not contain empty values", f.key)
		}
		if strings.Contains(v, ",") {
			return fmt.Errorf("filter %q: value %q contains the list separator", f.key, v)
		}
	}
	values.Add("filter", f.key+":in:"+strings.Join(f.values, ","))
	return nil
}

// rangeFilter renders one key:gte|lte:value filter.
type rangeFilter struct {
	key, op, value string
}

func (f rangeFilter) apply(values url.Values) error {
	if f.key == "" || f.value == "" {
		return fmt.Errorf("filter %q: key and value must be non-empty", f.key)
	}
	values.Add("filter", f.key+":"+f.op+":"+f.value)
	return nil
}

// rawFilter passes an explicit key:op:value through, for registry keys the
// SDK does not model. The server still allowlists key and op.
type rawFilter struct {
	key, op, value string
}

func (f rawFilter) apply(values url.Values) error {
	if f.key == "" || f.op == "" || f.value == "" {
		return fmt.Errorf("raw filter: key, op and value must be non-empty")
	}
	values.Add("filter", f.key+":"+f.op+":"+f.value)
	return nil
}

// timeRangeFilter renders one updated_at:gte|lte:RFC3339 filter, rejecting
// the zero time at apply time (constructors cannot error).
type timeRangeFilter struct {
	op string
	ts time.Time
}

func (f timeRangeFilter) apply(values url.Values) error {
	if f.ts.IsZero() {
		return fmt.Errorf("filter %q: timestamp must be non-zero", "updated_at")
	}
	values.Add("filter", "updated_at:"+f.op+":"+f.ts.Format(time.RFC3339))
	return nil
}

// reviewsFlag requests embedded needs_review queues (include_reviews=true).
// It selects rather than filters — layers with empty queues are omitted —
// and composes with the review filters.
type reviewsFlag struct{}

func (reviewsFlag) apply(values url.Values) error {
	values.Set("include_reviews", "true")
	return nil
}

// ReviewStatus restricts the listing to layers with a needs_review item in
// this status (e.g. "pending"). Matching items are embedded in each LayerRef.
func ReviewStatus(status string) LayerListFilter {
	return eqFilter{key: "review_status", value: status}
}

// ReviewStatusIn is ReviewStatus for several statuses at once.
func ReviewStatusIn(statuses ...string) LayerListFilter {
	return inFilter{key: "review_status", values: statuses}
}

// ReviewQueueReason restricts the listing to layers with a needs_review item
// carrying this queue reason.
func ReviewQueueReason(reason string) LayerListFilter {
	return eqFilter{key: "review_queue_reason", value: reason}
}

// LayerID restricts the listing to one layer ID.
func LayerID(id string) LayerListFilter {
	return eqFilter{key: "layer_id", value: id}
}

// LayerIDs restricts the listing to the given layer IDs.
func LayerIDs(ids ...string) LayerListFilter {
	return inFilter{key: "layer_id", values: ids}
}

// ProductRepoID restricts the listing to the layer of one product_repo.
func ProductRepoID(id string) LayerListFilter {
	return eqFilter{key: "product_repo_id", value: id}
}

// MerkleEpochAtLeast restricts the listing to layers at or after this
// merkle epoch (incremental sweeps).
func MerkleEpochAtLeast(epoch int) LayerListFilter {
	return rangeFilter{key: "merkle_epoch", op: "gte", value: strconv.Itoa(epoch)}
}

// MerkleEpochAtMost restricts the listing to layers at or before this
// merkle epoch.
func MerkleEpochAtMost(epoch int) LayerListFilter {
	return rangeFilter{key: "merkle_epoch", op: "lte", value: strconv.Itoa(epoch)}
}

// UpdatedSince restricts the listing to layers updated at or after this
// timestamp. Rendered as RFC 3339; the zero time is rejected at call time.
func UpdatedSince(ts time.Time) LayerListFilter {
	return timeRangeFilter{op: "gte", ts: ts}
}

// UpdatedUntil restricts the listing to layers updated at or before this
// timestamp. Rendered as RFC 3339; the zero time is rejected at call time.
func UpdatedUntil(ts time.Time) LayerListFilter {
	return timeRangeFilter{op: "lte", ts: ts}
}

// RawFilter passes an explicit key:op:value filter through to the ledger's
// allowlisted registry — the escape hatch for server-side keys this SDK
// version does not model.
func RawFilter(key, op, value string) LayerListFilter {
	return rawFilter{key: key, op: op, value: value}
}

// WithReviews embeds each layer's needs_review queue in the response.
func WithReviews() LayerListFilter {
	return reviewsFlag{}
}
