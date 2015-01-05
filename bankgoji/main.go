// @APIVersion 0.0.1
// @Title My Toy Bank API
// @Description This is a playground for experimentation
// @Contact joel.cooklin@gmail.com
// @TermsOfServiceUrl http://foo.bar
// @License BSD
// @LicenseUrl http://opensource.org/licenses/BSD-2-Clause

// @SubApi Banks API [/banks]

package bankgoji

import (
	//"fmt"
	log "github.com/Sirupsen/logrus"
	// "github.com/davecheney/profile"
	"github.com/jcooklin/webplay/common"
	"github.com/rs/cors"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	//"io"
	"net/http"
	"path/filepath"
)

type route struct {
	method  string
	path    string
	handler func(c web.C, w http.ResponseWriter, r *http.Request)
}

func init() {
	log.SetLevel(log.DebugLevel)
	//log.SetFlags(log.Ltime | log.Lshortfile)
	// cfg := profile.Config{
	// 	//MemProfile:     true,
	// 	CPUProfile:     true,
	// 	ProfilePath:    ".",  // store profiles in current directory
	// 	NoShutdownHook: true, // do not hook SIGINT
	// }
	// defer profile.Start(&cfg).Stop()
	// defer profile.Start(profile.CPUProfile).Stop()

}

var routes = []route{
	//{"GET", "/", root},
	{"GET", "/banks", getBanks},
	{"GET", "/banks/:bankId", getBank},
	{"GET", "/banks/:bankId/accounts", getBankAccounts},
	{"GET", "/banks/:bankId/accounts/:accountId", getBankAccount},
}

func main() {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://127.0.0.1:9000"},
	})
	//goji := LoadGoji(routes)
	goji.Use(c.Handler)
	Load()
	// goji.Get("/", root)
	// // @Title getBanks
	// // @Description retrieves banks
	// // @Success 200 {array} Bank
	// // @Router /banks [get]
	// goji.Get("/banks", getBanks)
	// goji.Get("/banks/:bankId", getBank)
	// goji.Get("/banks/:bankId/accounts", getBankAccounts)
	goji.Serve()
}

// Load loads all of the available routes.
func Load() {
	for _, route := range routes {
		switch route.method {
		case "GET":
			goji.Get(route.path, route.handler)
		default:
			panic("UNKNOWN HTTP METHOD: " + route.method)
		}
	}
	publicPath, _ := filepath.Abs("public")
	log.Info(publicPath)
	goji.Get("/*", http.FileServer(http.Dir(publicPath)))
}

func root(c web.C, w http.ResponseWriter, r *http.Request) {
	common.Index(w, r, "goji")
}

func getBankAccounts(c web.C, w http.ResponseWriter, r *http.Request) {
	common.GetBankAccounts(w, r, c.URLParams["bankId"])
}

func getBanks(c web.C, w http.ResponseWriter, r *http.Request) {
	common.GetBanks(w, r)
}

func getBank(c web.C, w http.ResponseWriter, r *http.Request) {
	common.GetBank(w, r, c.URLParams["bankId"])
}

func getBankAccount(c web.C, w http.ResponseWriter, r *http.Request) {
	common.GetBankAccount(w, r, c.URLParams["bankId"], c.URLParams["accountId"])
}
