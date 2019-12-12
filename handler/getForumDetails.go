package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func (h *DataBaseHandler) GetForumDetailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	if r.Method == "GET" {
		if h.usecase.IsForumInDB(slug) {
			h.usecase.GetForumInDB(slug)
			body, err := json.Marshal(h.usecase.GetForumInDB(slug))
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
