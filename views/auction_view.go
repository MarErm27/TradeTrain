package views

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"strconv"

	// "math/big"
	"net/http"
	// "strconv"
	"strings"
	"time"

	//"github.com/MarErm27/webTrainerWithUAdmin/models"
	"github.com/MarErm27/webTrainerWithUAdmin/models"
	"github.com/uadmin/uadmin"
)

type AuctionPageData struct {
	AuctionID string
	Price     string
}

// var participants = []map[string]interface{}{}

func getCurrentPosition(auctionID uint) int {
	// numberOfBetterOffers := []models.Offer{}
	// modelschemaBetter := uadmin.Schema["offer"]
	// uadmin.FilterList(&modelschemaBetter, "bid  DESC, id", true, 0, 1, &numberOfBetterOffers, "auction_id = ? AND participant = ?", auctionID, "Акционерное Общество \"Третий парк\"")

	// Initialize the Todo model
	numberOfBetterOffers := []models.Offer{}

	// Create a query in the sql variable to select all records in todos
	sql := `SELECT * FROM offers WHERE auction_id = ? AND participant = ? ORDER BY bid ASC, id ASC LIMIT 1`
	// Place it here
	db := uadmin.GetDB()
	// db.Raw(sql, auctionID, "Акционерное Общество \"Третий парк\"")
	// Store the query inside the Raw function in order to scan value to
	// the Todo model
	db.Raw(sql, auctionID, "Акционерное Общество \"Третий парк\"").Scan(&numberOfBetterOffers)

	if len(numberOfBetterOffers) > 0 {
		// uadmin.FilterList(&modelschemaBetter, "bid  DESC, id", true, 0, 100000, &numberOfBetterOffers, "auction_id = ? AND bid <= ?", auctionID, numberOfBetterOffers[0].Bid)
		// uadmin.Filter(&numberOfBetterOffers, "auction_id = ? AND participant = ?", auctionID, "Акционерное Общество \"Третий парк\"")

		sql = `SELECT * FROM offers WHERE auction_id = ? AND bid <= ? ORDER BY bid ASC, id ASC`
		// Place it here
		db := uadmin.GetDB()
		// db.Raw(sql, auctionID, "Акционерное Общество \"Третий парк\"")
		// Store the query inside the Raw function in order to scan value to
		// the Todo model
		db.Raw(sql, auctionID, numberOfBetterOffers[0].Bid).Scan(&numberOfBetterOffers)

		return len(numberOfBetterOffers)
	} else {
		return 0
	}
}

func round(input float64) float64 {
	if math.IsNaN(input) {
		return math.NaN()
	}
	sign := 1.0
	if input < 0 {
		sign = -1
		input *= -1
	}
	_, decimal := math.Modf(input)
	var rounded float64
	if decimal >= 0.5 {
		rounded = math.Ceil(input*100) / 100
	} else {
		rounded = math.Floor(input*100) / 100
	}
	return rounded * sign
}

func Round(input float64) float64 {
	if math.IsNaN(input) {
		return math.NaN()
	}
	sign := 1.0
	if input < 0 {
		sign = -1
		input *= -1
	}
	_, decimal := math.Modf(input)
	var rounded float64
	if decimal >= 0.5 {
		rounded = math.Ceil(input*100) / 100
	} else {
		rounded = math.Floor(input * 100)
	}
	return rounded * sign
}
func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func reduceLagging(lastBiddingStartNumber time.Duration, auctionID string) { //, wg *sync.WaitGroup) {
	// fmt.Println("statusIsGoingToBeChanged", biddingStartNumber, "auctionID: ", auctionID)
	<-time.After(lastBiddingStartNumber)

	// fmt.Println("changeAuctionStatusToActive: ", auctionID, biddingStartNumber)
	auction := models.Auction{}
	uadmin.Filter(&auction, "auction_id = ?", auctionID)

	rand.Seed(time.Now().UnixNano())
	var n int
	n = 1 + rand.Intn(2+1) // a ≤ n ≤ b
	t1 := time.Now()
	if n == 1 {
		// uadmin.Preload(&auction[0])

		lastOffer := models.Offer{}
		modelschema1 := uadmin.Schema["offer"]
		uadmin.FilterList(&modelschema1, "id", false, 0, 1, &lastOffer, "auction_id = ? AND best_bid = ?", auction.ID, true)
		var tm time.Time
		zeroTime := time.Time{}
		if lastOffer.BidTime == zeroTime {
			tm = auction.BiddingStart
		} else {
			tm = lastOffer.BidTime
		}

		rand.Seed(time.Now().UnixNano())
		randomN := 1 + rand.Intn(529-1+1)
		tmAdd := tm.Add(time.Second * time.Duration(randomN))
		diff := tmAdd.Sub(t1)
		auctionID := auction.AuctionID
		auction.HowLongOneComputerWillLag = diff
		uadmin.Save(&auction)
		// startLagAccessDenied = true

		go reduceLagging(diff, auctionID)

		// auctionTempName = "templates/noAccess.html"
		// uadmin.RenderHTML(w, r, auctionTempName, nil)

		// return
	} else {
		auction.OneComputerIsLagging = false
		uadmin.Save(&auction)
	}

}

