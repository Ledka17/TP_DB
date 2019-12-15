package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/Ledka17/TP_DB/model"
	"net/http"
	"strconv"
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

func (h *DataBaseHandler) GetForumThreadsHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	limit, _ := strconv.Atoi(queryParams["limit"][0])
	since := queryParams["since"][0]
	desc, _ := strconv.ParseBool(queryParams["desc"][0])

	slug := mux.Vars(r)["slug"]

	if r.Method == "GET" {
		if h.usecase.IsForumInDB(slug) {
			body, err := json.Marshal(h.usecase.GetThreadsForumInDB(slug, limit, since, desc))
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
