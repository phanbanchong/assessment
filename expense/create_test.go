//go:build unit

package expense

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

var (
	ExpenseWithIDJSON = `{
		"id": 1,
		"title": "title",
		"amount": 1,
		"note": "note",
		"tags": [
			"tag1",
			"tag2"
		]
	}`
	BadExpenseJSON = `{
		"title": "title",
		"amount": A,
		"note": "note",
		"tags": [
			"tag1",
			"tag2"
		]
	}`
	GoodExpenseJSON = `{
		"title": "title",
		"amount": 1,
		"note": "note",
		"tags": [
			"tag1",
			"tag2"
		]
	}`
)

func TestCreateExpenseHandler(t *testing.T) {
	t.Run("Create expense should be success", func(t *testing.T) {
		//Mock Database
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		mock.ExpectQuery("INSERT INTO expenses").
			WithArgs("title", 1.0, "note", pq.Array([]string{"tag1", "tag2"})).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))

		h := NewApplication(db)
		//Mock Echo Context
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(GoodExpenseJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if err = h.CreateExpenseHandler(c); err != nil {
			t.Errorf("should not return error but it got %v", err)
		}
	})

	t.Run("Create expense with invalid ID should got error", func(t *testing.T) {
		h := NewApplication(nil)
		//Mock Echo Context
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(ExpenseWithIDJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.CreateExpenseHandler(c)
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
	t.Run("Create expense bad request should be fail", func(t *testing.T) {
		h := NewApplication(nil)
		//Mock Echo Context
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(BadExpenseJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.CreateExpenseHandler(c)
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

	t.Run("Create expense with should be fail", func(t *testing.T) {
		//Mock Database
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		mock.ExpectQuery("INSERT INTO expenses").
			WithArgs("title", 1.0, "note", pq.Array([]string{"tag1", "tag2"})).
			WillReturnError(sql.ErrConnDone)
		h := NewApplication(db)

		//Mock Echo Context
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(GoodExpenseJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = h.CreateExpenseHandler(c)
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
