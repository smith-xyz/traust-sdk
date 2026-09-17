package skills_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/openshift/traust-sdk/go/v1/skills"
	"github.com/openshift/traust-sdk/go/v1/skills/skillstest"
	"github.com/openshift/traust-sdk/go/v1/types"
)

type rawProvider struct {
	payload []byte
}

func (p rawProvider) Execute(context.Context, skills.SkillMeta, []byte) ([]byte, error) {
	return bytes.Clone(p.payload), nil
}

func artifactValue[T any](t *testing.T, artifact types.Artifact[T]) T {
	t.Helper()
	value, err := artifact.Value()
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func TestRunPreservesResponseBytes(t *testing.T) {
	payload, err := json.MarshalIndent(skillstest.FixtureReport(), "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	payload = append(payload, ' ', '\n')
	artifact, err := skills.Scan.Run(
		context.Background(),
		rawProvider{payload: payload},
		skills.ScanInput{Repo: "https://github.com/org/repo", Ref: "main"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(artifact.Payload(), payload) {
		t.Fatal("skill response bytes changed")
	}
	if artifactValue(t, artifact).Title != "Fixture Security Assessment" {
		t.Fatal("typed value changed")
	}
}
