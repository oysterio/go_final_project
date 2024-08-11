package database

import (
	"database/sql"

	"go_final_project/tasks"
)

func GetTaskByID(id int, task *tasks.Task, db *sql.DB) (errText string, statusCode int, err error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	err = db.QueryRow(query, id).Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err == sql.ErrNoRows {
		errText = "Task not found"
		statusCode = 404
		return errText, statusCode, err
	} else if err != nil {
		errText = "Error executing db query"
		statusCode = 500
		return errText, statusCode, err
	}
	return "", 0, nil
}
