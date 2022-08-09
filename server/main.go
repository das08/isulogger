package main

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var (
	db           *sql.DB
	PostgresUser = os.Getenv("POSTGRES_USER")
	PostgresPass = os.Getenv("POSTGRES_PASSWORD")
)

type LogEntry struct {
	IsuconID  int       `json:"isucon_id" db:"isucon_id"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	Score     int       `json:"score" db:"score"`
}

func initializeDB() *sql.DB {
	_db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@isulogger-db:5432/isulogger?sslmode=disable", PostgresUser, PostgresPass))
	if err != nil {
		fmt.Println(err)
	}
	return _db
}

// Add a new log entry to the database
func insertLogEntry(entry *LogEntry) (bool, int) {
	var id int
	err := db.QueryRow("INSERT INTO entry(isucon_id, timestamp, score) VALUES($1,$2,$3) RETURNING id", entry.IsuconID, entry.Timestamp, entry.Score).Scan(&id)
	if err != nil {
		fmt.Println("Error: Create entry failed: ", err)
	}
	return id > 0, id
}

func main() {
	// Initialize DB connection
	db = initializeDB()
	defer db.Close()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORS())

	e.GET("/", hello)

	e.GET("/new", createLogEntry)

	e.Logger.Fatal(e.Start(":8082"))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func createLogEntry(c echo.Context) error {
	entry := new(LogEntry)
	entry.IsuconID = 1
	entry.Timestamp = time.Now()
	entry.Score = 777

	if ok, id := insertLogEntry(entry); !ok {
		return echo.ErrInternalServerError
	} else {
		return c.JSON(http.StatusOK, id)
	}
}
