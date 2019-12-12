package handler

import (
	"encoding/json"
	"github.com/Ledka17/TP_DB/model"
	"github.com/gorilla/mux"
	"net/http"
)

func (h *DataBaseHandler) CreateThreadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var thread model.Thread
		err := decoder.Decode(&thread)
		checkErr(err)

		if thread.Author == "" {
			writeWithError(w, 404)
			return
		}

		if h.usecase.IsThreadInDB(thread.Slug) || h.usecase.IsForumInDB(slug) {
			body, err := json.Marshal(h.usecase.GetThreadInDB(slug))
			checkErr(err)

			w.WriteHeader(409)
			w.Write(body)
			return
		}
		body, err := json.Marshal(h.usecase.CreateThreadInDB(slug, thread))
		checkErr(err)

		w.WriteHeader(201)
		w.Write(body)
		return
	}
	w.WriteHeader(400)
}
