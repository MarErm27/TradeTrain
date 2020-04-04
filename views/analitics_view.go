package views

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MarErm27/webTrainerWithUAdmin/models"
	"github.com/uadmin/uadmin"
)

// AnaliticsPanelPageData ...
type AnaliticsPanelPageData struct {
	TableNumber               int
	Table                     []map[string]interface{}
	BtnColor                  string
	ResponsibleFromGroupOne   string
	ResponsibleFromGroupTwo   string
	ResponsibleFromGroupThree string
}

// SummaryAnaliticsPanelPageData ...
type SummaryAnaliticsPanelPageData struct {
	SumOfOurInitialPrices    string
	SumOfNotOurInitialPrices string
	SummOfAllInitialPrices   string

	SumOurFallRub    string
	SumNotOurFallRub string
	SumAllFallRub    string

	SumOurAdditionalFall    string
	SumNotOurAdditionalFall string
	SumAllAdditionalFall    string

	SumOurPlanFallLimit    string
	SumNotOurPlanFallLimit string
	SumAllPlanFallLimit    string

	SumOurFallPercent    string
	SumNotOurFallPercent string
	SumAllFallPercent    string

	SumOurPlanFallLimitPercent    string
	SumNotOurPlanFallLimitPercent string
	SumAllPlanFallLimitPercent    string

	Economy            string
	AwerageFallEconomy string

	OurAverageStartLimit    string
	NotOurAverageStartLimit string
	AllAverageStartLimit    string

	OurAverageFall      string
	NotOurAverageFall   string
	AllAverageStartFall string

	OurDifAverageFall      string
	NotOurDifAverageFall   string
	AllDifAverageStartFall string

	OurVehiclesGor    string
	NotOurVehiclesGor string
	GORVehiclesGor    string

	OurVehiclesVIT    string
	NotOurVehiclesVIT string
	GORVehiclesVIT    string

	OurVehiclesZEL    string
	NotOurVehiclesZEL string
	GORVehiclesZEL    string

	OurVehiclesALL    string
	NotOurVehiclesALL string
	AllVehiclesALL    string

	Data []AnaliticsPanelPageData
}

