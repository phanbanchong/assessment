package expense

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func TestGetExpenseHandler(t *testing.T) {
	t.Run("Get expense by ID should be success", func(t *testing.T) {
		//Mock Database
		ID := 2
		var mock sqlmock.Sqlmock
		var err error
		db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
			AddRow(2, "expense 2", 2.0, "note 2", pq.Array([]string{"tag1", "tag2"}))

		mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1").
			ExpectQuery().
			WithArgs(ID).
			WillReturnRows(mockRows)

		//Mock Echo Context
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/expenses/"+strconv.Itoa(ID), nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(ID))

		if err = GetExpenseHandler(c); err != nil {
			t.Errorf("should not return error but it got %v", err)
		}

		resp := rec.Body.String()
		want := `{"id":2,"title":"expense 2","amount":2,"note":"note 2","tags":["tag1","tag2"]}` + "\n"
		if resp != want {
			t.Errorf("response error was not expected got: %s", resp)
		}

	})
	t.Run("Get expense by ID and ErrNoRows should got error", func(t *testing.T) {
		//Mock Database
		ID := 2
		var mock sqlmock.Sqlmock
		var err error
		db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1").
			ExpectQuery().
			WithArgs(ID).
			WillReturnError(sql.ErrNoRows)

		//Mock Echo Context
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/expenses/"+strconv.Itoa(ID), nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(ID))

		if err = GetExpenseHandler(c); err != nil {
			t.Errorf("should not return error but it got %v", err)
		}

		resp := rec.Body.String()
		want := `{"message":"Expense not found"}` + "\n"
		if resp != want {
			t.Errorf("response error was not expected got: %s", resp)
		}
	})
	t.Run("Get expense by ID anoter error from database should got error", func(t *testing.T) {
		//Mock Database
		ID := 2
		var mock sqlmock.Sqlmock
		var err error
		db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1").
			ExpectQuery().
			WithArgs(ID).
			WillReturnError(sql.ErrConnDone)

		//Mock Echo Context
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/expenses/"+strconv.Itoa(ID), nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(ID))

		if err = GetExpenseHandler(c); err != nil {
			t.Errorf("should not return error but it got %v", err)
		}

		resp := rec.Body.String()
		want := `{"message":"Unable to scan expense:sql: connection is already closed"}` + "\n"
		if resp != want {
			t.Errorf("response error was not expected got: %s", resp)
		}
	})
	t.Run("Get expense by invalid ID should got error", func(t *testing.T) {
		//Mock Echo Context
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/expenses/a", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("a")

		if err := GetExpenseHandler(c); err != nil {
			t.Errorf("should not return error but it got %v", err)
		}

		resp := rec.Body.String()
		want := `{"message":"ID is invalid"}` + "\n"
		if resp != want {
			t.Errorf("response error was not expected got: %s", resp)
		}
	})

	t.Run("Get all expense should be success", func(t *testing.T) {
		//Mock Database
		var mock sqlmock.Sqlmock
		var err error
		db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
			AddRow(1, "expense 2", 1.0, "note 1", pq.Array([]string{"tag1", "tag2"})).
			AddRow(2, "expense 2", 2.0, "note 2", pq.Array([]string{"tag1", "tag2"}))

		mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses ORDER BY id ASC").
			ExpectQuery().
			WillReturnRows(mockRows)

		//Mock Echo Context
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		if err = GetExpensesHandler(c); err != nil {
			t.Errorf("should not return error but it got %v", err)
		}

		resp := rec.Body.String()
		want := `[{"id":1,"title":"expense 2","amount":1,"note":"note 1","tags":["tag1","tag2"]},{"id":2,"title":"expense 2","amount":2,"note":"note 2","tags":["tag1","tag2"]}]` + "\n"
		if resp != want {
			t.Errorf("response error was not expected got: %s", resp)
		}
	})
	t.Run("Get all expense should got error", func(t *testing.T) {
		//Mock Database
		var mock sqlmock.Sqlmock
		var err error
		db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses ORDER BY id ASC").
			ExpectQuery().
			WillReturnError(sql.ErrConnDone)

		//Mock Echo Context
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		if err = GetExpensesHandler(c); err != nil {
			t.Errorf("should not return error but it got %v", err)
		}

		resp := rec.Body.String()
		want := `{"message":"Unable to get expenses from database:sql: connection is already closed"}` + "\n"
		if resp != want {
			t.Errorf("response error was not expected got: %s", resp)
		}
	})
}
