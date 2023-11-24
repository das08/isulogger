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

	_ "github.com/mattn/go-sqlite3"
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
	ID            int            `json:"id" db:"id"`
	ContestID     int            `json:"contest_id" db:"contest_id"`
	Timestamp     string         `json:"timestamp" db:"timestamp"`
	BranchName    string         `json:"branch_name" db:"branch_name"`
	Score         int            `json:"score" db:"score"`
	Message       string         `json:"message" db:"message"`
	AttachedFiles []AttachedFile `json:"attached_files"`
}

type AttachedFile struct {
	ID       int    `json:"id" db:"id"`
	EntryID  int    `json:"entry_id" db:"entry_id"`
	FileType string `json:"file_type" db:"file_type"`
	Source   string `json:"source" db:"source"`
	FilePath string `json:"file_path" db:"file_path"`
}

func initializeDB() *sql.DB {
	_db, err := sql.Open("sqlite3", fmt.Sprintf("file://%s?mode=rw", DB_PATH))
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
	if entry.Timestamp == "" {
		currentTime := time.Now()
		entry.Timestamp = currentTime.Format(time.RFC3339)
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
	result, err := db.Exec("INSERT INTO contest(contest_name) VALUES(?)", contestName)
	if err != nil {
		fmt.Println("Error: Create contest failed: ", err)
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

func selectLogEntry(contestID int, orderBy string) []LogEntry {
	var entry []LogEntry
	query := "SELECT id,timestamp,branch_name,score,message FROM entry WHERE contest_id = ? ORDER BY timestamp asc"
	if orderBy == "desc" {
		query = "SELECT id,timestamp,branch_name,score,message FROM entry WHERE contest_id = ? ORDER BY timestamp desc"
	}
	rows, err := db.Query(query, contestID)
	if err != nil {
		fmt.Println("Error: Get entry failed: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var e LogEntry
		e.ContestID = contestID
		err := rows.Scan(&e.ID, &e.Timestamp, &e.BranchName, &e.Score, &e.Message)
		if err != nil {
			fmt.Println("Error: Scan entry failed: ", err)
			continue
		}

		// Next, fetch the log files attached to this entry
		// TODO: this is a N+1 query, ISUCON CHANCE
		fileQuery := `SELECT id,file_type,source,file_path FROM attached_file WHERE entry_id=?`
		fileRows, err := db.Query(fileQuery, e.ID)
		if err != nil {
			fmt.Println("Error: Get attached file failed: ", err)
			continue
		}
		defer fileRows.Close()
		e.AttachedFiles = []AttachedFile{} // assign empty slice to avoid null
		for fileRows.Next() {
			var f AttachedFile
			f.EntryID = e.ID
			err := fileRows.Scan(&f.ID, &f.FileType, &f.Source, &f.FilePath)
			if err != nil {
				fmt.Println("Error: Scan attached file failed: ", err)
				continue
			}
			e.AttachedFiles = append(e.AttachedFiles, f)
		}

		entry = append(entry, e)
	}
	return entry
}

func deleteLogByID(entryID int) bool {
	result, err := db.Exec("DELETE FROM entry WHERE id = ?", entryID)
	if err != nil {
		fmt.Println("Error: Delete entry failed: ", err)
		return false
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Error: Get rows affected failed: ", err)
		return false
	}

	if rowsAffected == 0 {
		return false
	} else {
		return true
	}
}

func hasLatestEntry(contestID int, minutesAgo int) bool {
	t := time.Now().Add(-time.Duration(minutesAgo) * time.Minute)
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM entry WHERE contest_id = ? AND timestamp >= ?", contestID, t)
	if err != nil {
		fmt.Println("Error: Get entry failed: ", err)
	}
	if count > 0 {
		return true
	} else {
		return false
	}
}

func insertLogFileToLatest(contestID int, logType, source, logPath string) (bool, string) {
	// first find the latest entry_id
	var entryID int
	err := db.QueryRow(`SELECT id FROM entry WHERE contest_id=? ORDER BY timestamp DESC LIMIT 1`, contestID).Scan(&entryID)
	if err != nil {
		fmt.Println("Error: Get entry failed: ", err)
		return false, "Failed to find the latest entry"
	}

	return insertLogFileByID(contestID, entryID, logType, source, logPath)
}

func insertLogFileByID(contestID, entryID int, logType, source, logPath string) (bool, string) {
	// upsert the log file with (entry_id, log_type, source) as the key
	result, err := db.Exec(`INSERT INTO attached_file (entry_id, file_type, source, file_path) VALUES (?, ?, ?, ?) ON CONFLICT(entry_id, file_type, source) DO UPDATE SET file_path = ?`, entryID, logType, source, logPath, logPath)
	if err != nil {
		fmt.Println("Error: Create entry failed: ", err)
		return false, "Failed to insert log file"
	}

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
		return false, "Failed to update entry"
	}
}

func updateMessageByID(contestID, entryID int, message string) (bool, error) {
	result, err := db.Exec("UPDATE entry SET message = ? WHERE contest_id = ? AND id = ?", message, contestID, entryID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return (rowsAffected > 0), nil
}

func listSourcesUsedInContest(contestID int) ([]string, error) {
	rows, err := db.Query("SELECT DISTINCT(af.source) FROM attached_file af INNER JOIN entry e ON af.entry_id=e.id WHERE e.contest_id=? ORDER BY af.source ASC", contestID)
	if err != nil {
		return nil, err
	}

	var sources []string

	for rows.Next() {
		var source string
		err := rows.Scan(&source)
		if err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}

	return sources, nil
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
	e.GET("/entry/:contest_id/sources", listSources)

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
		Timestamp:  time.Now().Format(time.RFC3339),
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

	source := c.QueryParam("source")
	if source == "" {
		return c.String(http.StatusBadRequest, "Source is required")
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
		ok, id = insertLogFileByID(contestID, entryID, logType, source, file.Filename)
	} else {
		ok, id = insertLogFileToLatest(contestID, logType, source, file.Filename)
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

func listSources(c echo.Context) error {
	type returnData struct {
		Sources []string `json:"sources"`
	}

	contestIDRaw := c.Param("contest_id")
	if contestIDRaw == "" {
		return c.String(http.StatusBadRequest, "Contest ID is required")
	}
	contestID, err := strconv.Atoi(contestIDRaw)
	if err != nil {
		return c.String(http.StatusBadRequest, "Contest ID is invalid")
	}

	sources, err := listSourcesUsedInContest(contestID)
	if err != nil {
		c.Logger().Errorf("listSourcesUsedInContest: %v", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, returnData{Sources: sources})
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
