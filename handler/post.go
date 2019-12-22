package handler

import (
	"encoding/json"
	"github.com/Ledka17/TP_DB/model"
	"github.com/labstack/echo"
	"strconv"
)

func (h *DataBaseHandler) GetPostDetailsHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	checkErr(err)
	related := c.QueryParams()["related"]

	if h.usecase.IsPostInDB(id) {
		return c.JSON(200, h.usecase.GetPostInDB(id, related))
	}
	return writeWithError(c, 404, "post not found")
}

func (h *DataBaseHandler) ChangePostDetailsHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	checkErr(err)
	if h.usecase.IsPostInDB(id) {
		decoder := json.NewDecoder(c.Request().Body)
		var post model.PostUpdate
		err := decoder.Decode(&post)
		checkErr(err)

		return c.JSON(200, h.usecase.ChangePostInDB(id, post))
	}
	return writeWithError(c, 404, "post not found")
}