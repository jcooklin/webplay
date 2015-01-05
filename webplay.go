package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	// "github.com/davecheney/profile"
	"github.com/jcooklin/webplay/bankgoji"
	"github.com/jcooklin/webplay/bankhttprouter"
	"github.com/jcooklin/webplay/common"
	"github.com/olekukonko/tablewriter"
	"github.com/rs/cors"
	"github.com/zenazn/goji"
	"gopkg.in/jmcvetta/napping.v1"
	"os"
	"strconv"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

//todo import this shit from common package or top level
// type Account struct {
// 	Id      string  `json:"id"`
// 	Name    string  `json:"name"`
// 	Balance float64 `json:"balance"`
// }

// type Bank struct {
// 	Id       string             `json:"id"`
// 	Name     string             `json:"name"`
// 	Accounts map[string]Account `json:"accounts"`
// }

func main() {
	host := os.Getenv("HOST")
	if len(host) == 0 {
		host = "localhost"
	}
	port := os.Getenv("PORT")
	if (len(port)) == 0 {
		port = "8000"
	}
	url := fmt.Sprintf("http://%s:%s", host, port)
	app := cli.NewApp()
	app.Name = "webplay"
	app.Usage = "playin with go and web frameworks"
	// app.Action = func(c *cli.Context) {

	// }
	app.Commands = []cli.Command{
		{
			Name: "serve",
			Subcommands: []cli.Command{
				{
					Name: "goji",
					Action: func(c *cli.Context) {
						// cfg := profile.Config{
						// 	//MemProfile:     true,
						// 	CPUProfile:     true,
						// 	ProfilePath:    ".",  // store profiles in current directory
						// 	NoShutdownHook: true, // do not hook SIGINT
						// }
						// defer profile.Start(&cfg).Stop()

						cor := cors.New(cors.Options{
							AllowedOrigins: []string{"http://127.0.0.1:9000"},
						})
						goji.Use(cor.Handler)
						bankgoji.Load()
						goji.Serve()
					},
				},
				{
					Name: "httprouter",
					Action: func(c *cli.Context) {
						// router := httprouter.New()
						// router.GET("/", bankhttprouter.Index)
						// router.GET("/hello/:name", bankhttprouter.Hello)
						bankhttprouter.Load()
						bankhttprouter.Serve("")
						//log.Fatal(http.ListenAndServe(":"+port, router))
					},
				},
			},
		},
		{
			Name:  "bank",
			Usage: "bank operations",
			// Action: func(c *cli.Context) {
			// 	println("bank task: ", c.Args().First(), c.Args().Get(3))
			// },
			Subcommands: []cli.Command{
				{
					Name:        "list",
					Usage:       "list",
					Description: "list all banks",
					Action: func(c *cli.Context) {
						// res := struct {
						// 	Banks map[string]Bank
						// }{}
						res := map[string]common.Bank{}
						e := common.Error{}
						s := napping.Session{}
						resp, err := s.Get(url+"/banks", nil, &res, &e)
						if resp.Status() == 200 {
							table := tablewriter.NewWriter(os.Stdout)
							table.SetHeader([]string{"Id", "Name", "Num of Accounts", "Total Deposits"})
							//https://github.com/golang/go/issues/3117
							var bankTotals = make(map[string]*struct{ TotalDeposits float64 })
							for id, bank := range res {
								bankTotals[id] = &struct {
									TotalDeposits float64
								}{TotalDeposits: 0}
								for _, account := range bank.Accounts {
									bankTotals[id].TotalDeposits += account.Balance
								}
								log.Infof("bank %s has %.2f total deposits in %d accounts", bank.Name, bankTotals[id].TotalDeposits, len(bank.Accounts))
								table.Append([]string{bank.Id, bank.Name, strconv.Itoa(len(bank.Accounts)), strconv.FormatFloat(bankTotals[id].TotalDeposits, 'f', 2, 64)})
							}
							table.Render()
							// log.Println(resp.RawText())
						} else {
							log.Println(fmt.Sprintf("We got an error: %s", err.Error()))
						}
					},
				},
				{
					Name:        "show",
					Usage:       "show <bank id>",
					Description: "show bank details",
					Action: func(c *cli.Context) {
						//we need to require a positional parameter
						if len(c.Args()) < 1 {
							fmt.Println("Bank ID is required")
							fmt.Printf("Usage: %s\n", c.Command.Usage)
							os.Exit(1)
						}
						res := map[string]struct {
							Name          string
							TotalDeposits float64
						}{}
						e := common.Error{}
						s := napping.Session{}
						resp, err := s.Get(fmt.Sprintf(url+"/banks/%s", c.Args().First()), nil, &res, &e)
						fmt.Printf("show bank '%s' details here\n", c.Args().First())
						if resp.Status() == 200 {
							log.Println("we got a 200")
						} else {
							log.Println(fmt.Sprintf("We got an error: %s", err.Error()))
						}
					},
				},
				{
					Name:        "account",
					Usage:       "account <bank id> <account id>",
					Description: "Get a the account details at a bank",
					Action: func(c *cli.Context) {
						if len(c.Args()) < 2 {
							fmt.Println("a bankd ID and account ID are required")
							fmt.Printf("Usage: %s\n", c.Command.Usage)
							os.Exit(1)
						}
						res := common.Account{}
						e := common.Error{}
						s := napping.Session{}
						resp, err := s.Get(fmt.Sprintf(url+"/banks/%s/accounts/%s", c.Args().First(), c.Args().Get(1)), nil, &res, &e)
						fmt.Printf("show bank '%s' account '%s' details here\n", c.Args().First(), c.Args().Get(1))
						if resp.Status() == 200 {
							table := tablewriter.NewWriter(os.Stdout)
							table.SetHeader([]string{"Id", "Name", "Balance"})
							table.Append([]string{res.Id, res.Name, strconv.FormatFloat(res.Balance, 'f', 2, 64)})
							table.Render()
						} else {
							if err != nil {
								log.Println(fmt.Sprintf("We got an error: %s", err.Error()))
							} else {
								log.Printf("We got a %d response", resp.Status())
							}
						}
					},
				},
				{
					Name:        "accounts",
					Usage:       "accounts <bank id>",
					Description: "List a banks accounts",
					Action: func(c *cli.Context) {
						if len(c.Args()) < 1 {
							fmt.Println("Bank ID is required")
							fmt.Printf("Usage: %s\n", c.Command.Usage)
							os.Exit(1)
						}
						res := map[string]common.Account{}
						e := common.Error{}
						s := napping.Session{}
						resp, err := s.Get(fmt.Sprintf(url+"/banks/%s/accounts", c.Args().First()), nil, &res, &e)
						if resp.Status() == 200 {
							table := tablewriter.NewWriter(os.Stdout)
							table.SetHeader([]string{"Id", "Name", "Balance"})
							for _, account := range res {
								table.Append([]string{account.Id, account.Name, strconv.FormatFloat(account.Balance, 'f', 2, 64)})
							}
							table.Render()
						} else {
							if err != nil {
								log.Println(fmt.Sprintf("We got an error: %s", err.Error()))
							} else {
								log.Printf("We got a %d response", resp.Status())
							}
						}
					},
				},
			},
		},
	}

	app.Run(os.Args)
}
