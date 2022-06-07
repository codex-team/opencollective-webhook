package main

type Transaction struct {
	DateTime    string `csv:"datetime"`    // Transaction Date time
	Description string `csv:"description"` // Transaction information
	Amount      string `csv:"amount"`      // Transaction amount
	To          string `csv:"accountName"` // Transaction destination
}
