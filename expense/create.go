package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateExpenseHandler(c echo.Context) error {
	exp := Expense{}
	err := c.Bind(&exp)
	if exp.ID != 0 {
		return c.JSON(http.StatusBadRequest, Error{Message: "Field ID is invalid"})
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}

	exp, err = CreateExpense(db, exp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, exp)
}
