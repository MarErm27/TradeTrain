package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/MarErm27/webTrainerWithUAdmin/models"
	"github.com/MarErm27/webTrainerWithUAdmin/views"
	"github.com/uadmin/uadmin"
)

func showMainPage(w http.ResponseWriter, r *http.Request) {
	// data := TodoPageData{
	// 	PageTitle: "My TODO list",
	// 	// Todos:     Auctions,
	// }
	mainPageTempName := "templates/mainPage.html"
	uadmin.RenderHTML(w, r, mainPageTempName, nil)
}
func showSSL(w http.ResponseWriter, r *http.Request) {
	mainPageTempName := ".well-known/acme-challenge/b4f_4q3edwdA8n9UONaJ2U913snkTb_RcwLeel_-Yqs"
	uadmin.RenderHTML(w, r, mainPageTempName, nil)
}

func redirect(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target,
		// see comments below and consider the codes 308, 302, or 301
		http.StatusTemporaryRedirect)
}
func showResume(w http.ResponseWriter, r *http.Request) {

	f, err := os.Open("Ermeshev_Marat.pdf")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	defer f.Close()

	//Set header
	w.Header().Set("Content-type", "application/pdf")

	//Stream to response
	if _, err := io.Copy(w, f); err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
	}
}

func main() {
	// go http.ListenAndServe(":80", http.HandlerFunc(redirect)) //включить перед обновлением сервера
	uadmin.Register(
		models.Auction{},
		models.AuctionVersion{},
		models.Offer{},
		models.Participant{},
		models.AllAuctionsDataAndSettings{},
		models.Set{},
	)
	// setting := uadmin.Setting{}
	// uadmin.Get(&setting, "code = ?", "uAdmin.RootURL")
	// setting.ParseFormValue([]string{"/admin/"})
	// setting.Save()
	uadmin.RegisterInlines(
		models.Auction{},
		map[string]string{
			"participant":    "AuctionID",
			"offer":          "AuctionID",
			"auctionversion": "AuctionID",
		},
	)

	uadmin.Port = 443
	uadmin.RootURL = "/admin/"
	setting := uadmin.Setting{}
	uadmin.Update(&setting, "Value", strconv.Itoa(uadmin.Port), "code = ?", "uAdmin.Port")
	uadmin.Update(&setting, "Value", uadmin.RootURL, "code = ?", "uAdmin.RootURL")
	uadmin.SiteName = "Auctions and analitics"

	http.HandleFunc("/", http.HandlerFunc(showMainPage))
	http.HandleFunc("/.well-known/acme-challenge/", http.HandlerFunc(showSSL))
	http.HandleFunc("/resume", http.HandlerFunc(showResume))

	http.HandleFunc("/http_handler/", views.HTTPHandler)
	uadmin.StartServer() //выключить перед обновлением сервера
	// uadmin.StartSecureServer("mycert.pem", "private.pem") //включить перед обновлением сервера
}
