package agent

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"charm.land/fantasy"
	"github.com/charmbracelet/crush/internal/config"
	"github.com/stretchr/testify/require"
)

func TestBuildCircuitProviderUsesAzureDeploymentEndpoint(t *testing.T) {
	t.Parallel()

	var gotPath string
	var gotAuthorization string
	var gotAPIKey string
	var gotBody struct {
		User string `json:"user"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.String()
		gotAuthorization = r.Header.Get("Authorization")
		gotAPIKey = r.Header.Get("Api-Key")
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "chatcmpl-1",
			"object": "chat.completion",
			"created": 1,
			"model": "gpt-5-nano",
			"choices": [
				{
					"index": 0,
					"message": {
						"role": "assistant",
						"content": "hello"
					},
					"finish_reason": "stop"
				}
			],
			"usage": {
				"prompt_tokens": 1,
				"completion_tokens": 1,
				"total_tokens": 2
			}
		}`))
	}))
	defer srv.Close()

	cfg, err := config.Init(t.TempDir(), "", false)
	require.NoError(t, err)

	coord := &coordinator{cfg: cfg}
	provider, err := coord.buildCircuitProvider(
		srv.URL,
		"access-token",
		map[string]string{"Authorization": "Bearer should-be-removed"},
		map[string]string{
			"apiVersion": "2025-04-01-preview",
			"appkey":     "app-key",
		},
	)
	require.NoError(t, err)

	model, err := provider.LanguageModel(t.Context(), "gpt-5-nano")
	require.NoError(t, err)

	resp, err := model.Generate(t.Context(), fantasy.Call{
		Prompt: fantasy.Prompt{fantasy.NewUserMessage("say hello")},
	})
	require.NoError(t, err)

	require.Equal(t, "/openai/deployments/gpt-5-nano/chat/completions?api-version=2025-04-01-preview", gotPath)
	require.Empty(t, gotAuthorization)
	require.Equal(t, "access-token", gotAPIKey)
	require.Equal(t, `{"appkey":"app-key"}`, gotBody.User)
	require.Equal(t, "hello", resp.Content.Text())
}
