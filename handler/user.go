package handler

import (
	"encoding/json"
	"github.com/Ledka17/TP_DB/model"
	"github.com/labstack/echo"
)

func (h *DataBaseHandler) GetUserProfileHandler(c echo.Context) error {
	nickname := c.Param("nickname")
	if h.usecase.IsUserInDB(nickname, "") {
		return c.JSON(200, h.usecase.GetUserInDB(nickname, ""))
	}
	return writeWithError(c, 404)
}

func (h *DataBaseHandler) ChangeUserProfileHandler(c echo.Context) error {
	nickname := c.Param("nickname")
	decoder := json.NewDecoder(c.Request().Body)
	var userUpdate model.UserUpdate
	err := decoder.Decode(&userUpdate)
	checkErr(err)

	if h.usecase.IsUserInDB(nickname, "") {
		if h.usecase.IsUserInDB("", userUpdate.Email) {
			return writeWithError(c, 409)
		}

		return c.JSON(200, h.usecase.GetUserInDB(nickname, ""))
	}
	return writeWithError(c, 404)
}

func (h *DataBaseHandler) CreateUserHandler(c echo.Context) error {
	nickname := c.Param("nickname")

	decoder := json.NewDecoder(c.Request().Body)
	var user model.User
	err := decoder.Decode(&user)
	checkErr(err)

	if h.usecase.IsUserInDB(nickname, user.Email) {
		return c.JSON(409, h.usecase.GetUserInDB(nickname, user.Email))
	}

	return c.JSON(201, h.usecase.СreateUserInDB(nickname, user))
}

