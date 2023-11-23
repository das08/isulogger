package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"
)

var (
	db        *sql.DB
	DB_PATH   string = os.Getenv("DB_PATH")
	secretKey string
)

type Contest struct {
	ContestID   int    `json:"contest_id" db:"contest_id"`
	ContestName string `json:"contest_name" db:"contest_name"`
}

type LogEntry struct {
	ID            int       `json:"id" db:"id"`
	ContestID     int       `json:"contest_id" db:"contest_id"`
	Timestamp     time.Time `json:"timestamp" db:"timestamp"`
	BranchName    string    `json:"branch_name" db:"branch_name"`
	Score         int       `json:"score" db:"score"`
	Message       string    `json:"message" db:"message"`
	AccessLogPath string    `json:"access_log_path" db:"access_log_path"`
	SlowLogPath   string    `json:"slow_log_path" db:"slow_log_path"`
	ImagePath     string    `json:"image_path" db:"image_path"`
}

func initializeDB() *sql.DB {
	_db, err := sql.Open("sqlite3", fmt.Sprintf("file://%s", DB_PATH))
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

	result, err := db.Exec("INSERT INTO entry(contest_id, timestamp, branch_name, score, message) VALUES(?,?,?,?,?)", entry.ContestID, entry.Timestamp, entry.BranchName, entry.Score, entry.Message)
	if err != nil {
		fmt.Println("Error: Create entry failed: ", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error: Get last insert id failed: ", err)
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
		err := rows.Scan(&e.ID, &e.ContestID, &e.Timestamp, &e.BranchName, &e.Score, &e.Message, &e.AccessLogPath, &e.SlowLogPath, &e.ImagePath)
		if err != nil {
			fmt.Println("Error: Scan entry failed: ", err)
		} else {
			entry = append(entry, e)
		}
	}
	return entry
}

func deleteLogByID(entryID int) bool {
	var count int
	err := db.QueryRow("DELETE FROM entry WHERE id = $1 RETURNING id", entryID).Scan(&count)
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		fmt.Println("Error: Delete entry failed: ", err)
		return false
	}
	return true
}

func hasLatestEntry(contestID int, minutesAgo int) bool {
	t := time.Now().Add(-time.Duration(minutesAgo) * time.Minute)
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM entry WHERE contest_id = $1 AND timestamp >= $2", contestID, t).Scan(&count)
	if err != nil {
		fmt.Println("Error: Get entry failed: ", err)
	}
	if count > 0 {
		return true
	} else {
		return false
	}
}

func insertLogFileToLatest(contestID int, logType, logPath string) (bool, string) {
	var id int
	var err error
	switch logType {
	case "access":
		err = db.QueryRow("UPDATE entry SET access_log_path = $1 WHERE id IN (SELECT id FROM entry WHERE contest_id = $2 ORDER BY id DESC LIMIT 1) RETURNING id", logPath, contestID).Scan(&id)

	case "slow":
		err = db.QueryRow("UPDATE entry SET slow_log_path = $1 WHERE id IN (SELECT id FROM entry WHERE contest_id = $2 ORDER BY id DESC LIMIT 1) RETURNING id", logPath, contestID).Scan(&id)
	}

	if err != nil {
		fmt.Println("Error: Create entry failed: ", err)
	}

	if id > 0 {
		return true, fmt.Sprintf("%d", id)
	} else {
		return false, "Failed to update entry"
	}
}

func insertLogFileByID(contestID, entryID int, logType, logPath string) (bool, string) {
	var id int
	var err error
	switch logType {
	case "access":
		err = db.QueryRow("UPDATE entry SET access_log_path = $1 WHERE contest_id = $2 AND id = $3 RETURNING id", logPath, contestID, entryID).Scan(&id)
	case "slow":
		err = db.QueryRow("UPDATE entry SET slow_log_path = $1 WHERE contest_id = $2 AND id = $3 RETURNING id", logPath, contestID, entryID).Scan(&id)
	}

	if err != nil {
		fmt.Println("Error: Create entry failed: ", err)
	}

	if id > 0 {
		return true, fmt.Sprintf("%d", id)
	} else {
		return false, "Failed to update entry"
	}
}

