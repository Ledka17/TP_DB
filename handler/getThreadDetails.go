package handler

import (
	"TP_DB/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func (h *DataBaseHandler) GetThreadDetails(w http.ResponseWriter, r *http.Request) {
	slugOrId := mux.Vars(r)["slug_or_id"]

	if r.Method == "GET" {
		if h.usecase.IsThreadInDB(slugOrId) {
			body, err := json.Marshal(h.usecase.GetThreadInDB(slugOrId))
			checkErr(err)

			w.Write(body)
			w.WriteHeader(200)
			return
		}

		writeWithError(w, 404)
		return
	}

	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var threadUpdate model.ThreadUpdate
		err := decoder.Decode(&threadUpdate)
		checkErr(err)

		if h.usecase.IsThreadInDB(slugOrId) {
			body, err := json.Marshal(h.usecase.ChangeThreadInDB(threadUpdate, slugOrId))
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
