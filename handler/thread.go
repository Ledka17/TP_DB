package handler

import (
	"encoding/json"
	"github.com/Ledka17/TP_DB/model"
	"github.com/labstack/echo"
	"strconv"
)

func (h *DataBaseHandler) CreateThreadHandler(c echo.Context) error {
	slug := c.Param("slug")
	decoder := json.NewDecoder(c.Request().Body)
	var thread model.Thread
	err := decoder.Decode(&thread)
	checkErr(err)

	if thread.Author == "" || !h.usecase.IsUserInDB(thread.Author, "") || !h.usecase.IsForumInDB(slug) {
		return writeWithError(c, 404, "user not found")
	}

	if h.usecase.IsThreadInDB(thread.Slug) {
		return c.JSON(409, h.usecase.GetThreadInDB(thread.Slug))
	}
	return c.JSON(201, h.usecase.CreateThreadInDB(slug, thread))
}

func (h *DataBaseHandler) CreateThreadPosts(c echo.Context) error {
	decoder := json.NewDecoder(c.Request().Body)
	var posts []model.Post
	err := decoder.Decode(&posts)
	checkErr(err)

	slugOrId := c.Param("slug_or_id")

	if h.usecase.IsThreadInDB(slugOrId) {
		if h.usecase.CheckParentPost(posts) {
			return c.JSON(201, h.usecase.CreatePostsInDB(posts))
		}
		return writeWithError(c, 409, "have a conflicts in posts")
	}

	return writeWithError(c, 404, "thread not found")
}

func (h *DataBaseHandler) GetThreadDetails(c echo.Context) error {
	slugOrId := c.Param("slug_or_id")

	if h.usecase.IsThreadInDB(slugOrId) {
		return c.JSON(200, h.usecase.GetThreadInDB(slugOrId))
	}

	return writeWithError(c, 404, "thread not found")
}

func (h *DataBaseHandler) ChangeThreadDetails(c echo.Context) error {
	slugOrId := c.Param("slug_or_id")
	decoder := json.NewDecoder(c.Request().Body)
	var threadUpdate model.ThreadUpdate
	err := decoder.Decode(&threadUpdate)
	checkErr(err)

	if h.usecase.IsThreadInDB(slugOrId) {
		return c.JSON(200, h.usecase.ChangeThreadInDB(threadUpdate, slugOrId))
	}

	return writeWithError(c, 404, "thread not found")
}

func (h *DataBaseHandler) GetThreadPosts(c echo.Context) error {
	slugOrId := c.Param("slug_or_id")
	vars := c.QueryParams()
	limit, _ := strconv.Atoi(vars["limit"][0])
	since, _ := strconv.Atoi(vars["since"][0])
	sort := vars["sort"][0]
	desc, _ := strconv.ParseBool(vars["desc"][0])

	if h.usecase.IsThreadInDB(slugOrId) {
		return c.JSON(200, h.usecase.GetPostsInDB(slugOrId, limit, since, sort, desc))
	}
	return writeWithError(c, 404, "thread not found")
}

func (h *DataBaseHandler) VoteOnThread(c echo.Context) error {
	slugOrId := c.Param("slug_or_id")
	decoder := json.NewDecoder(c.Request().Body)
	var vote model.Vote
	err := decoder.Decode(&vote)
	checkErr(err)

	if h.usecase.IsThreadInDB(slugOrId) && h.usecase.IsUserInDB(vote.Nickname, "") {
		return c.JSON(200, h.usecase.VoteForThreadInDB(slugOrId, vote))
	}
	return writeWithError(c, 404, "thread or user not found")
}
