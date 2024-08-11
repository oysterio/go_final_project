package database

import (
	"database/sql"

	"go_final_project/tasks"
)

func AddTask(task *tasks.Task, db *sql.DB) (id *int64, errText string, err error) {
	var idTask int64
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		errText = "AddTaskHandler: Error executing db query"
		return nil, errText, err
	}

	// get new task ID
	idTask, err = res.LastInsertId()
	if err != nil {
		errText = "AddTaskHandler: Error executing db query"
		return nil, errText, err
	}
	id = &idTask
	return id, "", nil
}
