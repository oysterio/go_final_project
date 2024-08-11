package database

import (
	"database/sql"

	"go_final_project/tasks"
)

func EditTask(task *tasks.Task, db *sql.DB) (errText string, err error) {
	query := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat =? WHERE id = ?"
	_, err = db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.Id)
	if err != nil {
		errText = "Task not found"
		return errText, err
	}
	return "", nil
}
