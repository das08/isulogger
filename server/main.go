package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/das08/isulogger/parser"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"
)

var (
	db           *sql.DB
	PostgresUser = os.Getenv("POSTGRES_USER")
	PostgresPass = os.Getenv("POSTGRES_PASSWORD")
)

type Contest struct {
	ContestID   int    `json:"contest_id" db:"contest_id"`
	ContestName string `json:"contest_name" db:"contest_name"`
}

type LogEntry struct {
	ID            int       `json:"id" db:"id"`
	ContestID     int       `json:"contest_id" db:"contest_id"`
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
	if entry.ContestID == 0 {
		return false, ""
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	var id int
	err := db.QueryRow("INSERT INTO entry(contest_id, timestamp, score, message) VALUES($1,$2,$3,$4) RETURNING id", entry.ContestID, entry.Timestamp, entry.Score, entry.Message).Scan(&id)
	if err != nil {
		fmt.Println("Error: Create entry failed: ", err)
	}
	if id > 0 {
		return true, fmt.Sprintf("%d", id)
	} else {
		return false, ""
	}
}

func insertContest(contestName string) (bool, string) {
	var id int
	err := db.QueryRow("INSERT INTO contest(contest_name) VALUES($1) RETURNING contest_id", contestName).Scan(&id)
	if err != nil {
		fmt.Println("Error: Create contest failed: ", err)
	}
	if id > 0 {
		return true, fmt.Sprintf("%d", id)
	} else {
		return false, ""
	}
}

func selectContest() []Contest {
	var contest []Contest
	rows, err := db.Query("SELECT * FROM contest ORDER BY contest_id DESC ")
	if err != nil {
		fmt.Println("Error: Get contest failed: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c Contest
		err := rows.Scan(&c.ContestID, &c.ContestName)
		if err != nil {
			fmt.Println("Error: Scan contest failed: ", err)
		} else {
			contest = append(contest, c)
		}
	}
	return contest
}

func selectLogEntry(ContestID int, orderBy string) []LogEntry {
	var entry []LogEntry
	query := "SELECT * FROM entry WHERE contest_id = $1 ORDER BY timestamp asc"
	if orderBy == "desc" {
		query = "SELECT * FROM entry WHERE contest_id = $1 ORDER BY timestamp desc"
	}
	rows, err := db.Query(query, ContestID)
	if err != nil {
		fmt.Println("Error: Get entry failed: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var e LogEntry
		err := rows.Scan(&e.ID, &e.ContestID, &e.Timestamp, &e.Score, &e.Message, &e.AccessLogPath, &e.SlowLogPath, &e.ImagePath)
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

	e.GET("/new_log", createLogEntry)
	e.POST("/new_contest", createContest)
	e.GET("/get", getLogEntry)
	e.GET("/get_contest", getContest)

	e.GET("/parsed/alp/:id", parsedAlp)

	e.Static("/log", "log")

	e.Logger.Fatal(e.Start(":8082"))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func createLogEntry(c echo.Context) error {
	_contestID := c.QueryParam("contest_id")
	_score := c.QueryParam("score")
	message := c.QueryParam("message")

	contestID, err := strconv.Atoi(_contestID)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid contest_id")
	}
	score, err := strconv.Atoi(_score)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid score")
	}

	entry := LogEntry{}
	entry.ContestID = contestID
	entry.Timestamp = time.Now()
	entry.Score = score
	entry.Message = message

	if ok, id := insertLogEntry(&entry); !ok {
		return echo.ErrInternalServerError
	} else {
		return c.JSON(http.StatusOK, id)
	}
}

func createContest(c echo.Context) error {
	type postData struct {
		ContestName string `json:"contest_name"`
	}
	data := postData{}
	if err := c.Bind(&data); err != nil {
		return echo.ErrInternalServerError
	}

	fmt.Println("contestname: ", data.ContestName)

	if data.ContestName == "" {
		return echo.ErrBadRequest
	}

	if ok, id := insertContest(data.ContestName); !ok {
		return echo.ErrInternalServerError
	} else {
		return c.JSON(http.StatusOK, id)
	}
}

func getLogEntry(c echo.Context) error {
	contestIDRaw := c.QueryParam("contest_id")
	sortRaw := c.QueryParam("sort")
	orderBy := "desc"
	if contestIDRaw == "" {
		return echo.ErrBadRequest
	}
	if sortRaw == "asc" {
		orderBy = "asc"
	}

	contestID, err := strconv.Atoi(contestIDRaw)
	if err != nil {
		return echo.ErrBadRequest
	}
	entry := selectLogEntry(contestID, orderBy)
	return c.JSON(http.StatusOK, entry)
}

func parsedAlp(c echo.Context) error {
	id := c.Param("id")

	var alpPath string
	err := db.QueryRow("SELECT access_log_path FROM entry WHERE id = $1", id).Scan(&alpPath)
	if err != nil {
		return echo.ErrNotFound
	}

	// FIXME: do not hardcode /log path
	alpString, err := ioutil.ReadFile("log/" + alpPath)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	alp, err := parser.ParseAlpData(string(alpString))
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, alp)
}

func getContest(c echo.Context) error {
	contest := selectContest()
	return c.JSON(http.StatusOK, contest)
}
