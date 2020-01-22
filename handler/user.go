package handler

import (
	"encoding/json"
	"github.com/Ledka17/TP_DB/model"
	"github.com/labstack/echo"
)

func (h *DataBaseHandler) GetUserProfileHandler(c echo.Context) error {
	nickname := c.Param("nickname")

	user := h.usecase.GetUserInDB(nickname)
	emptyUser := model.User{}

	if user != emptyUser {
		return c.JSON(200, user)
	}
	return writeWithError(c, 404, "user not found")
}

func (h *DataBaseHandler) ChangeUserProfileHandler(c echo.Context) error {
	nickname := c.Param("nickname")
	decoder := json.NewDecoder(c.Request().Body)
	var userUpdate model.UserUpdate
	err := decoder.Decode(&userUpdate)
	checkErr(err)

	users := h.usecase.GetUsersInDB(nickname, userUpdate.Email)

	if len(users) == 2 || len(users) == 1 && users[0].Email == userUpdate.Email {
		return writeWithError(c, 409, "user already exists")
	}
	if len(users) == 1 {
		return c.JSON(200, h.usecase.ChangeUserInDB(nickname, userUpdate))
	}
	return writeWithError(c, 404, "user not found")
}

func (h *DataBaseHandler) CreateUserHandler(c echo.Context) error {
	nickname := c.Param("nickname")

	decoder := json.NewDecoder(c.Request().Body)
	var user model.User
	err := decoder.Decode(&user)
	checkErr(err)

	foundUsers := h.usecase.GetUsersInDB(nickname, user.Email)
	if len(foundUsers) != 0 {
		return c.JSON(409, foundUsers)
	}

	return c.JSON(201, h.usecase.СreateUserInDB(nickname, user))
}

