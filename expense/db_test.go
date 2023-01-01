package expense

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

func TestInitDB(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Now we execute our method
	expectDB := InitDB()
	if expectDB != nil && expectDB != db {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// Make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestCreateExpense(t *testing.T) {
	exp := Expense{
		Title:  "title",
		Amount: 1,
		Note:   "note",
		Tags:   []string{"tag1", "tag2"},
	}
	db = &sql.DB{}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("INSERT INTO expenses").
		WithArgs("title", 1.0, "note", pq.Array([]string{"tag1", "tag2"})).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))

	// Now we execute our method
	if _, err = CreateExpense(db, exp); err != nil {
		t.Errorf("error was not expected while insert expense: %s", err)
	}

	// Make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetExpense(t *testing.T) {
	ID := 2

	db = &sql.DB{}
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
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

	// Now we execute our method
	if _, err = GetExpenseByID(db, ID); err != nil {
		t.Errorf("error was not expected while select expense: %s", err)
	}

	// Make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
