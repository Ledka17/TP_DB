package handler

import (
	"encoding/json"
	"github.com/Ledka17/TP_DB/forum"
	"github.com/Ledka17/TP_DB/model"
	"github.com/labstack/echo"
	"log"
)

type DataBaseHandler struct {
	usecase forum.Repository
}

func NewHandler(e *echo.Echo, usecase forum.Repository) {
	handler := DataBaseHandler{usecase: usecase}

	e.POST("/forum/create", handler.CreateForumHandler)
	e.POST("/forum/:slug/create", handler.CreateThreadHandler)
	e.GET("/forum/:slug/details", handler.GetForumDetailsHandler)
	e.GET("/forum/:slug/threads", handler.GetForumThreadsHandler)
	e.GET("/forum/:slug/users", handler.GetForumUsersHandler)

	e.POST("/thread/:slug_or_id/create", handler.CreateThreadPosts)
	e.GET("/thread/:slug_or_id/details", handler.GetThreadDetails)
	e.POST("/thread/:slug_or_id/details", handler.ChangeThreadDetails)
	e.GET("/thread/:slug_or_id/posts", handler.GetThreadPosts)
	e.POST("/thread/:slug_or_id/vote", handler.VoteOnThread)

	e.GET("/post/:id/details", handler.GetPostDetailsHandler)
	e.POST("/post/:id/details", handler.ChangePostDetailsHandler)

	e.POST("/user/:nickname/create", handler.CreateUserHandler)
	e.GET("/user/:nickname/profile", handler.GetUserProfileHandler)
	e.POST("/user/:nickname/profile", handler.ChangeUserProfileHandler)

	e.POST("/service/clear", handler.ClearDBHandler)
	e.GET("/service/status", handler.GetServiceStatusHandler)
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func writeWithError(c echo.Context, statusCode int) error {
	errorRes := model.Error{"Ошибка."}
	body, err := json.Marshal(errorRes)
	checkErr(err)

	return c.JSON(statusCode, body)
}