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
	foundUser := h.usecase.GetUserInDB(forum.User)
	emptyUser := model.User{}
	if foundUser == emptyUser {
		return writeWithError(c, 404, "User not found")
	}
	foundForum := h.usecase.GetForumInDB(forum.Slug)
	emptyForum := model.Forum{}
	if foundForum != emptyForum {
		return c.JSON(409, foundForum)
	}
	forum.User = foundUser.Nickname
	return c.JSON(201, h.usecase.CreateForumInDB(forum))
}

func (h *DataBaseHandler) GetForumDetailsHandler(c echo.Context) error {
	slug := c.Param("slug")
	foundForum := h.usecase.GetForumInDB(slug)
	emptyForum := model.Forum{}

	if foundForum != emptyForum {
		return c.JSON(200, foundForum)
	}
	return writeWithError(c, 404, "forum not found")
}

func (h *DataBaseHandler) GetForumThreadsHandler(c echo.Context) error {
	limit := -1
	since := ""
	desc := false
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
	since := ""
	desc := false
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
		users := h.usecase.GetForumUsersInDB(forumSlug, limit, since, desc)
		return c.JSON(200, users)
	}
	return writeWithError(c, 404, "forum not found")
}
