package handler

import (
	"github.com/Ledka17/TP_DB/forum"
	"github.com/labstack/echo"
	"log"
)

type DataBaseHandler struct {
	usecase forum.Repository
}

func NewHandler(e *echo.Echo, usecase forum.Repository) {
	handler := DataBaseHandler{usecase: usecase}

	e.POST(ForumPath + "/create", handler.CreateForumHandler)
	e.POST(ForumPath + "/:slug/create", handler.CreateThreadHandler)
	e.GET(ForumPath + "/:slug/details", handler.GetForumDetailsHandler)
	e.GET(ForumPath + "/:slug/threads", handler.GetForumThreadsHandler)
	e.GET(ForumPath + "/:slug/users", handler.GetForumUsersHandler)

	e.POST(ThreadPath + "/:slug_or_id/create", handler.CreateThreadPosts)
	e.GET(ThreadPath + "/:slug_or_id/details", handler.GetThreadDetails)
	e.POST(ThreadPath + "/:slug_or_id/details", handler.ChangeThreadDetails)
	e.GET(ThreadPath + "/:slug_or_id/posts", handler.GetThreadPosts)
	e.POST(ThreadPath + "/:slug_or_id/vote", handler.VoteOnThread)

	e.GET(PostPath + "/:id/details", handler.GetPostDetailsHandler)
	e.POST(PostPath + "/:id/details", handler.ChangePostDetailsHandler)

	e.POST(UserPath + "/:nickname/create", handler.CreateUserHandler)
	e.GET(UserPath + "/:nickname/profile", handler.GetUserProfileHandler)
	e.POST(UserPath + "/:nickname/profile", handler.ChangeUserProfileHandler)

	e.POST(ServicePath + "/clear", handler.ClearDBHandler)
	e.GET(ServicePath + "/status", handler.GetServiceStatusHandler)
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func writeWithError(c echo.Context, statusCode int) error {
	errorRes := "Ошибка."
	return c.JSON(statusCode, errorRes)
}