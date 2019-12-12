package handler

import (
	"TP_DB/model"
	"encoding/json"
	"net/http"
)

func (h *DataBaseHandler) CreateForumHandler (w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var forum model.Forum
		err := decoder.Decode(&forum)
		checkErr(err)

		statusCode := 201
		body := []byte("")

		if forum.User == "" {
			statusCode = 404
			errorRes := model.Error{"Владелец форума не найден."}
			body, err = json.Marshal(errorRes)
			checkErr(err)
		} else {
			if h.usecase.IsForumInDB(forum.Slug) {
				statusCode = 409
				body, err = json.Marshal(h.usecase.GetForumInDB(forum.Slug))
				checkErr(err)
			} else {
				statusCode = 201
				body, err = json.Marshal(h.usecase.CreateForumInDB(forum))
				checkErr(err)
			}
		}

		w.WriteHeader(statusCode)
		w.Write(body)
		return
	}
	w.WriteHeader(400)
}
