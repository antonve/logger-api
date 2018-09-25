package models

import (
	"errors"
	"fmt"
	"strconv"

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
	UserID   uint64         `json:"user_id" db:"user_id"`
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

// ByLanguage only keeps the logs for a certain language
func (logCollection *LogCollection) ByLanguage(language enums.Language) {
	filteredLogs := make([]Log, 0)

	for _, log := range logCollection.Logs {
		if log.Language == language {
			filteredLogs = append(filteredLogs, log)
		}
	}

	logCollection.Logs = filteredLogs
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
		WHERE
			user_id = $1 AND
		  deleted = FALSE
	`, userID)

	return err
}

// GetAllWithFilters returns all logs with filters applied
func (logCollection *LogCollection) GetAllWithFilters(filters map[string]interface{}) error {
	db := GetDatabase()
	defer db.Close()

	where := "DELETED = FALSE"

	for filter, value := range filters {
		switch filter {
		case "user_id":
			filters["user_id"] = value.(uint64)
			if filters["user_id"] != 0 {
				where = where + " AND user_id = :user_id"
			}
		case "date":
			if value != "" {
				where = where + " AND date = :date"
			}
		case "from":
			if value != "" {
				where = where + " AND date >= :from"
			}
		case "until":
			if value != "" {
				where = where + " AND date <= :until"
			}
		case "language":
			if value != "" {
				where = where + " AND language = :language"
			}
		case "page":
			page, err := strconv.ParseUint(value.(string), 10, 64)
			if err != nil || page <= 0 {
				page = 1
			}

			filters["page"] = (page - 1) * 30
		}
	}

	if _, ok := filters["page"]; !ok {
		filters["page"] = 0
	}

	query := `
		WITH filtered_logs AS (
			SELECT
				unnest(ids) AS id
			FROM (
				SELECT
					array_agg(id) AS ids
				FROM logs
				WHERE ` + where + `
				GROUP BY date
				ORDER BY date DESC
				OFFSET :page
				LIMIT 30
			) AS agg_ids
		)
		SELECT
			id,
			user_id,
			language,
			to_char(date, 'YYYY-MM-DD') AS date,
			duration,
			activity,
			notes
		FROM logs l
		WHERE EXISTS (
			SELECT 1
			FROM filtered_logs fl
			WHERE fl.id = l.id
		)
		ORDER BY date DESC, language
	`

	rows, err := db.NamedQuery(query, filters)
	defer rows.Close()

	if err != nil {
		return err
	}

	for rows.Next() {
		var log Log
		err = rows.StructScan(&log)

		if err != nil {
			return err
		}

		logCollection.Logs = append(logCollection.Logs, log)
	}

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
			user_id,
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
		VALUES (:user_id, :language, :date, :duration, :activity, :notes)
		RETURNING id
	`
	rows, err := db.NamedQuery(query, log)
	defer rows.Close()

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
			language = :language,
			date = :date,
			duration = :duration,
			activity = :activity,
			notes = :notes
		WHERE
			id = :id AND
			user_id = :user_id AND
			deleted = FALSE
	`
	result, err := db.NamedExec(query, log)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if rows == 0 {
		err = fmt.Errorf("no log found with id %d for user %d", log.ID, log.UserID)
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
