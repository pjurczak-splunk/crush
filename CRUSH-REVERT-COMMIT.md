# Crush Commit Review Session — Revert of Parallel-Session Work / Circuit Provider

Session date: Wed Aug 05 2026
Repo: `github.com/charmbracelet/crush` (terminal AI coding assistant)
Review tool: Circuit (LLM gateway) with `claudeOpus4.8`
Constraint: review-only — no changes or fixes without explicit approval.

---

## 1. Overview

Two commits were reviewed:

| Commit | Message | Role |
|--------|---------|------|
| `aeb3ea39477c9e7ebc7363891a45e84d90f509ac` | `feat: add Circuit provider support` | Parent — adds a first-class Circuit provider type (Azure-style deployment endpoint) |
| `b27bd4d4fa84405d28d5fe2caebd77295367f7de` | `chore: revert parallel-session config work and downgrade deps` | HEAD — reverts parallel-session config code and downgrades dependencies (keeps Circuit) |

Commit order in history:

```
b27bd4d4 chore: revert parallel-session config work and downgrade deps   <- HEAD
aeb3ea39 feat: add Circuit provider support
b83944c4 chore: bump mcp sdk to 1.7.0 (#3447)
9c8ffe44 chore(legal): @Qalipso has signed the CLA
aa57f561 chore: auto-update files
```

Two review questions were asked via Circuit:

1. **b27bd4d4**: Why was `README.md` modified in a "revert + downgrade deps" commit?
2. **aeb3ea39**: Is it safe to remove `b27bd4d4`, and will Circuit provider support continue to work?

---

## 2. Review Question 1 — b27bd4d4: Why was README.md modified?

### 2.1 Commit stats

```
commit b27bd4d4fa84405d28d5fe2caebd77295367f7de
Author: Przemysław Jurczak <pjurczak@splunk.com>
Date:   Fri Jul 31 13:18:01 2026 +0200

    chore: revert parallel-session config work and downgrade deps

 README.md                                    |   9 --
 go.mod                                       |  34 +++---
 go.sum                                       |  80 +++++++-------
 internal/agent/agent.go                      |  10 ++
 internal/agent/coordinator.go                |  38 ++++---
 internal/agent/hyper/provider.json           |  54 +++-------
 internal/agent/tools/mcp/init.go             |   1 -
 internal/cmd/login.go                        |  11 +-
 internal/config/config.go                    |   6 --
 internal/config/load.go                      |   7 --
 internal/config/load_test.go                 |  19 ----
 internal/config/model_pin_test.go            |  39 -------
 internal/config/refresh_singleflight_test.go | 152 +++-----------------------
 internal/config/store.go                     | 154 +++++----------------------
 internal/oauth/mcp/handler.go                |  51 +--------
 internal/shell/coreutils.go                  |   8 ++
 internal/shell/dispatch.go                   |   2 +-
 internal/shell/dispatch_test.go              |  16 ++-
 internal/shell/run.go                        |  16 ++-
 internal/ui/dialog/mcp_auth.go               |   6 +-
 internal/ui/dialog/notifications.go          |  12 +--
 internal/ui/model/mcp.go                     |   2 +-
 internal/ui/model/ui.go                      |  11 +-
 internal/ui/notification/native.go           |  12 +--
 internal/ui/styles/quickstyle.go             |   1 -
 internal/ui/styles/styles.go                 |   1 -
 26 files changed, 197 insertions(+), 555 deletions(-)
```

### 2.2 The README change

Removed this install paragraph (about illumos/OpenIndiana and Oracle Solaris):

```
-On illumos (OpenIndiana, OmniOS), the command above works as-is. Only native
-OS notifications are unavailable there; terminal-based notifications (OSC) and
-the terminal bell still work. On Oracle Solaris, add `-tags sqlite3_dotlk` so
-the local database uses dot-file locking:
-
-```
-go install -tags sqlite3_dotlk github.com/charmbracelet/crush@latest
-```
```

### 2.3 Related code changes in the same commit (relevant subset)

- `internal/ui/notification/native.go`: dropped per-platform `defaultNotifyFunc`;
  now imports and calls `github.com/gen2brain/beeep` directly (`beeep.Notify`),
  plus sets `beeep.AppName = "Crush"`. The old comment said beeep's dbus
  dependency does not build on illumos/solaris.
