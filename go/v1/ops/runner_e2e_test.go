//go:build e2e && (darwin || linux)

package ops

import (
	"bufio"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
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

// processGroupCommand runs argv without a shell; cancel kills the whole process group, so a lost
// lease stops descendants too.
func processGroupCommand(ctx context.Context, name string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error { return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) }
	return cmd
}

type processExecutor struct{}

func (processExecutor) Execute(ctx context.Context, execution Execution) (ProducedReport, error) {
	output, err := processGroupCommand(ctx, execution.Argv[0], execution.Argv[1:]...).Output()
	if err != nil {
		return ProducedReport{}, err
	}
	return ProducedReport{Payload: output, Result: api.JobResult{Outcome: api.OutcomeSucceeded}}, nil
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
	submitted, err := person.SubmitWithResponse(ctx, api.SubmitRequest{
		Order:  api.OrderRequest{Asset: "command_result", Priority: 1},
		Config: &[]api.ConfigDigest{{Name: pack.ref(), Digest: pack.digest}},
	})
	if err != nil {
		t.Fatal(err)
	}
	want, err := api.Decode(submitted.HTTPResponse, submitted.Body, submitted.JSON200)
	if err != nil {
		t.Fatal(err)
	}

	runner, err := NewRunner(RunnerConfig{
		Client: e2eClient(t, "TRAUST_E2E_RUNNER_TOKEN_FILE"), Pack: pack, Executor: processExecutor{},
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

func TestE2ERunnerLocalCommandStopsDescendantsOnLeaseLoss(t *testing.T) {
	python, err := exec.LookPath("python3")
	if err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(t.TempDir(), "descendant-survived")
	const script = `import subprocess, sys, time
subprocess.Popen([sys.executable, "-c", "import pathlib,sys,time; print('child-ready', flush=True); time.sleep(0.5); pathlib.Path(sys.argv[1]).write_text('survived')", sys.argv[1]])
time.sleep(10)`
	var server *queuedServer
	runner, mock := queuedRunner(t, executorFunc(func(ctx context.Context, _ Execution) (ProducedReport, error) {
		command := processGroupCommand(ctx, python, "-c", script, marker)
		stdout, err := command.StdoutPipe()
		if err != nil {
			return ProducedReport{}, err
		}
		if err := command.Start(); err != nil {
			return ProducedReport{}, err
		}
		defer func() {
			if command.ProcessState == nil {
				_ = command.Cancel()
				_ = command.Wait()
			}
		}()
		ready, err := bufio.NewReader(stdout).ReadString('\n')
		if err != nil || ready != "child-ready\n" {
			return ProducedReport{}, errors.New("child process did not start")
		}
		server.mu.Lock()
		server.heartbeatStatus = 409
		server.mu.Unlock()
		return ProducedReport{}, command.Wait()
	}))
	server = mock
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	order, err := runner.RunOnce(ctx)
	var lost *LeaseLostError
	if order == nil || !errors.As(err, &lost) {
		t.Fatalf("order=%+v err=%v", order, err)
	}
	events, _ := server.snapshot()
	if strings.Contains(strings.Join(events, ","), "complete") {
		t.Fatalf("completed after lease loss: %v", events)
	}
	time.Sleep(700 * time.Millisecond)
	if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("descendant survived lease loss: %v", err)
	}
}
