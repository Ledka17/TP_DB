package handler

import (
	"TP_DB/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func (h *DataBaseHandler) VoteOnThread(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		slugOrId := mux.Vars(r)["slug_or_id"]

		if h.usecase.IsThreadInDB(slugOrId) {
			decoder := json.NewDecoder(r.Body)
			var vote model.Vote
			err := decoder.Decode(&vote)
			checkErr(err)

			body, err := json.Marshal(h.usecase.VoteForThreadInDB(slugOrId, vote))
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