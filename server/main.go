package main

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

var (
	db           *sql.DB
	PostgresUser = os.Getenv("POSTGRES_USER")
	PostgresPass = os.Getenv("POSTGRES_PASSWORD")
)

type LogEntry struct {
	ID            int       `json:"id" db:"id"`
	IsuconID      int       `json:"isucon_id" db:"isucon_id"`
	Timestamp     time.Time `json:"timestamp" db:"timestamp"`
	Score         int       `json:"score" db:"score"`
	Message       string    `json:"message" db:"message"`
	AccessLogPath string    `json:"access_log_path" db:"access_log_path"`
	SlowLogPath   string    `json:"slow_log_path" db:"slow_log_path"`
	ImagePath     string    `json:"image_path" db:"image_path"`
}

func initializeDB() *sql.DB {
	_db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@isulogger-db:5432/isulogger?sslmode=disable", PostgresUser, PostgresPass))
	if err != nil {
		fmt.Println(err)
	}
	return _db
}

// Add a new log entry to the database
func insertLogEntry(entry *LogEntry) (bool, string) {
	if entry.IsuconID == 0 {
		return false, ""
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	var id int
	err := db.QueryRow("INSERT INTO entry(isucon_id, timestamp, score) VALUES($1,$2,$3) RETURNING id", entry.IsuconID, entry.Timestamp, entry.Score).Scan(&id)
	if err != nil {
		fmt.Println("Error: Create entry failed: ", err)
	}
	if id > 0 {
		return true, fmt.Sprintf("%d", id)
	} else {
		return false, ""
	}
}

func selectLogEntry(isuconID int, orderBy string) []LogEntry {
	var entry []LogEntry
	query := "SELECT * FROM entry WHERE isucon_id = $1 ORDER BY timestamp asc"
	if orderBy == "desc" {
		query = "SELECT * FROM entry WHERE isucon_id = $1 ORDER BY timestamp desc"
	}
	rows, err := db.Query(query, isuconID)
	if err != nil {
		fmt.Println("Error: Get entry failed: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var e LogEntry
		err := rows.Scan(&e.ID, &e.IsuconID, &e.Timestamp, &e.Score, &e.Message, &e.AccessLogPath, &e.SlowLogPath, &e.ImagePath)
		if err != nil {
			fmt.Println("Error: Scan entry failed: ", err)
		} else {
			entry = append(entry, e)
		}
	}
	return entry
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
	e.GET("/get", getLogEntry)

	e.Static("/log", "log")

	e.Logger.Fatal(e.Start(":8082"))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func createLogEntry(c echo.Context) error {
	entry := LogEntry{}
	entry.IsuconID = 1
	entry.Timestamp = time.Now()
	entry.Score = 777

	if ok, id := insertLogEntry(&entry); !ok {
		return echo.ErrInternalServerError
	} else {
		return c.JSON(http.StatusOK, id)
	}
}

func getLogEntry(c echo.Context) error {
	isuconIDRaw := c.QueryParam("isucon_id")
	sortRaw := c.QueryParam("sort")
	orderBy := "desc"
	if isuconIDRaw == "" {
		return echo.ErrBadRequest
	}
	if sortRaw == "asc" {
		orderBy = "asc"
	}

	isuconID, err := strconv.Atoi(isuconIDRaw)
	if err != nil {
		return echo.ErrBadRequest
	}
	entry := selectLogEntry(isuconID, orderBy)
	return c.JSON(http.StatusOK, entry)
}
