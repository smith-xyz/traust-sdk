# traust-sdk

Typed SDKs for invoking security harness skills with pluggable execution backends.

## Setup

```bash
make setup   # enable git hooks (once per clone)
make test              # network-free Go unit tests
make test-integration  # storage: SQLite + optional local PostgreSQL
make status            # per-language versions and tags
```

Each language has its own `{lang}/VERSION` and git tag (`{tag_prefix}/vX.Y.Z` in `ci/languages.json`).

## Language SDKs

| Language | Path | Module |
|---|---|---|
| Go | [`go/`](go/) | `github.com/openshift/traust-sdk/go` |

## Architecture

```
traust-contracts        ← schemas, enums (source of truth)
traust-sdk              ← typed SDKs that consumers import (this repo)
```

The SDK is opinionated on **contracts** (input/output shapes validated against
schemas from `traust-contracts`) and unopinionated on **execution** (you
implement a Provider interface to run skills however you want).

## Regenerating from contracts

```bash
cd go/
make generate   # fetches and verifies the pinned canonical contracts commit
make test       # network-free unit tests
```

For storage integration coverage, `make test-integration` always exercises SQLite
and attempts PostgreSQL at one fixed local test DSN. Start that database with:

```bash
podman run --name traust-postgres --rm -d -e POSTGRES_USER=traust -e POSTGRES_PASSWORD=traust-test-only -e POSTGRES_DB=traust_test -p 127.0.0.1:5432:5432 -v traust-postgres-data:/var/lib/postgresql/data docker.io/library/postgres:16
```

These are fixed local test-only credentials, **not deployed credentials or
secrets**; the suite does not support a DSN override. An absent, unreachable, or
authentication-rejecting PostgreSQL skips that portion with a message. Once a
connection succeeds, Open, initialization, and protocol failures fail the test.

## License

Apache License 2.0 — see [`LICENSE`](LICENSE).
