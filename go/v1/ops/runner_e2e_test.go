//go:build e2e && (darwin || linux)

package ops

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/traust-security/traust-sdk/go/v1/api"
)

// Runs against traust-server built from core at CORE_REF; `make e2e` starts it.
func e2eClient(t *testing.T, tokens string) *api.BoundClient {
	t.Helper()
	url, path := os.Getenv("TRAUST_E2E_URL"), os.Getenv(tokens)
	if url == "" || path == "" {
		t.Fatalf("TRAUST_E2E_URL and %s are required; run make e2e", tokens)
	}
	client, err := api.NewBoundClient(api.ClientConfig{BaseURL: url, Tokens: api.FileTokenSource{Path: path}})
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestE2ERunnerRunsASubmittedPackCommand(t *testing.T) {
	raw, err := os.ReadFile("testdata/local-command.yaml")
	if err != nil {
		t.Fatal(err)
	}
	pack, err := LoadPack(raw)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	person := e2eClient(t, "TRAUST_E2E_TOKEN_FILE")
	want, err := person.Submit(ctx, api.SubmitRequest{
		Order:  api.OrderRequest{Asset: "command_result", Priority: 1},
		Config: &[]api.ConfigDigest{pack.Pin()},
	})
	if err != nil {
		t.Fatal(err)
	}

	runner, err := NewRunner(RunnerConfig{
		Client: e2eClient(t, "TRAUST_E2E_RUNNER_TOKEN_FILE"), Pack: pack, Executor: ProcessExecutor{},
		Assets: []string{"command_result"}, LeaseSeconds: 30, PollInterval: time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	order, err := runner.RunOnce(ctx)
	if err != nil || order == nil || *order.Id != *want.Id || order.State != api.RunOrderStateSucceeded {
		t.Fatalf("order=%+v err=%v", order, err)
	}
	if idle, err := runner.RunOnce(ctx); err != nil || idle != nil {
		t.Fatalf("drained queue returned order=%+v err=%v", idle, err)
	}
}
