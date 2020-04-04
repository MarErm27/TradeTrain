package models

import (
	"time"

	"github.com/uadmin/uadmin"
)

// Auction model ...
type Auction struct {
	uadmin.Model
	Name                        string
	Organizer                   string
	NotificationPublication     time.Time
	ApplicationDeadline         time.Time
	BiddingStart                time.Time
	BiddingFinish               time.Time
	Status                      string
	AuctionID                   string  `uadmin:"required;search"`
	InitialPrice                float64 `uadmin:"required;search"`
	HalfPercentOfInitialPrice   float64
	FivePercentsOfInitialPrice  float64
	CurrentPrice                float64
	CurrentFall                 float64
	YourFallLimitToWin          float64
	YourFallLimitToLoose        float64
	YourLastOffer               float64
	YourCurrentPosition         int
	TimeToFinish                time.Time
	CurrentBestOfferParticipant string
	CurrentBestOfferTime        time.Time
	TotalOffersSubmitted        int
	EstimatedDatesAndEndTime    time.Time
	OurAuction                  bool `uadmin:"search"`
	PersonResponsibleForOur     string
	NumberOfCollection          int
	NumberInCollection          int

	FallLimit                       float64
	PlanFallLimit                   float64
	PersonResponsibleForLimit       string
	OurFall                         float64
	Fall                            float64
	MinutesToNextStepOverLimit      float64
	TimeWhenWeGetMinutesBeforeLimit time.Time
	AreWeInvolved                   bool
	PersonResponsibleForInvolve     string
	LotNumber                       int
	Priority                        int
	GOR                             float64
	VIT                             float64
	ZEL                             float64
	NumberOfVehiclesBK              float64
	NumberOfVehiclesSK1             float64
	NumberOfVehiclesSK2             float64
	NumberOfVehicles                float64
	WeWon                           bool
	AuctionIsFinished               bool
	OneComputerIsLagging            bool
	HowLongOneComputerWillLag       time.Duration
}

func (a *Auction) Save() {
	// Multiply the Number and the Cost to get the value of the Sum
	a.HalfPercentOfInitialPrice = a.InitialPrice * 0.005
	a.FivePercentsOfInitialPrice = a.InitialPrice * 0.05
	a.CurrentPrice = a.InitialPrice
	// a.BidDuration = time.Now()
	uadmin.Save(a)
}
