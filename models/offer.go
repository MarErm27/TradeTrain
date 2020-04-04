package models

import (
	"time"

	"github.com/uadmin/uadmin"
)

// Offer model ...
type Offer struct {
	uadmin.Model
	Name        string
	Auction     Auction // <-- Category Model
	AuctionID   uint    // <-- CategoryID
	Bid         float64
	Participant string
	BidTime     time.Time
	BestBid     bool
}
