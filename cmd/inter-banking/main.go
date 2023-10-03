package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/agiacomolli/go-inter"
)

var (
	certFile  = flag.String("cert", "cert.crt", "signed certificate file")
	keyFile   = flag.String("key", "cert.key", "certificate private key file")
	tokenData = flag.String("token", "", "user token")
	start     = flag.String("start", "", "start date")
	end       = flag.String("end", "", "end date (not required)")
	help      = flag.Bool("help", false, "display this help message")
)

func main() {
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *tokenData == "" {
		fmt.Println("token is required")
		os.Exit(1)
	}

	if *start == "" {
		fmt.Println("start date is required")
		os.Exit(1)
	}

	startDate, err := time.Parse(time.DateOnly, *start)
	if err != nil {
		fmt.Printf("could not parse start date: %s\n", err)
		os.Exit(1)
	}

	var endDate time.Time
	if *end == "" {
		endDate = time.Now()
	} else {
		endDate, err = time.Parse(time.DateOnly, *end)
		if err != nil {
			fmt.Printf("could not parse start date: %s\n", err)
			os.Exit(1)
		}
	}

	cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		fmt.Printf("could not parse certificate files: %s\n", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	client := inter.NewClient(cert)

	token := inter.TokenFromString(*tokenData)

	banking := inter.NewBanking(client, token)

	transactions, err := banking.Transactions(ctx, startDate, endDate)
	if err != nil {
		fmt.Printf("could not get transactions: %s\n", err)
		os.Exit(1)
	}

	var payload strings.Builder
	fmt.Fprintf(&payload, "Statements from %s to %s\n\n",
		startDate.Format(time.DateOnly), endDate.Format(time.DateOnly))

	tw := tabwriter.NewWriter(&payload, 5, 1, 2, ' ', 0)
	fmt.Fprintln(tw, "Date\t     Value \tOperation\tType\tTitle\tDescription")

	for _, v := range transactions {
		fmt.Fprintf(tw, "%s\t%10.2f\t%s\t%s\t%s\t%s\t\n",
			v.Date.Format(time.DateOnly), v.Value,
			v.Operation, v.Type, v.Title, v.Description)
	}
	tw.Flush()

	fmt.Println(payload.String())
}
