package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"time"
)

type TransactionsVariables struct {
	CollectiveId int `json:"CollectiveId"`
	DateFrom string `json:"dateFrom"`
	DateTo string `json:"dateTo"`
}

type TransactionsQuery struct {
	OperationName string `json:"operationName"`
	Variables TransactionsVariables `json:"variables"`
	Query string `json:"query"`
}

//const maximumQuery  = "query Transactions($CollectiveId: Int!, $type: String, $limit: Int, $offset: Int, $dateFrom: String, $dateTo: String) {\n  allTransactions(CollectiveId: $CollectiveId, type: $type, limit: $limit, offset: $offset, dateFrom: $dateFrom, dateTo: $dateTo) {\n    id\n    uuid\n    description\n    createdAt\n    type\n    amount\n    currency\n    hostCurrency\n    hostCurrencyFxRate\n    netAmountInCollectiveCurrency\n    hostFeeInHostCurrency\n    platformFeeInHostCurrency\n    paymentProcessorFeeInHostCurrency\n    paymentMethod {\n      service\n      type\n      name\n      data\n      __typename\n    }\n    collective {\n      id\n      slug\n      type\n      name\n      __typename\n    }\n    fromCollective {\n      id\n      name\n      slug\n      path\n      image\n      __typename\n    }\n    usingVirtualCardFromCollective {\n      id\n      slug\n      name\n      __typename\n    }\n    host {\n      id\n      slug\n      name\n      currency\n      hostFeePercent\n      __typename\n    }\n    ... on Expense {\n      category\n      attachment\n      __typename\n    }\n    ... on Order {\n      createdAt\n      subscription {\n        interval\n        __typename\n      }\n      __typename\n    }\n    refundTransaction {\n      id\n      uuid\n      description\n      createdAt\n      type\n      amount\n      currency\n      hostCurrency\n      hostCurrencyFxRate\n      netAmountInCollectiveCurrency\n      hostFeeInHostCurrency\n      platformFeeInHostCurrency\n      paymentProcessorFeeInHostCurrency\n      paymentMethod {\n        service\n        type\n        name\n        data\n        __typename\n      }\n      collective {\n        id\n        slug\n        type\n        name\n        __typename\n      }\n      fromCollective {\n        id\n        name\n        slug\n        path\n        image\n        __typename\n      }\n      usingVirtualCardFromCollective {\n        id\n        slug\n        name\n        __typename\n      }\n      host {\n        id\n        slug\n        name\n        currency\n        hostFeePercent\n        __typename\n      }\n      ... on Expense {\n        category\n        attachment\n        __typename\n      }\n      ... on Order {\n        createdAt\n        subscription {\n          interval\n          __typename\n        }\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n}\n"
const query = "query Transactions($CollectiveId: Int!, $type: String, $limit: Int, $offset: Int, $dateFrom: String, $dateTo: String) {\n  allTransactions(CollectiveId: $CollectiveId, type: $type, limit: $limit, offset: $offset, dateFrom: $dateFrom, dateTo: $dateTo) {\n    id\n    type\n    amount\n    currency\n    netAmountInCollectiveCurrency\n    collective {\n      name\n    }\n    fromCollective {\n      name\n      path\n    }\n  }\n}\n"
const graphQlURL = "https://opencollective.com/api/graphql"
const currentStateFilename = ".opencollective-current.json"
const webhookURL = "https://notify.bot.codex.so/u/"

func loadCurrentState() (GraphQlData, error) {
	if _, err := os.Stat(currentStateFilename); os.IsNotExist(err) {
		return GraphQlData{}, err
	}

	jsonFile, err := os.Open(currentStateFilename)
	if err != nil {
		log.Fatalf("File open error: %v", err)
	}
	defer jsonFile.Close()

	data, _ := ioutil.ReadAll(jsonFile)
	var state = &GraphQlData{}
	if err := json.Unmarshal(data, &state); err != nil {
		log.Fatalf("JSON file unmarshalling error: %v", err)
	}

	return *state, nil
}

func saveCurrentState(body []byte) {
	if err := ioutil.WriteFile(currentStateFilename, body, 0644); err != nil {
		log.Fatalf("JSON saving error: %v", err)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Provide CodeX.Bot webhook token as command argument")
	}

	token := os.Args[1]

	gqlQuery := TransactionsQuery{
		OperationName: "Transactions",
		Variables: TransactionsVariables{
			CollectiveId: 37258,
			DateFrom: "2019-02-28T21:00:00.000Z",
			DateTo: time.Now().Format("2006-01-02T15:04:05.000Z"),
		},
		Query: query,
	}

	request, err := json.Marshal(&gqlQuery)
	if err != nil {
		log.Fatalf("JSON Marshalling error: %v", err)
	}

	req, err := http.NewRequest("POST", graphQlURL, bytes.NewBuffer(request))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("HTTP POST error: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Printf("Invalid status code: %d", resp.StatusCode)
		fmt.Println("response Headers:", resp.Header)
		fmt.Println("response Body:", string(body))
	}

	var response = &GraphQlData{}
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatalf("JSON Unmarshalling error: %v", err)
	}

	currentState, err := loadCurrentState()

	// Not a first run
	if err == nil {
		ids := make(map[int]bool)
		for _, transaction := range currentState.Data.AllTransactions {
			ids[transaction.Id] = true
		}

		newTransactions := []Transaction{}
		if len(currentState.Data.AllTransactions) < len(response.Data.AllTransactions) {
			for _, transaction := range response.Data.AllTransactions {
				if _, ok := ids[transaction.Id]; !ok {
					newTransactions = append(newTransactions, transaction)
				}
			}
		}

		sort.Sort(Transactions(newTransactions))
		for _, transaction := range newTransactions {
			data := url.Values{}
			data.Set("message", fmt.Sprintf("💰 %d$ donation to %s from %s", transaction.Amount / 100, transaction.Collective.Name, transaction.FromCollective.Name))
			_, err := MakeHTTPRequest("POST", fmt.Sprintf("%s%s", webhookURL, token), []byte(data.Encode()), map[string]string{
				"Content-Type": "application/x-www-form-urlencoded",
			})
			if err != nil {
				log.Fatalf("Webhook error: %v", err)
			}
		}
	}

	saveCurrentState(body)
}