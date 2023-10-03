package inter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type OAuth struct {
	client *Client
}

func NewOAuth(client *Client) *OAuth {
	return &OAuth{
		client: client,
	}
}

func (o *OAuth) Authorize(ctx context.Context, clientID, clientSecret string, scopes ...string) (Token, error) {
	endpoint := fmt.Sprintf("%s/oauth/v2/token", o.client.apiBaseUrl)

	form := url.Values{}
	form.Add("scope", strings.Join(scopes, " "))
	form.Add("grant_type", "client_credentials")
	form.Add("client_id", clientID)
	form.Add("client_secret", clientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint,
		bytes.NewBufferString(form.Encode()))
	if err != nil {
		return Token{}, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := o.client.Do(req)
	if err != nil {
		return Token{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Token{}, err
	}

	if resp.StatusCode != http.StatusOK {
		e, err := parseApiOAuthResponseError(data)
		if err != nil {
			return Token{}, err
		}

		return Token{}, errors.New(e.ErrorTitle)
	}

	return parseApiResponseToken(data)
}

type apiOAuthResponseError struct {
	Error      string `json:"error"`
	ErrorTitle string `json:"error_title"`
}

func parseApiOAuthResponseError(d []byte) (apiOAuthResponseError, error) {
	var e apiOAuthResponseError

	err := json.Unmarshal(d, &e)
	if err != nil {
		return apiOAuthResponseError{}, err
	}

	return e, nil
}

type apiToken struct {
	Data      string `json:"access_token"`
	Type      string `json:"token_type"`
	ExpiresIn int    `json:"expires_in"`
	Scopes    string `json:"scope"`
}

func parseApiResponseToken(d []byte) (Token, error) {
	var tmp apiToken

	err := json.Unmarshal(d, &tmp)
	if err != nil {
		return Token{}, err
	}

	return Token{
		Data:      tmp.Data,
		Type:      tmp.Type,
		Scopes:    strings.Split(tmp.Scopes, " "),
		ExpiresAt: time.Now().Add(time.Duration(tmp.ExpiresIn) * time.Second),
	}, nil
}
