package health

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func GetHealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, Health{Status: "UP", Timestamp: time.Now().Format(time.RFC3339)})
}
