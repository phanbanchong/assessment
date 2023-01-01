package expense

import (
	"database/sql"
	"os"

	"github.com/labstack/gommon/log"

	"github.com/lib/pq"
)

var db *sql.DB

func InitDB() *sql.DB {
	var err error

	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connection to database error", err)
	}

	command := "CREATE TABLE IF NOT EXISTS expenses ( id SERIAL PRIMARY KEY, title TEXT, amount FLOAT, note TEXT, tags TEXT[]);"

	_, err = db.Exec(string(command))
	if err != nil {
		log.Fatal("Unable to create table", err)
	} else {
		log.Printf("Executed: %s", command)
	}
	return db
}

func GetExpenseByID(db *sql.DB, id int) (Expense, error) {
	exp := Expense{}
	stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")
	if err != nil {
		return exp, err
	}

	rows := stmt.QueryRow(id)
	if rows.Err() != nil {
		return exp, rows.Err()
	}
	err = rows.Scan(&exp.ID, &exp.Title, &exp.Amount, &exp.Note, pq.Array(&exp.Tags))
	return exp, err
}

func CreateExpense(db *sql.DB, exp Expense) (Expense, error) {
	row := db.QueryRow("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id", exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags))
	err := row.Scan(&exp.ID)
	if err != nil {
		log.Errorf("Error insert expense error: %v", err)
		return exp, err
	}
	return exp, nil
}
