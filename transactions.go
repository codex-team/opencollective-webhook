package main

type Collective struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Slug string `json:"slug"`
}

type Subscription struct {
	Interval string `json:"interval"`
}

type Transaction struct {
	Id int `json:"id"`
	Type string `json:"type"`
	Amount int `json:"amount"`
	Currency string `json:"currency"`
	NetAmountInCollectiveCurrency int `json:"netAmountInCollectiveCurrency"`
	Collective Collective `json:"collective"`
	FromCollective Collective `json:"fromCollective"`
	Subscription Subscription `json:"subscription"`
}

type Transactions []Transaction

func (a Transactions) Len() int           { return len(a) }
func (a Transactions) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Transactions) Less(i, j int) bool { return a[i].Id < a[j].Id }

type GraphQlDataTransactions struct {
	AllTransactions Transactions `json:"allTransactions"`
}

type GraphQlData struct {
	Data GraphQlDataTransactions `json:"data"`
}