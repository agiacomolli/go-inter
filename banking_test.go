package inter

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestApiTransactionTypeMap(t *testing.T) {
	t.Run("returns unknow if type not found", func(t *testing.T) {
		transactionTypeStr := "wix"
		transactionType := apiTransactionTypeMap[transactionTypeStr]
		require.Equal(t, transactionType.String(), "unknow")
	})

	t.Run("returns the type correctly", func(t *testing.T) {
		transactionTypeStr := "PIX"
		transactionType := apiTransactionTypeMap[transactionTypeStr]
		require.Equal(t, transactionType, PixTransactionType)
	})
}

func TestApiTransactionOperationMap(t *testing.T) {
	t.Run("returns invalid if type not found", func(t *testing.T) {
		transactionOperationStr := "x"
		transactionOperation := apiTransactionOperationMap[transactionOperationStr]
		require.Equal(t, transactionOperation.String(), "invalid")
	})

	t.Run("returns the type correctly", func(t *testing.T) {
		transactionOperationStr := "C"
		transactionOperation := apiTransactionOperationMap[transactionOperationStr]
		require.Equal(t, transactionOperation, CreditTransactionOperation)
	})
}

func TestParseApiTransactions(t *testing.T) {
	t.Run("returns an error if data is invalid", func(t *testing.T) {
		data := []byte(`{"transacoes": {}}`)

		_, err := parseApiTransactions(data)
		require.Error(t, err)
	})

	t.Run("correctly parses input data", func(t *testing.T) {
		var (
			transactionDate        = "2022-02-02"
			transactionType        = "PIX"
			transactionOperation   = "C"
			transactionValue       = "123.45"
			transactionTitle       = "title"
			transactionDescription = "description"
		)

		data := []byte(fmt.Sprintf(`{
	"transacoes": [{
		"dataEntrada": "%s",
		"tipoTransacao": "%s",
		"tipoOperacao": "%s",
		"valor": "%s",
		"titulo": "%s",
		"descricao": "%s"
	}]
}`, transactionDate, transactionType, transactionOperation, transactionValue,
			transactionTitle, transactionDescription))

		transactionDateTime, err := time.Parse(time.DateOnly, transactionDate)
		require.NoError(t, err)

		transactionFloatValue, err := strconv.ParseFloat(transactionValue, 32)
		require.NoError(t, err)

		want := []Transaction{
			Transaction{
				Date:        transactionDateTime,
				Type:        apiTransactionTypeMap[transactionType],
				Operation:   apiTransactionOperationMap[transactionOperation],
				Value:       float32(transactionFloatValue),
				Title:       transactionTitle,
				Description: transactionDescription,
			},
		}

		got, err := parseApiTransactions(data)
		require.NoError(t, err)
		require.Equal(t, got, want)
	})
}

func TestBankingTransactions(t *testing.T) {
	t.Run("returns an error on context cancelation", func(t *testing.T) {
		client := NewClient(tls.Certificate{})

		banking := NewBanking(client, Token{})

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := banking.Transactions(ctx, time.Now(), time.Now())
		require.ErrorIs(t, err, context.Canceled)
	})

	var response string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, response)
	}))
	defer ts.Close()

	t.Run("returns an error if data is invalid", func(t *testing.T) {
		response = `{"transactions": "}`

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
			transactionDate        = "2022-02-02"
			transactionType        = "PIX"
			transactionOperation   = "C"
			transactionValue       = "123.45"
			transactionTitle       = "title"
			transactionDescription = "description"
		)

		response = fmt.Sprintf(`{
	"transacoes": [{
		"dataEntrada": "%s",
		"tipoTransacao": "%s",
		"tipoOperacao": "%s",
		"valor": "%s",
		"titulo": "%s",
		"descricao": "%s"
	}]
}`, transactionDate, transactionType, transactionOperation, transactionValue,
			transactionTitle, transactionDescription)

		transactionDateTime, err := time.Parse(time.DateOnly, transactionDate)
		require.NoError(t, err)

		transactionFloatValue, err := strconv.ParseFloat(transactionValue, 32)
		require.NoError(t, err)

		want := []Transaction{
			Transaction{
				Date:        transactionDateTime,
				Type:        apiTransactionTypeMap[transactionType],
				Operation:   apiTransactionOperationMap[transactionOperation],
				Value:       float32(transactionFloatValue),
				Title:       transactionTitle,
				Description: transactionDescription,
			},
		}

		client := NewClient(tls.Certificate{})
		client.apiBaseUrl = ts.URL

		banking := NewBanking(client, Token{})

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		got, err := banking.Transactions(ctx, time.Now(), time.Now())
		require.NoError(t, err)
		require.Equal(t, got, want)
	})
}
