package inter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type TransactionType int

const (
	PixTransactionType = TransactionType(iota + 1)
	PagamentoTransactionType
	TransferenciaTransactionType
)

func (t TransactionType) String() string {
	switch t {
	case PixTransactionType:
		return "pix"
	case PagamentoTransactionType:
		return "pagamento"
	case TransferenciaTransactionType:
		return "transferencia"
	}

	return "unknow"
}

type TransactionOperation int

const (
	CreditTransactionOperation = TransactionOperation(iota + 1)
	DebitTransactionOperation
)

func (t TransactionOperation) String() string {
	switch t {
	case CreditTransactionOperation:
		return "credit"
	case DebitTransactionOperation:
		return "debit"
	}

	return "invalid"
}

type Transaction struct {
	Date        time.Time
	Type        TransactionType
	Operation   TransactionOperation
	Value       float32
	Title       string
	Description string
}

type Banking struct {
	client *Client
	token  Token
}

func NewBanking(client *Client, token Token) *Banking {
	return &Banking{
		client: client,
		token:  token,
	}
}

func (b *Banking) Transactions(ctx context.Context, start, end time.Time) ([]Transaction, error) {
	endpoint := fmt.Sprintf("%s/banking/v2/extrato", b.client.apiBaseUrl)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return []Transaction{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", b.token.Data))

	q := url.Values{}
	q.Add("dataInicio", start.Format(time.DateOnly))
	q.Add("dataFim", end.Format(time.DateOnly))

	req.URL.RawQuery = q.Encode()

	resp, err := b.client.Do(req)
	if err != nil {
		return []Transaction{}, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Transaction{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []Transaction{}, errors.New(string(data))
	}

	return parseApiTransactions(data)
}

type apiTransaction struct {
	Date        string `json:"dataEntrada"`
	Type        string `json:"tipoTransacao"`
	Operation   string `json:"tipoOperacao"`
	Value       string `json:"valor"`
	Title       string `json:"titulo"`
	Description string `json:"descricao"`
}

type apiTransactions struct {
	Transactions []apiTransaction `json:"transacoes"`
}

var apiTransactionTypeMap map[string]TransactionType = map[string]TransactionType{
	"PIX":           PixTransactionType,
	"PAGAMENTO":     PagamentoTransactionType,
	"TRANSFERENCIA": TransferenciaTransactionType,
}

var apiTransactionOperationMap map[string]TransactionOperation = map[string]TransactionOperation{
	"C": CreditTransactionOperation,
	"D": DebitTransactionOperation,
}

func transactionFromApi(a apiTransaction) (Transaction, error) {
	date, err := time.Parse(time.DateOnly, a.Date)
	if err != nil {
		return Transaction{}, err
	}

	value, err := strconv.ParseFloat(a.Value, 32)
	if err != nil {
		return Transaction{}, err
	}

	return Transaction{
		Date:        date,
		Type:        apiTransactionTypeMap[a.Type],
		Operation:   apiTransactionOperationMap[a.Operation],
		Value:       float32(value),
		Title:       strings.TrimSpace(a.Title),
		Description: strings.TrimSpace(a.Description),
	}, nil
}

func parseApiTransactions(d []byte) ([]Transaction, error) {
	var tmp apiTransactions

	err := json.Unmarshal(d, &tmp)
	if err != nil {
		return []Transaction{}, err
	}

	transactions := make([]Transaction, 0, len(tmp.Transactions))

	for _, v := range tmp.Transactions {
		t, err := transactionFromApi(v)
		if err != nil {
			return []Transaction{}, err
		}

		transactions = append(transactions, t)
	}

	return transactions, nil
}
