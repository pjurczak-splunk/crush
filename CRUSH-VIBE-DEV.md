feat: add Circuit provider support

feat started
#29JUL
codex resume 019fac0e-63a8-7562-ab4a-7ca4e63c8471

codebuff --continue 2026-07-29T11-58-22.645Z

feat completed
PATH="/private/tmp/crush-dev:$PATH" crush

??JUL opencode -s ses_04db87cf5ffezzoqiNrwb4dcku

Review CRUSH-POST-PATCH-RUN.md internal/agent/circuit_provider_test.go internal/discover/circuit.go internal/discover/circuit_test.go analyse and summarize the change (first commit), next review the files that are changes to be committed also analyse and summarize the change (second commit), do not commit, not change anything

Here's the Opus 4.8 analysis of both commit groups:                                                                                                                                                         
                                                                                                                                                                                                                 
     Commit 1 — CRUSH-POST-PATCH-RUN.md (Setup Runbook)                                                                                                                                                          
                                                                                                                                                                                                                 
     Runbook documenting end-to-end integration testing: building a patched dev binary with proper Go toolchain, setting CIRCUIT_API_APP_KEY/                                                                    
     CIRCUIT_API_ACCESS_TOKEN env vars, and running a smoke test (go run . run -m circuit/gpt-5-nano "say hello in one sentence") against https://chat-ai.cisco.com (                                            
     API 2025-04-01-preview). Successful output confirms the Circuit provider works.                                                                                                                             
                                                                                                                                                                                                                 
     Issues: Go toolchain mismatch (1.22.5 local vs 1.26.5 required) not pinned; macOS-specific paths limit reproducibility.
     
     Commit 2 — Code Changes                                                                                                                                                                                     
                                                                                                                                                                                                                 
     A/B) internal/discover/circuit.go + circuit_test.go — No-op enricher registering "circuit" as a known custom provider. Decision is sound (Circuit deployment                                                
     IDs are explicit, no discovery endpoint). ⚠️ Both files exist on disk but are gitignored and unstaged — they won't land in the commit.                                                                      
                                                                                                                                                                                                                 
     C) internal/agent/circuit_provider_test.go — Does not exist. The most logic-heavy new code (buildCircuitProvider) has zero test coverage.                                                                   
                                                                                                                                                                                                                 
     D) coordinator.go — buildCircuitProvider constructs a fantasy.Provider using OpenAI SDK's Azure compatibility layer: strips Authorization, sets api-key,                                                    
     injects user via WithJSONSet("user", ...). Validates api_key, base_url, and either appkey/user. Azure alias renamed azure→fantasyazure to avoid collision.                                                  
     Debug logging added around model creation.                                                                                                                                                                  
                                                                                                                                                                                                                 
     E) crush.json — Clean Circuit provider block with env-var references, single gpt-5-nano model.                                                                                                              
                                                                                                                                                                                                                 
     F) MCP refactoring — Unrelated removal of orphaned token logic, config simplification, handler streamlining.                                                                                                
                                                                                                                                                                                                                 
     Key Risks (prioritized) 
 │Priority                        │Issue                                                              │Recommendation                                           │                                            
     ├────────────────────────────────┼───────────────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────┤                                            
     │High                            │Discover enricher + test unstaged/gitignored                       │Stage if runtime-dependent                               │                                            
     ├────────────────────────────────┼───────────────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────┤                                            
     │High                            │No test for buildCircuitProvider                                   │Add unit tests for auth/validation                       │                                            
     ├────────────────────────────────┼───────────────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────┤                                            
     │Medium                          │Debug logging near credentials                                     │Ensure token redaction                                   │                                            
     ├────────────────────────────────┼───────────────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────┤                                            
     │Medium                          │Feature + MCP cleanup in one changeset                             │Split into separate commits                              │                                            
     ├────────────────────────────────┼───────────────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────┤                                            
     │Low                             │fantasyazure alias unclear                                         │Rename or add comment                                    │                                            
     └────────────────────────────────┴───────────────────────────────────────────────────────────────────┴─────────────────────────────────────────────────────────┘                                            
                                                                                                                                                                                                                 
Would you like Opus 4.8 to draft the missing buildCircuitProvider test cases or go deeper on any area?                                                                                                      
   

#30JUL
dont push, review and split the staged changes into two clean commits, first one for the Circuit provider,
next commit only these changes and the second part keep staged OK?

codebuff --continue 2026-07-30T09-34-14.283Z

opencode -s ses_048a5240bffel6RAMteoQZNGoT

chore: revert parallel-session config work and downgrade deps
feat: add Circuit provider support

#4AUG

reivew file CRUSH-CURRENT-STATE.md
opencode -s ses_048a5240bffel6RAMteoQZNGoT

CRUSH-REVERT-COMMIT.md
opencode -s ses_02d96bc7dffeqwn6ege7KGEqrm
