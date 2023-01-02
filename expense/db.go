package expense

import (
	"database/sql"
	"os"

	"github.com/labstack/gommon/log"

	"github.com/lib/pq"
)

func InitDB() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
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
		log.Errorf("Insert expense error: %v", err)
		return exp, err
	}
	return exp, nil
}

func UpdateExpense(db *sql.DB, exp Expense) (Expense, error) {
	stmt, err := db.Prepare("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id = $1")
	if err != nil {
		return exp, err
	}

	if _, err := stmt.Exec(exp.ID, exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags)); err != nil {
		log.Errorf("Update expense error: %v", err)
		return exp, err
	}
	return exp, nil
}

func GetExpenses(db *sql.DB) ([]Expense, error) {
	expenses := []Expense{}
	stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses ORDER BY id ASC")
	if err != nil {
		return expenses, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return expenses, err
	}

	for rows.Next() {
		expense := Expense{}
		err := rows.Scan(&expense.ID, &expense.Title, &expense.Amount, &expense.Note, pq.Array(&expense.Tags))
		if err != nil {
			return expenses, err
		}
		expenses = append(expenses, expense)
	}
	return expenses, nil
}
