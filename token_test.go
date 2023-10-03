package inter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestToken(t *testing.T) {
	t.Run("creates token from string", func(t *testing.T) {
		data := "token-data"

		token := TokenFromString(data)
		require.Equal(t, token.Data, data)
	})

	t.Run("returns true for a non-expired token", func(t *testing.T) {
		token := Token{
			ExpiresAt: time.Now().Add(time.Second),
		}
		require.True(t, token.Valid())
	})

	t.Run("returns false for an expired token", func(t *testing.T) {
		token := Token{
			ExpiresAt: time.Now().Add(-time.Second),
		}
		require.False(t, token.Valid())
	})
}
