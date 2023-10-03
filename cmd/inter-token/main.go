package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/agiacomolli/go-inter"
)

var (
	certFile     = flag.String("cert", "cert.crt", "signed certificate file")
	keyFile      = flag.String("key", "cert.key", "certificate private key file")
	clientID     = flag.String("client-id", "", "client identification")
	clientSecret = flag.String("client-secret", "", "client secret")
	scopes       = flag.String("scopes", "", "comma-separated client scopes")
	outputFormat = flag.String("output-format", "token", "output format [token|info|json]")
	help         = flag.Bool("help", false, "display this help message")
)

func main() {
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *clientID == "" {
		fmt.Println("client-id is required")
		os.Exit(1)
	}

	if *clientSecret == "" {
		fmt.Println("client-secret is required")
		os.Exit(1)
	}

	cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		fmt.Printf("could not parse certificate files: %s\n", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	client := inter.NewClient(cert)

	oauth := inter.NewOAuth(client)

	token, err := oauth.Authorize(ctx, *clientID, *clientSecret,
		strings.Split(*scopes, ",")...)
	if err != nil {
		fmt.Printf("could not authorize: %s\n", err)
		os.Exit(1)
	}

	var output strings.Builder

	switch *outputFormat {
	case "token":
		_, err = fmt.Fprint(&output, token.Data)
	case "info":
		err = writeInfoOutput(&output, token)
	case "json":
		err = writeJsonOutput(&output, token)
	default:
		fmt.Println("invalid output format")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("could not print output: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(output.String())
}

func writeInfoOutput(w io.Writer, t inter.Token) error {
	_, err := fmt.Fprintf(w, `token     %s
type      %s
expires   %s
scopes    %s`,
		t.Data, t.Type, t.ExpiresAt.Format(time.RFC3339),
		strings.Join(t.Scopes, " "))

	return err
}

type tokenOutput struct {
	Data      string   `json:"data"`
	Type      string   `json:"type"`
	ExpiresAt string   `json:"expires_at"`
	Scopes    []string `json:"scopes"`
}

func writeJsonOutput(w io.Writer, t inter.Token) error {
	tmp := tokenOutput{
		Data:      t.Data,
		Type:      t.Type,
		ExpiresAt: t.ExpiresAt.Format(time.RFC3339),
		Scopes:    t.Scopes,
	}

	b, err := json.MarshalIndent(tmp, "", "\t")
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(w, string(b))

	return err
}
