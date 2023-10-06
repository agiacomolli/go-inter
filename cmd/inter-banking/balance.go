package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/agiacomolli/go-inter"
)

var (
	date        string
	dateUsage   = "balance date"
	defaultDate = ""

	outputFormat        string
	outputFormatUsage   = "output format"
	defaultOutputFormat = "short"
)

func balanceCommand(ctx context.Context, banking *inter.Banking, args []string) {
	flag := flag.NewFlagSet("balance", flag.ExitOnError)

	flag.StringVar(&date, "d", defaultDate, dateUsage)
	flag.StringVar(&date, "date", defaultDate, dateUsage)
	flag.StringVar(&outputFormat, "f", defaultOutputFormat, outputFormatUsage)
	flag.StringVar(&outputFormat, "output-format", defaultOutputFormat, outputFormatUsage)

	flag.Usage = mainUsage
	flag.Parse(args)

	var balanceDate time.Time
	if date == "" {
		balanceDate = time.Now()
	} else {
		var err error
		balanceDate, err = time.Parse(time.DateOnly, date)
		if err != nil {
			fmt.Printf("could not parse balance date: %s\n", err)
			os.Exit(1)
		}
	}

	balance, err := banking.Balance(ctx, balanceDate)
	if err != nil {
		fmt.Printf("could not get balance: %s\n", err)
		os.Exit(1)
	}

	switch outputFormat {
	case "short":
		fmt.Printf("%.2f\n", balance.Available)
	case "full":
		var payload strings.Builder
		fmt.Fprintf(&payload, "Balances at %s\n\n", date)
		tw := tabwriter.NewWriter(&payload, 5, 1, 2, ' ', 0)
		fmt.Fprintf(tw, "Available\t%10.2f\n", balance.Available)
		fmt.Fprintf(tw, "Limit\t%10.2f\n", balance.Limit)
		fmt.Fprintf(tw, "On hold\t%10.2f\n", balance.CheckOnHold)
		fmt.Fprintf(tw, "Judicially blocked\t%10.2f\n", balance.JudiciallyBlocked)
		fmt.Fprintf(tw, "Administratively blocked\t%10.2f", balance.AdministrativelyBlocked)
		tw.Flush()
		fmt.Println(payload.String())
	default:
		fmt.Println("invalid output format")
		os.Exit(1)
	}

}
