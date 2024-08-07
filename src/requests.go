// Package requests makes requests to Open Collective and telegram
package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gocarina/gocsv"
)

const openCollectiveURL string = "https://rest.opencollective.com/v2/editorjs/transactions.txt?" +
	"kind=ADDED_FUNDS,BALANCE_TRANSFER,CONTRIBUTION,EXPENSE,PLATFORM_TIP&" +
	"includeGiftCardTransactions=1&includeIncognitoTransactions=1&includeChildrenTransactions=1"

const webhookContentType string = "application/x-www-form-urlencoded"

const webhookURL = "https://notify.bot.codex.so/u/"

// GetTransactionsFromOpenCollective gets data from Open Collective api and returns a list of Transaction objects.
func GetTransactionsFromOpenCollective() []*Transaction {
	// Get CSV from Open Collective.
	resp, err := http.Get(openCollectiveURL)
	if err != nil {
		log.Fatalln(err)
	}

	var transactions []*Transaction

	// Parse response to Transaction struct, creates the list of objects.
	err = gocsv.Unmarshal(resp.Body, &transactions)
	if err != nil {
		log.Fatalln(err)
	}

	return transactions
}

// SendToChat gets the list of new transactions and sends webhook to notifier about them.
func SendToChat(transactions []*Transaction, webhookToken string) {
	var message string
	for i := 0; i < len(transactions); i++ {
		amount, err := strconv.ParseFloat(transactions[i].Amount, 64)

		if err != nil {
			log.Fatal(err)
		}

		// Check if amount is less than zero.
		if amount > 0 {
			message = message + fmt.Sprintf("💰 %d$ %s to %s \n\n", int(amount), transactions[i].Description, transactions[i].To)
		}
	}
	if len(message) > 0 {
		data := url.Values{}
		data.Set("message", message)

		_, err := http.Post(fmt.Sprintf("%s%s", webhookURL, webhookToken), webhookContentType,
			strings.NewReader(data.Encode()))

		if err != nil {
			log.Fatal(err)
		}
	}
}
