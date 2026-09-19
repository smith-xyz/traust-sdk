// Code generated from traust-contracts v0.18.2. DO NOT EDIT.

package types

import (
	"github.com/traust-security/traust-sdk/go/v1/enums"
)

type ThreatRegister struct {
	ByProduct                  map[string]interface{}   `json:"by_product,omitempty"`
	DocVariance                map[string]interface{}   `json:"doc_variance,omitempty"`
	Meta                       *ThreatRegisterMeta      `json:"meta,omitempty"`
	Models                     *int                     `json:"models,omitempty"`
	ModelsSkippedNonconforming *int                     `json:"models_skipped_nonconforming,omitempty"`
	Note                       *string                  `json:"note,omitempty"`
	OwnershipCuts              map[string]interface{}   `json:"ownership_cuts,omitempty"`
	QuickWins                  []map[string]interface{} `json:"quick_wins,omitempty"`
	Scoring                    *string                  `json:"scoring,omitempty"`
	TenantBoundaries           []map[string]interface{} `json:"tenant_boundaries,omitempty"`
	Threats                    []Threat                 `json:"threats"`
	Totals                     map[string]interface{}   `json:"totals,omitempty"`
	Updated                    *string                  `json:"updated,omitempty"`
	Version                    *int                     `json:"version,omitempty"`
}

type Threat struct {
	Actors              []string               `json:"actors,omitempty"`
	Asset               *string                `json:"asset,omitempty"`
	Controls            *string                `json:"controls,omitempty"`
	Evidence            []string               `json:"evidence,omitempty"`
	Id                  string                 `json:"id"`
	Impact              enums.ThreatImpact     `json:"impact"`
	IsolationBoundaries []string               `json:"isolation_boundaries,omitempty"`
	IsolationDimensions []string               `json:"isolation_dimensions,omitempty"`
	Key                 string                 `json:"key"`
	Likelihood          enums.ThreatLikelihood `json:"likelihood"`
	Linddun             *bool                  `json:"linddun,omitempty"`
	Model               string                 `json:"model"`
	Product             *string                `json:"product,omitempty"`
	Score               *int                   `json:"score,omitempty"`
	Status              enums.ThreatStatus     `json:"status"`
	SubjectId           *string                `json:"subject_id,omitempty"`
	Surface             *string                `json:"surface,omitempty"`
	Threat              string                 `json:"threat"`
}

type ThreatRegisterMeta struct {
	Generated                  *string `json:"generated,omitempty"`
	Key                        *string `json:"key,omitempty"`
	Models                     int     `json:"models"`
	ModelsSkippedNonconforming *int    `json:"models_skipped_nonconforming,omitempty"`
	Root                       *string `json:"root,omitempty"`
	Scoring                    *string `json:"scoring,omitempty"`
	TenantBoundaryCount        *int    `json:"tenant_boundary_count,omitempty"`
	ThreatCount                int     `json:"threat_count"`
}
