package common

import (
	"encoding/json"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"io"
	"math/rand"
	"time"
)

type Account struct {
	Id      string  `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

type Bank struct {
	Id       string             `json:"id"`
	Name     string             `json:"name"`
	Accounts map[string]Account `json:"accounts"`
}

type Error struct {
	Message string
	Errors  []struct {
		Resource string
		Field    string
		Code     string
	}
}

var bankId = pseudo_uuid()
var Banks = map[string]Bank{
	bankId: Bank{
		Id:       bankId,
		Name:     "ACME",
		Accounts: make(map[string]Account),
	},
}

func (a Bank) Write(w io.Writer) {
	slcA, _ := json.Marshal(a)
	fmt.Fprint(w, string(slcA))
}

func (a Bank) WriteSummary(w io.Writer) {
	summary, _ := json.Marshal(a.GetSummary())
	w.Write(summary)
}

func (a Bank) GetSummary() interface{} {
	var totalDeposits float64
	for _, account := range a.Accounts {
		totalDeposits += account.Balance
	}
	return struct {
		Name          string
		TotalDeposits float64
	}{
		a.Name,
		totalDeposits,
	}
}

func (a Account) Write(w io.Writer) {
	slcA, _ := json.Marshal(a)
	fmt.Fprint(w, string(slcA))
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	for _, bank := range Banks {
		for i := 0; i < 100; i++ {
			id := pseudo_uuid()
			bank.Accounts[id] = Account{
				Id:      id,
				Name:    randomdata.FullName(randomdata.RandomGender),
				Balance: rand.Float64() * 99999,
			}
		}
	}
}
