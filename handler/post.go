package handler

import (
	"encoding/json"
	"github.com/Ledka17/TP_DB/model"
	"github.com/labstack/echo"
	"strconv"
	"strings"
)

func (h *DataBaseHandler) GetPostDetailsHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	checkErr(err)
	related := strings.Split(c.QueryParam("related"), ",")

	post := h.usecase.GetPostInDB(id)
	emptyPost := model.Post{}
	if post != emptyPost {
		postFull := model.PostFull{
			Author: nil,
			Forum:  nil,
			Post:   &post,
			Thread: nil,
		}

		if checkInRelated("user", related) {
			user := h.usecase.GetUserInDB(post.Author)
			postFull.Author = &user
		}
		if checkInRelated("forum", related) {
			forum := h.usecase.GetForumInDB(post.Forum)
			postFull.Forum = &forum
		}
		if checkInRelated("thread", related) {
			thread := h.usecase.GetThreadById(int(post.ThreadId))
			postFull.Thread = &thread
		}
		return c.JSON(200, postFull)
	}
	return writeWithError(c, 404, "post not found")
}

func (h *DataBaseHandler) ChangePostDetailsHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	checkErr(err)
	decoder := json.NewDecoder(c.Request().Body)
	var post model.PostUpdate
	err = decoder.Decode(&post)
	checkErr(err)

	changedPost, _ := h.usecase.ChangePostInDB(id, post)
	if changedPost.Author == "" {
		return writeWithError(c, 404, "post not found")
	}

	return c.JSON(200, changedPost)
}

func checkInRelated(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}