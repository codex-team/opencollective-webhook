package main

import (
	"github.com/joho/godotenv" //nolint:gofumpt
	"log"
	"os"
	"strconv"
	"time" //nolint:gofumpt
)

func main() {
	// load values from .env into the system.
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	var startDateString string

	// Get environment variables.
	startDateString, startDateExists := os.LookupEnv("START_DATE")
	if !startDateExists {
		log.Fatal("No START_DATE in environment variables")
	}

	webhookToken, webhookTokenExists := os.LookupEnv("TOKEN")

	if !webhookTokenExists {
		log.Fatal("No Codex Bot TOKEN in environment variables")
	}

	periodicityString, exists := os.LookupEnv("PERIODICITY")

	if !exists {
		log.Fatal("No PERIODICITY in environment variables")
	}

	// Convert string to integer.
	periodicity, err := strconv.Atoi(periodicityString)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s Environment variables initiated", time.Now().Format(dateFormat))

	for {
		// Get all transactions from Open Collective.
		var transactions = GetTransactionsFromOpenCollective()

		log.Printf("%s Got %d transactions from open collective", time.Now().Format(dateFormat), len(transactions))

		// Get transactions, which was after certain date.
		sortedTransactions := sortTransactionsByStartDate(startDateString, transactions)

		// Check, if there is no new transactions.
		if len(sortedTransactions) > 0 {
			// Send messages to telegram about new transactions.
			SendToChat(sortedTransactions, webhookToken)

			log.Printf("%s Webhook(s) with %d new transaction(s) was(were) sent",
				time.Now().Format(dateFormat), len(sortedTransactions))
		} else {
			log.Printf("%s There was no new transactions", time.Now().Format(dateFormat))
		}
		// Change start date to date of the last transaction.
		startDateString = transactions[0].DateTime

		// Make delay.
		time.Sleep(time.Duration(periodicity) * time.Minute)
	}
}
