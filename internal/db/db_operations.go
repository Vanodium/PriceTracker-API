package db_operations

import (
	"database/sql"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type CreateTrackerRequest struct {
	TrackerId   int64
	UserId      int64
	TrackerUrl  string
	CssSelector string
	StartPrice  string
}

func connectDatabase() *sql.DB {
	database, err := sql.Open("sqlite3", "../../internal/db/database.db")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database connected")
	return database
}

func AddUser(userHash string, userEmail string) error {
	db := connectDatabase()
	defer db.Close()

	if userEmailExists(userEmail) {
		return updateUserHash(userEmail, userHash)
	}
	operation, err := db.Prepare("INSERT INTO Users (user_hash, user_email) VALUES (?, ?)")
	if err != nil {
		return err
	}
	operation.Exec(userHash, userEmail)
	log.Println("New user added")
	return nil
}

func userEmailExists(userEmail string) bool {
	db := connectDatabase()
	defer db.Close()

	operation := "SELECT user_hash FROM Users WHERE user_email = ?"
	var hash string
	row := db.QueryRow(operation, userEmail).Scan(&hash)
	return row != sql.ErrNoRows
}

func updateUserHash(userEmail string, userHash string) error {
	db := connectDatabase()
	defer db.Close()

	operation, err := db.Prepare("UPDATE Users SET user_hash = ? WHERE user_email = ?")
	if err != nil {
		return err
	}
	operation.Exec(userHash, userEmail)
	log.Printf("Token changed for user %s", userEmail)
	return nil
}

func UserExists(userHash string) bool {
	db := connectDatabase()
	defer db.Close()

	operation := "SELECT user_email FROM Users WHERE user_hash = ?"
	var email string
	row := db.QueryRow(operation, userHash).Scan(&email)
	return row != sql.ErrNoRows
}

func GetUserId(userHash string) int64 {
	db := connectDatabase()
	defer db.Close()

	operation := "SELECT user_id FROM Users WHERE user_hash = ?"
	var userId int64
	db.QueryRow(operation, userHash).Scan(&userId)
	log.Printf("Found user %d", userId)
	return userId
}

func GetEmailById(userId int64) string {
	db := connectDatabase()
	defer db.Close()

	operation := "SELECT user_email FROM Users WHERE user_id = ?"
	var userEmail string
	db.QueryRow(operation, userId).Scan(&userEmail)
	return userEmail
}

func GetUserTrackers(userId int64) [][]string {
	db := connectDatabase()
	defer db.Close()

	rows, err := db.Query("SELECT tracker_id, tracker_url FROM Trackers WHERE user_id = ?", userId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var trackers [][]string
	var trackerId, trackerUrl string
	for rows.Next() {
		err = rows.Scan(&trackerId, &trackerUrl)
		if err != nil {
			log.Fatal(err)
		}
		trackers = append(trackers, []string{trackerId, trackerUrl})
	}
	return trackers
}

func AddTracker(userId int64, trackerUrl string, cssSelector string) (int64, error) {
	db := connectDatabase()
	defer db.Close()

	operation, err := db.Prepare("INSERT INTO Trackers (user_id, tracker_url, css_selector) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	result, _ := operation.Exec(userId, trackerUrl, cssSelector)

	trackerId, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	log.Printf("Tracker %d added", trackerId)
	return trackerId, nil
}

func DeleteTracker(userId int64, trackerId int64) error {
	db := connectDatabase()
	defer db.Close()
	operation, err := db.Prepare("DELETE FROM Trackers WHERE user_id = ? AND tracker_id = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = operation.Exec(userId, trackerId)
	if err != nil {
		return err
	}
	log.Printf("Tracker %d deleted", trackerId)
	return nil
}

func GetTrackerById(trackerId int64) (int64, string, string) {
	db := connectDatabase()
	defer db.Close()

	var userId int64
	var trackerUrl, cssSelector string
	db.QueryRow("SELECT user_id, css_selector, path FROM Trackers WHERE tracker_id = ?", trackerId).Scan(&userId, &trackerUrl, &cssSelector)
	return userId, trackerUrl, cssSelector
}

func GetOldPrices(currentDate int64) [][]string {
	db := connectDatabase()
	defer db.Close()

	latestCheckDate := currentDate - 3600
	rows, err := db.Query("SELECT tracker_id, last_price FROM Prices WHERE date < ?", latestCheckDate)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var recordings [][]string
	var trackerId int64
	var lastPrice string
	for rows.Next() {
		err = rows.Scan(&trackerId, &lastPrice)
		if err != nil {
			log.Fatal(err)
		}
		recordings = append(recordings, []string{strconv.FormatInt(trackerId, 10), lastPrice})
	}
	return recordings
}

func AddPrice(trackerId int64, newPrice string, currentDate int64) error {
	db := connectDatabase()
	defer db.Close()

	operation, err := db.Prepare("INSERT INTO Prices (tracker_id, last_price, date) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	operation.Exec(trackerId, newPrice, currentDate)
	log.Printf("Price added for tracker %d", trackerId)
	return nil
}

func UpdatePrice(trackerId int64, newPrice string) error {
	db := connectDatabase()
	defer db.Close()

	operation, err := db.Prepare("UPDATE Prices SET last_price = ? WHERE tracker_id = ?")
	if err != nil {
		return err
	}
	operation.Exec(newPrice, trackerId)
	log.Printf("Price changed for tracker %d", trackerId)
	return nil
}

func UpdatePriceDate(trackerId int64, currentDate int64) error {
	db := connectDatabase()
	defer db.Close()

	operation, err := db.Prepare("UPDATE Prices SET date = ? WHERE tracker_id = ?")
	if err != nil {
		return err
	}
	operation.Exec(currentDate, trackerId)
	log.Printf("Price for %d checked %d", trackerId, currentDate)
	return nil
}
