package expense

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func TestUpdateExpenseHandler(t *testing.T) {
	t.Run("Update expense should be success", func(t *testing.T) {
		exp := Expense{
			ID:     1,
			Title:  "title",
			Amount: 1,
			Note:   "note",
			Tags:   []string{"tag1", "tag2"},
		}
		//Mock Database
		var mock sqlmock.Sqlmock
		var err error
		db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		mock.ExpectPrepare("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id = $1").
			ExpectExec().
			WithArgs(exp.ID, exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		//Mock Echo Context
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/expenses/3", strings.NewReader(GoodExpenseJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(exp.ID))

		if err = UpdateExpenseHandler(c); err != nil {
			t.Errorf("should not return error but it got %v", err)
		}

		resp := rec.Body.String()
		want := `{"id":1,"title":"title","amount":1,"note":"note","tags":["tag1","tag2"]}` + "\n"
		if resp != want {
			t.Errorf("response error was not expected got: %s", resp)
		}
	})

	t.Run("Update expense with invalid ID should got error", func(t *testing.T) {
		//Mock Echo Context
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/expenses", strings.NewReader(ExpenseWithIDJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("a")

		err := UpdateExpenseHandler(c)
		if err != nil {
			t.Errorf("should not return error but it got %v", err)
		}
		if c.Response().Status != http.StatusBadRequest {
			t.Errorf("should status bed request but it got %v", c.Response().Status)
		}

		resp := rec.Body.String()
		want := `{"message":"Field ID is invalid"}` + "\n"
		if resp != want {
			t.Errorf("response error was not expected got: %s", resp)
		}
	})
	t.Run("Update expense bad request should be fail", func(t *testing.T) {
		//Mock Echo Context
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/expenses", strings.NewReader(BadExpenseJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(3))

		err := UpdateExpenseHandler(c)
		if err != nil {
			t.Errorf("should not return error but it got %v", err)
		}
		if c.Response().Status != http.StatusBadRequest {
			t.Errorf("should status bed request but it got %v", c.Response().Status)
		}

		resp := rec.Body.String()
		want := `{"message":"code=400, message=Syntax error: offset=35, error=invalid character 'A' looking for beginning of value, internal=invalid character 'A' looking for beginning of value"}` + "\n"
		if resp != want {
			t.Errorf("response error was not expected got: %s", resp)
		}
	})

	t.Run("Update expense and database loss should be fail", func(t *testing.T) {
		exp := Expense{
			ID:     1,
			Title:  "title",
			Amount: 1,
			Note:   "note",
			Tags:   []string{"tag1", "tag2"},
		}
		//Mock Database
		var mock sqlmock.Sqlmock
		var err error
		db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		mock.ExpectPrepare("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id = $1").
			ExpectExec().
			WithArgs(exp.ID, exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags)).
			WillReturnError(sql.ErrConnDone)

		//Mock Echo Context
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/expenses/3", strings.NewReader(GoodExpenseJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(exp.ID))

		err = UpdateExpenseHandler(c)
		if err != nil {
			t.Errorf("should not return error but it got %v", err)
		}
		if c.Response().Status != http.StatusInternalServerError {
			t.Errorf("should status bed request but it got %v", c.Response().Status)
		}

		resp := rec.Body.String()
		want := `{"message":"sql: connection is already closed"}` + "\n"
		if resp != want {
			t.Errorf("response error was not expected got: %s", resp)
		}
	})
}
