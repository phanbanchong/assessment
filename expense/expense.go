package expense

import "database/sql"

type handler struct {
	DB *sql.DB
}

func NewApplication(db *sql.DB) *handler {
	return &handler{db}
}

type Expense struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

type Error struct {
	Message string `json:"message"`
}
