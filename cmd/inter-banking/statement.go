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
	start        string
	startUsage   = "start date"
	defaultStart = ""

	end        string
	endUsage   = "end date"
	defaultEnd = ""
)

func statementCommand(ctx context.Context, banking *inter.Banking, args []string) {
	flag := flag.NewFlagSet("statement", flag.ExitOnError)

	flag.StringVar(&start, "s", defaultStart, startUsage)
	flag.StringVar(&start, "start-date", defaultStart, startUsage)
	flag.StringVar(&end, "e", defaultEnd, endUsage)
	flag.StringVar(&end, "end-date", defaultEnd, endUsage)

	flag.Usage = mainUsage
	flag.Parse(args)

	if start == "" {
		fmt.Println("start date is required")
		os.Exit(1)
	}

	startDate, err := time.Parse(time.DateOnly, start)
	if err != nil {
		fmt.Printf("could not parse start date: %s\n", err)
		os.Exit(1)
	}

	var endDate time.Time
	if end == "" {
		endDate = time.Now()
	} else {
		endDate, err = time.Parse(time.DateOnly, end)
		if err != nil {
			fmt.Printf("could not parse start date: %s\n", err)
			os.Exit(1)
		}
	}

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