- `internal/ui/model/ui.go`: removed `notification.NativeSupported` gating in
  `selectNotificationBackend`; the "prefer OSC where native is unavailable
  (illumos/solaris)" branch is gone. Local-session fallback now checks
  `runtime.GOOS == "darwin"` only.
- `internal/ui/dialog/notifications.go`: removed hiding of the `native`
  notification option when `!notification.NativeSupported`.
- `internal/shell/run.go` / `coreutils.go`: `execMiddleware` type alias moved
  from `run.go` to `coreutils.go`; `coreutils.go` now imports
  `mvdan.cc/sh/v3/interp`. New comment in `run.go` still says
  `coreUtilsExecHandler` is nil on illumos/solaris.
- `internal/ui/styles/quickstyle.go` / `styles.go`: removed `NeedsAuthIcon`;
  `mcp.go` now uses `ErrorIcon` for `StateNeedsAuth`.

### 2.4 Prompt sent to Circuit (claudeOpus4.8)

Review of commit b27bd4d4, single question: why was README.md modified, and is
the README removal consistent with (or contradicted by) the code changes in the
same commit (removal of `NativeSupported` gating, unconditional `beeep` import,
`execMiddleware` alias move into `coreutils.go`). The prompt embedded the
README + shell + UI notification diff hunks.

### 2.5 Circuit (claudeOpus4.8) findings — Question 1

**(a) Most likely reason the README was touched.**

The README removal is *consistent with* — and in fact directly caused by — the
code changes in this commit. This commit removes the platform-abstraction layer
that made illumos/Solaris buildable:

- `notification.NativeSupported` gating is deleted from both `ui/model/ui.go`
  (backend selection) and `ui/dialog/notifications.go` (menu option hiding).
- `native.go` no longer uses a per-platform `defaultNotifyFunc` (which was a
  no-op on illumos/solaris where beeep's dbus dep doesn't build); it now calls
  `beeep.Notify` directly and imports `github.com/gen2brain/beeep`
  unconditionally.
- The `execMiddleware` type alias moved into `coreutils.go`, which imports
  `mvdan.cc/sh/v3/interp`.

Since native OS notification support and the special install path for
illumos/Solaris no longer exist as first-class supported configurations, the
README paragraph documenting them was removed to keep docs in sync with the
code. The README edit is *not* unrelated — it's the documentation half of
dropping that platform-specific handling.

**(b) Consistent or contradicted?**

Logically consistent with the code direction (dropping the illumos/solaris
carve-out from both docs and code together). However, technically *contradicted
by build reality*, which is the suspicious part:

- The old comments explicitly said beeep's dbus dependency **does not build on
  illumos/solaris**, and native notifications there needed the `NativeSupported`
  gating.
- `native.go` now imports and calls `beeep` unconditionally, and
  `ui.go`/`notifications.go` no longer fall back to OSC. If beeep still doesn't
  build on illumos/solaris, **removing the README instructions doesn't fix
  that — it just hides it**; the package may now fail to compile on those
  platforms rather than gracefully degrading.
- The commit title says "revert parallel-session config work and downgrade
  deps." Removing platform gating, changing notification backend selection, and
  moving a type alias are **not obviously part of a config revert or dep
  downgrade**.

**Flags for human reviewer (Q1):**

1. **Scope mismatch vs. commit title.** The notification/gating refactor and the
   README removal don't fit "revert parallel-session config work and downgrade
   deps." Confirm these belong in this commit (possible squash/leftover).
2. **Potential build regression on illumos/solaris.** Verify whether beeep (and
   its dbus dep) now builds there. If not, the unconditional `beeep` import in
   `native.go` and removal of `NativeSupported` may break the build there —
   directly the situation the removed README paragraph and old comments warned
   about.
3. **Inconsistent internal signaling.** `run.go`'s new comment still says
   `coreUtilsExecHandler` is nil on illumos/solaris, yet the notification path
   no longer treats those platforms specially.
4. **The dropped `sqlite3_dotlk` Solaris instruction** was unrelated to
   notifications; confirm the SQLite dot-locking build path is still
   valid/intended, since nothing in this diff touches it.

---

## 3. Review Question 2 — aeb3ea39: Is it safe to remove b27bd4d4? Will Circuit keep working?

### 3.1 Commit stats

