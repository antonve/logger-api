package models

import (
	"errors"
	"fmt"

	"github.com/antonve/logger-api/models/enums"
	"github.com/jmoiron/sqlx/types"
)

// LogCollection array of logs
type LogCollection struct {
	Logs []Log `json:"logs"`
}

// Log model
type Log struct {
	ID       uint64         `json:"id" db:"id"`
	UserID   uint64         `json:"user_id" db:"id"`
	Language enums.Language `json:"language" db:"language"`
	Date     string         `json:"date" db:"date"`
	Duration uint64         `json:"duration" db:"duration"`
	Activity enums.Activity `json:"activity" db:"activity"`
	Notes    types.JSONText `json:"notes" db:"notes"`
}

// Length returns the amount of logs in the collection
func (logCollection *LogCollection) Length() int {
	return len(logCollection.Logs)
}

// Validate the Log model
func (log *Log) Validate() error {
	if log.UserID == 0 {
		return errors.New("invalid `UserID` supplied")
	}
	if log.Date == "" {
		return errors.New("invalid `Date` supplied")
	}
	if log.Duration == 0 {
		return errors.New("invalid `Duration` supplied")
	}
	if len(log.Activity) == 0 || !log.Activity.IsValid() {
		return errors.New("invalid `Activity` supplied")
	}

	return nil
}

// GetAll returns all logs
func (logCollection *LogCollection) GetAll() error {
	db := GetDatabase()
	defer db.Close()

	err := db.Select(&logCollection.Logs, `
		SELECT
			id,
			user_id,
			language,
			to_char(date, 'YYYY-MM-DD') AS date,
			duration,
			activity,
			notes
		FROM logs
		WHERE deleted = FALSE
	`)

	return err
}

// GetAllFromUser returns all logs from a certain user
func (logCollection *LogCollection) GetAllFromUser(userID uint64) error {
	db := GetDatabase()
	defer db.Close()

	stmt, err := db.Preparex(`
		SELECT
			id,
			user_id,
			language,
			to_char(date, 'YYYY-MM-DD') AS date,
			duration,
			activity,
			notes
		FROM logs
    WHERE
      user_id = $1 AND
		  deleted = FALSE
	`)
	if err != nil {
		return err
	}

	err = stmt.Get(&logCollection.Logs, userID)

	return err
}

// Get a log by id
func (logCollection *LogCollection) Get(id uint64) (*Log, error) {
	db := GetDatabase()
	defer db.Close()

	// Init log
	log := Log{}

	// Get log
	stmt, err := db.Preparex(`
		SELECT
			id,
			user_id
			language,
			to_char(date, 'YYYY-MM-DD') AS date,
			duration,
			activity,
			notes
		FROM logs
		WHERE
			id = $1 AND
			deleted = FALSE
	`)
	if err != nil {
		return nil, err
	}

	stmt.Get(&log, id)
	if log.ID == 0 {
		return nil, fmt.Errorf("no log found with id %v", id)
	}

	return &log, nil
}

// Add a log to the database
func (logCollection *LogCollection) Add(log *Log) (uint64, error) {
	db := GetDatabase()
	defer db.Close()

	query := `
		INSERT INTO logs (user_id, language, date, duration, activity, notes)
		VALUES (:language, :date, :duration, :activity, :notes)
		RETURNING id
	`
	rows, err := db.NamedQuery(query, log)

	if err != nil {
		return 0, err
	}

	var id uint64
	if rows.Next() {
		rows.Scan(&id)
	}

	return id, nil
}

// Update a log
func (logCollection *LogCollection) Update(log *Log) error {
	db := GetDatabase()
	defer db.Close()

	query := `
		UPDATE logs
		SET
			user_id = :user_id,
			language = :language,
			date = :date,
			duration = :duration,
			activity = :activity,
			notes = :notes
		WHERE
			id = :id AND
			deleted = FALSE
	`
	result, err := db.NamedExec(query, log)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if rows == 0 {
		err = fmt.Errorf("no log found with id %v", log.ID)
	}

	return err
}

// Delete a log
func (logCollection *LogCollection) Delete(log *Log) error {
	db := GetDatabase()
	defer db.Close()

	query := `
		UPDATE logs
		SET deleted = TRUE
		WHERE
			id = :id AND
			deleted = FALSE
	`
	result, err := db.NamedExec(query, log)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if rows == 0 {
		err = fmt.Errorf("no log found with id %v or it has already been deleted", log.ID)
	}

	return err
}

// IsOwner checks the owner
func (log *Log) IsOwner(userID uint64) bool {
	if log.UserID == userID {
		return true
	}

	return false
}
