package database

import (
	"database/sql"
)

func CheckTaskExistence(id string, db *sql.DB) (errText string, statusCode int, err error) {
	var idTask int
	err = db.QueryRow("SELECT id FROM scheduler WHERE id = ?", id).Scan(&idTask)
	if err == sql.ErrNoRows {
		errText = "EditTaskHandler: Task not found"
		statusCode = 404
		return errText, statusCode, err
	} else if err != nil {
		errText = "EditTaskHandler: Error checking task existence"
		statusCode = 500
		return errText, statusCode, err
	}
	return "", 0, nil
}
