package models

import (
	"github.com/uadmin/uadmin"
)

// AllAuctionsDataAndSettings model ...
type AllAuctionsDataAndSettings struct {
	uadmin.Model
	Name                   string
	AuctionsPerParticipant int
	NumberOfParticipants   int
	
}