```
commit aeb3ea39477c9e7ebc7363891a45e84d90f509ac
Author: Przemysław Jurczak <pjurczak@splunk.com>
Date:   Fri Jul 31 13:17:02 2026 +0200

    feat: add Circuit provider support

 README.md                                     |  48 ++++++++++++
 crush.json                                    |  19 +++++
 internal/agent/circuit_provider_test.go       |  83 +++++++++++++++++++++
 internal/agent/coordinator.go                 | 101 ++++++++++++++++++++++++++
 internal/discover/circuit.go                  |  22 ++++++
 internal/discover/circuit_test.go             |  15 ++++
 internal/skills/builtin/crush-config/SKILL.md |  40 +++++++++-
 schema.json                                   |   1 +
 8 files changed, 328 insertions(+), 1 deletion(-)
```

### 3.2 What aeb3ea39 adds

- **`internal/agent/coordinator.go`**: `buildCircuitProvider(baseURL, apiKey,
  headers, options)`; `resolveProviderOptionString(any)` (via
  `c.cfg.Resolve`); `providerOptionsToStringMap(map[string]any)`; const
  `circuitDefaultAPIVersion = "2025-04-01-preview"`; new `case "circuit":` in
  `buildProvider`. Uses `charm.land/fantasy/providers/openai` options plus
  `github.com/openai/openai-go/v3/azure` (`openaiazure.WithEndpoint`) and
  `github.com/openai/openai-go/v3/option` (`WithJSONSet("user", userValue)`).
  Strips `Authorization` header, sets `api-key` header from the access token.
- **`internal/agent/circuit_provider_test.go`**: httptest-based end-to-end test
  asserting the Azure deployment path, `api-key` header, empty `Authorization`,
  `user` body field, and generated response text.
- **`internal/discover/circuit.go`** + test: registers a no-op
  `circuitEnricher` so `circuit` is a known custom provider type.
- **`crush.json` / README.md / SKILL.md / schema.json**: docs + config wiring
  for the `circuit` provider type.

### 3.3 Dependency state at each commit

`go.mod` at aeb3ea39:

```
github.com/mattn/go-isatty v0.0.24
github.com/modelcontextprotocol/go-sdk v1.7.0
github.com/openai/openai-go/v3 v3.46.0
github.com/posthog/posthog-go v1.22.0
github.com/pressly/goose/v3 v3.27.3
```

Deps changed by b27bd4d4 (`git diff aeb3ea39..b27bd4d4 -- go.mod`):

```
- github.com/mattn/go-isatty v0.0.24            + v0.0.23
- github.com/modelcontextprotocol/go-sdk v1.7.0 + v1.7.0-pre.3
- github.com/openai/openai-go/v3 v3.46.0        + v3.44.0
- github.com/posthog/posthog-go v1.22.0          + v1.19.0
- github.com/pressly/goose/v3 v3.27.3           + v3.27.2
```

The Circuit code was written and tested against **openai-go v3.46.0**; b27bd4d4
downgraded to v3.44.0 while keeping the Circuit code.

### 3.4 Empirical verification performed (this session)

| Check | At HEAD (b27bd4d4, v3.44.0) | At aeb3ea39 (v3.46.0) |
|-------|------------------------------|------------------------|
| `TestBuildCircuitProviderUsesAzureDeploymentEndpoint` | `ok` | `ok` |
| `TestCircuitIsRegisteredCustomProvider` | `ok` | `ok` |
| `go build ./...` | — | `BUILD OK` |
| `go test ./internal/config/...` | — | `ok` |
| `go test ./internal/shell/...` | — | `ok` |

Commands used:

```
go test ./internal/agent -run TestBuildCircuitProviderUsesAzureDeploymentEndpoint -count=1
go test ./internal/discover -run TestCircuitIsRegisteredCustomProvider -count=1
git worktree add <tmp> aeb3ea39    # target state, without touching current checkout
go build ./... && go test ./internal/config/... ./internal/shell/...
git worktree remove <tmp>
```

Conclusion of verification: Circuit tests pass under **both** openai-go v3.44.0
(HEAD) and v3.46.0 (target state), and the full target tree builds cleanly.

### 3.5 Prompt sent to Circuit (claudeOpus4.8)

Single question, full diff of aeb3ea39 embedded as context. Question: is it
safe to remove b27bd4d4 (go back to aeb3ea39) and will Circuit provider support
continue to work? Included the empirically verified facts (tests pass at both
HEAD and aeb3ea39) and asked for residual risks + what else gets re-introduced.

