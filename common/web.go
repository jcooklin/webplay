package common

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request, s string) {
	fmt.Fprintf(w, "Welcome to Bank Webplay (%s)", s)
}

// GetBanks returns a map of Banks.
// The bank Id is the key of the map.
func GetBanks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	banks := make(map[string]interface{})
	for index, bank := range Banks {
		// banks[index] = bank.GetSummary()
		banks[index] = bank
	}
	slcA, _ := json.Marshal(banks)
	w.Write(slcA)
}

// GetBank returns a bank provided it's Id.
func GetBank(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")
	// id := c.URLParams["bankId"]
	bank, ok := Banks[id]
	if !ok {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	bank.WriteSummary(w)
}

// GetBankAccounts returns a map dict of accounts.
// The account number is the key.
func GetBankAccounts(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")
	bank, ok := Banks[id]
	if !ok {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	slcA, _ := json.Marshal(bank.Accounts)
	w.Write(slcA)
}

func GetBankAccount(w http.ResponseWriter, r *http.Request, bankId string, accountId string) {
	w.Header().Set("Content-Type", "application/json")
	account, ok := Banks[bankId].Accounts[accountId]
	if !ok {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	slcA, _ := json.Marshal(account)
	w.Write(slcA)
}
