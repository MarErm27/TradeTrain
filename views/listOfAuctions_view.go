package views

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MarErm27/webTrainerWithUAdmin/models"
	"github.com/uadmin/uadmin"
)

// type AuctionList struct {
// 	Title        string
// 	InitialPrice string
// 	// Done         bool
// }

// type TodoPageData struct {
// 	PageTitle string
// 	Todos     []AuctionList
// }
// type TodoPageData struct {
// 	PageTitle string
// 	Todos     []map[string]interface{}{}
// }
func addSpacesToPrice(price string) string {
	var s []string
	var positionOfAChar int
	str := strings.Split(price, "")
	for i := len(str) - 1; i >= 0; i-- {

		if string(str[i]) == "." {
			positionOfAChar = 0
		}
		s = append(s, string(str[i]))
		if string(str[i]) == "-" {
			positionOfAChar = 0
			continue
		}
		if positionOfAChar == 3 {
			s = append(s, " ")
			positionOfAChar = 0
		}

		positionOfAChar += 1

	}
	var stringWithSpaces string
	for i := len(s) - 1; i >= 0; i-- {

		stringWithSpaces += s[i]

	}

	return stringWithSpaces
}

// type UniqueRand struct {
// 	generated map[int]bool
// }

// var u *UniqueRand
// var u = make([]map[bool]UniqueRand, 0)

// var u map[int]bool

func getRandUnicNumberInRange(u map[int]bool, numberOfParticipants int) int {
	for {
		i := rand.Intn(numberOfParticipants-1+1) + 1
		// u.generated[i] = true
		if !u[i] {
			u[i] = true
			return i
		} else {
			return getRandUnicNumberInRange(u, numberOfParticipants)
		}
	}
}

func createParticipants(title string) {
	var rangeOfOffer int
	var minRangeOfOffer int
	var offerOfParticipant int
	var minOfferOfParticipant float64
	var maxOffer int
	var maxLimit int
	var minLimit int
	var sumToAddToLargestOffOfParticipantToGeMaxOfferOfUser int
	var sumToSubtractFromLargestOffOfParticipantToGeMinOfferOfUser int
	var YourFallLimitToWin float64
	var YourFallLimitToLoose float64
	rand.Seed(time.Now().UnixNano())
	numberOfParticipants := 2 + rand.Intn(7-2+1)
	u := make(map[int]bool)
	// var u map[int]bool
	// UniqueRand = nil
	for i := 1; i < numberOfParticipants; i++ {
		rand.Seed(time.Now().UnixNano())
		participant := models.Participant{}
		rangeOfOffer = i * 3
		minRangeOfOffer = rangeOfOffer - 2
		// rand.Intn(max-min+1) + min
		offerOfParticipant = rand.Intn(rangeOfOffer-minRangeOfOffer+1) + minRangeOfOffer
		if i == 1 {
			minOfferOfParticipant = float64(offerOfParticipant)
		}

		// participant.Name = strconv.Itoa(1 + rand.Intn(numberOfParticipants-1+1))
		participant.Name = strconv.Itoa(getRandUnicNumberInRange(u, numberOfParticipants))
		participant.AuctionID = title
		participant.OfferOfParticipant = float64(offerOfParticipant)
		uadmin.Save(&participant)
	}
	maxOffer = offerOfParticipant
	maxLimit = 22 - maxOffer

	minLimit = maxOffer
	// if minLimit == 0 {
	// 	minLimit = 1
	// }
	auction := models.Auction{}
	modelschema := uadmin.Schema["auction"]
	uadmin.FilterList(&modelschema, "id", true, 0, 1, &auction, "auction_id = ?", title)
	// if len(auction) > 0 {

	if auction.AuctionID != "" {
		sumToAddToLargestOffOfParticipantToGeMaxOfferOfUser = rand.Intn(maxLimit-1+1) + 1
		sumToSubtractFromLargestOffOfParticipantToGeMinOfferOfUser = rand.Intn(minLimit-1+1) + 1
		YourFallLimitToWin = float64(maxOffer + sumToAddToLargestOffOfParticipantToGeMaxOfferOfUser)
		YourFallLimitToLoose = float64(maxOffer - sumToSubtractFromLargestOffOfParticipantToGeMinOfferOfUser)
		if YourFallLimitToWin == 22 {
			YourFallLimitToWin = 21.99
		}
		if YourFallLimitToWin <= 21 {
			if rand.Intn(2-1+1)+1 == 1 {
				YourFallLimitToWin += 0.5
			}
		}

		if YourFallLimitToLoose >= 1 {
			if rand.Intn(2-1+1)+1 == 1 {
				YourFallLimitToLoose -= 0.5
			}
		}

		if YourFallLimitToLoose <= 0 {
			YourFallLimitToLoose = 0.5
		} else if YourFallLimitToLoose <= minOfferOfParticipant {
			minOfferOfParticipant -= 0.5
		}

		auction.YourFallLimitToWin = YourFallLimitToWin
		auction.YourFallLimitToLoose = YourFallLimitToLoose
		uadmin.Save(&auction)
	}
}

