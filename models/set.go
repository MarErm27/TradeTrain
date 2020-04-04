package models

import (
	"github.com/uadmin/uadmin"
)

// Sets model ...
type Set struct {
	uadmin.Model
	Name                      string
	ResponsibleFromGroupOne   string
	ResponsibleFromGroupTwo   string
	ResponsibleFromGroupThree string
}
