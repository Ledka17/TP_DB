package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

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
