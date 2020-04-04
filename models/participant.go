package models

import (
	"github.com/uadmin/uadmin"
)

// Participant model ...
type Participant struct {
	uadmin.Model
	Name               string
	AuctionID          string
	OfferOfParticipant float64
}
