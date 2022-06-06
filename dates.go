package main

import (
	"log"
	"time"
)

const dateFormat string = "2006-01-02T15:04:05"

// parseDateFromString get string and parse it to Date.
func parseDateFromString(str string) time.Time {
	date, err := time.Parse(dateFormat, str)
	if err != nil {
		log.Fatal(err)
	}

	return date
}

// compareDate get two dates and return true, if the first date is after the second.
func compareDate(firstDate time.Time, secondDate time.Time) bool {
	return firstDate.After(secondDate)
}

// sortTransactionsByStartDate gets start date and the list of transactions, and returns transactions,
// which are after the date.
func sortTransactionsByStartDate(startDateString string, transactions []*Transaction) []*Transaction {
	var startDate = parseDateFromString(startDateString)

	var sortedTransactions []*Transaction

	// Loop with all getting transactions, from end to start.
	for i := len(transactions) - 1; i >= 0; i-- {
		// Get Date of transaction.
		trTime := parseDateFromString(transactions[i].DateTime)
		// Check if transaction is after the start date.
		if compareDate(trTime, startDate) { // Add transaction to the list of new transactions.
			sortedTransactions = append(sortedTransactions, transactions[i])
		}
	}

	return sortedTransactions
}
