package inter

import (
	"time"
)

type Token struct {
	Data      string
	Type      string
	Scopes    []string
	ExpiresAt time.Time
}

func (t Token) Valid() bool {
	return time.Now().Before(t.ExpiresAt)
}

func TokenFromString(s string) Token {
	return Token{Data: s}
}
