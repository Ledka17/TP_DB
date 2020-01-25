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

	if thread.Author == "" {
		return writeWithError(c, 404, "user is empty")
	}

	foundUser := h.usecase.GetUserInDB(thread.Author)
	emptyUser := model.User{}
	foundForum := h.usecase.GetForumInDB(slug)
	emptyForum := model.Forum{}
	if foundForum == emptyForum || foundUser == emptyUser {
		return writeWithError(c, 404, "user or forum not found")
	}

	if thread.Slug != "" {
		foundThread := h.usecase.GetThreadInDB(thread.Slug)
		emptyThread := model.Thread{}
		if foundThread != emptyThread {
			return c.JSON(409, foundThread)
		}
	}
	thread.Forum = foundForum.Slug
	thread.ForumId = foundForum.Id
	thread.UserId = foundUser.Id
	return c.JSON(201, h.usecase.CreateThreadInDB(thread))
}

func (h *DataBaseHandler) CreateThreadPosts(c echo.Context) error {
	decoder := json.NewDecoder(c.Request().Body)
	var posts []model.Post
	err := decoder.Decode(&posts)
	checkErr(err)

	slugOrId := c.Param("slug_or_id")

	posts, err = h.usecase.CreatePostsInDB(posts, slugOrId)
	if err != nil {
		if err.Error() == "conflicts in posts" {
			return writeWithError(c, 409, err.Error())
		}
		return writeWithError(c, 404, err.Error())
		//if h.usecase.CheckParentPost(posts, slugOrId) {
		//	return c.JSON(201, h.usecase.CreatePostsInDB(posts, foundThread))
		//}
		//return writeWithError(c, 409, "have a conflicts in posts")
	}

	return c.JSON(201, posts)
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

	foundThread := h.usecase.GetThreadInDB(slugOrId)
	emptyThread := model.Thread{}
	if foundThread != emptyThread {
		return c.JSON(200, h.usecase.ChangeThreadInDB(threadUpdate, foundThread))
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

	foundThread := h.usecase.GetThreadInDB(slugOrId)
	emptyThread := model.Thread{}
	if foundThread != emptyThread {
		return c.JSON(200, h.usecase.GetPostsInDB(foundThread.Id, limit, since, sort, desc))
	}
	return writeWithError(c, 404, "thread not found")
}

func (h *DataBaseHandler) VoteOnThread(c echo.Context) error {
	slugOrId := c.Param("slug_or_id")
	decoder := json.NewDecoder(c.Request().Body)
	var vote model.Vote
	err := decoder.Decode(&vote)
	checkErr(err)

	err = h.usecase.VoteForThreadInDB(slugOrId, vote)
	if err != nil {
		return writeWithError(c, 404, err.Error())
	}
	thread := h.usecase.GetThreadInDB(slugOrId)
	return c.JSON(200, thread)
}
