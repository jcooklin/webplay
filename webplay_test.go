// +build go1.4

package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/jcooklin/webplay/bankgoji"
	"github.com/jcooklin/webplay/bankhttprouter"
	"github.com/jcooklin/webplay/common"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web/middleware"
	"gopkg.in/jmcvetta/napping.v1"
	// "io/ioutil"
	"github.com/davecheney/profile"
	"os"
	"strconv"
	"testing"
)

var url string

func TestMain(m *testing.M) {
	//log.SetOutput(ioutil.Discard)
	//log.SetOutput(os.Stdout)
	goji.DefaultMux.Abandon(middleware.Logger)
	host := os.Getenv("HOST")
	if len(host) == 0 {
		host = "localhost"
	}
	port := os.Getenv("PORT")
	if (len(port)) == 0 {
		port = "8000"
	}
	url = fmt.Sprintf("http://%s:%s", host, port)
	bankgoji.Load()
	log.Info("Testing goji")
	cfg := profile.Config{
		//MemProfile:     true,
		CPUProfile:     true,
		ProfilePath:    ".",  // store profiles in current directory
		NoShutdownHook: true, // do not hook SIGINT
	}
	p := profile.Start(&cfg)
	go goji.Serve()
	res := m.Run()
	if res > 0 {
		os.Exit(res)
	}
	bankhttprouter.Load()
	log.Info("Testing httprouter")
	go bankhttprouter.Serve("8001")
	url = fmt.Sprintf("http://%s:%s", host, "8001")
	res = m.Run()
	p.Stop()
	os.Exit(res)
}

func TestBankList(t *testing.T) {
	res := map[string]common.Bank{}
	e := common.Error{}
	s := napping.Session{}
	resp, err := s.Get(url+"/banks", nil, &res, &e)
	if resp.Status() != 200 {
		if err == nil {
			log.Error(err.Error())
		}
		t.Error(fmt.Sprintf("we expected 200 but got %s", strconv.Itoa(resp.Status())))
	}
}

func TestBankListAccount(t *testing.T) {
	res := map[string]common.Bank{}
	e := common.Error{}
	s := napping.Session{}
	resp, err := s.Get(url+"/banks", nil, &res, &e)
	if resp.Status() != 200 {
		if err != nil {
			log.Error(err.Error())
		}
		t.Error(fmt.Sprintf("we expected 200 but got %s", strconv.Itoa(resp.Status())))
	}

	bankIds := make([]string, 0, len(res))
	for key := range res {
		bankIds = append(bankIds, key)
	}

	accountsResult := map[string]common.Account{}
	accountsResp, accountsErr := s.Get(fmt.Sprintf("%s/banks/%s/accounts", url, bankIds[0]), nil, &accountsResult, &e)
	if accountsResp.Status() != 200 {
		if err != nil {
			log.Error(accountsErr.Error())
		}
		t.Error(fmt.Sprintf("we expected 200 but got %s for %s:",
			strconv.Itoa(accountsResp.Status()),
			fmt.Sprintf("%s/banks/%s/accounts", url, bankIds[0])))
	}
}

func BenchmarkBankList(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res := map[string]common.Bank{}
		e := common.Error{}
		s := napping.Session{}
		resp, err := s.Get(url+"/banks", nil, &res, &e)
		if resp.Status() != 200 {
			if err != nil {
				log.Error(err.Error())
			}
			b.Error(fmt.Sprintf("we expected a 200 response but got %s", strconv.Itoa(resp.Status())))
		}
		for bankId, bank := range res {
			for accountId := range bank.Accounts {
				res := common.Account{}
				resp, err := s.Get(fmt.Sprintf("%s/banks/%s/accounts/%s", url, bankId, accountId), nil, &res, &e)
				if resp.Status() != 200 {
					if err != nil {
						log.Error(err.Error())
					}
					b.Error(fmt.Sprintf("we expected a 200 response but got %s", strconv.Itoa(resp.Status())))
				}
			}
		}
	}
}
