package handler

import (
	"encoding/json"
	"github.com/Ledka17/TP_DB/model"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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

func (h *DataBaseHandler) CreateThreadPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var posts []model.Post
		err := decoder.Decode(&posts)
		checkErr(err)

		slugOrId := mux.Vars(r)["slug_or_id"]
		statusCode := 404

		if h.usecase.IsThreadInDB(slugOrId) {
			if h.usecase.CheckParentPost(posts) {
				body, err := json.Marshal(h.usecase.CreatePostsInDB(posts))
				checkErr(err)
				w.Write(body)
				w.WriteHeader(200)
				return
			}
			statusCode = 409
		}

		writeWithError(w, statusCode)
		return
	}
	w.WriteHeader(400)
}

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

func (h *DataBaseHandler) GetThreadPosts(w http.ResponseWriter, r *http.Request) {
	slugOrId := mux.Vars(r)["slug_or_id"]
	vars := r.URL.Query()
	limit, _ := strconv.Atoi(vars["limit"][0])
	since, _ := strconv.Atoi(vars["since"][0])
	sort := vars["sort"][0]
	desc, _ := strconv.ParseBool(vars["desc"][0])


	if r.Method == "GET" {
		if h.usecase.IsThreadInDB(slugOrId) {
			body, err := json.Marshal(h.usecase.GetPostsInDB(slugOrId, limit, since, sort, desc))
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
