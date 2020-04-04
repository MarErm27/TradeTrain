package models

import (
	"fmt"
	"time"

	"github.com/uadmin/uadmin"
)

// AuctionVersion !
type AuctionVersion struct {
	uadmin.Model
	Auction                   Auction
	AuctionID                 uint
	Number                    int `uadmin:"help:version number"`
	Date                      time.Time
	FallLimit                 float64
	PersonResponsibleForLimit string
}

// Returns the version number
func (a AuctionVersion) String() string {
	return fmt.Sprint(a.Number)
}
