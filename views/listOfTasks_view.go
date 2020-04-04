package views

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/MarErm27/webTrainerWithUAdmin/models"
	"github.com/uadmin/uadmin"
)

// AnaliticsPanelPageData ...
type ListOfTasksPageData struct {
	Images  [][]string
	Results []map[string]interface{}
}

// ListOfTasksHandler !
func ListOfTasksHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is not logged in
	// if uadmin.IsAuthenticated(r) == nil {
	// 	http.Redirect(w, r, "/admin", http.StatusSeeOther)
	// 	return
	// }
	// r.URL.Path creates a new path called /todo
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/ListOfTasks")

	switch r.Method {
	case "GET":
		var taskListTempName string

		taskListTempName = "templates/listOfTasks.html"

		auctions := []models.Auction{}
		uadmin.All(&auctions)
		if len(auctions) == 0 {
			createAuctions()
			uadmin.All(&auctions)
		}
		uadmin.AdminPage("number_of_collection", true, 0, len(auctions), &auctions, "")
		results := []map[string]interface{}{}
		var previousNumber int
		for i := range auctions {
			uadmin.Preload(&auctions[i])
			if auctions[i].NumberOfCollection != previousNumber {
				s := strings.Split(strconv.Itoa(auctions[i].NumberOfCollection), "")
				fmt.Println("path", r.URL.Path)
				results = append(results, map[string]interface{}{
					"SliceOfNumber":      s,
					"NumberOfCollection": auctions[i].NumberOfCollection,
				})
				previousNumber = auctions[i].NumberOfCollection
			}
		}

		// data := TodoPageData{
		// 	PageTitle: "My TODO list",
		// 	Todos:     results,
		// }
		// Pass TodoList data object to the specified HTML path

		// sumData = ListOfTasksPageData{
		// 	Images:  images,
		// 	Results: results,
		// }
		uadmin.RenderHTML(w, r, taskListTempName, results)

	case "POST":

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}