func updateMessageByID(contestID, entryID int, message string) (bool, error) {
	result, err := db.Exec("UPDATE entry SET message = $1 WHERE contest_id = $2 AND id = $3", message, contestID, entryID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return (rowsAffected > 0), nil
}

func main() {
	// Initialize DB connection
	db = initializeDB()
	defer db.Close()

	e := echo.New()

	// API secret key
	secretKey = os.Getenv("SECRET_KEY")
	if secretKey == "" {
		e.Logger.Warnf("SECRET_KEY is not set")
		e.Logger.Warnf("SECRET_KEY is \"isulogger\"")
		secretKey = "isulogger"
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(SharedKeyAuthMiddleware)
	e.Use(middleware.CORS())

	e.GET("/", hello)

	e.POST("/contest", createContest)
	e.POST("/entry", createLogEntry)
	e.DELETE("/entry/:entry_id", deleteLogEntry)
	e.POST("/entry/:contest_id/:log_type", uploadLogFile)

	e.PUT("/entry/:contest_id/:entry_id/message", updateMessage)

	e.GET("/entry", getLogEntry)
	e.GET("/contest", getContest)
	e.GET("/check/entry", checkLatestEntry)

	e.Static("/log", "log")

	e.Logger.Fatal(e.Start(":8082"))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func createLogEntry(c echo.Context) error {
	type postData struct {
		ContestID  int    `json:"contest_id"`
		BranchName string `json:"branch_name"`
		Score      int    `json:"score"`
		Message    string `json:"message"`
	}
	var p postData
	if err := c.Bind(&p); err != nil {
		return echo.ErrInternalServerError
	}
	if p.ContestID == 0 {
		return c.String(http.StatusBadRequest, "Contest ID is required")
	}

	entry := LogEntry{
		ContestID:  p.ContestID,
		Timestamp:  time.Now(),
		BranchName: p.BranchName,
		Score:      p.Score,
		Message:    p.Message,
	}

	if ok, id := insertLogEntry(&entry); !ok {
		return echo.ErrInternalServerError
	} else {
		return c.JSON(http.StatusOK, id)
	}
}

func deleteLogEntry(c echo.Context) error {
	entryIDraw := c.Param("entry_id")
	entryID, err := strconv.Atoi(entryIDraw)
	if err != nil {
		return c.String(http.StatusBadRequest, "Entry ID is invalid")
	}

	if ok := deleteLogByID(entryID); !ok {
		return echo.ErrInternalServerError
	} else {
		return c.JSON(http.StatusOK, "OK")
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
		return c.String(http.StatusBadRequest, "Contest Name is required")
	}

	if ok, id := insertContest(data.ContestName); !ok {
		return echo.ErrInternalServerError
	} else {
		return c.JSON(http.StatusOK, id)
	}
}

func uploadLogFile(c echo.Context) error {
	contestIDRaw := c.Param("contest_id")
	if contestIDRaw == "" {
		return c.String(http.StatusBadRequest, "Contest ID is required")
	}
	contestID, err := strconv.Atoi(contestIDRaw)
	if err != nil {
		return c.String(http.StatusBadRequest, "Contest ID is invalid")
	}
	logType := c.Param("log_type")
	if logType == "" {
		return c.String(http.StatusBadRequest, "Log Type is required")
	}
	if logType != "access" && logType != "slow" {
		return c.String(http.StatusBadRequest, "Log Type is invalid")
	}

	file, err := c.FormFile("log")
	if err != nil {
		return echo.ErrInternalServerError
	}

	entryIDRaw := c.FormValue("entry_id")
	entryID, err := strconv.Atoi(entryIDRaw)
	if err != nil && entryIDRaw != "" {
		return c.String(http.StatusBadRequest, "Entry ID is invalid")
	}

	src, err := file.Open()
	if err != nil {
		return echo.ErrInternalServerError
	}
	defer src.Close()
	dst, err := os.Create("./log/" + file.Filename)
	if err != nil {
		return echo.ErrInternalServerError
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return echo.ErrInternalServerError
	}

	var ok bool
	var id string
	if entryID > 0 {
		ok, id = insertLogFileByID(contestID, entryID, logType, file.Filename)
	} else {
		ok, id = insertLogFileToLatest(contestID, logType, file.Filename)
	}

	if !ok {
		return c.JSON(http.StatusInternalServerError, id)
	} else {
		return c.JSON(http.StatusOK, id)
	}
}

func updateMessage(c echo.Context) error {
	type postData struct {
		Message string `json:"message"`
	}

	contestIDRaw := c.Param("contest_id")
	if contestIDRaw == "" {
		return c.String(http.StatusBadRequest, "Contest ID is required")
	}
	contestID, err := strconv.Atoi(contestIDRaw)
	if err != nil {
		return c.String(http.StatusBadRequest, "Contest ID is invalid")
	}

	entryIDRaw := c.Param("entry_id")
	if entryIDRaw == "" {
		return c.String(http.StatusBadRequest, "Entry ID is required")
	}
	entryID, err := strconv.Atoi(entryIDRaw)
	if err != nil {
		return c.String(http.StatusBadRequest, "Entry ID is invalid")
	}

	var p postData
	if err := c.Bind(&p); err != nil {
		return echo.ErrInternalServerError
	}

	updated, err := updateMessageByID(contestID, entryID, p.Message)
	if err != nil {
		c.Logger().Errorf("updateMessageById: %v", err)
		return echo.ErrInternalServerError
	} else if !updated {
		return c.String(http.StatusNotFound, "contest or entry not found")
	} else {
		return c.NoContent(http.StatusNoContent)
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
		return c.String(http.StatusBadRequest, "Contest ID is invalid")
	}
	entry := selectLogEntry(contestID, orderBy)
	return c.JSON(http.StatusOK, entry)
}

func getContest(c echo.Context) error {
	contest := selectContest()
	return c.JSON(http.StatusOK, contest)
}

func checkLatestEntry(c echo.Context) error {
	contestIDRaw := c.QueryParam("contest_id")
	minutesAgoRaw := c.QueryParam("minutes_ago")
	if contestIDRaw == "" || minutesAgoRaw == "" {
		return echo.ErrBadRequest
	}
	contestID, err := strconv.Atoi(contestIDRaw)
	if err != nil {
		return c.String(http.StatusBadRequest, "Contest ID is invalid")
	}
	minutesAgo, err := strconv.Atoi(minutesAgoRaw)
	if err != nil {
		return c.String(http.StatusBadRequest, "Minutes ago is invalid")
	}
	if hasLatestEntry(contestID, minutesAgo) {
		return c.String(http.StatusOK, "true")
	} else {
		return c.String(http.StatusOK, "false")
	}
}