func createAuctions() {
	then := time.Date(
		2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	rand.Seed(then.UnixNano())

	// Auctions := []AuctionList{}

	initialPriceOfAllAuctions := 160000000000
	averagePrice := initialPriceOfAllAuctions / 40
	currentMaximumNumber := 42
	randInt := 1
	var randInitPrice int64

	var weHaveToCreateParticipants bool
	participants := []models.Participant{}
	uadmin.All(&participants)
	if len(participants) == 0 {
		weHaveToCreateParticipants = true
	}

	for i := 1; i <= currentMaximumNumber; i++ {
		rand.Seed(time.Now().UnixNano())
		randomN := 1 + rand.Intn(2-1+1)
		var plusMinusK int
		if randomN == 1 {
			plusMinusK = -1
		} else {
			plusMinusK = 1
		}

		rand.Seed(then.UnixNano())
		auction := models.Auction{}
		randInt = randInt + currentMaximumNumber + rand.Intn(currentMaximumNumber-i+1)

		randInitPrice = int64(averagePrice*100+rand.Intn((averagePrice*60)-i+1)*plusMinusK) / 100

		rPrice := float64(randInitPrice)
		// rPrice := SliceOfPrices[i-1]
		// rPrice := Lots[i-1]["InitialPrice"]
		title := "0" + strconv.Itoa(372200004520000002+randInt)

		auction.AuctionID = title
		// auction.InitialPrice, _ = rPrice.(float64)
		auction.InitialPrice = rPrice
		auction.Name = "Услуги (работы) по регулярным перевозкам пассажиров автобусами в городском и пригородном сообщении"
		auction.Organizer = "СБП ГКУ \"Организатор перевозок\""
		auction.BiddingStart = time.Now()
		// auction.CurrentPrice, _ = rPrice.(float64)
		auction.CurrentPrice = rPrice
		auction.HalfPercentOfInitialPrice = auction.InitialPrice * 0.005
		auction.FivePercentsOfInitialPrice = auction.InitialPrice * 0.05
		auction.Status = "Прием заявок"
		// auction.LotNumber = Lots[i-1]["LotNumber"].(int)
		// auction.Priority = Lots[i-1]["Priority"].(int)
		// auction.OurAuction = Lots[i-1]["OurAuction"].(bool)
		// auction.GOR = Lots[i-1]["GOR"].(float64)
		// auction.VIT = Lots[i-1]["VIT"].(float64)
		// auction.ZEL = Lots[i-1]["ZEL"].(float64)
		// auction.NumberOfVehiclesBK = Lots[i-1]["NumberOfVehiclesBK"].(float64)
		// auction.NumberOfVehiclesSK1 = Lots[i-1]["NumberOfVehiclesSK1"].(float64)
		// auction.NumberOfVehiclesSK2 = Lots[i-1]["NumberOfVehiclesSK2"].(float64)
		// auction.NumberOfVehicles = Lots[i-1]["NumberOfVehicles"].(float64)
		// auction.FallLimit = Lots[i-1]["FallLimit"].(float64)
		// auction.PlanFallLimit = Lots[i-1]["PlanFallLimit"].(float64)

		auction.AreWeInvolved = true
		uadmin.Save(&auction)
		if weHaveToCreateParticipants {
			createParticipants(title)
		}
		// break
	}
	fmt.Println("AuctionsISCreated")
}

// listOfAuctionsHandler !
func listOfAuctionsHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is not logged in
	// if uadmin.IsAuthenticated(r) == nil {
	// 	http.Redirect(w, r, "/admin", http.StatusSeeOther)
	// 	return
	// }

	// r.URL.Path creates a new path called /todo
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/ListOfAuctions")

	rand.Seed(time.Now().UnixNano())
	var auctionListTempName string
	n := 1 + rand.Intn(1000-1+1) // a ≤ n ≤ b
	if n == 1 {
		auctionListTempName = "templates/noInternet.html"
		uadmin.RenderHTML(w, r, auctionListTempName, nil)
		return
	} else {
		auctionListTempName = "templates/listOfAuctions.html"
	}

	auctions := []models.Auction{}
	uadmin.All(&auctions)
	if len(auctions) == 0 {
		createAuctions()
		uadmin.All(&auctions)
	}
	// fmt.Println("LotNumber", "Priority", "AuctionID", "InitialPrice") //распечатка для Александра

	results := []map[string]interface{}{}
	for i := range auctions {
		uadmin.Preload(&auctions[i])
		var status bool
		if auctions[i].Status == "Прием заявок" {
			status = true
		}

		// fmt.Println(auctions[i].LotNumber, auctions[i].Priority, auctions[i].AuctionID, fmt.Sprintf("%.2f", auctions[i].InitialPrice)) //распечатка для Александра
		// Assigns the string of interface in each Todo fields
		results = append(results, map[string]interface{}{
			"Title":        auctions[i].AuctionID,
			"InitialPrice": addSpacesToPrice(fmt.Sprintf("%.2f", auctions[i].InitialPrice)),
			"Status":       auctions[i].Status,
			"StatusBool":   status,
		})
	}

	// data := TodoPageData{
	// 	PageTitle: "My TODO list",
	// 	Todos:     results,
	// }
	// Pass TodoList data object to the specified HTML path
	uadmin.RenderHTML(w, r, auctionListTempName, results)
	// uadmin.RenderHTML(w, r, auctionListTempName, results)
}
