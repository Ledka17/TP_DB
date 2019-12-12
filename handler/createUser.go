package handler

import (
	"encoding/json"
	"github.com/Ledka17/TP_DB/model"
	"github.com/gorilla/mux"
	"net/http"
)

func (h *DataBaseHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nickname := mux.Vars(r)["nickname"]

		decoder := json.NewDecoder(r.Body)
		var user model.User
		err := decoder.Decode(&user)
		checkErr(err)

		if h.usecase.IsUserInDB(nickname, user.Email) {
			body, err := json.Marshal(h.usecase.GetUserInDB(nickname, user.Email))
			checkErr(err)

			w.WriteHeader(409)
			w.Write(body)
			return
		}

		body, err := json.Marshal(h.usecase.СreateUserInDB(nickname, user))
		checkErr(err)

		w.WriteHeader(201)
		w.Write(body)
		return
	}
	w.WriteHeader(400)
}
