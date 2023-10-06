package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/agiacomolli/go-inter"
)

var (
	certFile        string
	certFileUsage   = "signed certificate file"
	defaultCertFile = "cert.crt"

	keyFile        string
	keyFileUsage   = "certificate private key file"
	defaultKeyFile = "cert.key"

	tokenData        string
	tokenDataUsage   = "user token"
	defaultTokenData = ""
)

func main() {
	flag.StringVar(&certFile, "c", defaultCertFile, certFileUsage)
	flag.StringVar(&certFile, "cert", defaultCertFile, certFileUsage)
	flag.StringVar(&keyFile, "k", defaultKeyFile, keyFileUsage)
	flag.StringVar(&keyFile, "key", defaultKeyFile, keyFileUsage)
	flag.StringVar(&tokenData, "t", defaultTokenData, tokenDataUsage)
	flag.StringVar(&tokenData, "token", defaultTokenData, tokenDataUsage)

	flag.Usage = mainUsage
	flag.Parse()

	if tokenData == "" {
		fmt.Println("token is required")
		os.Exit(1)
	}
	token := inter.TokenFromString(tokenData)

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		fmt.Printf("could not parse certificate files: %s\n", err)
		os.Exit(1)
	}
	client := inter.NewClient(cert)

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("no subcommand set")
		os.Exit(1)
	}

	cmd, args := args[0], args[1:]

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	banking := inter.NewBanking(client, token)

	switch cmd {
	case "balance":
		balanceCommand(ctx, banking, args)
	case "statement":
		statementCommand(ctx, banking, args)
	default:
		fmt.Println("command not found:", cmd)
		os.Exit(1)
	}
}

func mainUsage() {
	fmt.Fprintf(flag.CommandLine.Output(),
		`Usage: inter-banking [OPTION...] <COMMAND>

  -h, --help                 give this help list
  -c, --cert                 signed certificate file (default 'cert.crt')
  -k, --key                  certificate private key file (default 'cert.key')
  -t, --token                personal user token


balance                      get account balance

  -d, --date                 balance date in the format YYYY-MM-DD (defaults to
                             today)
  -f, --output-format        the output format used to show balance; can be
                             'short' (default) or 'full'

statement                    fetch account statements

  -s, --start-date           statements start date in the format YYYY-MM-DD
  -e, --end-date             statements end date in the format YYYY-MM-DD (defaults to
                             today)
`)
}
