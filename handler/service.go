package handler

import (
	"github.com/labstack/echo"
)

func (h *DataBaseHandler) ClearDBHandler(c echo.Context) error {
	h.usecase.CleanUp()
	return c.JSON(200, "ok")
}

func (h *DataBaseHandler) GetServiceStatusHandler(c echo.Context) error {
	return c.JSON(200, h.usecase.GetStatusDB())
}
