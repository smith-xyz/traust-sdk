package ops

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/traust-security/traust-sdk/go/v1/api"
)

// RunnerConfig binds a pinned pack and the assets this host runs.
type RunnerConfig struct {
	Client            *api.BoundClient
	Pack              *Pack
	Executor          Executor
	Assets            []string
	LeaseSeconds      int
	HeartbeatInterval time.Duration
	RequestTimeout    time.Duration
	PollInterval      time.Duration
}

// Runner claims queued orders for its assets and runs each one's pinned pack command.
type Runner struct {
	config RunnerConfig
}

// NewRunner refuses assets a fleet-wide runner cannot run from its pinned pack.
func NewRunner(cfg RunnerConfig) (*Runner, error) {
	if cfg.Client == nil || cfg.Pack == nil || cfg.Executor == nil {
		return nil, errors.New("ops: runner needs a Client, a Pack and an Executor")
	}
	if len(cfg.Assets) == 0 || cfg.LeaseSeconds <= 0 || cfg.LeaseSeconds > maxLeaseSeconds {
		return nil, fmt.Errorf("ops: runner needs assets and a lease of 1..%d seconds", maxLeaseSeconds)
	}
	lease := leaseDuration(cfg.LeaseSeconds)
	cfg.HeartbeatInterval = cmp.Or(cfg.HeartbeatInterval, lease/3)
	cfg.RequestTimeout = cmp.Or(cfg.RequestTimeout, defaultRequestTimeout)
	if cfg.HeartbeatInterval <= 0 || cfg.HeartbeatInterval >= lease || cfg.RequestTimeout <= 0 || cfg.PollInterval <= 0 {
		return nil, errors.New("ops: runner needs positive request and poll intervals and a heartbeat shorter than the lease")
	}
	for i, name := range cfg.Assets {
		if slices.Contains(cfg.Assets[:i], name) {
			return nil, fmt.Errorf("ops: asset %q listed twice", name)
		}
		if _, err := cfg.Pack.command(name); err != nil {
			return nil, err
		}
	}
	cfg.Assets = slices.Clone(cfg.Assets)
	return &Runner{config: cfg}, nil
}

// Run claims and runs orders until ctx is done, waiting PollInterval after an empty claim.
func (r *Runner) Run(ctx context.Context) error {
	for {
		order, err := r.RunOnce(ctx)
		if err != nil {
			return err
		}
		if order != nil {
			continue
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(r.config.PollInterval):
		}
	}
}

// RunOnce claims and runs one order. It returns nil, nil when the queue has nothing for this runner.
func (r *Runner) RunOnce(ctx context.Context) (*api.RunOrder, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	order, err := r.claim(ctx)
	if err != nil || order == nil {
		return nil, err
	}
	l, err := hold(r.config.Client, order, r.config.LeaseSeconds, r.config.HeartbeatInterval, r.config.RequestTimeout)
	if err != nil {
		return order, err
	}
	command, err := r.config.Pack.command(order.Asset)
	if err != nil || !slices.Contains(r.config.Assets, order.Asset) {
		return r.finish(ctx, l, refused(api.RunOrderReasonAssetNotInPack))
	}
	if order.RepositoryId != nil {
		return r.finish(ctx, l, refused(api.RunOrderReasonBadRequest))
	}
	if err := l.renew(ctx); err != nil {
		return l.snapshot(), err
	}
	pinned, err := supervise(ctx, l, func(read context.Context) (bool, error) { return r.pinned(read, l) })
	if err != nil {
		return l.snapshot(), err
	}
	if !pinned {
		return r.finish(ctx, l, refused(api.RunOrderReasonBadRequest))
	}
	if err := l.renew(ctx); err != nil {
		return l.snapshot(), err
	}
	produced, err := l.run(ctx, r.config.Executor, Execution{Argv: command.resolve(l.id())})
	if lost(err) || ctx.Err() != nil {
		return l.snapshot(), err
	}
	if err == nil {
		err = checkOutcome(produced.Result)
	}
	if err != nil {
		finished, finishErr := r.finish(ctx, l, api.JobResult{Outcome: api.OutcomeFailed, Reason: reason(api.RunOrderReasonRunnerCrashed)})
		return finished, errors.Join(err, finishErr)
	}
	return r.finish(ctx, l, produced.Result)
}

func (r *Runner) claim(ctx context.Context) (*api.RunOrder, error) {
	ctx, cancel := context.WithTimeout(ctx, r.config.RequestTimeout)
	defer cancel()
	assets := slices.Clone(r.config.Assets)
	order, err := r.config.Client.Claim(ctx, api.ClaimRequest{Assets: &assets, LeaseSeconds: r.config.LeaseSeconds})
	if err != nil {
		return nil, fmt.Errorf("ops: claim: %w", err)
	}
	return order, nil
}

// pinned reports whether core planned the order with exactly this runner's pack. It runs under
// supervise, so a slow read keeps the lease and a lost lease cancels it.
func (r *Runner) pinned(ctx context.Context, l *lease) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.config.RequestTimeout)
	defer cancel()
	manifest, err := r.config.Client.GetOrderManifest(ctx, l.id())
	if err != nil {
		return false, fmt.Errorf("ops: order %s manifest: %w", l.id(), err)
	}
	if manifest.Config == nil {
		return false, nil
	}
	var digests []string
	for _, config := range *manifest.Config {
		if config.Name == r.config.Pack.ref() {
			digests = append(digests, config.Digest)
		}
	}
	return len(digests) == 1 && digests[0] == r.config.Pack.digest, nil
}

func (r *Runner) finish(ctx context.Context, l *lease, result api.JobResult) (*api.RunOrder, error) {
	completed, err := l.complete(ctx, result)
	order := l.latest(completed)
	return &order, err
}

func refused(why api.RunOrderReason) api.JobResult {
	return api.JobResult{Outcome: api.OutcomeRefused, Reason: reason(why)}
}

func reason(why api.RunOrderReason) *api.RunOrderReason { return &why }
