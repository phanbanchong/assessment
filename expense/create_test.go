package expense

import (
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
		var mock sqlmock.Sqlmock
		var err error
		db, mock, err = sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		mock.ExpectQuery("INSERT INTO expenses").
			WithArgs("title", 1.0, "note", pq.Array([]string{"tag1", "tag2"})).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))

		//Mock Echo Context
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(GoodExpenseJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if err = CreateExpenseHandler(c); err != nil {
			t.Errorf("should not return error but it got %v", err)
		}
	})

	t.Run("Create expense with ID should be fail", func(t *testing.T) {
		//Mock Database
		var mock sqlmock.Sqlmock
		var err error
		db, mock, err = sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		mock.ExpectQuery("INSERT INTO expenses").
			WithArgs("title", 1.0, "note", pq.Array([]string{"tag1", "tag2"})).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))

		//Mock Echo Context
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(ExpenseWithIDJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = CreateExpenseHandler(c)
		if err != nil {
			t.Errorf("should not return error but it got %v", err)
		}
		if c.Response().Status != http.StatusBadRequest {
			t.Errorf("should status bed request but it got %v", c.Response().Status)
		}
	})
	t.Run("Create expense bad request should be fail", func(t *testing.T) {
		//Mock Database
		var mock sqlmock.Sqlmock
		var err error
		db, mock, err = sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		mock.ExpectQuery("INSERT INTO expenses").
			WithArgs("title", 1.0, "note", pq.Array([]string{"tag1", "tag2"})).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))

		//Mock Echo Context
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(BadExpenseJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = CreateExpenseHandler(c)
		if err != nil {
			t.Errorf("should not return error but it got %v", err)
		}
		if c.Response().Status != http.StatusBadRequest {
			t.Errorf("should status bed request but it got %v", c.Response().Status)
		}
	})
}
