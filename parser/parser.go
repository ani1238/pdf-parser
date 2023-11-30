package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ledongthuc/pdf"
)

type Transactions []Transaction

var TotalTransactions []Transaction

func matchFormat(s string) bool {
	return regexp.MustCompile(`^[0-9\-\/]+$`).MatchString(s)
}

type Transaction struct {
	TxnDate    time.Time
	TxnDesc    string
	TxnAmount  float64
	TxnBalance float64
}

func parseTransaction(s string) (*Transaction, error) {
	re := regexp.MustCompile(`^(\d{1,2}\/\d{1,2}\/\d{2,4})\s(\S.*?)\s(\d{1,3}(?:,\d{3})*\.\d{2})\s(\d{1,3}(?:,\d{3})*\.\d{2})`)
	isFormat := re.MatchString(s)
	if isFormat {
		// fmt.Println(s)
		re1 := regexp.MustCompile(`(\d{1,2}\/\d{1,2}\/\d{2,4})`)
		date := re1.FindStringSubmatch(s)
		actualDate, err := time.Parse("2/1/2006", date[0])
		if err != nil {
			return nil, err
		}
		re2 := regexp.MustCompile(`\s\S+`)
		description := ""
		desc := re2.FindAllStringSubmatch(s, -1)
		for i, v := range desc {
			if i == len(desc)-1 || i == len(desc)-2 {
				break
			}
			description += v[0] + " "
		}
		re3 := regexp.MustCompile(`(\d{1,3}(?:,\d{3})*\.\d{2})`)
		money := re3.FindAllStringSubmatch(s, -1)
		amt, _ := strconv.ParseFloat(strings.Replace(money[0][0], ",", "", -1), 8)
		bal, _ := strconv.ParseFloat(strings.Replace(money[1][0], ",", "", -1), 8)
		// fmt.Println(date[1], desc, description, amt, bal)
		transactionConcrete := Transaction{TxnDate: actualDate, TxnDesc: description, TxnAmount: amt, TxnBalance: bal}
		// fmt.Println(&transactionConcrete)
		return &transactionConcrete, nil
	}
	return nil, fmt.Errorf("The given string is not a transaction")
}

func ParsePdf(path string) error {
	f, r, err := pdf.Open(path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return err
	}
	totalPage := r.NumPage()
	// fmt.Println("num of pages: ", totalPage)
	// transactionRows := make([]int, 0)

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		rows, _ := p.GetTextByRow()
		for _, row := range rows {
			// fmt.Printf(">>>> row: %d ", row.Position)
			transactionString := ""
			for _, word := range row.Content {
				// fmt.Printf(word.S)
				transactionString += word.S

			}
			transactionVar, _ := parseTransaction(transactionString)
			if transactionVar != nil && err == nil {
				TotalTransactions = append(TotalTransactions, *transactionVar)
				// fmt.Println(*transactionVar)
			}

		}
	}
	// fmt.Println(transactionRows)
	return nil
}

func GetAllTransactions() (Transactions, error) {
	if len(TotalTransactions) > 0 {
		return TotalTransactions, nil
	}
	return nil, fmt.Errorf("There are no transactions")
}

func GetTransactionsByDate(startDate time.Time, endDate time.Time) (Transactions, error) {
	var resultTransactions Transactions

	for _, val := range TotalTransactions {
		fmt.Println(startDate, endDate, val.TxnDate)
		if (val.TxnDate.After(startDate) && val.TxnDate.Before(endDate)) || val.TxnDate == startDate || val.TxnDate == endDate {
			fmt.Println(val.TxnDate)
			resultTransactions = append(resultTransactions, val)
		}
	}
	if len(resultTransactions) == 0 {
		return nil, fmt.Errorf("There are no transactions within given dates")
	}
	return resultTransactions, nil
}

func GetBalanceByDate(requiredDate time.Time) (*float64, error) {
	var resultBalance float64 = -1.0 //assuming balance cannot go negative
	for _, val := range TotalTransactions {
		resultBalance = val.TxnBalance
		if val.TxnDate.After(requiredDate) {
			break
		}
	}
	if resultBalance == -1.0 {
		return nil, fmt.Errorf("Date out of range of data")
	}
	return &resultBalance, nil
}
