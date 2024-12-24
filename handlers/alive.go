package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func AliveCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "I'm alive!"})
}
