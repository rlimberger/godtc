package main

import (
	"flag"
	"fmt"
	"github.com/rlimberger/godtc"
	"log"
	"strings"
)

func main() {

	// required arguments
	tradeAccount := flag.String("sc_account", "", "Sierrachart trade account")
	numberOfDays := flag.Int("sc_days", 0, "Number of days to request fills for")

	// make sure required arguments are specified
	flag.Parse()
	if *tradeAccount == "" || *numberOfDays == 0 {
		panic("Missing argument(s). Please specify username, password, account and days.")
	}

	// create SierraChart DTC client and logon to the local DTC server
	c, err := godtc.NewClient()
	if err != nil {
		panic(err)
	}

	// request historical fills from SierraChart
	log.Printf("Requesting %d day(s) of historical fills for trade account `%s` from SierraChart\n", *numberOfDays, *tradeAccount)
	fills, err := c.RequestHistoricalFills(*tradeAccount, *numberOfDays)
	if err != nil {
		panic(err)
	}
	log.Printf("Received %d fills from Sierrachart\n", len(fills))

	//convert SierraChart fills to Tradervue executions
	var executions []tradervue.Execution
	for _, fill := range fills {
		tve, err := dtc2tv.ExecutionFromDTCOrderFill(fill)
		if err != nil {
			panic(fmt.Sprintf("Error during conversion %s", err.Error()))
		}
		executions = append(executions, tve)
	}

	//import into TV
	log.Printf("Importing %d executions into Tradervue...\n", len(executions))
	tags := strings.Split(*tagsRaw, ",")
	err = tradervue.Import(executions, *username, *password, tags, accountTag)
	if err != nil {
		panic(fmt.Sprintf("Tradervue import failed with %s", err.Error()))
	}
}
