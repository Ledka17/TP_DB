package handler

import (
	"encoding/json"
	"github.com/Ledka17/TP_DB/model"
	"github.com/gorilla/mux"
	"net/http"
)

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
