package bankhttprouter

import (
	// log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/jcooklin/webplay/common"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"net/http"
	"os"
	"path/filepath"
)

type route struct {
	method  string
	path    string
	handler func(w http.ResponseWriter, r *http.Request, p httprouter.Params)
}

var routes = []route{
	{"GET", "/", Index},
	{"GET", "/banks", getBanks},
	{"GET", "/banks/:bankId", getBank},
	{"GET", "/banks/:bankId/accounts", getBankAccounts},
	{"GET", "/banks/:bankId/accounts/:accountId", getBankAccount},
}

var router *httprouter.Router

func init() {
	router = httprouter.New()
}

// Load loads all of the available routes.
func Load() {
	//router := httprouter.New()
	for _, route := range routes {
		switch route.method {
		case "GET":
			router.GET(route.path, route.handler)
		default:
			panic("UNKNOWN HTTP METHOD: " + route.method)
		}
	}
}

// Serve will start the net/http server using port 8000.
// The env variable PORT can be used to override the port.
func Serve(port string) {
	cor := cors.New(cors.Options{
		AllowedOrigins: []string{"http://127.0.0.1:9000"},
	})
	if port == "" {
		port = os.Getenv("PORT")
	}
	if (len(port)) == 0 {
		port = "8000"
	}
	// log.Infof("port: %s", port)
	//n := negroni.Classic()
	n := negroni.New()
	n.Use(cor)
	publicPath, _ := filepath.Abs("public")
	router.ServeFiles("/api/*filepath", http.Dir(publicPath))
	// n.UseHandler(http.FileServer(http.Dir("swagger.json")))
	n.UseHandler(router)
	n.Run(":" + port)
	//log.Fatal(http.ListenAndServe(":"+port, router))
}

// Index is the default route.
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	common.Index(w, r, "httprouter")
}

func getBanks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	common.GetBanks(w, r)
}

func getBank(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	common.GetBank(w, r, p.ByName("bankId"))
}

func getBankAccounts(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	common.GetBankAccounts(w, r, p.ByName("bankId"))
}

func getBankAccount(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	common.GetBankAccount(w, r, p.ByName("bankId"), p.ByName("accountId"))
}
