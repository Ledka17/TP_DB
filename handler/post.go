package handler

import (
	"encoding/json"
	"github.com/Ledka17/TP_DB/model"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (h *DataBaseHandler) PostDetailsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	checkErr(err)
	related := r.URL.Query()["related"]

	if r.Method == "GET" {
		if h.usecase.IsPostInDB(id) {
			body, err := json.Marshal(h.usecase.GetPostInDB(id, related))
			checkErr(err)

			w.WriteHeader(200)
			w.Write(body)
			return
		}
		writeWithError(w, 404)
		return
	}

	if r.Method == "POST" {
		if h.usecase.IsPostInDB(id) {
			decoder := json.NewDecoder(r.Body)
			var post model.PostUpdate
			err := decoder.Decode(&post)
			checkErr(err)

			body, err := json.Marshal(h.usecase.ChangePostInDB(id, post))
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

