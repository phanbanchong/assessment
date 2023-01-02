package expense

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *handler) UpdateExpenseHandler(c echo.Context) error {
	exp := Expense{}
	err := c.Bind(&exp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}
	exp.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Field ID is invalid"})
	}

	exp, err = UpdateExpense(h.DB, exp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, exp)
}
