package handler

import (
	"encoding/json"
	"github.com/Ledka17/TP_DB/model"
	"github.com/labstack/echo"
	"strconv"
)

func (h *DataBaseHandler) CreateForumHandler (c echo.Context) error {
	decoder := json.NewDecoder(c.Request().Body)
	var forum model.Forum
	err := decoder.Decode(&forum)
	checkErr(err)
	if forum.User == "" || !h.usecase.IsUserInDB(forum.User, "") {
		return writeWithError(c, 404, "User not found")
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
	return writeWithError(c, 404, "forum not found")
}

func (h *DataBaseHandler) GetForumThreadsHandler(c echo.Context) error {
	limit := -1
	since := ""
	desc := false
	// TODO get params
	if c.QueryParam("limit") != "" {
		limit, _ = strconv.Atoi(c.QueryParam("limit"))
	}
	if c.QueryParam("since") != "" {
		since = c.QueryParam("since")
	}
	if c.QueryParam("desc") != "" {
		desc, _ = strconv.ParseBool(c.QueryParam("desc"))
	}

	forumSlug := c.Param("slug")

	if h.usecase.IsForumInDB(forumSlug) {
		return c.JSON(200, h.usecase.GetThreadsForumInDB(forumSlug, limit, since, desc))
	}
	return writeWithError(c, 404, "forum not found")
}

func (h *DataBaseHandler) GetForumUsersHandler(c echo.Context) error {
	limit := -1
	since := 0
	desc := false
	// TODO get params
	if c.QueryParam("limit") != "" {
		limit, _ = strconv.Atoi(c.QueryParam("limit"))
	}
	if c.QueryParam("since") != "" {
		since, _ = strconv.Atoi(c.QueryParam("since"))
	}
	if c.QueryParam("desc") != "" {
		desc, _ = strconv.ParseBool(c.QueryParam("desc"))
	}

	forumSlug := c.Param("slug")

	if h.usecase.IsForumInDB(forumSlug) {
		users := h.usecase.GetForumUsersInDB(forumSlug, limit, since, desc)
		return c.JSON(200, users)
	}
	return writeWithError(c, 404, "forum not found")
}
