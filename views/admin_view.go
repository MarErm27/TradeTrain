package views

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/MarErm27/webTrainerWithUAdmin/models"
	"github.com/uadmin/uadmin"
)

var wg sync.WaitGroup

// AdminHandler !
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is not logged in
	// if uadmin.IsAuthenticated(r) == nil {
	// 	http.Redirect(w, r, "/admin", http.StatusSeeOther)
	// 	return
	// }
	// r.URL.Path creates a new path called /todo
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/Admin")

	switch r.Method {
	case "GET":
		var taskListTempName string

		taskListTempName = "templates/AdminPage.html"

		uadmin.RenderHTML(w, r, taskListTempName, nil)

	case "POST":
		//Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		// if err := r.ParseForm(); err != nil {
		// 	fmt.Fprintf(w, "ParseForm() err: %v", err)
		// 	return
		// }
		// fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		// name := r.FormValue("name")
		// address := r.FormValue("address")
		// fmt.Fprintf(w, "Name = %s\n", name)
		// fmt.Fprintf(w, "Address = %s\n", address)

		// //case method is get nomenclature list - get nomenclature for dropdown menu in a table
		Method := r.FormValue("Method")
		switch Method {

		case "CleanDB":
			offers := []models.Offer{}
			uadmin.All(&offers)
			for i := range offers {
				uadmin.Delete(&offers[i])
			}
			participants := []models.Participant{}
			uadmin.All(&participants)
			for i := range participants {
				uadmin.Delete(&participants[i])
			}

			auctions := []models.Auction{}
			uadmin.All(&auctions)
			for i := range auctions {
				uadmin.Delete(&auctions[i])
			}

			fmt.Println("CleanDB")
		case "setAllAuctionStartTime":
			UnixTimeToSet := r.FormValue("UnixTimeToSet")

			u, err := strconv.ParseInt(UnixTimeToSet, 10, 64)
			if err != nil {
				panic(err)
			}
			tm := time.Unix(u, 0)

			auctions := []models.Auction{}
			uadmin.All(&auctions)

			t1 := time.Now()
			for i := range auctions {
				wg.Add(1)
				uadmin.Preload(&auctions[i])

				auctions[i].BiddingStart = tm
				auctions[i].BiddingFinish = tm
				mtx.Lock()
				uadmin.Save(&auctions[i])
				mtx.Unlock()
				rand.Seed(time.Now().UnixNano())
				randomN := 1 + rand.Intn(599-1+1)

				diffToChangeStatus := tm.Sub(t1)
				tmAdd := tm.Add(time.Second * time.Duration(randomN))

				diff := tmAdd.Sub(t1)
				auctionID := auctions[i].AuctionID
				go changeAuctionStatusToActive(diffToChangeStatus, auctionID)
				go setParticipantAction(diff, auctionID, tm)
			}
			wg.Wait()
			fmt.Println("setAllAuctionStartTime", UnixTimeToSet)
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

var mtx sync.Mutex

func changeAuctionStatusToActive(biddingStartNumber time.Duration, auctionID string) { //, wg *sync.WaitGroup) {
	defer wg.Done()
	// fmt.Println("statusIsGoingToBeChanged", biddingStartNumber, "auctionID: ", auctionID)
	<-time.After(biddingStartNumber)

	// fmt.Println("changeAuctionStatusToActive: ", auctionID, biddingStartNumber)
	auction := []models.Auction{}
	uadmin.Filter(&auction, "auction_id = ?", auctionID)
	auction[0].Status = "Проведение аукциона"
	mtx.Lock()
	uadmin.Save(&auction[0])
	mtx.Unlock()
}

func randFloats(min, max float64, n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = min + rand.Float64()*(max-min)
	}
	return res
}

func setParticipantAction(randomNumber time.Duration, auctionID string, tm time.Time) {
	// fmt.Println("setParticipantAction", randomNumber, "auctionID: ", auctionID)
	<-time.After(randomNumber)

	results := []map[string]interface{}{}
	auction := []models.Auction{}
	uadmin.Filter(&auction, "auction_id = ?", auctionID)
	if tm != auction[0].BiddingStart {
		fmt.Println("This auction has different action running. No need for another one.")
		return
	}

	offers := []models.Offer{}
	modelschema := uadmin.Schema["offer"]
	uadmin.FilterList(&modelschema, "id", false, 0, 1, &offers, "auction_id = ? AND best_bid = ?", auction[0].ID, true)
	// uadmin.FilterList(&modelschema, "id", false, 0, 1, &offers, "auction_id = ?", auction[0].ID)
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
	var LastOfferTime time.Time
	var LastBestOrNotParticipant string

	var thisIsAFirstOffer bool
	if len(results) != 0 {
		LastOfferTime, _ = results[0]["BidTime"].(time.Time)
		LastBestOrNotParticipant = offers[0].Participant

	} else {
		thisIsAFirstOffer = true
	}

	offers = []models.Offer{}
	modelschema = uadmin.Schema["offer"]
	uadmin.FilterList(&modelschema, "id", false, 0, 1, &offers, "auction_id = ?", auction[0].ID)

	currentFall := auction[0].CurrentFall

	lastBestOffers := []models.Offer{}
	modelschema = uadmin.Schema["offer"]
	uadmin.FilterList(&modelschema, "id", false, 0, 1, &lastBestOffers, "auction_id = ? AND best_bid = ?", auction[0].ID, true)

	var LastBestParticipant string
	if len(lastBestOffers) > 0 {
		LastBestParticipant = lastBestOffers[0].Participant
	}

	participants := []models.Participant{}
	uadmin.Filter(&participants, "auction_id = ? AND name <> ? AND offer_of_participant >= ?", auctionID, LastBestParticipant, currentFall+0.5)

	// var sliceToShuffle []map[string]interface{}{}

	var OfferAuctionID string
	var OfferBidFloat64 float64
	var OfferParticipant string
	for i := range participants {
		uadmin.Preload(&participants[i])
		if participants[i].OfferOfParticipant > (currentFall + 0.5) {
			OfferAuctionID = participants[i].AuctionID
			OfferBidFloat64 = participants[i].OfferOfParticipant
			OfferParticipant = participants[i].Name
		}
		break
	}

	if OfferAuctionID == "" {
		participant := models.Participant{}
		modelschema := uadmin.Schema["participant"]
		uadmin.FilterList(&modelschema, "offer_of_participant", false, 0, 1, &participant, "auction_id = ? AND offer_of_participant >= ?", auctionID, currentFall+0.5)
		if participant.OfferOfParticipant > (currentFall + 0.5) {
			randomN := 1 + rand.Intn(599-1+1)
			t1 := time.Now()
			BidTime := time.Now().Add(time.Second * time.Duration(randomN))
			timeToChangeStatusOfAuction := LastOfferTime.Add(time.Minute * time.Duration(10))
			diffBetweenEndOfAuctionAndBitTime := timeToChangeStatusOfAuction.Sub(t1)
			if diffBetweenEndOfAuctionAndBitTime > 0 {
				durationToEndOfAuction := timeToChangeStatusOfAuction.Sub(t1)
				randomN = 1 + rand.Intn(int(durationToEndOfAuction.Seconds())+1)
				BidTime = time.Now().Add(time.Second * time.Duration(randomN))
				diff := BidTime.Sub(t1)
				go setParticipantAction(diff, auctionID, tm)
			} else {
				diff := timeToChangeStatusOfAuction.Sub(t1)
				go changeAuctionStatusToOver(diff, auctionID, LastBestOrNotParticipant)
			}
		} else {
			t1 := time.Now()
			timeToChangeStatusOfAuction := LastOfferTime.Add(time.Minute * time.Duration(10))
			diff := timeToChangeStatusOfAuction.Sub(t1)

			go changeAuctionStatusToOver(diff, auctionID, LastBestOrNotParticipant)
		}
		return
	}

	offer := models.Offer{}
	if LastBestOrNotParticipant != OfferParticipant &&
		OfferBidFloat64 >= (currentFall+0.5) {

		offer.Auction = auction[0]
		offer.Auction.AuctionID = OfferAuctionID
		var offerOfParticipant float64
		chanceToPlaceARandomOffer := 1 + rand.Intn(4-1+1)
		priceFrom := auction[0].CurrentPrice - auction[0].HalfPercentOfInitialPrice
		priceTo := auction[0].CurrentPrice - auction[0].FivePercentsOfInitialPrice
		rand.Seed(time.Now().UnixNano())
		calculatedBigOfferOfParticipant := randFloats(priceTo, priceFrom, 1)[0]
		if chanceToPlaceARandomOffer == 1 &&
			currentFall < 20 &&
			calculatedBigOfferOfParticipant > (auction[0].InitialPrice-((OfferBidFloat64*auction[0].InitialPrice)/100)) {
			offerOfParticipant = calculatedBigOfferOfParticipant
		} else {
			offerOfParticipant = priceFrom
		}
		offer.Bid = offerOfParticipant

		offer.Participant = OfferParticipant
		if (OfferBidFloat64-currentFall) > 0 || thisIsAFirstOffer {
			offer.BestBid = true
			auction[0].CurrentPrice = offerOfParticipant

			auction[0].YourCurrentPosition = getCurrentPosition(auction[0].ID) + 1
		} else {
			offer.BestBid = false
		}

		offer.BidTime = time.Now()

		auction[0].CurrentFall = 100 - auction[0].CurrentPrice*100/auction[0].InitialPrice
		// Store input in the Friend model
		uadmin.Save(&offer)
	}
	uadmin.Save(&auction[0])

	randomN := 1 + rand.Intn(599-1+1)
	t1 := time.Now()
	BidTime := offer.BidTime.Add(time.Second * time.Duration(randomN))

	diff := BidTime.Sub(t1)
	go setParticipantAction(diff, auctionID, tm)

}

func changeAuctionStatusToOver(randomNumber time.Duration, auctionID string, lastBestOrNotParticipant string) {
	<-time.After(randomNumber)

	auction := []models.Auction{}
	uadmin.Filter(&auction, "auction_id = ?", auctionID)

	offers := []models.Offer{}
	modelschema := uadmin.Schema["offer"]
	uadmin.FilterList(&modelschema, "id", false, 0, 1, &offers, "auction_id = ?", auction[0].ID)

	var LastOfferTime time.Time
	var LastBestOrNotParticipantFinalCheck string

	if len(offers) > 0 {
		LastOfferTime = offers[0].BidTime
		LastBestOrNotParticipantFinalCheck = offers[0].Participant
	}

	if LastBestOrNotParticipantFinalCheck == lastBestOrNotParticipant {
		auction[0].Status = "Подведение итогов"
		auction[0].BiddingFinish = LastOfferTime.Add(time.Minute * time.Duration(10))
		uadmin.Save(&auction[0])
	} else {

		t1 := time.Now()
		timeToChangeStatusOfAuction := LastOfferTime.Add(time.Minute * time.Duration(10))
		diff := timeToChangeStatusOfAuction.Sub(t1)
		go changeAuctionStatusToOver(diff, auctionID, LastBestOrNotParticipantFinalCheck)
	}
}