### 3.6 Circuit (claudeOpus4.8) findings — Question 2

**(a) Is it safe to remove b27bd4d4 for Circuit support?**

**Yes.** The Circuit code is self-contained and dependency-consistent at
aeb3ea39:

- `buildCircuitProvider` uses `openaiazure.WithEndpoint(baseURL, apiVersion)`.
  At aeb3ea39 the pinned version is **v3.46.0**, exactly the version the
  Circuit code was written and compiled against. Removing b27bd4d4 reverts
  openai-go back to v3.46.0 in lockstep with the code — **no version skew**.
  (The downgrade to v3.44.0 introduced by b27bd4d4 is the latent risk if
  `WithEndpoint`'s signature differed between versions; going back to v3.46.0
  eliminates that concern.)
- No coupling to parallel-session work: nothing in `coordinator.go`'s Circuit
  path, `internal/discover/circuit.go`, or either test references OAuth refresh,
  model pinning, or MCP token pruning. A partial/clean revert will not break
  Circuit.

**(b) Residual risks beyond the tests:**

1. **go.mod / go.sum consistency.** Reverting b27bd4d4 must restore both the
   openai-go bump (→ v3.46.0) and the mcp go-sdk change (v1.7.0-pre.3 →
   v1.7.0) atomically. Confirm `go.sum` hashes and run `go mod verify` +
   `go mod tidy` — a stale `go.sum` is the most likely silent breakage.
2. **Runtime behavior of `WithJSONSet`/`WithEndpoint` between SDK versions.**
   The test was re-run at v3.46.0 in this session and passed, so this risk is
   now retired.
3. **Transitive dependents of v3.46.0.** Other providers in `coordinator.go`
   (Azure, OpenAI-compat, etc.) also move to v3.46.0. Covered by the full
   build + agent/config/shell test runs at aeb3ea39 in this session.

**(c) What else removing b27bd4d4 re-introduces:**

Removing b27bd4d4 is a *full revert of the revert*, so everything it undid
comes back. The reviewer must consciously accept:

- **Parallel-session config work returns** (biggest item): cross-process OAuth
  refresh hardening (`withRefreshLock`, `credentialWriteLockDeadline`,
  `newerDiskToken`/`usableDiskToken`, `pinPreferredModelLocked`, MCP orphaned
  OAuth token pruning, singleflight/rotating-token tests). This is functional
  code, not just Circuit-adjacent.
- **Dep downgrades are undone:** openai-go → v3.46.0, mcp go-sdk → v1.7.0.
- **Shell change:** `execMiddleware` alias relocation (back into `run.go`).
- **Notification change:** `NativeSupported` removal is undone (comes back).
- **README/docs edits** carried by b27bd4d4 are reverted.

**Recommendation (Circuit + session):** If the goal is *only* to align Circuit
with the SDK version it was written for, prefer a **targeted change** (bump
openai-go back to v3.46.0 + adjust go.sum) rather than removing b27bd4d4
wholesale — that avoids silently resurrecting the parallel-session work, the
shell alias move, and the notification change, which are unrelated to Circuit
and may have been intentionally reverted. Reserve the full removal of
b27bd4d4 for the case where the parallel-session work is genuinely wanted back.

---

## 4. Consolidated takeaways

1. The README edit in b27bd4d4 is the doc counterpart to dropping
   illumos/solaris special-casing, but the underlying code change may be a
   functional/build regression on those platforms and looks out of scope for the
   commit's stated purpose.
2. Circuit provider support is safe under both dependency states (v3.44.0 at
   HEAD and v3.46.0 at aeb3ea39); verified by test runs in this session.
3. Removing b27bd4d4 is safe *for Circuit* but re-introduces the parallel-session
   config work (OAuth refresh hardening, model pinning, MCP token pruning) plus
   the pre-downgrade deps — decide intentionally, or use a targeted dep restore
   instead.

## 5. Follow-ups not yet done

- Full review of aeb3ea39's other aspects (config handling, provider_options
  resolution, error handling) — not requested.
- Review of b27bd4d4 questions 2–8 (deps downgrade consistency, shell
  alias/build tags, dispatch_test.go comment/code mismatch, `NeedsAuthIcon`
  dangling refs, removed-symbol leftovers incl. `cmp.Or` in login.go,
  `refreshOAuthTokenLocked` scope, hyper/provider.json model removal) — pending.
