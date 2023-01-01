package expense

import (
	"database/sql"
	"log"
	"os"

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

func CreateExpense(db *sql.DB, exp Expense) (Expense, error) {
	row := db.QueryRow("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id", exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags))
	err := row.Scan(&exp.ID)
	if err != nil {
		log.Fatal("Error insert expense", err)
		return exp, err
	}
	return exp, nil
}
