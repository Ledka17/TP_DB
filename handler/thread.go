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
		return writeWithError(c, 404, "user or forum not found")
	}

	foundThread := h.usecase.GetThreadInDB(thread.Slug)
	emptyThread := model.Thread{}
	if foundThread != emptyThread {
		return c.JSON(409, foundThread)
	}
	return c.JSON(201, h.usecase.CreateThreadInDB(slug, thread))
}

func (h *DataBaseHandler) CreateThreadPosts(c echo.Context) error {
	decoder := json.NewDecoder(c.Request().Body)
	var posts []model.Post
	err := decoder.Decode(&posts)
	checkErr(err)

	slugOrId := c.Param("slug_or_id")

	nicknames := make([]string, 0, len(posts))
	for _, post := range posts {
		nicknames = append(nicknames, post.Author)
	}

	if h.usecase.IsThreadInDB(slugOrId) && h.usecase.IsUsersInDB(nicknames) { // TODO проверка наличия авторов в базе
		if h.usecase.CheckParentPost(posts, slugOrId) {
			return c.JSON(201, h.usecase.CreatePostsInDB(posts, slugOrId))
		}
		return writeWithError(c, 409, "have a conflicts in posts")
	}

	return writeWithError(c, 404, "thread or users not found")
}

func (h *DataBaseHandler) GetThreadDetails(c echo.Context) error {
	slugOrId := c.Param("slug_or_id")

	foundThread := h.usecase.GetThreadInDB(slugOrId)
	emptyThread := model.Thread{}
	if foundThread != emptyThread {
		return c.JSON(200, foundThread)
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

	limit := -1
	since := -1
	sort := ""
	desc := false
	if c.QueryParam("limit") != "" {
		limit, _ = strconv.Atoi(c.QueryParam("limit"))
	}
	if c.QueryParam("since") != "" {
		since, _ = strconv.Atoi(c.QueryParam("since"))
	}
	if c.QueryParam("sort") != "" {
		sort = c.QueryParam("sort")
	}
	if c.QueryParam("desc") != "" {
		desc, _ = strconv.ParseBool(c.QueryParam("desc"))
	}

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

	foundThread := h.usecase.GetThreadInDB(slugOrId)
	emptyThread := model.Thread{}
	if foundThread != emptyThread && h.usecase.IsUserInDB(vote.Nickname, "") {
		return c.JSON(200, h.usecase.VoteForThreadInDB(foundThread, vote))
	}
	return writeWithError(c, 404, "thread or user not found")
}
