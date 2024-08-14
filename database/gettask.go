package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"go_final_project/tasks"
)

func (d Database) GetTaskByID(id int) (tasks.Task, int, error) {
	var t tasks.Task
	log.Printf("%v", id)
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	err := d.db.QueryRow(query, id).Scan(&t.Id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err == sql.ErrNoRows {
		return t, http.StatusNotFound, nil
	} else if err != nil {
		return t, http.StatusInternalServerError, fmt.Errorf("failed to retrieve task: %w", err)
	}
	return t, http.StatusOK, nil
}
