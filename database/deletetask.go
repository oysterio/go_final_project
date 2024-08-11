package database

import (
	"database/sql"
)

func DeleteTask(id int, db *sql.DB) (errText string, err error) {
	query := "DELETE FROM scheduler WHERE id = ?"
	result, err := db.Exec(query, id)
	if err != nil {
		errText = "Error deleting task"
		return errText, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		errText = "Unable to determine the number of affected rows"
		return errText, err
	} else if rowsAffected == 0 {
		errText = "Task not found"
		return errText, err
	}
	return "", nil
}
