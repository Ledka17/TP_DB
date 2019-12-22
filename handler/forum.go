package handler

import (
	"encoding/json"
	"github.com/Ledka17/TP_DB/model"
	"github.com/labstack/echo"
	"log"
	"strconv"
)

func (h *DataBaseHandler) CreateForumHandler (c echo.Context) error {
	decoder := json.NewDecoder(c.Request().Body)
	var forum model.Forum
	err := decoder.Decode(&forum)
	checkErr(err)

	log.Print(forum)

	if forum.User == "" || !h.usecase.IsUserInDB(forum.User, "") {
		return c.JSON(404, "User not found")
	}
	if h.usecase.IsForumInDB(forum.Slug) {
		return c.JSON(409, h.usecase.GetForumInDB(forum.Slug))
	}
	return c.JSON(201, h.usecase.CreateForumInDB(forum))
}

func (h *DataBaseHandler) GetForumDetailsHandler(c echo.Context) error {
	slug := c.Param("slug")
	if h.usecase.IsForumInDB(slug) {
		return c.JSON(200, h.usecase.GetForumInDB(slug))
	}
	return writeWithError(c, 404)
}

func (h *DataBaseHandler) GetForumThreadsHandler(c echo.Context) error {
	queryParams := c.QueryParams()
	limit, _ := strconv.Atoi(queryParams["limit"][0])
	since := queryParams["since"][0]
	desc, _ := strconv.ParseBool(queryParams["desc"][0])

	forumSlug := c.Param("slug")

	if h.usecase.IsForumInDB(forumSlug) {
		return c.JSON(200, h.usecase.GetThreadsForumInDB(forumSlug, limit, since, desc))
	}
	return writeWithError(c, 404)
}

func (h *DataBaseHandler) GetForumUsersHandler(c echo.Context) error {
	queryParams := c.QueryParams()
	limit, _ := strconv.Atoi(queryParams["limit"][0])
	since := queryParams["since"][0]
	desc, _ := strconv.ParseBool(queryParams["desc"][0])

	slug := c.Param("slug")

	if h.usecase.IsForumInDB(slug) {
		return c.JSON(200, h.usecase.GetForumUsersInDB(slug, limit, since, desc))
	}
	return writeWithError(c, 404)
}
