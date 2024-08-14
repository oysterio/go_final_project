// package tasks provides tools for working on scheduler tasks
package database

import (
	"database/sql"
	"fmt"
	"time"

	"go_final_project/constants"
	"go_final_project/tasks"
)

// GetSearchQuery get search query and provide SQL query for search execution
func (d Database) GetTasks(searchStr string) (taskList []tasks.Task, err error) {

	var t tasks.Task
	var searchParam string
	var query string
	var rows *sql.Rows

	switch {
	case searchStr != "":
		searchDate, err := time.Parse("02.01.2006", searchStr)
		if err == nil {
			// get tasks by date
			searchParam = searchDate.Format(constants.DateFormat)
			query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?"
			rows, err = d.db.Query(query, searchParam, constants.TaskLimit)
			if err != nil {
				return taskList, fmt.Errorf("failed to get task by date: %w", err)
			}
		} else {
			// get tasks by title/comment
			searchParam = "%" + searchStr + "%"
			query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE $1 OR comment LIKE $1 ORDER BY date LIMIT $2"
			rows, err = d.db.Query(query, searchParam, constants.TaskLimit)
			if err != nil {
				return taskList, fmt.Errorf("failed to get task by content: %w", err)
			}
		}
		defer rows.Close()
	default:
		// get tasks w/o condition
		query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?"
		rows, err = d.db.Query(query, constants.TaskLimit)
		if err != nil {
			return taskList, fmt.Errorf("failed to get tasks: %w", err)
		}
		defer rows.Close()
	}

	if err = rows.Err(); err != nil {
		return taskList, fmt.Errorf("failed to iterate over rows: %w", err)
	}

	for rows.Next() {
		var id int64
		if err = rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return taskList, fmt.Errorf("failed to scan tasks: %w", err)
		}
		t.Id = fmt.Sprint(id)
		taskList = append(taskList, t)
	}
	return taskList, nil
}
