package inter

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestApiOauthResponseError(t *testing.T) {
	t.Run("returns an error if data is invalid", func(t *testing.T) {
		data := []byte(`{
	"error: "error description",
	"error_title": "error title"
}`)

		_, err := parseApiOAuthResponseError(data)
		require.Error(t, err)
	})

	t.Run("correctly parses input data", func(t *testing.T) {
		want := apiOAuthResponseError{
			Error:      "response error description",
			ErrorTitle: "response error title",
		}

		data := []byte(fmt.Sprintf(`{
	"error": "%s",
	"error_title": "%s"
}`, want.Error, want.ErrorTitle))

		got, err := parseApiOAuthResponseError(data)
		require.NoError(t, err)
		require.Equal(t, got, want)
	})
}

func TestOAuthAutorize(t *testing.T) {
	t.Run("returns an error on context cancelation", func(t *testing.T) {
		client := NewClient(tls.Certificate{})

		oauth := NewOAuth(client)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := oauth.Authorize(ctx, "client-id", "client-secret")
		require.ErrorIs(t, err, context.Canceled)
	})

	var response string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, response)
	}))
	defer ts.Close()

	t.Run("returns an error if data is invalid", func(t *testing.T) {
		response = `{"access_token}`

		client := NewClient(tls.Certificate{})
		client.apiBaseUrl = ts.URL

		oauth := NewOAuth(client)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		_, err := oauth.Authorize(ctx, "client-id", "client-secret")
		require.Error(t, err)
	})

	t.Run("returns the created token", func(t *testing.T) {
		var (
			tokenData       = "test-token-123"
			tokenType       = "test-token-type"
			tokenExpiresSec = 300
			tokenScopes     = []string{"test-token-scope"}
		)

		response = fmt.Sprintf(`{
	"access_token": "%s",
	"token_type": "%s",
	"expires_in": %d,
	"scope": "%s"
}`, tokenData, tokenType, tokenExpiresSec, strings.Join(tokenScopes, " "))

		want := Token{
			Data:      tokenData,
			Type:      tokenType,
			Scopes:    tokenScopes,
			ExpiresAt: time.Now().Add(time.Duration(tokenExpiresSec) * time.Second),
		}

		client := NewClient(tls.Certificate{})
		client.apiBaseUrl = ts.URL

		oauth := NewOAuth(client)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		got, err := oauth.Authorize(ctx, "client-id", "client-secret")
		require.NoError(t, err)
		compareToken(t, got, want)
	})
}

func compareToken(t *testing.T, a, b Token) {
	t.Helper()

	require.Equal(t, a.Data, b.Data)
	require.Equal(t, a.Type, b.Type)
	require.Equal(t, a.Scopes, b.Scopes)
	require.WithinDuration(t, a.ExpiresAt, b.ExpiresAt, time.Second)
}
