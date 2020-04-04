package views

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MarErm27/webTrainerWithUAdmin/models"
	"github.com/uadmin/uadmin"
)

// TaskHandler !
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is not logged in
	// if uadmin.IsAuthenticated(r) == nil {
	// 	http.Redirect(w, r, "/admin", http.StatusSeeOther)
	// 	return
	// }
	// r.URL.Path creates a new path called /todo
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/Task")

	var taskListTempName string

	taskListTempName = "templates/task.html"
	switch r.Method {
	case "GET":
		keys, ok := r.URL.Query()["NumberOfCollection"]

		if !ok || len(keys[0]) < 1 {
			log.Println("Url Param 'NumberOfCollection' is missing")
			return
		}
		key := keys[0]

		auctions := []models.Auction{}
		uadmin.Filter(&auctions, "number_of_collection = ?", key)
		// if len(auctions) == 0 {
		// 	createAuctions()
		// 	uadmin.Filter(&auction, "number_of_collection = ?", key)
		// }
		// results := AnaliticsPanelPageData{}
		data := []AnaliticsPanelPageData{}
		result := []map[string]interface{}{}
		sumData := SummaryAnaliticsPanelPageData{}
		// var data
		var auctionNumber int
		var tableNumber int
		tableNumber = 1
		lenOfAuctions := len(auctions)
		var sumOfAllInitialPrices float64
		for i := range auctions {
			uadmin.Preload(&auctions[i])
			sumOfAllInitialPrices += auctions[i].InitialPrice

			// Assigns the string of interface in each Todo fields
			var PulseTR string
			var RedTR bool
			if (auctions[i].Fall >= auctions[i].FallLimit || auctions[i].OurFall >= auctions[i].FallLimit) &&
				(auctions[i].Fall > 0 || auctions[i].OurFall > 0) {
				PulseTR = "pulse"
				RedTR = true
			}
			var showTimer bool
			if auctions[i].MinutesToNextStepOverLimit == 0 {
				showTimer = true
			}
			// var AreWeInvolved bool
			// if auctions[i].OurFall > 0 || auctions[i].Fall > 0 {
			// 	AreWeInvolved = true
			// }
			result = append(result, map[string]interface{}{
				"Title":                           auctions[i].AuctionID,
				"InitialPrice":                    addSpacesToPrice(fmt.Sprintf("%.2f", auctions[i].InitialPrice)),
				"OurAuction":                      auctions[i].OurAuction,
				"FallLimit":                       auctions[i].FallLimit,
				"PlanFallLimit":                   auctions[i].PlanFallLimit,
				"PersonResponsibleForLimit":       auctions[i].PersonResponsibleForLimit,
				"OurFall":                         auctions[i].OurFall,
				"Fall":                            auctions[i].Fall,
				"PulseTR":                         PulseTR,
				"RedTR":                           RedTR,
				"MinutesToNextStepOverLimit":      auctions[i].MinutesToNextStepOverLimit,
				"TimeWhenWeGetMinutesBeforeLimit": auctions[i].TimeWhenWeGetMinutesBeforeLimit.Unix(),
				"ShowTimer":                       showTimer,
				"WeWon":                           auctions[i].WeWon,
				"AuctionIsFinished":               auctions[i].AuctionIsFinished,
				"AreWeInvolved":                   auctions[i].AreWeInvolved,
				"LotNumber":                       auctions[i].LotNumber,
			})
			auctionNumber += 1
			if auctionNumber == 8 || (i+1) == lenOfAuctions {
				set := models.Set{}
				uadmin.Filter(&set, "name = ?", key)
				data = append(data, AnaliticsPanelPageData{
					TableNumber:               auctions[i].NumberOfCollection,
					Table:                     result,
					ResponsibleFromGroupOne:   set.ResponsibleFromGroupOne,
					ResponsibleFromGroupTwo:   set.ResponsibleFromGroupTwo,
					ResponsibleFromGroupThree: set.ResponsibleFromGroupThree,
				})
				tableNumber += 1
				auctionNumber = 0
				result = []map[string]interface{}{}
			}
			// fmt.Println("auctionNumber", auctionNumber)
		}

		sumData = SummaryAnaliticsPanelPageData{
			SummOfAllInitialPrices: addSpacesToPrice(fmt.Sprintf("%.2f", sumOfAllInitialPrices)),
			Data:                   data,
		}
		// data := TodoPageData{
		// 	Results: results,
		// 	Todos:     results,
		// }
		// Pass TodoList data object to the specified HTML path
		uadmin.RenderHTML(w, r, taskListTempName, sumData)

	case "POST":
		Method := r.FormValue("Method")
		AuctionID := r.FormValue("AuctionID")

		res := map[string]interface{}{}
		switch Method {

		case "SaveOurFall":
			OurFall := r.FormValue("OurFall")
			auction := models.Auction{}
			uadmin.Filter(&auction, "auction_id = ?", AuctionID)

			s, err := strconv.ParseFloat(OurFall, 64)
			if err != nil {
				res["status"] = "ERROR"
				res["err_msg"] = "OurFall " + OurFall + " can't be parsed to float64."
				uadmin.ReturnJSON(w, r, res)
				return
			}

			auction.OurFall = s
			uadmin.Save(&auction)
			fmt.Println("SaveOurFall")
		case "SaveFall":
			Fall := r.FormValue("Fall")
			auction := models.Auction{}
			uadmin.Filter(&auction, "auction_id = ?", AuctionID)
			s, err := strconv.ParseFloat(Fall, 64)
			if err != nil {
				res["status"] = "ERROR"
				res["err_msg"] = "Fall " + Fall + " can't be parsed to float64."
				uadmin.ReturnJSON(w, r, res)
				return
			}
			auction.Fall = s
			uadmin.Save(&auction)
			fmt.Println("SaveFall")
		case "tellHowManyMinutesLeftToNextStepOverTheLimit":
			Minutes := r.FormValue("Minutes")
			auction := models.Auction{}
			uadmin.Filter(&auction, "auction_id = ?", AuctionID)
			s, err := strconv.ParseFloat(Minutes, 64)
			if err != nil {
				res["status"] = "ERROR"
				res["err_msg"] = "Minutes " + Minutes + " can't be parsed to float64."
				uadmin.ReturnJSON(w, r, res)
				return
			}
			auction.MinutesToNextStepOverLimit = s
			auction.TimeWhenWeGetMinutesBeforeLimit = time.Now()
			uadmin.Save(&auction)
		case "tellThatWeWon":
			auction := models.Auction{}
			uadmin.Filter(&auction, "auction_id = ?", AuctionID)
			auction.AuctionIsFinished = true
			auction.WeWon = true
			uadmin.Save(&auction)
		case "cancelWon":
			auction := models.Auction{}
			uadmin.Filter(&auction, "auction_id = ?", AuctionID)
			auction.AuctionIsFinished = false
			auction.WeWon = false
			uadmin.Save(&auction)
		case "tellThatWeLoosed":
			auction := models.Auction{}
			uadmin.Filter(&auction, "auction_id = ?", AuctionID)
			auction.AuctionIsFinished = true
			auction.WeWon = false
			uadmin.Save(&auction)
		case "cancelLoose":
			auction := models.Auction{}
			uadmin.Filter(&auction, "auction_id = ?", AuctionID)
			auction.AuctionIsFinished = false
			auction.WeWon = false
			uadmin.Save(&auction)
		case "responsibleFromGroupOne":
			User := r.FormValue("User")
			keys, ok := r.URL.Query()["NumberOfCollection"]

			if !ok || len(keys[0]) < 1 {
				log.Println("Url Param 'NumberOfCollection' is missing")
				return
			}
			key := keys[0]
			// fmt.Println("key", key)
			set := models.Set{}
			uadmin.Filter(&set, "name = ?", key)
			set.ResponsibleFromGroupOne = User
			uadmin.Save(&set)
		case "responsibleFromGroupTwo":
			User := r.FormValue("User")
			keys, ok := r.URL.Query()["NumberOfCollection"]

			if !ok || len(keys[0]) < 1 {
				log.Println("Url Param 'NumberOfCollection' is missing")
				return
			}
			key := keys[0]
			// fmt.Println("key", key)
			set := models.Set{}
			uadmin.Filter(&set, "name = ?", key)
			set.ResponsibleFromGroupTwo = User
			uadmin.Save(&set)
		case "responsibleFromGroupThree":
			User := r.FormValue("User")
			keys, ok := r.URL.Query()["NumberOfCollection"]

			if !ok || len(keys[0]) < 1 {
				log.Println("Url Param 'NumberOfCollection' is missing")
				return
			}
			key := keys[0]
			// fmt.Println("key", key)
			set := models.Set{}
			uadmin.Filter(&set, "name = ?", key)
			set.ResponsibleFromGroupThree = User
			uadmin.Save(&set)
		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
	// Pass TodoList data object to the specified HTML path
	// uadmin.RenderHTML(w, r, taskListTempName, nil)
}
