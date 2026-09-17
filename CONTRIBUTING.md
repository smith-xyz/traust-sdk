# Contributing

## Setup

```bash
git clone <repo-url> && cd traust-sdk
make setup
```

Enables git hooks. One time per clone. Go SDK requires **Go 1.26+** (`toolchain` in `go/go.mod` pins the patched release for `govulncheck`).

Code generation fetches the canonical `traust-contracts` repository at the immutable commit pinned in `go/Makefile`; no sibling checkout is required. The verified checkout is reused from the user cache. See [README.md](README.md) for repository/ref/cache overrides.

## Commit messages

Conventional commits required. Format: `type(scope): subject`

Types: `feat`, `fix`, `perf`, `refactor`, `docs`, `test`, `chore`, `ci`, `build`, `style`, `revert`

Use the language id as scope when the change is SDK-specific:

```
feat(go): add harness provider interface
fix(go): correct enum code generation
ci: update shared gates              ← no language release required
```

## Hooks

| Hook | Runs | Speed |
|------|------|-------|
| `commit-msg` | subject format check | instant |
| `pre-commit` | gofmt on staged `.go` | <1s |
| `pre-push` | `go vet` + `go test ./...` | ~10s |

Bypass: `--no-verify` on commit or push. CI still enforces.

## Releasing

Each language versions independently. Tag format: `{tag_prefix}/vX.Y.Z` (see `ci/languages.json`).

If your PR has release-worthy commits (`feat`, `fix`, `perf`, or `!`) **for a language**:

```bash
make bump go minor
# add to go/CHANGELOG.md:
#   ## [X.Y.Z]
#   ### Added
#   - ...
make check-release
git add go/VERSION go/CHANGELOG.md
git commit -m "chore(go): release X.Y.Z"
```

Only bump languages touched by the PR.

CI/infra-only PRs: use `ci:` / `chore:` commits — no VERSION bump.

Registry: `ci/languages.json` (must stay aligned with `ci/lang/*.yml`).

## CI pipeline

Stages fan out across languages in parallel:

| Stage | Jobs (parallel) |
|-------|-----------------|
| `lint` | `lint:<lang>`, … + gates (MR only) |
| `test` | `test:<lang>`, … |
| `security` | `security:<lang>`, … |
| `release` | `release:tag` — tags each language with a pending `{lang}/VERSION` bump |

Jobs auto-enable when a language's marker file exists (`ci/languages.json`). Add a language: extend the table in `AGENTS.md`, register in `ci/languages.json`, add `ci/lang/<lang>.yml` + `ci/lang/<lang>/ci.py`.

**Main:** `release.py release-if-ready --all` compares each `{lang}/VERSION` to latest `{tag_prefix}/v*` tag.

## Running tests

```bash
make test              # network-free unit tests
make test-integration  # SQLite plus optional local PostgreSQL storage coverage
make security          # govulncheck + gosec (required in CI)
make -C go test
make -C go check-drift # fetches/verifies the pinned contracts commit
```

`make test-integration` uses one fixed local PostgreSQL test DSN. Start it with:

```bash
podman run --name traust-postgres --rm -d -e POSTGRES_USER=traust -e POSTGRES_PASSWORD=traust-test-only -e POSTGRES_DB=traust_test -p 127.0.0.1:5432:5432 -v traust-postgres-data:/var/lib/postgresql/data docker.io/library/postgres:16
```

These are fixed local test-only credentials, **not deployed credentials or
secrets**, and no DSN override is supported. If PostgreSQL is absent, unreachable,
or rejects authentication, its test reports a skip while SQLite still runs. Once a
connection succeeds, Open, initialization, and protocol failures fail the suite.

## Architecture

See [README.md](README.md) for module layout and contract regeneration.
