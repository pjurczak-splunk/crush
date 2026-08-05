
# create custom provider with claude & codex 

(charm ~/pj-code/charm-starter) codex resume 019fac0e-63a8-7562-ab4a-7ca4e63c8471

1. The installed Go here is 1.22.5, while this repo requires go 1.26.5.
go version (go installation)

2. 
export CIRCUIT_API_APP_KEY=egai-prd-cx-541213691-coding-1769512278368
export CIRCUIT_API_ACCESS_TOKEN='...oauth access token...'

3. get token 
go to https://ai-chat.cisco.com/bridgeit-platform/home
generate access token using client_id and client secret
sedning request to url -X POST \
    https://id.cisco.com/oauth2/default/v1/token \
    -H 'Accept: */*' \
    -H 'Content-Type: application/x-www-form-urlencoded' \
    -H 'Authorization: Basic 'not-sure-what-token-it-is' \
    -d 'grant_type=client_credentials'

4. test provider
go run . run -m circuit/gpt-5-nano "say hello in one sentence"

go run . run --verbose -m circuit/gpt-5-nano "say hello in one sentence"
Wed Jul 29:~/pj-code/charm-starter/v0.87.0-circuit-run $ vim crush.json
INFO Running in non-interactive mode
INFO Configured Circuit provider base_url=https://chat-ai.cisco.com api_version=2025-04-01-preview has_user=true has_access_token=true
INFO Configured Circuit provider base_url=https://chat-ai.cisco.com api_version=2025-04-01-preview has_user=true has_access_token=true
INFO Overriding large model for non-interactive run provider=circuit model=gpt-5-nano
WARN Using large model as small model for unknown provider provider=circuit model=gpt-5-nano
INFO Configured Circuit provider base_url=https://chat-ai.cisco.com api_version=2025-04-01-preview has_user=true has_access_token=true
INFO Configured Circuit provider base_url=https://chat-ai.cisco.com api_version=2025-04-01-preview has_user=true has_access_token=true
INFO Configured Circuit provider base_url=https://chat-ai.cisco.com api_version=2025-04-01-preview has_user=true has_access_token=true
INFO Configured Circuit provider base_url=https://chat-ai.cisco.com api_version=2025-04-01-preview has_user=true has_access_token=true
INFO Configured Circuit provider base_url=https://chat-ai.cisco.com api_version=2025-04-01-preview has_user=true has_access_token=true
INFO Configured Circuit provider base_url=https://chat-ai.cisco.com api_version=2025-04-01-preview has_user=true has_access_token=true
INFO Configured Circuit provider base_url=https://chat-ai.cisco.com api_version=2025-04-01-preview has_user=true has_access_token=true
INFO Configured Circuit provider base_url=https://chat-ai.cisco.com api_version=2025-04-01-preview has_user=true has_access_token=true
INFO Configured Circuit provider base_url=https://chat-ai.cisco.com api_version=2025-04-01-preview has_user=true has_access_token=true
INFO Configured Circuit provider base_url=https://chat-ai.cisco.com api_version=2025-04-01-preview has_user=true has_access_token=true
INFO Created session for non-interactive run session_id=6d4f6cb6-1bef-4c5f-a861-b46d2747eb34
INFO Configured Circuit provider base_url=https://chat-ai.cisco.com api_version=2025-04-01-preview has_user=true has_access_token=true
INFO Configured Circuit provider base_url=https://chat-ai.cisco.com api_version=2025-04-01-preview has_user=true has_access_token=true
INFO Configured Circuit provider base_url=https://chat-ai.cisco.com api_version=2025-04-01-preview has_user=true has_access_token=true
INFO Configured Circuit provider base_url=https://chat-ai.cisco.com api_version=2025-04-01-preview has_user=true has_access_token=true
INFO ModelProvider called provider=circuit model=gpt-5-nano
ERRO Title generation hit token limit with small model; trying next
ERRO Title generation hit token limit with large model; trying next
Hello.INFO Skill turn summary component=skills session_id=6d4f6cb6-1bef-4c5f-a861-b46d2747eb34 prompt_len=25 active_total=4 loaded_total=0 loaded_this_turn=[]

WARN shutdown timeout exceeded, some messages may be lost
WARN shutdown timeout: 1 batches still in flight
ERRO Failed to flush PostHog events error="shutdown timeout: context deadline exceeded"
WARN sending request - Post "https://data.charm.land/batch/": context canceled
ERRO 6 messages dropped: shutdown timeout


# fixing with codebuff

# merging to branch on crush
To continue this session later, run:
codebuff --continue 2026-07-29T11-58-22.645Z
Session: 216 credits · Remaining: 278 credits · Cycle: Aug 13

        go run . models
        go run . run -m circuit/gpt-5-nano "say hello in one sentence"
        Hello.


# run crush dev version

Build your patched copy into a separate binary, then put that directory ahead of Homebrew on PATH.

# BUILD

build not creating binary 'GOTOOLCHAIN=local /opt/homebrew/opt/go/bin/go build -o /private/tmp/crush-dev/crush .' but it returns go: go.mod requires go >= 1.26.5 (running go 1.22.5; GOTOOLCHAIN=local)' how to fix it, not change anything yet


• Your /opt/homebrew/opt/go/bin/go wrapper reports go1.26.5, but when GOTOOLCHAIN=local is set it reveals the real local toolchain is go1.22.5. That is why the build refuses go.mod’s go >= 1.26.5.

  Do not use GOTOOLCHAIN=local here. Let Go download/use the required toolchain:

  mkdir -p /private/tmp/crush-dev
  GOCACHE=/private/tmp/go-cache GOMODCACHE=/private/tmp/gomodcache \
    /opt/homebrew/opt/go/bin/go build -o /private/tmp/crush-dev/crush .

  Then run your patched binary:

  PATH="/private/tmp/crush-dev:$PATH" crush --version
  PATH="/private/tmp/crush-dev:$PATH" crush run --verbose -m circuit/gpt-5-nano "say hello in one sentence"

  If that still tries to write under /Users/pjurczak/go/pkg/mod or fails downloading the toolchain, use:

  GOTOOLCHAIN=auto GOCACHE=/private/tmp/go-cache GOMODCACHE=/private/tmp/gomodcache \
    /opt/homebrew/opt/go/bin/go build -o /private/tmp/crush-dev/crush .

  Permanent fix is to upgrade Homebrew Go so the local binary is actually new enough:

  brew update
  brew upgrade go
  go version

  Expected after upgrade: go version go1.26.5 darwin/arm64 or newer without relying on auto toolchain download.

ISSUES:
crush dev binary works with circuit custom provider works from ~/pj-code/charm-starter/crush:circuit-patch directory
try global config: crush dev binary (PATH="/private/tmp/crush-dev:$PATH" crush) does not show circuit custom provider from /tmp
https://github.com/pjurczak-splunk/crush#configuration
