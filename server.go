package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/phanbanchong/assessment/expense"
	"github.com/phanbanchong/assessment/health"
)

func ContextDB(db *sql.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	}
}

func main() {
	db := expense.InitDB()
	h := expense.NewApplication(db)
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", health.GetHealthHandler)

	e.GET("/expenses", h.GetExpensesHandler)
	e.GET("/expenses/:id", h.GetExpenseHandler)
	e.POST("/expenses", h.CreateExpenseHandler)
	e.PUT("/expenses/:id", h.UpdateExpenseHandler)

	go func(e *echo.Echo) {
		if err := e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server")
		}
	}(e)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
