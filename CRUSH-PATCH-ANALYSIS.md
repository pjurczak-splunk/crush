# Crush Patch Analysis

## Original request

> Review CRUSH-POST-PATCH-RUN.md internal/agent/circuit_provider_test.go internal/discover/circuit.go internal/discover/circuit_test.go analyse and summarize the change (first commit), next review the files that are changes to be committed also analyse and summarize the change (second commit), do not commit, not change anything

---

Analysis of the working-tree changes, which bundle **two logically distinct efforts** into the staging area.

---

## First commit — Circuit custom provider

Feature work described in `CRUSH-POST-PATCH-RUN.md` and validated by the new test files.

### What `CRUSH-POST-PATCH-RUN.md` documents
It's a run journal, not code. It records the manual bring-up of a custom **Circuit** provider (Cisco's `chat-ai.cisco.com` Azure-style endpoint):

- Env setup (`CIRCUIT_API_APP_KEY`, `CIRCUIT_API_ACCESS_TOKEN` obtained via an OAuth client-credentials call to `id.cisco.com`).
- Smoke test: `go run . run -m circuit/gpt-5-nano "say hello..."` → responds `Hello.` ✅
- A toolchain snag (Homebrew Go wrapper reports 1.26.5 but `GOTOOLCHAIN=local` reveals real 1.22.5; fix is to *not* pin `local` and let Go fetch the required toolchain, or `brew upgrade go`).
- A known open issue: the dev binary shows the circuit provider when run from the repo dir but **not** from `/tmp` (config discovery is working-dir relative).

### The actual code enabling Circuit
(in `internal/agent/coordinator.go`, `internal/discover/circuit.go`, `crush.json`, `README.md`, `SKILL.md`, `schema.json`)

- **`buildCircuitProvider`** — a new provider builder. Circuit is Azure-deployment-shaped but authenticates differently:
  - Uses the OpenAI SDK with `openaiazure.WithEndpoint(baseURL, apiVersion)` so the request path becomes `/openai/deployments/<model>/chat/completions?api-version=...`.
  - **Strips** any `Authorization` header and sends the access token as the **`api-key`** header instead.
  - Injects the app key into the request **body's `user` field** — either from `provider_options.user` (verbatim JSON string) or synthesized from `provider_options.appkey` as `{"appkey":"..."}`.
  - Resolves `$`-env references in `provider_options` at build time via `resolveProviderOptionString` → `cfg.Resolve`.
  - Validates that `api_key`, `base_url`, and an app key/user are all present.
  - Default API version constant `2025-04-01-preview`.
- Wired into `buildProvider`'s type switch (`case "circuit":`).
- **`internal/discover/circuit.go`** — registers `"circuit"` as a known custom provider type via a no-op `circuitEnricher` (models must be declared explicitly; no discovery endpoint). This registration is what makes `IsKnownCustomProvider("circuit")` and the schema enum include circuit.
- **Tests**: `circuit_provider_test.go` spins up an httptest server and asserts the deployment path, empty `Authorization`, `api-key: access-token`, and `user: {"appkey":"app-key"}` body. `circuit_test.go` asserts circuit is registered/known.
- **Config/docs**: `crush.json` adds a live circuit provider block; `README.md`, `SKILL.md`, and `schema.json` document/enumerate the new type.

Also mixed into this area of `coordinator.go`/`agent.go` are **added `slog.Debug` logging statements** around model creation and the agent stream lifecycle, plus a cosmetic gofumpt-style rewrap of an `append(...)` call and an `azure` → `fantasyazure` import alias (to disambiguate from the newly imported `github.com/openai/openai-go/v3/azure`).

---

## Second commit — dependency downgrades, illumos revert, and misc cleanup

A separate, unrelated set of changes:

- **`go.mod` / `go.sum`** — downgrades a batch of dependencies (Azure SDK, openai-go, posthog-go, goose, brotli, klauspost/compress, mattn/go-isatty, modelcontextprotocol/go-sdk to a `-pre.3`, various `golang.org/x/*` and `modernc.org/*` modules, etc.). Effectively reverts recent dependabot bumps.
- **`README.md`** — removes the illumos/Solaris install notes section.
- **Auth / config plumbing** — changes across `internal/config/store.go`, `internal/oauth/mcp/handler.go`, `internal/cmd/login.go`, `internal/config/refresh_singleflight_test.go`, `internal/config/load.go`, `internal/config/load_test.go` (login/refresh single-flight behavior so parallel sessions don't invalidate each other's login).
- **Shell** — `internal/shell/coreutils.go`, `dispatch.go`, `dispatch_test.go`, `run.go` (multi-call coreutils / binary passthrough test fix).
- **MCP** — `internal/agent/tools/mcp/init.go` bump to match the go-sdk `-pre.3`.
- **UI** — `internal/ui/model/mcp.go` uses `ErrorIcon` instead of `NeedsAuthIcon`; the `NeedsAuthIcon` style is removed from `styles.go`/`quickstyle.go`; `dialog/mcp_auth.go` and `dialog/notifications.go` gofumpt formatting; `notification/native.go` and `model/ui.go` touched.
- **`internal/agent/hyper/provider.json`** — tweaks model pricing and drops `kimi-k3` and `qwen3.7-flash`.

---

## Notes / flags

- The **unstaged** `.gitignore` change replaces the whole file with `*` (ignore everything) — separate from both commits and almost certainly **not** something to keep.
- The two changes are entangled in `README.md` and `go.mod`/`go.sum` (the same files carry both the Circuit docs and the dep downgrades/illumos revert). For two clean commits, split those files with `git add -p`.

_No files were changed or committed while producing this analysis._

---

## Suggested followups

- **Split into 2 commits** — Split the staged changes into two clean commits: one for the Circuit provider and one for the dependency downgrades and misc cleanup. Use `git add -p` to separate the entangled `README.md` and `go.mod`/`go.sum`. Don't push.
- **Fix .gitignore** — Revert the unstaged `.gitignore` change that replaces the whole file with a single `*` back to the original ignore rules.
- **Run Circuit tests** — Run the Circuit provider tests (`circuit_provider_test.go` and `circuit_test.go`) and report whether they pass.

# Review CRUSH-POST-PATCH-RUN.md internal/agent/circuit_provider_test.go internal/discover/circuit.go internal/discover/circuit_test.go analyse and summarize the change (first commit), next review the files that are changes to be committed also analyse and summarize the change (second commit), do not commit, not change anything 