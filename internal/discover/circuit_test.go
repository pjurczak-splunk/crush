package discover

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCircuitIsRegisteredCustomProvider(t *testing.T) {
	t.Parallel()

	require.NotNil(t, GetEnricher("circuit"))
	require.True(t, IsKnownCustomProvider("circuit"))
	require.Contains(t, RegisteredProviderTypes(), "circuit")
}
