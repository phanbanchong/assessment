package expense

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func GetExpenseHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorf("error: %v", err)
		return c.JSON(http.StatusInternalServerError, Error{Message: "ID is invalid"})
	}
	exp, err := GetExpenseByID(db, id)
	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Error{Message: "Expense not found"})
	case nil:
		return c.JSON(http.StatusOK, exp)
	default:
		return c.JSON(http.StatusInternalServerError, Error{Message: "Unable to scan expense:" + err.Error()})
	}
}

func GetExpensesHandler(c echo.Context) error {
	expenses, err := GetExpenses(db)
	if err != nil {
		log.Errorf("Unable to get expenses from db:" + err.Error())
		return c.JSON(http.StatusInternalServerError, Error{Message: "Unable to get expenses from database:" + err.Error()})
	}
	return c.JSON(http.StatusOK, expenses)
}
