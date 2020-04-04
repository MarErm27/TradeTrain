package views

import (
	"net/http"
	"strings"
)

// HTTPHandler !
func HTTPHandler(w http.ResponseWriter, r *http.Request) {
	// r.URL.Path creates a new path called /http_handler
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/http_handler")
	if strings.HasPrefix(r.URL.Path, "/ListOfAuctions") {
		listOfAuctionsHandler(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/Auction") {
		AuctionHandler(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/Task") {
		TaskHandler(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/Analytics") {
		AnalyticsHandler(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/Admin") {
		AdminHandler(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/ListOfTasks") {
		ListOfTasksHandler(w, r)
		return
	}
}
