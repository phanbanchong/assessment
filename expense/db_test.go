//go:build unit

package expense

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

func TestCreateExpense(t *testing.T) {
	exp := Expense{
		Title:  "title",
		Amount: 1,
		Note:   "note",
		Tags:   []string{"tag1", "tag2"},
	}
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

func TestUpdateExpense(t *testing.T) {
	exp := Expense{
		ID:     3,
		Title:  "title",
		Amount: 3.0,
		Note:   "note",
		Tags:   []string{"tag1", "tag2"},
	}
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id = $1").
		ExpectExec().
		WithArgs(exp.ID, exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Now we execute our method
	if _, err = UpdateExpense(db, exp); err != nil {
		t.Errorf("error was not expected while select expense: %s", err)
	}

	// Make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