// AuctionHandler !
func AuctionHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is not logged in
	// if uadmin.IsAuthenticated(r) == nil {
	// 	http.Redirect(w, r, "/admin", http.StatusSeeOther)
	// 	return
	// }
	// r.URL.Path creates a new path called /todo
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/Auction")

	switch r.Method {
	case "GET":
		keys, ok := r.URL.Query()["AuctionID"]

		if !ok || len(keys[0]) < 1 {
			log.Println("Url Param 'AuctionID' is missing")
			return
		}

		rand.Seed(time.Now().UnixNano())
		var auctionTempName string
		n := 1 + rand.Intn(1000-1+1) // a ≤ n ≤ b
		if n == 1 {
			auctionTempName = "templates/noInternet.html"
			uadmin.RenderHTML(w, r, auctionTempName, nil)
			return
		} else {
			auctionTempName = "templates/auction.html"
		}

		// Query()["key"] will return an array of items,
		// we only want the single item.
		key := keys[0]

		// Assigns a map as a string of interface to store any types of values
		results := []map[string]interface{}{}
		// Fetches all object in the database
		auction := []models.Auction{}

		modelschema := uadmin.Schema["auction"]
		uadmin.FilterList(&modelschema, "id", true, 0, 1, &auction, "auction_id = ?", key)

		if len(auction) == 0 {
			fmt.Fprintf(w, "Sorry, auction "+key+" not found.")
			return
		}
		var startLagAccessDenied bool
		if !auction[0].OneComputerIsLagging {
			rand.Seed(time.Now().UnixNano())
			n = 1 + rand.Intn(40+1) // a ≤ n ≤ b
			t1 := time.Now()
			if n == 5 {
				// uadmin.Preload(&auction[0])

				lastOffer := models.Offer{}
				modelschema1 := uadmin.Schema["offer"]
				uadmin.FilterList(&modelschema1, "id", false, 0, 1, &lastOffer, "auction_id = ? AND best_bid = ?", auction[0].ID, true)
				var tm time.Time
				zeroTime := time.Time{}
				if lastOffer.BidTime == zeroTime {
					tm = auction[0].BiddingStart
				} else {
					tm = lastOffer.BidTime
				}

				rand.Seed(time.Now().UnixNano())
				randomN := 1 + rand.Intn(529-1+1)
				tmAdd := tm.Add(time.Second * time.Duration(randomN))
				diff := tmAdd.Sub(t1)
				auctionID := auction[0].AuctionID
				auction[0].OneComputerIsLagging = true
				auction[0].HowLongOneComputerWillLag = diff
				uadmin.Save(&auction[0])
				startLagAccessDenied = true

				go reduceLagging(diff, auctionID)

				// auctionTempName = "templates/noAccess.html"
				// uadmin.RenderHTML(w, r, auctionTempName, nil)

				// return
			}
		}

		// Accesses and fetches data from another model

		offers := []models.Offer{}
		uadmin.Filter(&offers, "auction_id = ?", auction[0].ID)
		// uadmin.FilterList(&modelschema, "id", false, 0, 1, &offers, "auction_id = ?", auction[0].ID)
		var thereWasOffersBefore bool
		if len(offers) > 0 {
			thereWasOffersBefore = true
		}
		uadmin.Filter(&offers, "auction_id = ? AND participant = ?", auction[0].ID, "Акционерное Общество \"Третий парк\"")

		var thereWasYourOffer bool
		if len(offers) > 0 {
			thereWasYourOffer = true
		}
		uadmin.Filter(&offers, "auction_id = ?", auction[0].ID)
		offersTable := []map[string]interface{}{}
		for i := range offers {
			uadmin.Preload(&offers[i])
			offersTable = append(offersTable, map[string]interface{}{
				"Bid":         addSpacesToPrice(fmt.Sprintf("%.2f", offers[i].Bid)),
				"Participant": offers[i].Participant,
				"BidTime":     offers[i].BidTime.Unix(),
				"BestBid":     offers[i].BestBid,
			})
		}

		lastOffer := models.Offer{}
		modelschema1 := uadmin.Schema["offer"]
		uadmin.FilterList(&modelschema1, "id", false, 0, 1, &lastOffer, "auction_id = ? AND best_bid = ?", auction[0].ID, true)

		for i := range auction {
			uadmin.Preload(&auction[i])
			// Assigns the string of interface in each Todo fields
			results = append(results, map[string]interface{}{
				"AuctionID":                    auction[i].AuctionID,
				"Name":                         auction[i].Name,
				"InitialPrice":                 addSpacesToPrice(fmt.Sprintf("%.2f", auction[i].InitialPrice)),
				"BidDuration":                  time.Now().Unix() - auction[i].BiddingStart.Unix(),
				"CurrentServerTime":            time.Now().Unix(),
				"CurrentFall":                  auction[0].CurrentFall,
				"CurrentPrice":                 addSpacesToPrice(fmt.Sprintf("%.2f", auction[0].CurrentPrice)),
				"thereWasOffersBefore":         thereWasOffersBefore,
				"thereWasYourOffer":            thereWasYourOffer,
				"YourLastOffer":                auction[i].YourLastOffer,
				"YourCurrentPosition":          auction[i].YourCurrentPosition,
				"AmountOfOffersThatWerePosted": len(offers),
				"offersTable":                  offersTable,
				"lastOfferTime":                lastOffer.BidTime.Unix(),
				"BestPretendentName":           lastOffer.Participant,
				"BestPretendentBid":            lastOffer.Bid,
				"Status":                       auction[i].Status,
				"BiddingDuration":              fmtDuration(auction[i].BiddingFinish.Sub(auction[i].BiddingStart)),
				"BiddingStart":                 auction[i].BiddingStart.Unix(),
				"OneComputerIsLagging":         auction[0].OneComputerIsLagging,
				"StartLagAccessDenied":         startLagAccessDenied,
			})
			break
		}

		// Prints the todo in JSON format
		// uadmin.ReturnJSON(w, r, results)
		uadmin.RenderHTML(w, r, auctionTempName, results[0])

	case "POST":
		Method := r.FormValue("Method")
		switch Method {
		case "PlaceAnOffer":
			// var ItIsABestOffer bool
			OfferAuctionID := r.FormValue("AuctionID")
			OfferBid := r.FormValue("Offer")
			OfferParticipant := r.FormValue("Participant")
			results := []map[string]interface{}{}
			res := map[string]interface{}{}

			auction := []models.Auction{}
			uadmin.Filter(&auction, "auction_id = ?", OfferAuctionID)

			if auction[0].Status == "Подведение итогов" {
				res["status"] = "ERROR"
				res["err_msg"] = "Auction is over. Your offer is not accepted."
				uadmin.ReturnJSON(w, r, res)
				return
			}
			offers := []models.Offer{}
			modelschema := uadmin.Schema["offer"]
			uadmin.FilterList(&modelschema, "id", false, 0, 1, &offers, "auction_id = ? AND best_bid = ?", auction[0].ID, true)
			for i := range offers {
				uadmin.Preload(&offers[i])
				// Assigns the string of interface in each Todo fields
				results = append(results, map[string]interface{}{

					"LastParticipant": offers[i].Participant,
					"LastBid":         offers[i].Bid,
					"IsItABestBid":    offers[i].BestBid,
					"BidTime":         offers[i].BidTime,
				})
			}
			// Validate if the friendName variable is empty.
			if OfferAuctionID == "" {
				res["status"] = "ERROR"
				res["err_msg"] = "AuctionID is required."
				uadmin.ReturnJSON(w, r, res)
				return

			}
			var LastParticipant string

			var floatZero float64
			LastBid := floatZero
			var thisIsAFirstOffer bool
			if len(results) != 0 {
				LastParticipant, _ = results[0]["LastParticipant"].(string)
				//IsItABestBid, _ := results[0]["IsItABestBid"].(bool)
				LastBid = results[0]["LastBid"].(float64)

			} else {
				thisIsAFirstOffer = true
			}
			OfferBidFloat64, err := strconv.ParseFloat(OfferBid, 64)
			if err != nil {
				fmt.Println("err: ", err)
			}
			if LastParticipant == OfferParticipant && OfferBidFloat64 == LastBid {
				res["status"] = "ERROR"
				res["err_msg"] = "This offer has a same price as a previous one a given from a same participant. It will not be recorded/"
				uadmin.ReturnJSON(w, r, res)
				return
			} else {
				fmt.Println("test")
			}

			offer := models.Offer{}
			if OfferBidFloat64 >= (auction[0].CurrentPrice - auction[0].FivePercentsOfInitialPrice) { //
				// Store input into the Name, Email, and Password fields
				offer.Auction = auction[0]
				offer.Auction.AuctionID = OfferAuctionID
				offer.Bid = OfferBidFloat64
				offer.Participant = OfferParticipant
				if (LastBid-OfferBidFloat64) > 0 || thisIsAFirstOffer {
					// ItIsABestOffer = true
					offer.BestBid = true
					auction[0].CurrentPrice = OfferBidFloat64
					auction[0].CurrentFall = 100 - auction[0].CurrentPrice*100/auction[0].InitialPrice
					uadmin.Save(&offer)
					auction[0].YourCurrentPosition = getCurrentPosition(auction[0].ID)

				} else {
					offer.BestBid = false
					// numberOfBetterOffers := []models.Offer{}
					// modelschemaBetter := uadmin.Schema["offer"]
					// uadmin.FilterList(&modelschemaBetter, "id", false, 0, 1, &numberOfBetterOffers, "auction_id = ? AND bid <= ?", auction[0].ID, OfferBidFloat64)
					uadmin.Save(&offer)
					auction[0].YourCurrentPosition = getCurrentPosition(auction[0].ID)
				}

				offer.BidTime = time.Now()

				// Store input in the Friend model
				uadmin.Save(&offer)

			} else {
				res["status"] = "ERROR"
				res["err_msg"] = "Your bid is already best."
				uadmin.ReturnJSON(w, r, res)
				return
			}

			if OfferParticipant == "Акционерное Общество \"Третий парк\"" {
				auction[0].YourLastOffer = OfferBidFloat64
				// auction[0].YourCurrentPosition = 1
			}
			uadmin.Save(&auction[0])

			// offers1 := []models.Offer{}
			// modelschema1 := uadmin.Schema["offer"]
			// uadmin.FilterList(&modelschema1, "id", false, 0, 1, &offers1, "auction_id = ?", auction[0].ID)

			// resToShowToClient := []map[string]interface{}{}

			// resToShowToClient = append(resToShowToClient, map[string]interface{}{
			// 	"CurrentFall":                  auction[0].CurrentFall,
			// 	"ItWasABestOffer":              ItIsABestOffer,
			// 	"CurrentPrice":                 auction[0].CurrentPrice,
			// 	"YourLastOffer":                auction[0].YourLastOffer,
			// 	"YourCurrentPosition":          auction[0].YourCurrentPosition,
			// 	"AmountOfOffersThatWerePosted": len(offers),
			// })
			// uadmin.ReturnJSON(w, r, resToShowToClient)

			UpdateAuctionID := r.FormValue("AuctionID")
			UpdateLastBidTime := r.FormValue("LastBidTime")
			// i, err := strconv.ParseInt(UpdateLastBidTime, 10, 64)
			// if err != nil {
			// 	panic(err)
			// }
			// tm := time.Unix(i, 0)
			auction = []models.Auction{}
			uadmin.Filter(&auction, "auction_id = ?", UpdateAuctionID)
			var tm time.Time
			if UpdateLastBidTime == "" {
				tm = auction[0].BiddingStart
			} else {
				i, err := strconv.ParseInt(UpdateLastBidTime, 10, 64)
				if err != nil {
					panic(err)
				}
				tm = time.Unix(i, 0)
			}

			offers = []models.Offer{}
			uadmin.Filter(&offers, "auction_id = ?", auction[0].ID)
			var thereWasOffersBefore bool
			if len(offers) > 0 {
				thereWasOffersBefore = true
			}
			uadmin.Filter(&offers, "auction_id = ? AND participant = ?", auction[0].ID, "Акционерное Общество \"Третий парк\"")
			var thereWasYourOffer bool
			if len(offers) > 0 {
				thereWasYourOffer = true
			}
			uadmin.Filter(&offers, "auction_id = ?  AND bid_time > ?", auction[0].ID, tm)
			offersTable := []map[string]interface{}{}
			for i := range offers {
				uadmin.Preload(&offers[i])
				offersTable = append(offersTable, map[string]interface{}{
					"Bid":         addSpacesToPrice(fmt.Sprintf("%.2f", offers[i].Bid)),
					"Participant": offers[i].Participant,
					"BidTime":     offers[i].BidTime.Unix(),
					"BestBid":     offers[i].BestBid,
				})
			}
			lastOffer := models.Offer{}
			modelschema1 := uadmin.Schema["offer"]
			uadmin.FilterList(&modelschema1, "id", false, 0, 1, &lastOffer, "auction_id = ? AND best_bid = ?", auction[0].ID, true)

			results = []map[string]interface{}{}
			results = append(results, map[string]interface{}{
				"AuctionID":                    auction[0].AuctionID,
				"Name":                         auction[0].Name,
				"BidDuration":                  time.Now().Unix() - auction[0].BiddingStart.Unix(),
				"CurrentServerTime":            time.Now().Unix(),
				"CurrentFall":                  auction[0].CurrentFall,
				"CurrentPrice":                 addSpacesToPrice(fmt.Sprintf("%.2f", auction[0].CurrentPrice)),
				"thereWasOffersBefore":         thereWasOffersBefore,
				"thereWasYourOffer":            thereWasYourOffer,
				"YourLastOffer":                auction[0].YourLastOffer,
				"YourCurrentPosition":          auction[0].YourCurrentPosition,
				"AmountOfOffersThatWerePosted": len(offers),
				"offersTable":                  offersTable,
				"lastOfferTime":                lastOffer.BidTime.Unix(),
				"BestPretendentName":           lastOffer.Participant,
				"BestPretendentBid":            lastOffer.Bid,
				"Status":                       auction[0].Status,
				"BiddingDuration":              fmtDuration(auction[0].BiddingFinish.Sub(auction[0].BiddingStart)),
				"BiddingStart":                 auction[0].BiddingStart.Unix(),
			})
			uadmin.ReturnJSON(w, r, results)
			// fmt.Println("PageUpdate", UpdateAuctionID, UpdateLastBidTime, tm)

		case "PageUpdate":

			UpdateAuctionID := r.FormValue("AuctionID")
			UpdateLastBidTime := r.FormValue("LastBidTime")
			auction := []models.Auction{}
			uadmin.Filter(&auction, "auction_id = ?", UpdateAuctionID)
			var tm time.Time
			if UpdateLastBidTime == "" {
				tm = auction[0].BiddingStart
			} else {
				i, err := strconv.ParseInt(UpdateLastBidTime, 10, 64)
				if err != nil {
					panic(err)
				}
				tm = time.Unix(i, 0)
			}

			offers := []models.Offer{}
			uadmin.Filter(&offers, "auction_id = ?", auction[0].ID)
			var thereWasOffersBefore bool
			if len(offers) > 0 {
				thereWasOffersBefore = true
			}
			uadmin.Filter(&offers, "auction_id = ? AND participant = ?", auction[0].ID, "Акционерное Общество \"Третий парк\"")
			var thereWasYourOffer bool
			if len(offers) > 0 {
				thereWasYourOffer = true
			}
			uadmin.Filter(&offers, "auction_id = ?  AND bid_time > ?", auction[0].ID, tm)
			offersTable := []map[string]interface{}{}
			for i := range offers {
				uadmin.Preload(&offers[i])
				offersTable = append(offersTable, map[string]interface{}{
					"Bid":         addSpacesToPrice(fmt.Sprintf("%.2f", offers[i].Bid)),
					"Participant": offers[i].Participant,
					"BidTime":     offers[i].BidTime.Unix(),
					"BestBid":     offers[i].BestBid,
				})
			}
			lastOffer := models.Offer{}
			modelschema1 := uadmin.Schema["offer"]
			uadmin.FilterList(&modelschema1, "id", false, 0, 1, &lastOffer, "auction_id = ? AND best_bid = ?", auction[0].ID, true)

			results := []map[string]interface{}{}
			results = append(results, map[string]interface{}{
				"AuctionID":                    auction[0].AuctionID,
				"Name":                         auction[0].Name,
				"BidDuration":                  time.Now().Unix() - auction[0].BiddingStart.Unix(),
				"CurrentServerTime":            time.Now().Unix(),
				"CurrentFall":                  auction[0].CurrentFall,
				"CurrentPrice":                 addSpacesToPrice(fmt.Sprintf("%.2f", auction[0].CurrentPrice)),
				"thereWasOffersBefore":         thereWasOffersBefore,
				"thereWasYourOffer":            thereWasYourOffer,
				"YourLastOffer":                auction[0].YourLastOffer,
				"YourCurrentPosition":          auction[0].YourCurrentPosition,
				"AmountOfOffersThatWerePosted": len(offers),
				"offersTable":                  offersTable,
				"lastOfferTime":                lastOffer.BidTime.Unix(),
				"BestPretendentName":           lastOffer.Participant,
				"BestPretendentBid":            lastOffer.Bid,
				"Status":                       auction[0].Status,
				"BiddingDuration":              fmtDuration(auction[0].BiddingFinish.Sub(auction[0].BiddingStart)),
				"BiddingStart":                 auction[0].BiddingStart.Unix(),
			})
			uadmin.ReturnJSON(w, r, results)
			// fmt.Println("PageUpdate", UpdateAuctionID, UpdateLastBidTime, tm)
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func Pow(a *big.Float, e uint64) *big.Float {
	result := Zero().Copy(a)
	for i := uint64(0); i < e-1; i++ {
		result = Mul(result, a)
	}
	return result
}

func Root(a *big.Float, n uint64) *big.Float {
	limit := Pow(New(2), 256)
	n1 := n - 1
	n1f, rn := New(float64(n1)), Div(New(1.0), New(float64(n)))
	x, x0 := New(1.0), Zero()
	_ = x0
	for {
		potx, t2 := Div(New(1.0), x), a
		for b := n1; b > 0; b >>= 1 {
			if b&1 == 1 {
				t2 = Mul(t2, potx)
			}
			potx = Mul(potx, potx)
		}
		x0, x = x, Mul(rn, Add(Mul(n1f, x), t2))
		if Lesser(Mul(Abs(Sub(x, x0)), limit), x) {
			break
		}
	}
	return x
}

func Abs(a *big.Float) *big.Float {
	return Zero().Abs(a)
}

func New(f float64) *big.Float {
	r := big.NewFloat(f)
	r.SetPrec(256)
	return r
}

func Div(a, b *big.Float) *big.Float {
	return Zero().Quo(a, b)
}

func Zero() *big.Float {
	r := big.NewFloat(0.0)
	r.SetPrec(256)
	return r
}

func Mul(a, b *big.Float) *big.Float {
	return Zero().Mul(a, b)
}

func Add(a, b *big.Float) *big.Float {
	return Zero().Add(a, b)
}

func Sub(a, b *big.Float) *big.Float {
	return Zero().Sub(a, b)
}

func Lesser(x, y *big.Float) bool {
	return x.Cmp(y) == -1
}