// AnalyticsHandler !
func AnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is not logged in
	// if uadmin.IsAuthenticated(r) == nil {
	// 	http.Redirect(w, r, "/admin", http.StatusSeeOther)
	// 	return
	// }
	// r.URL.Path creates a new path called /todo
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/Analytics")

	var taskListTempName string

	taskListTempName = "templates/analyticsPanel.html"

	// // Pass TodoList data object to the specified HTML path
	// uadmin.RenderHTML(w, r, taskListTempName, nil)

	switch r.Method {
	case "GET":
		auctions := []models.Auction{}
		uadmin.All(&auctions)
		if len(auctions) == 0 {
			createAuctions()
			uadmin.All(&auctions)
		}
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
		var sumOfOurInitialPrices float64
		var sumOfNotOurInitialPrices float64

		var sumOurFallRub float64
		var sumNotOurFallRub float64
		var sumAllFallRub float64

		var sumOurAdditionalFall float64
		var sumNotOurAdditionalFall float64
		var sumAllAdditionalFall float64

		var sumOurPlanFallLimit float64
		var sumNotOurPlanFallLimit float64
		var sumAllPlanFallLimit float64

		var sumOurAreaVehicles float64
		var sumNotOurAreaVehicles float64
		var sumAllAreaVehicles float64

		var sumGOROurAreaVehicles float64
		var sumGORNotOurAreaVehicles float64
		var sumGORAllAreaVehicles float64

		var sumVITOurAreaVehicles float64
		var sumVITNotOurAreaVehicles float64
		var sumVITAllAreaVehicles float64

		var sumZELOurAreaVehicles float64
		var sumZELNotOurAreaVehicles float64
		var sumZELAllAreaVehicles float64

		var economy float64
		var sumOurFallRubEconomy float64

		var sumOurPlanFallLimitEconomy float64
		var allInitialPricesWhereWeWon float64
		var fallDownWhereWeWon float64

		var ourFallLimits float64
		var allOurStartFalls float64
		var notOurFallLimits float64
		var allFallLimits float64

		// var thereWasRedTRInASer bool
		timeLayout := "15:04:05"
		for i := range auctions {
			uadmin.Preload(&auctions[i])

			ver := []models.AuctionVersion{}
			uadmin.Filter(&ver, "auction_id = ?", auctions[i].ID)
			vers := []map[string]interface{}{}
			// for j := range ver {
			for j := len(ver) - 1; j >= 0; j-- {
				vers = append(vers, map[string]interface{}{
					"FallLimit":                 ver[j].FallLimit,
					"PersonResponsibleForLimit": ver[j].PersonResponsibleForLimit,
					"Date":                      ver[j].Date.Format(timeLayout),
				})
			}
			sumOfAllInitialPrices += auctions[i].InitialPrice

			var status bool
			if auctions[i].Status == "Прием заявок" {
				status = true
			}

			var PulseTR string
			var RedTR bool
			var ourFallIsBigger string
			var fallIsBigger string
			if (auctions[i].OurFall < auctions[i].Fall || auctions[i].OurFall >= auctions[i].FallLimit) && !auctions[i].AuctionIsFinished {
				ourFallIsBigger = "bg-danger text-light"
			} else if !auctions[i].AuctionIsFinished {
				ourFallIsBigger = "bg-success text-light"
			}
			if (auctions[i].Fall >= auctions[i].FallLimit) && !auctions[i].AuctionIsFinished {
				fallIsBigger = "bg-danger text-light"
			}
			if (auctions[i].Fall >= auctions[i].FallLimit || auctions[i].OurFall >= auctions[i].FallLimit) &&
				(auctions[i].Fall > 0 || auctions[i].OurFall > 0) &&
				(!auctions[i].WeWon && !auctions[i].AuctionIsFinished) {
				PulseTR = "pulse"
				RedTR = true
				// "bg-danger text-light"
			}

			var showTimer bool
			if auctions[i].MinutesToNextStepOverLimit == 0 {
				showTimer = true
			}
			// Assigns the string of interface in each Todo fields

			var fall float64

			if auctions[i].OurFall > auctions[i].Fall {
				fall = auctions[i].OurFall
			} else {
				fall = auctions[i].Fall
			}

			// var fallLimit float64
			// if auctions[i].WeWon && auctions[i].AuctionIsFinished {
			// 	fallLimit = fall
			// } else {
			// 	fallLimit = auctions[i].FallLimit
			// }

			var fallRub float64
			fallRub = auctions[i].InitialPrice - fall*auctions[i].InitialPrice/100
			var fallRubPlus05 float64
			fallRubPlus05 = auctions[i].InitialPrice - (fall+0.5)*auctions[i].InitialPrice/100
			var fallLimitRubForEconomy float64
			fallLimitRubForEconomy = (auctions[i].InitialPrice - auctions[i].FallLimit*auctions[i].InitialPrice/100)
			// var currentFall float64
			var fallLimitRub float64
			if auctions[i].WeWon && auctions[i].AuctionIsFinished {
				fallLimitRub = fallRub
				// currentFall = fall
			} else {
				fallLimitRub = (auctions[i].InitialPrice - auctions[i].FallLimit*auctions[i].InitialPrice/100)
				// currentFall = auctions[i].FallLimit
			}
			var planFallLimitRub float64
			planFallLimitRub = (auctions[i].InitialPrice - auctions[i].PlanFallLimit*auctions[i].InitialPrice/100)
			var additionalFall float64
			additionalFall = fallLimitRub - planFallLimitRub
			var doWeHaveInformationAboutCurrentFall bool
			if auctions[i].OurFall != 0 || auctions[i].Fall != 0 {
				doWeHaveInformationAboutCurrentFall = true
			}

			var colorOfFalls string
			if (auctions[i].OurFall < auctions[i].Fall && fall >= auctions[i].FallLimit) &&
				(!auctions[i].WeWon && !auctions[i].AuctionIsFinished) {
				colorOfFalls = "bg-danger text-light"
			} else {
				colorOfFalls = "bg-success text-light"

			}

			if auctions[i].WeWon {
				economy += auctions[i].InitialPrice - auctions[i].InitialPrice*fall
			}

			// var AreWeInvolved bool

			// if (auctions[i].OurFall > 0 || auctions[i].Fall > 0) {
			// 	AreWeInvolved = true
			// }
			sumAllFallRub += fallRub
			sumAllAdditionalFall += additionalFall
			sumAllPlanFallLimit += planFallLimitRub

			sumAllAreaVehicles += auctions[i].GOR*auctions[i].NumberOfVehicles + auctions[i].VIT*auctions[i].NumberOfVehicles + auctions[i].ZEL*auctions[i].NumberOfVehicles
			sumGORAllAreaVehicles += auctions[i].GOR * auctions[i].NumberOfVehicles
			sumVITAllAreaVehicles += auctions[i].VIT * auctions[i].NumberOfVehicles
			sumZELAllAreaVehicles += auctions[i].ZEL * auctions[i].NumberOfVehicles
			allFallLimits += fallLimitRub
			if auctions[i].OurAuction {
				sumOfOurInitialPrices += auctions[i].InitialPrice
				sumOurFallRub += fallRub
				sumOurAdditionalFall += additionalFall
				sumOurPlanFallLimit += planFallLimitRub

				ourFallLimits += fallLimitRub
				// if auctions[i].WeWon && auctions[i].AuctionIsFinished {
				// 	ourFallLimits +=
				// }else{
				// 	ourFallLimits += fallLimitRub
				// }

				sumOurAreaVehicles += auctions[i].GOR*auctions[i].NumberOfVehicles + auctions[i].VIT*auctions[i].NumberOfVehicles + auctions[i].ZEL*auctions[i].NumberOfVehicles
				sumGOROurAreaVehicles += auctions[i].GOR * auctions[i].NumberOfVehicles
				sumVITOurAreaVehicles += auctions[i].VIT * auctions[i].NumberOfVehicles
				sumZELOurAreaVehicles += auctions[i].ZEL * auctions[i].NumberOfVehicles
				allOurStartFalls += auctions[i].InitialPrice - (auctions[i].InitialPrice*auctions[i].PlanFallLimit)/100
				// allOurStartFalls += PlanFallLimitRub
			} else {
				sumOfNotOurInitialPrices += auctions[i].InitialPrice
				sumNotOurFallRub += fallRub
				sumNotOurAdditionalFall += additionalFall
				sumNotOurPlanFallLimit += planFallLimitRub

				notOurFallLimits += fallLimitRub
				// if auctions[i].WeWon && auctions[i].AuctionIsFinished {
				// 	notOurFallLimits +=
				// }else{
				// 	notOurFallLimits += fallLimitRub
				// }

				sumNotOurAreaVehicles += auctions[i].GOR*auctions[i].NumberOfVehicles + auctions[i].VIT*auctions[i].NumberOfVehicles + auctions[i].ZEL*auctions[i].NumberOfVehicles
				sumGORNotOurAreaVehicles += auctions[i].GOR * auctions[i].NumberOfVehicles
				sumVITNotOurAreaVehicles += auctions[i].VIT * auctions[i].NumberOfVehicles
				sumZELNotOurAreaVehicles += auctions[i].ZEL * auctions[i].NumberOfVehicles
			}
			if auctions[i].WeWon && auctions[i].AuctionIsFinished {
				allInitialPricesWhereWeWon += auctions[i].InitialPrice
				fallDownWhereWeWon += auctions[i].InitialPrice - fallRub

				sumOurFallRubEconomy += fallRub
				// sumOurPlanFallLimitEconomy += planFallLimitRub
				sumOurPlanFallLimitEconomy += fallLimitRubForEconomy

				// fallLimitRub

			} else {

			}
			var branch string
			if auctions[i].GOR == 1 {
				branch += "ГОР"
			}
			if auctions[i].VIT == 1 {
				if branch == "" {
					branch += "ВИТ"
				} else {
					branch += ";ВИТ"
				}
			}
			if auctions[i].ZEL == 1 {
				if branch == "" {
					branch += "ЗЕЛ"
				} else {
					branch += ";ЗЕЛ"
				}
			}

			// auctions[i].InitialPrice - (auctions[i].InitialPrice * PlanFallLimit)

			result = append(result, map[string]interface{}{
				"Title":                     auctions[i].AuctionID,
				"InitialPrice":              addSpacesToPrice(fmt.Sprintf("%.2f", auctions[i].InitialPrice)),
				"HalfPercentOfInitialPrice": addSpacesToPrice(fmt.Sprintf("%.2f", auctions[i].HalfPercentOfInitialPrice)),
				"Status":                    auctions[i].Status,
				"StatusBool":                status,
				"OurAuction":                auctions[i].OurAuction,
				"PersonResponsibleForOur":   auctions[i].PersonResponsibleForOur,
				"FallLimit":                 auctions[i].FallLimit,
				"FallLimitRub":              addSpacesToPrice(fmt.Sprintf("%.2f", fallLimitRub)),
				"PlanFallLimit":             auctions[i].PlanFallLimit,
				"PlanFallLimitRub":          addSpacesToPrice(fmt.Sprintf("%.2f", planFallLimitRub)),

				"EconomyPercent": fmt.Sprintf("%.2f", auctions[i].FallLimit-fall),
				// "EconomyPercent":                      fmt.Sprintf("%.2f", auctions[i].PlanFallLimit-fall),
				// "EconomyRub":                          addSpacesToPrice(fmt.Sprintf("%.2f", auctions[i].InitialPrice*(auctions[i].PlanFallLimit-fall)/100)),
				"EconomyRub": addSpacesToPrice(fmt.Sprintf("%.2f", auctions[i].InitialPrice*(auctions[i].FallLimit-fall)/100)),

				"AdditionalFallPercent":               auctions[i].FallLimit - auctions[i].PlanFallLimit,
				"AdditionalFallRub":                   addSpacesToPrice(fmt.Sprintf("%.2f", additionalFall)),
				"PersonResponsibleForLimit":           auctions[i].PersonResponsibleForLimit,
				"OurFall":                             auctions[i].OurFall,
				"Fall":                                auctions[i].Fall,
				"FallRub":                             addSpacesToPrice(fmt.Sprintf("%.2f", fallRub)),
				"fallRubPlus05":                       addSpacesToPrice(fmt.Sprintf("%.2f", fallRubPlus05)),
				"NMCKMinusFallRub":                    addSpacesToPrice(fmt.Sprintf("%.2f", auctions[i].InitialPrice-fallRub)),
				"NMCKMinusFallRubPlus05":              addSpacesToPrice(fmt.Sprintf("%.2f", auctions[i].InitialPrice-fallRubPlus05)),
				"DoWeHaveInformationAboutCurrentFall": doWeHaveInformationAboutCurrentFall,
				"PulseTR":                             PulseTR,
				"RedTR":                               RedTR,
				"OurFallIsBigger":                     ourFallIsBigger,
				"FallIsBigger":                        fallIsBigger,
				"MinutesToNextStepOverLimit":          auctions[i].MinutesToNextStepOverLimit,
				"TimeWhenWeGetMinutesBeforeLimit":     auctions[i].TimeWhenWeGetMinutesBeforeLimit.Unix(),
				"ShowTimer":                           showTimer,
				"AreWeInvolved":                       auctions[i].AreWeInvolved,
				"PersonResponsibleForInvolve":         auctions[i].PersonResponsibleForInvolve,
				"LotNumber":                           auctions[i].LotNumber,
				"Priority":                            auctions[i].Priority,
				"WeWon":                               auctions[i].WeWon,
				"AuctionIsFinished":                   auctions[i].AuctionIsFinished,
				"NumberOfVehiclesBK":                  auctions[i].NumberOfVehiclesBK,
				"NumberOfVehiclesSK1":                 auctions[i].NumberOfVehiclesSK1,
				"NumberOfVehiclesSK2":                 auctions[i].NumberOfVehiclesSK2,
				"NumberOfVehicles":                    auctions[i].NumberOfVehicles,
				"GOR":                                 auctions[i].GOR,
				"VIT":                                 auctions[i].VIT,
				"ZEL":                                 auctions[i].ZEL,
				"Branch":                              branch,
				"ColorOfFalls":                        colorOfFalls,
				"Vers":                                vers,
			})

			auctionNumber += 1
			// if RedTR {
			// 	thereWasRedTRInASer = true

			// }

			bootstrap4olors := []string{"primary", "secondary", "success", "danger", "warning", "info", "light"}

			if auctionNumber == 8 || (i+1) == lenOfAuctions {

				set := models.Set{}
				uadmin.Filter(&set, "name = ?", tableNumber)
				data = append(data, AnaliticsPanelPageData{
					TableNumber: tableNumber,
					Table:       result,
					// BtnColor:                  thereWasRedTRInASer,
					BtnColor: bootstrap4olors[tableNumber-1],

					ResponsibleFromGroupOne:   set.ResponsibleFromGroupOne,
					ResponsibleFromGroupTwo:   set.ResponsibleFromGroupTwo,
					ResponsibleFromGroupThree: set.ResponsibleFromGroupThree,
				})
				tableNumber += 1
				auctionNumber = 0
				result = []map[string]interface{}{}
				// thereWasRedTRInASer = false
			}

			// fmt.Println("auctionNumber", auctionNumber)
		}

		sumData = SummaryAnaliticsPanelPageData{

			SumOfOurInitialPrices:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOfOurInitialPrices)),
			SumOfNotOurInitialPrices: addSpacesToPrice(fmt.Sprintf("%.2f", sumOfNotOurInitialPrices)),
			SummOfAllInitialPrices:   addSpacesToPrice(fmt.Sprintf("%.2f", sumOfAllInitialPrices)),

			SumOurFallRub:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOfOurInitialPrices-sumOurFallRub)),
			SumNotOurFallRub: addSpacesToPrice(fmt.Sprintf("%.2f", sumOfNotOurInitialPrices-sumNotOurFallRub)),
			SumAllFallRub:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOfAllInitialPrices-sumAllFallRub)),

			SumOurFallPercent:    addSpacesToPrice(fmt.Sprintf("%.2f", 100-sumOurFallRub*100/sumOfOurInitialPrices)),
			SumNotOurFallPercent: addSpacesToPrice(fmt.Sprintf("%.2f", 100-sumNotOurFallRub*100/sumOfNotOurInitialPrices)),
			SumAllFallPercent:    addSpacesToPrice(fmt.Sprintf("%.2f", 100-sumAllFallRub*100/sumOfAllInitialPrices)),

			SumOurPlanFallLimit:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOurPlanFallLimit)),
			SumNotOurPlanFallLimit: addSpacesToPrice(fmt.Sprintf("%.2f", sumNotOurPlanFallLimit)),
			SumAllPlanFallLimit:    addSpacesToPrice(fmt.Sprintf("%.2f", sumAllPlanFallLimit)),

			SumOurAdditionalFall:    addSpacesToPrice(fmt.Sprintf("%.2f", -1*(sumOurPlanFallLimit-sumOurFallRub))),
			SumNotOurAdditionalFall: addSpacesToPrice(fmt.Sprintf("%.2f", -1*(sumNotOurPlanFallLimit-sumNotOurFallRub))),
			SumAllAdditionalFall:    addSpacesToPrice(fmt.Sprintf("%.2f", -1*(sumAllPlanFallLimit-sumAllFallRub))),

			// SumOurPlanFallLimit:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOurPlanFallLimit)),
			// SumNotOurPlanFallLimit: addSpacesToPrice(fmt.Sprintf("%.2f", sumNotOurPlanFallLimit)),
			// SumAllPlanFallLimit:    addSpacesToPrice(fmt.Sprintf("%.2f", sumAllPlanFallLimit)),

			// SumOurPlanFallLimitPercent:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOurFallRub*100/sumOurAdditionalFall)),
			// SumNotOurPlanFallLimitPercent: addSpacesToPrice(fmt.Sprintf("%.2f", sumNotOurFallRub*100/sumNotOurAdditionalFall)),
			// SumAllPlanFallLimitPercent:    addSpacesToPrice(fmt.Sprintf("%.2f", sumAllFallRub*100/sumAllPlanFallLimit)),

			SumOurPlanFallLimitPercent:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOurPlanFallLimit*100/sumOfOurInitialPrices-sumOurFallRub*100/sumOfOurInitialPrices)),
			SumNotOurPlanFallLimitPercent: addSpacesToPrice(fmt.Sprintf("%.2f", sumNotOurPlanFallLimit*100/sumOfNotOurInitialPrices-sumNotOurFallRub*100/sumOfNotOurInitialPrices)),
			SumAllPlanFallLimitPercent:    addSpacesToPrice(fmt.Sprintf("%.2f", sumAllPlanFallLimit*100/sumOfAllInitialPrices-sumAllFallRub*100/sumOfAllInitialPrices)),

			// Economy: addSpacesToPrice(fmt.Sprintf("%.2f", economy)),
			Economy:            addSpacesToPrice(fmt.Sprintf("%.2f", sumOurFallRubEconomy-sumOurPlanFallLimitEconomy)),
			AwerageFallEconomy: addSpacesToPrice(fmt.Sprintf("%.2f", fallDownWhereWeWon/allInitialPricesWhereWeWon*100)),

			OurAverageStartLimit:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOfOurInitialPrices-sumOurPlanFallLimit)),
			NotOurAverageStartLimit: addSpacesToPrice(fmt.Sprintf("%.2f", sumOfNotOurInitialPrices-sumNotOurPlanFallLimit)),
			AllAverageStartLimit:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOfAllInitialPrices-sumAllPlanFallLimit)),

			OurAverageFall:      addSpacesToPrice(fmt.Sprintf("%.2f", sumOfOurInitialPrices-ourFallLimits)),
			NotOurAverageFall:   addSpacesToPrice(fmt.Sprintf("%.2f", sumOfNotOurInitialPrices-notOurFallLimits)),
			AllAverageStartFall: addSpacesToPrice(fmt.Sprintf("%.2f", sumOfAllInitialPrices-allFallLimits)),

			OurDifAverageFall:      addSpacesToPrice(fmt.Sprintf("%.2f", (sumOfOurInitialPrices-sumOurPlanFallLimit)-(sumOfOurInitialPrices-ourFallLimits))),
			NotOurDifAverageFall:   addSpacesToPrice(fmt.Sprintf("%.2f", (sumOfNotOurInitialPrices-sumNotOurPlanFallLimit)-(sumOfNotOurInitialPrices-notOurFallLimits))),
			AllDifAverageStartFall: addSpacesToPrice(fmt.Sprintf("%.2f", (sumOfAllInitialPrices-sumAllPlanFallLimit)-(sumOfAllInitialPrices-allFallLimits))),

			// ALLAverageStartLimit:
			// sumAllAreaVehicles

			OurVehiclesGor:    fmt.Sprintf("%.2f", ((sumGOROurAreaVehicles * 100) / sumOurAreaVehicles)),
			NotOurVehiclesGor: fmt.Sprintf("%.2f", ((sumGORNotOurAreaVehicles * 100) / sumNotOurAreaVehicles)),
			GORVehiclesGor:    fmt.Sprintf("%.2f", ((sumGORAllAreaVehicles * 100) / sumAllAreaVehicles)),

			OurVehiclesVIT:    fmt.Sprintf("%.2f", ((sumVITOurAreaVehicles * 100) / sumOurAreaVehicles)),
			NotOurVehiclesVIT: fmt.Sprintf("%.2f", ((sumVITNotOurAreaVehicles * 100) / sumNotOurAreaVehicles)),
			GORVehiclesVIT:    fmt.Sprintf("%.2f", ((sumVITAllAreaVehicles * 100) / sumAllAreaVehicles)),

			OurVehiclesZEL:    fmt.Sprintf("%.2f", ((sumZELOurAreaVehicles * 100) / sumOurAreaVehicles)),
			NotOurVehiclesZEL: fmt.Sprintf("%.2f", ((sumZELNotOurAreaVehicles * 100) / sumNotOurAreaVehicles)),
			GORVehiclesZEL:    fmt.Sprintf("%.2f", ((sumZELAllAreaVehicles * 100) / sumAllAreaVehicles)),

			// GORVehiclesALL:    fmt.Sprintf("%.2f", ((sumGORAllAreaVehicles * 100) / sumAllAreaVehicles)),
			OurVehiclesALL:    fmt.Sprintf("%.2f", ((sumOurAreaVehicles * 100) / sumAllAreaVehicles)),
			NotOurVehiclesALL: fmt.Sprintf("%.2f", ((sumNotOurAreaVehicles * 100) / sumAllAreaVehicles)),
			AllVehiclesALL:    fmt.Sprintf("%.2f", ((sumAllAreaVehicles * 100) / sumAllAreaVehicles)),

			// sumOurAreaVehicles += auctions[i].GOR*auctions[i].NumberOfVehicles + auctions[i].VIT*auctions[i].NumberOfVehicles + auctions[i].ZEL*auctions[i].NumberOfVehicles
			// sumGOROurAreaVehicles += auctions[i].GOR * auctions[i].NumberOfVehicles
			// sumVITOurAreaVehicles += auctions[i].VIT * auctions[i].NumberOfVehicles
			// sumZELOurAreaVehicles += auctions[i].ZEL * auctions[i].NumberOfVehicles

			// sumNotOurAreaVehicles += auctions[i].GOR*auctions[i].NumberOfVehicles + auctions[i].VIT*auctions[i].NumberOfVehicles + auctions[i].ZEL*auctions[i].NumberOfVehicles
			// 	sumGORNotOurAreaVehicles += auctions[i].GOR * auctions[i].NumberOfVehicles
			// 	sumVITNotOurAreaVehicles += auctions[i].VIT * auctions[i].NumberOfVehicles
			// 	sumZELNotOurAreaVehicles += auctions[i].ZEL * auctions[i].NumberOfVehicles

			// sumAllAreaVehicles += auctions[i].GOR*auctions[i].NumberOfVehicles + auctions[i].VIT*auctions[i].NumberOfVehicles + auctions[i].ZEL*auctions[i].NumberOfVehicles
			// sumGORAllAreaVehicles += auctions[i].GOR * auctions[i].NumberOfVehicles
			// sumVITAllAreaVehicles += auctions[i].VIT * auctions[i].NumberOfVehicles
			// sumZELAllAreaVehicles += auctions[i].ZEL * auctions[i].NumberOfVehicles

			Data: data,
		}
		// data := TodoPageData{
		// 	Results: results,
		// 	Todos:     results,
		// }
		// Pass TodoList data object to the specified HTML path
		uadmin.RenderHTML(w, r, taskListTempName, sumData)

	case "POST":
		r.ParseForm() //Parse url parameters passed, then parse the response packet for the POST body (request body)
		// attention: If you do not call ParseForm method, the following data can not be obtained form
		//fmt.Println(r.Form) // print information on server side.
		Method := r.FormValue("Method")
		// res := map[string]interface{}{}
		switch Method {

		case "SavePreSettings":

			TableDataString := r.Form["TableData"][0]
			TableDataRows := strings.Split(TableDataString, "\n")
			var TableData [][]string
			var numberOfTable []string
			for _, j := range TableDataRows {
				TableDataCells := strings.Split(j, "&")
				if len(TableDataCells) == 1 {
					numberOfTable = TableDataCells
				} else {
					TableDataCells = TableDataCells[:len(TableDataCells)-1]
					data := append(numberOfTable, TableDataCells...)
					TableData = append(TableData, data)
				}
			}

			fmt.Println(TableData)
			var auctionNumber int
			res := map[string]interface{}{}
			lenOfAuctions := len(TableData)
			for i, j := range TableData {
				auction := []models.Auction{}
				uadmin.Filter(&auction, "auction_id = ?", j[1])
				if j[3] == "true" {
					auction[0].OurAuction = true
				} else if j[3] == "false" {
					auction[0].OurAuction = false
				}
				numberOfCollection, err := strconv.Atoi(j[0])
				if err != nil {
					res["status"] = "ERROR"
					res["err_msg"] = "Number of collection " + j[0] + " can't be parsed to int."
					uadmin.ReturnJSON(w, r, res)
					return
				}
				auction[0].NumberOfCollection = numberOfCollection
				auctionNumber += 1
				if auctionNumber == 8 || (i+1) == lenOfAuctions {
					auctionNumber = 0
				}
				auction[0].NumberInCollection = auctionNumber
				initialPrice, err := strconv.ParseFloat(strings.ReplaceAll(j[2], " ", ""), 64)
				if err != nil {
					res["status"] = "ERROR"
					res["err_msg"] = "Initial price " + j[2] + " can't be parsed to float64."
					uadmin.ReturnJSON(w, r, res)
					return
				}

				auction[0].InitialPrice = initialPrice
				fallLimit, err := strconv.ParseFloat(strings.ReplaceAll(j[4], " ", ""), 64)
				if err != nil {
					res["status"] = "ERROR"
					res["err_msg"] = "Fall limit " + j[4] + " can't be parsed to float64."
					uadmin.ReturnJSON(w, r, res)
					return
				}

				auction[0].FallLimit = fallLimit
				planFallLimit, err := strconv.ParseFloat(strings.ReplaceAll(j[5], " ", ""), 64)

				if err != nil {
					res["status"] = "ERROR"
					res["err_msg"] = "Fall limit " + j[5] + " can't be parsed to float64."
					uadmin.ReturnJSON(w, r, res)
					return
				}
				auction[0].PlanFallLimit = planFallLimit

				auction[0].PersonResponsibleForLimit = j[6]
				uadmin.Save(&auction[0])
			}
			fmt.Println("path", r.URL.Path)
			fmt.Println("scheme", r.URL.Scheme)
			fmt.Println(r.Form["url_long"])
			for k, v := range r.Form {
				fmt.Println("key:", k)
				fmt.Println("val:", strings.Join(v, ""))
			}
			fmt.Fprintf(w, "Hello SavePreSettings!")
		case "setActPermission":
			User := r.FormValue("User")
			AuctionID := r.FormValue("AuctionID")
			AreWeInvolvedAjax := r.FormValue("AreWeInvolved")
			auction := models.Auction{}
			uadmin.Filter(&auction, "auction_id = ?", AuctionID)
			var AreWeInvolved bool
			if AreWeInvolvedAjax == "1" {
				AreWeInvolved = true
			}
			auction.AreWeInvolved = AreWeInvolved
			// auction.AuthorOfInvolved = User
			auction.PersonResponsibleForInvolve = User
			uadmin.Save(&auction)
			fmt.Fprintf(w, "Hello setActPermission!")
			// lastVer := models.InvolvedVersion{}
			// modelschema := uadmin.Schema["involvedversion"]
			// uadmin.FilterList(&modelschema, "id", false, 0, 1, &lastVer, "auction_id = ?", auction.ID)
			// if lastVer.AreWeInvolved != auction.AreWeInvolved {

			// 	ver := models.InvolvedVersion{}
			// 	ver.Date = time.Now()
			// 	ver.AuctionID = auction.ID
			// 	ver.AreWeInvolved = auction.AreWeInvolved
			// 	ver.PersonResponsibleForInvolve = auction.PersonResponsibleForInvolve

			// 	// Counts the version number based on the DocumentID and increment it
			// 	// by 1
			// 	ver.Number = uadmin.Count([]models.InvolvedVersion{}, "auction_id = ?", auction.ID) + 1

			// 	// Save the document version
			// 	uadmin.Save(&ver)
			//fmt.Fprintf(w, "Hello SaveChanges!")
			// } else {
			// 	fmt.Fprintf(w, "Act permission didn't changed!")
			// }

		case "ChangeFallLimit":
			res := map[string]interface{}{}
			User := r.FormValue("User")
			AuctionID := r.FormValue("AuctionID")
			FallLimit := r.FormValue("FallLimit")

			auction := models.Auction{}
			uadmin.Filter(&auction, "auction_id = ?", AuctionID)
			s, err := strconv.ParseFloat(FallLimit, 64)
			if err != nil {
				res["status"] = "ERROR"
				res["err_msg"] = "Fall " + FallLimit + " can't be parsed to float64."
				uadmin.ReturnJSON(w, r, res)
				return
			}
			auction.FallLimit = s
			auction.PersonResponsibleForLimit = User
			uadmin.Save(&auction)
			fmt.Println("SaveFall")

			lastVer := models.AuctionVersion{}
			// uadmin.All(&lastVer)
			// uadmin.Filter(&lastVer, "auction_id = ?", auction.ID)
			modelschema := uadmin.Schema["auctionversion"]

			uadmin.FilterList(&modelschema, "id", false, 0, 1, &lastVer, "auction_id = ?", auction.ID)

			if lastVer.FallLimit != auction.FallLimit {

				ver := models.AuctionVersion{}
				ver.Date = time.Now()
				ver.AuctionID = auction.ID
				ver.FallLimit = auction.FallLimit
				ver.PersonResponsibleForLimit = auction.PersonResponsibleForLimit

				// Counts the version number based on the DocumentID and increment it
				// by 1
				ver.Number = uadmin.Count([]models.AuctionVersion{}, "auction_id = ?", auction.ID) + 1

				// Save the document version
				uadmin.Save(&ver)

				auctions := []models.Auction{}
				uadmin.All(&auctions)
				if len(auctions) == 0 {
					createAuctions()
					uadmin.All(&auctions)
				}
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
				var sumOfOurInitialPrices float64
				var sumOfNotOurInitialPrices float64

				var sumOurFallRub float64
				var sumNotOurFallRub float64
				var sumAllFallRub float64

				var sumOurAdditionalFall float64
				var sumNotOurAdditionalFall float64
				var sumAllAdditionalFall float64

				var sumOurPlanFallLimit float64
				var sumNotOurPlanFallLimit float64
				var sumAllPlanFallLimit float64

				var sumOurAreaVehicles float64
				var sumNotOurAreaVehicles float64
				var sumAllAreaVehicles float64

				var sumGOROurAreaVehicles float64
				var sumGORNotOurAreaVehicles float64
				var sumGORAllAreaVehicles float64

				var sumVITOurAreaVehicles float64
				var sumVITNotOurAreaVehicles float64
				var sumVITAllAreaVehicles float64

				var sumZELOurAreaVehicles float64
				var sumZELNotOurAreaVehicles float64
				var sumZELAllAreaVehicles float64

				var economy float64
				var sumOurFallRubEconomy float64

				var sumOurPlanFallLimitEconomy float64
				var allInitialPricesWhereWeWon float64
				var fallDownWhereWeWon float64

				var ourFallLimits float64
				var allOurStartFalls float64
				var notOurFallLimits float64
				var allFallLimits float64

				// var thereWasRedTRInASer bool
				timeLayout := "15:04:05"
				for i := range auctions {
					uadmin.Preload(&auctions[i])

					ver := []models.AuctionVersion{}
					uadmin.Filter(&ver, "auction_id = ?", auctions[i].ID)
					vers := []map[string]interface{}{}
					// for j := range ver {
					for j := len(ver) - 1; j >= 0; j-- {
						vers = append(vers, map[string]interface{}{
							"FallLimit":                 ver[j].FallLimit,
							"PersonResponsibleForLimit": ver[j].PersonResponsibleForLimit,
							"Date":                      ver[j].Date.Format(timeLayout),
						})
					}
					sumOfAllInitialPrices += auctions[i].InitialPrice

					var status bool
					if auctions[i].Status == "Прием заявок" {
						status = true
					}

					var PulseTR string
					var RedTR bool
					if (auctions[i].Fall >= auctions[i].FallLimit || auctions[i].OurFall >= auctions[i].FallLimit) &&
						(auctions[i].Fall > 0 || auctions[i].OurFall > 0) &&
						(!auctions[i].WeWon && !auctions[i].AuctionIsFinished) {
						PulseTR = "pulse"
						RedTR = true
					}

					var showTimer bool
					if auctions[i].MinutesToNextStepOverLimit == 0 {
						showTimer = true
					}
					// Assigns the string of interface in each Todo fields

					var fall float64

					if auctions[i].OurFall > auctions[i].Fall {
						fall = auctions[i].OurFall
					} else {
						fall = auctions[i].Fall
					}

					// var fallLimit float64
					// if auctions[i].WeWon && auctions[i].AuctionIsFinished {
					// 	fallLimit = fall
					// } else {
					// 	fallLimit = auctions[i].FallLimit
					// }

					var fallRub float64
					fallRub = auctions[i].InitialPrice - fall*auctions[i].InitialPrice/100
					var fallRubPlus05 float64
					fallRubPlus05 = auctions[i].InitialPrice - (fall+0.5)*auctions[i].InitialPrice/100
					var fallLimitRubForEconomy float64
					fallLimitRubForEconomy = (auctions[i].InitialPrice - auctions[i].FallLimit*auctions[i].InitialPrice/100)
					// var currentFall float64
					var fallLimitRub float64
					if auctions[i].WeWon && auctions[i].AuctionIsFinished {
						fallLimitRub = fallRub
						// currentFall = fall
					} else {
						fallLimitRub = (auctions[i].InitialPrice - auctions[i].FallLimit*auctions[i].InitialPrice/100)
						// currentFall = auctions[i].FallLimit
					}
					var planFallLimitRub float64
					planFallLimitRub = (auctions[i].InitialPrice - auctions[i].PlanFallLimit*auctions[i].InitialPrice/100)
					var additionalFall float64
					additionalFall = fallLimitRub - planFallLimitRub
					var doWeHaveInformationAboutCurrentFall bool
					if auctions[i].OurFall != 0 || auctions[i].Fall != 0 {
						doWeHaveInformationAboutCurrentFall = true
					}

					var colorOfFalls string
					if (auctions[i].OurFall < auctions[i].Fall && fall >= auctions[i].FallLimit) &&
						(!auctions[i].WeWon && !auctions[i].AuctionIsFinished) {
						colorOfFalls = "bg-danger text-light"
					} else {
						colorOfFalls = "bg-success text-light"

					}

					if auctions[i].WeWon {
						economy += auctions[i].InitialPrice - auctions[i].InitialPrice*fall
					}

					// var AreWeInvolved bool

					// if (auctions[i].OurFall > 0 || auctions[i].Fall > 0) {
					// 	AreWeInvolved = true
					// }
					sumAllFallRub += fallRub
					sumAllAdditionalFall += additionalFall
					sumAllPlanFallLimit += planFallLimitRub

					sumAllAreaVehicles += auctions[i].GOR*auctions[i].NumberOfVehicles + auctions[i].VIT*auctions[i].NumberOfVehicles + auctions[i].ZEL*auctions[i].NumberOfVehicles
					sumGORAllAreaVehicles += auctions[i].GOR * auctions[i].NumberOfVehicles
					sumVITAllAreaVehicles += auctions[i].VIT * auctions[i].NumberOfVehicles
					sumZELAllAreaVehicles += auctions[i].ZEL * auctions[i].NumberOfVehicles
					allFallLimits += fallLimitRub
					if auctions[i].OurAuction {
						sumOfOurInitialPrices += auctions[i].InitialPrice
						sumOurFallRub += fallRub
						sumOurAdditionalFall += additionalFall
						sumOurPlanFallLimit += planFallLimitRub

						ourFallLimits += fallLimitRub
						// if auctions[i].WeWon && auctions[i].AuctionIsFinished {
						// 	ourFallLimits +=
						// }else{
						// 	ourFallLimits += fallLimitRub
						// }

						sumOurAreaVehicles += auctions[i].GOR*auctions[i].NumberOfVehicles + auctions[i].VIT*auctions[i].NumberOfVehicles + auctions[i].ZEL*auctions[i].NumberOfVehicles
						sumGOROurAreaVehicles += auctions[i].GOR * auctions[i].NumberOfVehicles
						sumVITOurAreaVehicles += auctions[i].VIT * auctions[i].NumberOfVehicles
						sumZELOurAreaVehicles += auctions[i].ZEL * auctions[i].NumberOfVehicles
						allOurStartFalls += auctions[i].InitialPrice - (auctions[i].InitialPrice*auctions[i].PlanFallLimit)/100
						// allOurStartFalls += PlanFallLimitRub
					} else {
						sumOfNotOurInitialPrices += auctions[i].InitialPrice
						sumNotOurFallRub += fallRub
						sumNotOurAdditionalFall += additionalFall
						sumNotOurPlanFallLimit += planFallLimitRub

						notOurFallLimits += fallLimitRub
						// if auctions[i].WeWon && auctions[i].AuctionIsFinished {
						// 	notOurFallLimits +=
						// }else{
						// 	notOurFallLimits += fallLimitRub
						// }

						sumNotOurAreaVehicles += auctions[i].GOR*auctions[i].NumberOfVehicles + auctions[i].VIT*auctions[i].NumberOfVehicles + auctions[i].ZEL*auctions[i].NumberOfVehicles
						sumGORNotOurAreaVehicles += auctions[i].GOR * auctions[i].NumberOfVehicles
						sumVITNotOurAreaVehicles += auctions[i].VIT * auctions[i].NumberOfVehicles
						sumZELNotOurAreaVehicles += auctions[i].ZEL * auctions[i].NumberOfVehicles
					}
					if auctions[i].WeWon && auctions[i].AuctionIsFinished {
						allInitialPricesWhereWeWon += auctions[i].InitialPrice
						fallDownWhereWeWon += auctions[i].InitialPrice - fallRub

						sumOurFallRubEconomy += fallRub
						// sumOurPlanFallLimitEconomy += planFallLimitRub
						sumOurPlanFallLimitEconomy += fallLimitRubForEconomy

						// fallLimitRub

					} else {

					}
					var branch string
					if auctions[i].GOR == 1 {
						branch += "ГОР"
					}
					if auctions[i].VIT == 1 {
						if branch == "" {
							branch += "ВИТ"
						} else {
							branch += ";ВИТ"
						}
					}
					if auctions[i].ZEL == 1 {
						if branch == "" {
							branch += "ЗЕЛ"
						} else {
							branch += ";ЗЕЛ"
						}
					}

					// auctions[i].InitialPrice - (auctions[i].InitialPrice * PlanFallLimit)

					result = append(result, map[string]interface{}{
						"Title":                     auctions[i].AuctionID,
						"InitialPrice":              addSpacesToPrice(fmt.Sprintf("%.2f", auctions[i].InitialPrice)),
						"HalfPercentOfInitialPrice": addSpacesToPrice(fmt.Sprintf("%.2f", auctions[i].HalfPercentOfInitialPrice)),
						"Status":                    auctions[i].Status,
						"StatusBool":                status,
						"OurAuction":                auctions[i].OurAuction,
						"PersonResponsibleForOur":   auctions[i].PersonResponsibleForOur,
						"FallLimit":                 auctions[i].FallLimit,
						"FallLimitRub":              addSpacesToPrice(fmt.Sprintf("%.2f", fallLimitRub)),
						"PlanFallLimit":             auctions[i].PlanFallLimit,
						"PlanFallLimitRub":          addSpacesToPrice(fmt.Sprintf("%.2f", planFallLimitRub)),

						"EconomyPercent": fmt.Sprintf("%.2f", auctions[i].FallLimit-fall),
						// "EconomyPercent":                      fmt.Sprintf("%.2f", auctions[i].PlanFallLimit-fall),
						// "EconomyRub":                          addSpacesToPrice(fmt.Sprintf("%.2f", auctions[i].InitialPrice*(auctions[i].PlanFallLimit-fall)/100)),
						"EconomyRub": addSpacesToPrice(fmt.Sprintf("%.2f", auctions[i].InitialPrice*(auctions[i].FallLimit-fall)/100)),

						"AdditionalFallPercent":               auctions[i].FallLimit - auctions[i].PlanFallLimit,
						"AdditionalFallRub":                   addSpacesToPrice(fmt.Sprintf("%.2f", additionalFall)),
						"PersonResponsibleForLimit":           auctions[i].PersonResponsibleForLimit,
						"OurFall":                             auctions[i].OurFall,
						"Fall":                                auctions[i].Fall,
						"FallRub":                             addSpacesToPrice(fmt.Sprintf("%.2f", fallRub)),
						"fallRubPlus05":                       addSpacesToPrice(fmt.Sprintf("%.2f", fallRubPlus05)),
						"NMCKMinusFallRub":                    addSpacesToPrice(fmt.Sprintf("%.2f", auctions[i].InitialPrice-fallRub)),
						"NMCKMinusFallRubPlus05":              addSpacesToPrice(fmt.Sprintf("%.2f", auctions[i].InitialPrice-fallRubPlus05)),
						"DoWeHaveInformationAboutCurrentFall": doWeHaveInformationAboutCurrentFall,
						"PulseTR":                             PulseTR,
						"RedTR":                               RedTR,
						"MinutesToNextStepOverLimit":          auctions[i].MinutesToNextStepOverLimit,
						"TimeWhenWeGetMinutesBeforeLimit":     auctions[i].TimeWhenWeGetMinutesBeforeLimit.Unix(),
						"ShowTimer":                           showTimer,
						"AreWeInvolved":                       auctions[i].AreWeInvolved,
						"PersonResponsibleForInvolve":         auctions[i].PersonResponsibleForInvolve,
						"LotNumber":                           auctions[i].LotNumber,
						"Priority":                            auctions[i].Priority,
						"WeWon":                               auctions[i].WeWon,
						"AuctionIsFinished":                   auctions[i].AuctionIsFinished,
						"NumberOfVehiclesBK":                  auctions[i].NumberOfVehiclesBK,
						"NumberOfVehiclesSK1":                 auctions[i].NumberOfVehiclesSK1,
						"NumberOfVehiclesSK2":                 auctions[i].NumberOfVehiclesSK2,
						"NumberOfVehicles":                    auctions[i].NumberOfVehicles,
						"GOR":                                 auctions[i].GOR,
						"VIT":                                 auctions[i].VIT,
						"ZEL":                                 auctions[i].ZEL,
						"Branch":                              branch,
						"ColorOfFalls":                        colorOfFalls,
						"Vers":                                vers,
					})

					auctionNumber += 1
					// if RedTR {
					// 	thereWasRedTRInASer = true

					// }

					bootstrap4olors := []string{"primary", "secondary", "success", "danger", "warning", "info", "light"}

					if auctionNumber == 8 || (i+1) == lenOfAuctions {

						set := models.Set{}
						uadmin.Filter(&set, "name = ?", tableNumber)
						data = append(data, AnaliticsPanelPageData{
							TableNumber: tableNumber,
							Table:       result,
							// BtnColor:                  thereWasRedTRInASer,
							BtnColor: bootstrap4olors[tableNumber-1],

							ResponsibleFromGroupOne:   set.ResponsibleFromGroupOne,
							ResponsibleFromGroupTwo:   set.ResponsibleFromGroupTwo,
							ResponsibleFromGroupThree: set.ResponsibleFromGroupThree,
						})
						tableNumber += 1
						auctionNumber = 0
						result = []map[string]interface{}{}
						// thereWasRedTRInASer = false
					}

					// fmt.Println("auctionNumber", auctionNumber)
				}

				sumData = SummaryAnaliticsPanelPageData{

					SumOfOurInitialPrices:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOfOurInitialPrices)),
					SumOfNotOurInitialPrices: addSpacesToPrice(fmt.Sprintf("%.2f", sumOfNotOurInitialPrices)),
					SummOfAllInitialPrices:   addSpacesToPrice(fmt.Sprintf("%.2f", sumOfAllInitialPrices)),

					SumOurFallRub:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOfOurInitialPrices-sumOurFallRub)),
					SumNotOurFallRub: addSpacesToPrice(fmt.Sprintf("%.2f", sumOfNotOurInitialPrices-sumNotOurFallRub)),
					SumAllFallRub:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOfAllInitialPrices-sumAllFallRub)),

					SumOurFallPercent:    addSpacesToPrice(fmt.Sprintf("%.2f", 100-sumOurFallRub*100/sumOfOurInitialPrices)),
					SumNotOurFallPercent: addSpacesToPrice(fmt.Sprintf("%.2f", 100-sumNotOurFallRub*100/sumOfNotOurInitialPrices)),
					SumAllFallPercent:    addSpacesToPrice(fmt.Sprintf("%.2f", 100-sumAllFallRub*100/sumOfAllInitialPrices)),

					SumOurPlanFallLimit:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOurPlanFallLimit)),
					SumNotOurPlanFallLimit: addSpacesToPrice(fmt.Sprintf("%.2f", sumNotOurPlanFallLimit)),
					SumAllPlanFallLimit:    addSpacesToPrice(fmt.Sprintf("%.2f", sumAllPlanFallLimit)),

					SumOurAdditionalFall:    addSpacesToPrice(fmt.Sprintf("%.2f", -1*(sumOurPlanFallLimit-sumOurFallRub))),
					SumNotOurAdditionalFall: addSpacesToPrice(fmt.Sprintf("%.2f", -1*(sumNotOurPlanFallLimit-sumNotOurFallRub))),
					SumAllAdditionalFall:    addSpacesToPrice(fmt.Sprintf("%.2f", -1*(sumAllPlanFallLimit-sumAllFallRub))),

					// SumOurPlanFallLimit:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOurPlanFallLimit)),
					// SumNotOurPlanFallLimit: addSpacesToPrice(fmt.Sprintf("%.2f", sumNotOurPlanFallLimit)),
					// SumAllPlanFallLimit:    addSpacesToPrice(fmt.Sprintf("%.2f", sumAllPlanFallLimit)),

					// SumOurPlanFallLimitPercent:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOurFallRub*100/sumOurAdditionalFall)),
					// SumNotOurPlanFallLimitPercent: addSpacesToPrice(fmt.Sprintf("%.2f", sumNotOurFallRub*100/sumNotOurAdditionalFall)),
					// SumAllPlanFallLimitPercent:    addSpacesToPrice(fmt.Sprintf("%.2f", sumAllFallRub*100/sumAllPlanFallLimit)),

					SumOurPlanFallLimitPercent:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOurPlanFallLimit*100/sumOfOurInitialPrices-sumOurFallRub*100/sumOfOurInitialPrices)),
					SumNotOurPlanFallLimitPercent: addSpacesToPrice(fmt.Sprintf("%.2f", sumNotOurPlanFallLimit*100/sumOfNotOurInitialPrices-sumNotOurFallRub*100/sumOfNotOurInitialPrices)),
					SumAllPlanFallLimitPercent:    addSpacesToPrice(fmt.Sprintf("%.2f", sumAllPlanFallLimit*100/sumOfAllInitialPrices-sumAllFallRub*100/sumOfAllInitialPrices)),

					// Economy: addSpacesToPrice(fmt.Sprintf("%.2f", economy)),
					Economy:            addSpacesToPrice(fmt.Sprintf("%.2f", sumOurFallRubEconomy-sumOurPlanFallLimitEconomy)),
					AwerageFallEconomy: addSpacesToPrice(fmt.Sprintf("%.2f", fallDownWhereWeWon/allInitialPricesWhereWeWon*100)),

					OurAverageStartLimit:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOfOurInitialPrices-sumOurPlanFallLimit)),
					NotOurAverageStartLimit: addSpacesToPrice(fmt.Sprintf("%.2f", sumOfNotOurInitialPrices-sumNotOurPlanFallLimit)),
					AllAverageStartLimit:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOfAllInitialPrices-sumAllPlanFallLimit)),

					OurAverageFall:      addSpacesToPrice(fmt.Sprintf("%.2f", sumOfOurInitialPrices-ourFallLimits)),
					NotOurAverageFall:   addSpacesToPrice(fmt.Sprintf("%.2f", sumOfNotOurInitialPrices-notOurFallLimits)),
					AllAverageStartFall: addSpacesToPrice(fmt.Sprintf("%.2f", sumOfAllInitialPrices-allFallLimits)),

					OurDifAverageFall:      addSpacesToPrice(fmt.Sprintf("%.2f", (sumOfOurInitialPrices-sumOurPlanFallLimit)-(sumOfOurInitialPrices-ourFallLimits))),
					NotOurDifAverageFall:   addSpacesToPrice(fmt.Sprintf("%.2f", (sumOfNotOurInitialPrices-sumNotOurPlanFallLimit)-(sumOfNotOurInitialPrices-notOurFallLimits))),
					AllDifAverageStartFall: addSpacesToPrice(fmt.Sprintf("%.2f", (sumOfAllInitialPrices-sumAllPlanFallLimit)-(sumOfAllInitialPrices-allFallLimits))),

					// ALLAverageStartLimit:
					// sumAllAreaVehicles

					OurVehiclesGor:    fmt.Sprintf("%.2f", ((sumGOROurAreaVehicles * 100) / sumOurAreaVehicles)),
					NotOurVehiclesGor: fmt.Sprintf("%.2f", ((sumGORNotOurAreaVehicles * 100) / sumNotOurAreaVehicles)),
					GORVehiclesGor:    fmt.Sprintf("%.2f", ((sumGORAllAreaVehicles * 100) / sumAllAreaVehicles)),

					OurVehiclesVIT:    fmt.Sprintf("%.2f", ((sumVITOurAreaVehicles * 100) / sumOurAreaVehicles)),
					NotOurVehiclesVIT: fmt.Sprintf("%.2f", ((sumVITNotOurAreaVehicles * 100) / sumNotOurAreaVehicles)),
					GORVehiclesVIT:    fmt.Sprintf("%.2f", ((sumVITAllAreaVehicles * 100) / sumAllAreaVehicles)),

					OurVehiclesZEL:    fmt.Sprintf("%.2f", ((sumZELOurAreaVehicles * 100) / sumOurAreaVehicles)),
					NotOurVehiclesZEL: fmt.Sprintf("%.2f", ((sumZELNotOurAreaVehicles * 100) / sumNotOurAreaVehicles)),
					GORVehiclesZEL:    fmt.Sprintf("%.2f", ((sumZELAllAreaVehicles * 100) / sumAllAreaVehicles)),

					// GORVehiclesALL:    fmt.Sprintf("%.2f", ((sumGORAllAreaVehicles * 100) / sumAllAreaVehicles)),
					OurVehiclesALL:    fmt.Sprintf("%.2f", ((sumOurAreaVehicles * 100) / sumAllAreaVehicles)),
					NotOurVehiclesALL: fmt.Sprintf("%.2f", ((sumNotOurAreaVehicles * 100) / sumAllAreaVehicles)),
					AllVehiclesALL:    fmt.Sprintf("%.2f", ((sumAllAreaVehicles * 100) / sumAllAreaVehicles)),

					// sumOurAreaVehicles += auctions[i].GOR*auctions[i].NumberOfVehicles + auctions[i].VIT*auctions[i].NumberOfVehicles + auctions[i].ZEL*auctions[i].NumberOfVehicles
					// sumGOROurAreaVehicles += auctions[i].GOR * auctions[i].NumberOfVehicles
					// sumVITOurAreaVehicles += auctions[i].VIT * auctions[i].NumberOfVehicles
					// sumZELOurAreaVehicles += auctions[i].ZEL * auctions[i].NumberOfVehicles

					// sumNotOurAreaVehicles += auctions[i].GOR*auctions[i].NumberOfVehicles + auctions[i].VIT*auctions[i].NumberOfVehicles + auctions[i].ZEL*auctions[i].NumberOfVehicles
					// 	sumGORNotOurAreaVehicles += auctions[i].GOR * auctions[i].NumberOfVehicles
					// 	sumVITNotOurAreaVehicles += auctions[i].VIT * auctions[i].NumberOfVehicles
					// 	sumZELNotOurAreaVehicles += auctions[i].ZEL * auctions[i].NumberOfVehicles

					// sumAllAreaVehicles += auctions[i].GOR*auctions[i].NumberOfVehicles + auctions[i].VIT*auctions[i].NumberOfVehicles + auctions[i].ZEL*auctions[i].NumberOfVehicles
					// sumGORAllAreaVehicles += auctions[i].GOR * auctions[i].NumberOfVehicles
					// sumVITAllAreaVehicles += auctions[i].VIT * auctions[i].NumberOfVehicles
					// sumZELAllAreaVehicles += auctions[i].ZEL * auctions[i].NumberOfVehicles

					Data: data,
				}
				// data := TodoPageData{
				// 	Results: results,
				// 	Todos:     results,
				// }
				// Pass TodoList data object to the specified HTML path

				uadmin.ReturnJSON(w, r, sumData)
			} else {
				res["status"] = "ERROR"
				res["err_msg"] = "Fall " + FallLimit + " is same as it was before."
				uadmin.ReturnJSON(w, r, res)
				return
			}
		case "ChangeOurNotOur":
			User := r.FormValue("User")
			AuctionID := r.FormValue("AuctionID")
			OurFromAjax := r.FormValue("Our")
			auction := models.Auction{}
			uadmin.Filter(&auction, "auction_id = ?", AuctionID)
			var Our bool
			if OurFromAjax == "1" {
				Our = true
			}
			auction.OurAuction = Our
			// auction.AuthorOfInvolved = User
			auction.PersonResponsibleForOur = User
			uadmin.Save(&auction)
			fmt.Fprintf(w, "Hello ChangeOurNotOur!")

		case "pageUpdate":
			// res := map[string]interface{}{}

			// FallLimit := r.FormValue("FallLimit")

			auctions := []models.Auction{}
			uadmin.All(&auctions)

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
			var sumOfOurInitialPrices float64
			var sumOfNotOurInitialPrices float64

			var sumOurFallRub float64
			var sumNotOurFallRub float64
			var sumAllFallRub float64

			var sumOurAdditionalFall float64
			var sumNotOurAdditionalFall float64
			var sumAllAdditionalFall float64

			var sumOurPlanFallLimit float64
			var sumNotOurPlanFallLimit float64
			var sumAllPlanFallLimit float64

			var sumOurAreaVehicles float64
			var sumNotOurAreaVehicles float64
			var sumAllAreaVehicles float64

			var sumGOROurAreaVehicles float64
			var sumGORNotOurAreaVehicles float64
			var sumGORAllAreaVehicles float64

			var sumVITOurAreaVehicles float64
			var sumVITNotOurAreaVehicles float64
			var sumVITAllAreaVehicles float64

			var sumZELOurAreaVehicles float64
			var sumZELNotOurAreaVehicles float64
			var sumZELAllAreaVehicles float64

			var economy float64
			var sumOurFallRubEconomy float64

			var sumOurPlanFallLimitEconomy float64
			var allInitialPricesWhereWeWon float64
			var fallDownWhereWeWon float64

			var ourFallLimits float64
			var allOurStartFalls float64
			var notOurFallLimits float64
			var allFallLimits float64

			// var thereWasRedTRInASer bool
			timeLayout := "15:04:05"
			for i := range auctions {
				uadmin.Preload(&auctions[i])

				ver := []models.AuctionVersion{}
				uadmin.Filter(&ver, "auction_id = ?", auctions[i].ID)
				vers := []map[string]interface{}{}
				// for j := range ver {
				for j := len(ver) - 1; j >= 0; j-- {
					vers = append(vers, map[string]interface{}{
						"FallLimit":                 ver[j].FallLimit,
						"PersonResponsibleForLimit": ver[j].PersonResponsibleForLimit,
						"Date":                      ver[j].Date.Format(timeLayout),
					})
				}
				sumOfAllInitialPrices += auctions[i].InitialPrice

				var status bool
				if auctions[i].Status == "Прием заявок" {
					status = true
				}

				var PulseTR string
				var RedTR bool
				if (auctions[i].Fall >= auctions[i].FallLimit || auctions[i].OurFall >= auctions[i].FallLimit) &&
					(auctions[i].Fall > 0 || auctions[i].OurFall > 0) &&
					(!auctions[i].WeWon && !auctions[i].AuctionIsFinished) {
					PulseTR = "pulse"
					RedTR = true
				}

				var showTimer bool
				if auctions[i].MinutesToNextStepOverLimit == 0 {
					showTimer = true
				}
				// Assigns the string of interface in each Todo fields

				var fall float64

				if auctions[i].OurFall > auctions[i].Fall {
					fall = auctions[i].OurFall
				} else {
					fall = auctions[i].Fall
				}

				// var fallLimit float64
				// if auctions[i].WeWon && auctions[i].AuctionIsFinished {
				// 	fallLimit = fall
				// } else {
				// 	fallLimit = auctions[i].FallLimit
				// }

				var fallRub float64
				fallRub = auctions[i].InitialPrice - fall*auctions[i].InitialPrice/100
				var fallRubPlus05 float64
				fallRubPlus05 = auctions[i].InitialPrice - (fall+0.5)*auctions[i].InitialPrice/100
				var fallLimitRubForEconomy float64
				fallLimitRubForEconomy = (auctions[i].InitialPrice - auctions[i].FallLimit*auctions[i].InitialPrice/100)
				// var currentFall float64
				var fallLimitRub float64
				if auctions[i].WeWon && auctions[i].AuctionIsFinished {
					fallLimitRub = fallRub
					// currentFall = fall
				} else {
					fallLimitRub = (auctions[i].InitialPrice - auctions[i].FallLimit*auctions[i].InitialPrice/100)
					// currentFall = auctions[i].FallLimit
				}
				var planFallLimitRub float64
				planFallLimitRub = (auctions[i].InitialPrice - auctions[i].PlanFallLimit*auctions[i].InitialPrice/100)
				var additionalFall float64
				additionalFall = fallLimitRub - planFallLimitRub
				var doWeHaveInformationAboutCurrentFall bool
				if auctions[i].OurFall != 0 || auctions[i].Fall != 0 {
					doWeHaveInformationAboutCurrentFall = true
				}

				var colorOfFalls string
				if (auctions[i].OurFall < auctions[i].Fall && fall >= auctions[i].FallLimit) &&
					(!auctions[i].WeWon && !auctions[i].AuctionIsFinished) {
					colorOfFalls = "bg-danger text-light"
				} else {
					colorOfFalls = "bg-success text-light"

				}

				if auctions[i].WeWon {
					economy += auctions[i].InitialPrice - auctions[i].InitialPrice*fall
				}

				// var AreWeInvolved bool

				// if (auctions[i].OurFall > 0 || auctions[i].Fall > 0) {
				// 	AreWeInvolved = true
				// }
				sumAllFallRub += fallRub
				sumAllAdditionalFall += additionalFall
				sumAllPlanFallLimit += planFallLimitRub

				sumAllAreaVehicles += auctions[i].GOR*auctions[i].NumberOfVehicles + auctions[i].VIT*auctions[i].NumberOfVehicles + auctions[i].ZEL*auctions[i].NumberOfVehicles
				sumGORAllAreaVehicles += auctions[i].GOR * auctions[i].NumberOfVehicles
				sumVITAllAreaVehicles += auctions[i].VIT * auctions[i].NumberOfVehicles
				sumZELAllAreaVehicles += auctions[i].ZEL * auctions[i].NumberOfVehicles
				allFallLimits += fallLimitRub
				if auctions[i].OurAuction {
					sumOfOurInitialPrices += auctions[i].InitialPrice
					sumOurFallRub += fallRub
					sumOurAdditionalFall += additionalFall
					sumOurPlanFallLimit += planFallLimitRub

					ourFallLimits += fallLimitRub
					// if auctions[i].WeWon && auctions[i].AuctionIsFinished {
					// 	ourFallLimits +=
					// }else{
					// 	ourFallLimits += fallLimitRub
					// }

					sumOurAreaVehicles += auctions[i].GOR*auctions[i].NumberOfVehicles + auctions[i].VIT*auctions[i].NumberOfVehicles + auctions[i].ZEL*auctions[i].NumberOfVehicles
					sumGOROurAreaVehicles += auctions[i].GOR * auctions[i].NumberOfVehicles
					sumVITOurAreaVehicles += auctions[i].VIT * auctions[i].NumberOfVehicles
					sumZELOurAreaVehicles += auctions[i].ZEL * auctions[i].NumberOfVehicles
					allOurStartFalls += auctions[i].InitialPrice - (auctions[i].InitialPrice*auctions[i].PlanFallLimit)/100
					// allOurStartFalls += PlanFallLimitRub
				} else {
					sumOfNotOurInitialPrices += auctions[i].InitialPrice
					sumNotOurFallRub += fallRub
					sumNotOurAdditionalFall += additionalFall
					sumNotOurPlanFallLimit += planFallLimitRub

					notOurFallLimits += fallLimitRub
					// if auctions[i].WeWon && auctions[i].AuctionIsFinished {
					// 	notOurFallLimits +=
					// }else{
					// 	notOurFallLimits += fallLimitRub
					// }

					sumNotOurAreaVehicles += auctions[i].GOR*auctions[i].NumberOfVehicles + auctions[i].VIT*auctions[i].NumberOfVehicles + auctions[i].ZEL*auctions[i].NumberOfVehicles
					sumGORNotOurAreaVehicles += auctions[i].GOR * auctions[i].NumberOfVehicles
					sumVITNotOurAreaVehicles += auctions[i].VIT * auctions[i].NumberOfVehicles
					sumZELNotOurAreaVehicles += auctions[i].ZEL * auctions[i].NumberOfVehicles
				}
				if auctions[i].WeWon && auctions[i].AuctionIsFinished {
					allInitialPricesWhereWeWon += auctions[i].InitialPrice
					fallDownWhereWeWon += auctions[i].InitialPrice - fallRub

					sumOurFallRubEconomy += fallRub
					// sumOurPlanFallLimitEconomy += planFallLimitRub
					sumOurPlanFallLimitEconomy += fallLimitRubForEconomy

					// fallLimitRub

				} else {

				}
				var branch string
				if auctions[i].GOR == 1 {
					branch += "ГОР"
				}
				if auctions[i].VIT == 1 {
					if branch == "" {
						branch += "ВИТ"
					} else {
						branch += ";ВИТ"
					}
				}
				if auctions[i].ZEL == 1 {
					if branch == "" {
						branch += "ЗЕЛ"
					} else {
						branch += ";ЗЕЛ"
					}
				}

				// auctions[i].InitialPrice - (auctions[i].InitialPrice * PlanFallLimit)

				result = append(result, map[string]interface{}{
					"Title":                     auctions[i].AuctionID,
					"InitialPrice":              addSpacesToPrice(fmt.Sprintf("%.2f", auctions[i].InitialPrice)),
					"HalfPercentOfInitialPrice": addSpacesToPrice(fmt.Sprintf("%.2f", auctions[i].HalfPercentOfInitialPrice)),
					"Status":                    auctions[i].Status,
					"StatusBool":                status,
					"OurAuction":                auctions[i].OurAuction,
					"PersonResponsibleForOur":   auctions[i].PersonResponsibleForOur,
					"FallLimit":                 auctions[i].FallLimit,
					"FallLimitRub":              addSpacesToPrice(fmt.Sprintf("%.2f", fallLimitRub)),
					"PlanFallLimit":             auctions[i].PlanFallLimit,
					"PlanFallLimitRub":          addSpacesToPrice(fmt.Sprintf("%.2f", planFallLimitRub)),

					"EconomyPercent": fmt.Sprintf("%.2f", auctions[i].FallLimit-fall),
					// "EconomyPercent":                      fmt.Sprintf("%.2f", auctions[i].PlanFallLimit-fall),
					// "EconomyRub":                          addSpacesToPrice(fmt.Sprintf("%.2f", auctions[i].InitialPrice*(auctions[i].PlanFallLimit-fall)/100)),
					"EconomyRub": addSpacesToPrice(fmt.Sprintf("%.2f", auctions[i].InitialPrice*(auctions[i].FallLimit-fall)/100)),

					"AdditionalFallPercent":               auctions[i].FallLimit - auctions[i].PlanFallLimit,
					"AdditionalFallRub":                   addSpacesToPrice(fmt.Sprintf("%.2f", additionalFall)),
					"PersonResponsibleForLimit":           auctions[i].PersonResponsibleForLimit,
					"OurFall":                             auctions[i].OurFall,
					"Fall":                                auctions[i].Fall,
					"FallRub":                             addSpacesToPrice(fmt.Sprintf("%.2f", fallRub)),
					"fallRubPlus05":                       addSpacesToPrice(fmt.Sprintf("%.2f", fallRubPlus05)),
					"NMCKMinusFallRub":                    addSpacesToPrice(fmt.Sprintf("%.2f", auctions[i].InitialPrice-fallRub)),
					"NMCKMinusFallRubPlus05":              addSpacesToPrice(fmt.Sprintf("%.2f", auctions[i].InitialPrice-fallRubPlus05)),
					"DoWeHaveInformationAboutCurrentFall": doWeHaveInformationAboutCurrentFall,
					"PulseTR":                             PulseTR,
					"RedTR":                               RedTR,
					"MinutesToNextStepOverLimit":          auctions[i].MinutesToNextStepOverLimit,
					"TimeWhenWeGetMinutesBeforeLimit":     auctions[i].TimeWhenWeGetMinutesBeforeLimit.Unix(),
					"ShowTimer":                           showTimer,
					"AreWeInvolved":                       auctions[i].AreWeInvolved,
					"PersonResponsibleForInvolve":         auctions[i].PersonResponsibleForInvolve,
					"LotNumber":                           auctions[i].LotNumber,
					"Priority":                            auctions[i].Priority,
					"WeWon":                               auctions[i].WeWon,
					"AuctionIsFinished":                   auctions[i].AuctionIsFinished,
					"NumberOfVehiclesBK":                  auctions[i].NumberOfVehiclesBK,
					"NumberOfVehiclesSK1":                 auctions[i].NumberOfVehiclesSK1,
					"NumberOfVehiclesSK2":                 auctions[i].NumberOfVehiclesSK2,
					"NumberOfVehicles":                    auctions[i].NumberOfVehicles,
					"GOR":                                 auctions[i].GOR,
					"VIT":                                 auctions[i].VIT,
					"ZEL":                                 auctions[i].ZEL,
					"Branch":                              branch,
					"ColorOfFalls":                        colorOfFalls,
					"Vers":                                vers,
				})

				auctionNumber += 1
				// if RedTR {
				// 	thereWasRedTRInASer = true

				// }

				bootstrap4olors := []string{"primary", "secondary", "success", "danger", "warning", "info", "light"}

				if auctionNumber == 8 || (i+1) == lenOfAuctions {

					set := models.Set{}
					uadmin.Filter(&set, "name = ?", tableNumber)
					data = append(data, AnaliticsPanelPageData{
						TableNumber: tableNumber,
						Table:       result,
						// BtnColor:                  thereWasRedTRInASer,
						BtnColor: bootstrap4olors[tableNumber-1],

						ResponsibleFromGroupOne:   set.ResponsibleFromGroupOne,
						ResponsibleFromGroupTwo:   set.ResponsibleFromGroupTwo,
						ResponsibleFromGroupThree: set.ResponsibleFromGroupThree,
					})
					tableNumber += 1
					auctionNumber = 0
					result = []map[string]interface{}{}
					// thereWasRedTRInASer = false
				}

				// fmt.Println("auctionNumber", auctionNumber)
			}

			sumData = SummaryAnaliticsPanelPageData{

				SumOfOurInitialPrices:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOfOurInitialPrices)),
				SumOfNotOurInitialPrices: addSpacesToPrice(fmt.Sprintf("%.2f", sumOfNotOurInitialPrices)),
				SummOfAllInitialPrices:   addSpacesToPrice(fmt.Sprintf("%.2f", sumOfAllInitialPrices)),

				SumOurFallRub:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOfOurInitialPrices-sumOurFallRub)),
				SumNotOurFallRub: addSpacesToPrice(fmt.Sprintf("%.2f", sumOfNotOurInitialPrices-sumNotOurFallRub)),
				SumAllFallRub:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOfAllInitialPrices-sumAllFallRub)),

				SumOurFallPercent:    addSpacesToPrice(fmt.Sprintf("%.2f", 100-sumOurFallRub*100/sumOfOurInitialPrices)),
				SumNotOurFallPercent: addSpacesToPrice(fmt.Sprintf("%.2f", 100-sumNotOurFallRub*100/sumOfNotOurInitialPrices)),
				SumAllFallPercent:    addSpacesToPrice(fmt.Sprintf("%.2f", 100-sumAllFallRub*100/sumOfAllInitialPrices)),

				SumOurPlanFallLimit:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOurPlanFallLimit)),
				SumNotOurPlanFallLimit: addSpacesToPrice(fmt.Sprintf("%.2f", sumNotOurPlanFallLimit)),
				SumAllPlanFallLimit:    addSpacesToPrice(fmt.Sprintf("%.2f", sumAllPlanFallLimit)),

				SumOurAdditionalFall:    addSpacesToPrice(fmt.Sprintf("%.2f", -1*(sumOurPlanFallLimit-sumOurFallRub))),
				SumNotOurAdditionalFall: addSpacesToPrice(fmt.Sprintf("%.2f", -1*(sumNotOurPlanFallLimit-sumNotOurFallRub))),
				SumAllAdditionalFall:    addSpacesToPrice(fmt.Sprintf("%.2f", -1*(sumAllPlanFallLimit-sumAllFallRub))),

				// SumOurPlanFallLimit:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOurPlanFallLimit)),
				// SumNotOurPlanFallLimit: addSpacesToPrice(fmt.Sprintf("%.2f", sumNotOurPlanFallLimit)),
				// SumAllPlanFallLimit:    addSpacesToPrice(fmt.Sprintf("%.2f", sumAllPlanFallLimit)),

				// SumOurPlanFallLimitPercent:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOurFallRub*100/sumOurAdditionalFall)),
				// SumNotOurPlanFallLimitPercent: addSpacesToPrice(fmt.Sprintf("%.2f", sumNotOurFallRub*100/sumNotOurAdditionalFall)),
				// SumAllPlanFallLimitPercent:    addSpacesToPrice(fmt.Sprintf("%.2f", sumAllFallRub*100/sumAllPlanFallLimit)),

				SumOurPlanFallLimitPercent:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOurPlanFallLimit*100/sumOfOurInitialPrices-sumOurFallRub*100/sumOfOurInitialPrices)),
				SumNotOurPlanFallLimitPercent: addSpacesToPrice(fmt.Sprintf("%.2f", sumNotOurPlanFallLimit*100/sumOfNotOurInitialPrices-sumNotOurFallRub*100/sumOfNotOurInitialPrices)),
				SumAllPlanFallLimitPercent:    addSpacesToPrice(fmt.Sprintf("%.2f", sumAllPlanFallLimit*100/sumOfAllInitialPrices-sumAllFallRub*100/sumOfAllInitialPrices)),

				// Economy: addSpacesToPrice(fmt.Sprintf("%.2f", economy)),
				Economy:            addSpacesToPrice(fmt.Sprintf("%.2f", sumOurFallRubEconomy-sumOurPlanFallLimitEconomy)),
				AwerageFallEconomy: addSpacesToPrice(fmt.Sprintf("%.2f", fallDownWhereWeWon/allInitialPricesWhereWeWon*100)),

				OurAverageStartLimit:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOfOurInitialPrices-sumOurPlanFallLimit)),
				NotOurAverageStartLimit: addSpacesToPrice(fmt.Sprintf("%.2f", sumOfNotOurInitialPrices-sumNotOurPlanFallLimit)),
				AllAverageStartLimit:    addSpacesToPrice(fmt.Sprintf("%.2f", sumOfAllInitialPrices-sumAllPlanFallLimit)),

				OurAverageFall:      addSpacesToPrice(fmt.Sprintf("%.2f", sumOfOurInitialPrices-ourFallLimits)),
				NotOurAverageFall:   addSpacesToPrice(fmt.Sprintf("%.2f", sumOfNotOurInitialPrices-notOurFallLimits)),
				AllAverageStartFall: addSpacesToPrice(fmt.Sprintf("%.2f", sumOfAllInitialPrices-allFallLimits)),

				OurDifAverageFall:      addSpacesToPrice(fmt.Sprintf("%.2f", (sumOfOurInitialPrices-sumOurPlanFallLimit)-(sumOfOurInitialPrices-ourFallLimits))),
				NotOurDifAverageFall:   addSpacesToPrice(fmt.Sprintf("%.2f", (sumOfNotOurInitialPrices-sumNotOurPlanFallLimit)-(sumOfNotOurInitialPrices-notOurFallLimits))),
				AllDifAverageStartFall: addSpacesToPrice(fmt.Sprintf("%.2f", (sumOfAllInitialPrices-sumAllPlanFallLimit)-(sumOfAllInitialPrices-allFallLimits))),

				// ALLAverageStartLimit:
				// sumAllAreaVehicles

				OurVehiclesGor:    fmt.Sprintf("%.2f", ((sumGOROurAreaVehicles * 100) / sumOurAreaVehicles)),
				NotOurVehiclesGor: fmt.Sprintf("%.2f", ((sumGORNotOurAreaVehicles * 100) / sumNotOurAreaVehicles)),
				GORVehiclesGor:    fmt.Sprintf("%.2f", ((sumGORAllAreaVehicles * 100) / sumAllAreaVehicles)),

				OurVehiclesVIT:    fmt.Sprintf("%.2f", ((sumVITOurAreaVehicles * 100) / sumOurAreaVehicles)),
				NotOurVehiclesVIT: fmt.Sprintf("%.2f", ((sumVITNotOurAreaVehicles * 100) / sumNotOurAreaVehicles)),
				GORVehiclesVIT:    fmt.Sprintf("%.2f", ((sumVITAllAreaVehicles * 100) / sumAllAreaVehicles)),

				OurVehiclesZEL:    fmt.Sprintf("%.2f", ((sumZELOurAreaVehicles * 100) / sumOurAreaVehicles)),
				NotOurVehiclesZEL: fmt.Sprintf("%.2f", ((sumZELNotOurAreaVehicles * 100) / sumNotOurAreaVehicles)),
				GORVehiclesZEL:    fmt.Sprintf("%.2f", ((sumZELAllAreaVehicles * 100) / sumAllAreaVehicles)),

				// GORVehiclesALL:    fmt.Sprintf("%.2f", ((sumGORAllAreaVehicles * 100) / sumAllAreaVehicles)),
				OurVehiclesALL:    fmt.Sprintf("%.2f", ((sumOurAreaVehicles * 100) / sumAllAreaVehicles)),
				NotOurVehiclesALL: fmt.Sprintf("%.2f", ((sumNotOurAreaVehicles * 100) / sumAllAreaVehicles)),
				AllVehiclesALL:    fmt.Sprintf("%.2f", ((sumAllAreaVehicles * 100) / sumAllAreaVehicles)),

				// sumOurAreaVehicles += auctions[i].GOR*auctions[i].NumberOfVehicles + auctions[i].VIT*auctions[i].NumberOfVehicles + auctions[i].ZEL*auctions[i].NumberOfVehicles
				// sumGOROurAreaVehicles += auctions[i].GOR * auctions[i].NumberOfVehicles
				// sumVITOurAreaVehicles += auctions[i].VIT * auctions[i].NumberOfVehicles
				// sumZELOurAreaVehicles += auctions[i].ZEL * auctions[i].NumberOfVehicles

				// sumNotOurAreaVehicles += auctions[i].GOR*auctions[i].NumberOfVehicles + auctions[i].VIT*auctions[i].NumberOfVehicles + auctions[i].ZEL*auctions[i].NumberOfVehicles
				// 	sumGORNotOurAreaVehicles += auctions[i].GOR * auctions[i].NumberOfVehicles
				// 	sumVITNotOurAreaVehicles += auctions[i].VIT * auctions[i].NumberOfVehicles
				// 	sumZELNotOurAreaVehicles += auctions[i].ZEL * auctions[i].NumberOfVehicles

				// sumAllAreaVehicles += auctions[i].GOR*auctions[i].NumberOfVehicles + auctions[i].VIT*auctions[i].NumberOfVehicles + auctions[i].ZEL*auctions[i].NumberOfVehicles
				// sumGORAllAreaVehicles += auctions[i].GOR * auctions[i].NumberOfVehicles
				// sumVITAllAreaVehicles += auctions[i].VIT * auctions[i].NumberOfVehicles
				// sumZELAllAreaVehicles += auctions[i].ZEL * auctions[i].NumberOfVehicles

				Data: data,
			}
			// data := TodoPageData{
			// 	Results: results,
			// 	Todos:     results,
			// }
			// Pass TodoList data object to the specified HTML path

			uadmin.ReturnJSON(w, r, sumData)

		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}
