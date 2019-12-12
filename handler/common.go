package handler

import (
	"TP_DB/forum"
	"TP_DB/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type DataBaseHandler struct {
	usecase forum.Repository
}

func NewHandler(r *mux.Router, usecase forum.Repository) {
	handler := DataBaseHandler{usecase: usecase}

	r.HandleFunc("/forum/create", handler.CreateForumHandler)
	r.HandleFunc("/forum/{slug}/create", handler.CreateThreadHandler)
	r.HandleFunc("/forum/{slug}/details", handler.GetForumDetailsHandler)
	r.HandleFunc("/forum/{slug}/threads", handler.GetForumThreadsHandler)
	r.HandleFunc("/forum/{slug}/users", handler.GetForumUsersHandler)

	r.HandleFunc("/thread/{slug_or_id}/create", handler.CreateThreadPosts)
	r.HandleFunc("/thread/{slug_or_id}/details", handler.GetThreadDetails)
	r.HandleFunc("/thread/{slug_or_id}/posts", handler.GetThreadPosts)
	r.HandleFunc("/thread/{slug_or_id}/vote", handler.VoteOnThread)

	r.HandleFunc("/post/{id}/details", handler.PostDetailsHandler)

	r.HandleFunc("/user/{nickname}/create", handler.CreateUserHandler)
	r.HandleFunc("/user/{nickname}/profile", handler.UserProfileHandler)

	r.HandleFunc("/service/clear", handler.ClearDBHandler)
	r.HandleFunc("/service/status", handler.GetServiceStatusHandler)
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func writeWithError(w http.ResponseWriter, statusCode int) {
	errorRes := model.Error{"Ошибка."}
	body, err := json.Marshal(errorRes)
	checkErr(err)

	w.Write(body)
	w.WriteHeader(statusCode)
}