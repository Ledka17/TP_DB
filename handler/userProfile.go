package handler

import (
	"encoding/json"
	"github.com/Ledka17/TP_DB/model"
	"github.com/gorilla/mux"
	"net/http"
)

func (h *DataBaseHandler) UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	nickname := mux.Vars(r)["nickname"]
	if r.Method == "GET" {
		if h.usecase.IsUserInDB(nickname, "") {
			body, err := json.Marshal(h.usecase.GetUserInDB(nickname, ""))
			checkErr(err)

			w.WriteHeader(200)
			w.Write(body)
			return
		}
		writeWithError(w, 404)
		return
	}
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var userUpdate model.UserUpdate
		err := decoder.Decode(&userUpdate)
		checkErr(err)

		if h.usecase.IsUserInDB(nickname, "") {
			if h.usecase.IsUserInDB("", userUpdate.Email) {
				writeWithError(w, 409)
				return
			}
			body, err := json.Marshal(h.usecase.GetUserInDB(nickname, ""))
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
