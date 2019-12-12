package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (h *DataBaseHandler) GetForumUsersHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	limit, _ := strconv.Atoi(queryParams["limit"][0])
	since := queryParams["since"][0]
	desc, _ := strconv.ParseBool(queryParams["desc"][0])

	slug := mux.Vars(r)["slug"]

	if r.Method == "GET" {
		if h.usecase.IsForumInDB(slug) {
			body, err := json.Marshal(h.usecase.GetForumUsersInDB(slug, limit, since, desc))
			checkErr(err)

			w.WriteHeader(200)
			w.Write(body)
			return
		}
		writeWithError(w, 404)
		return
	}
	w.WriteHeader(400)
}
