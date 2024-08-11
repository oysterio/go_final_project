// package handlers provides API handlers
package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"go_final_project/tasks"
)

const taskLimit = 50

type TasksResponse struct {
	Tasks []tasks.Task `json:"tasks"`
}

// GetTasksHandler takes request and provide existing tasks from db
func GetTasksHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		SendErrorResponse(w, "GetTasksHandler: Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var taskList []tasks.Task
	var task tasks.Task
	var rows *sql.Rows
	var err error
	searchStr := r.FormValue("search")
	// get tasks with search query provided
	if searchStr != "" {
		searchParam, query := tasks.GetSearchQuery(searchStr)
		rows, err := db.Query(query, searchParam, taskLimit)
		if err != nil {
			SendErrorResponse(w, "GetTasksHandler: Error executing db query", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
	} else {
		// get tasks w/o condition
		query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT $1"
		rows, err = db.Query(query, taskLimit)
		if err != nil {
			SendErrorResponse(w, "GetTasksHandler: Error executing db query", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
	}

	// scan provided tasks
	if err := rows.Err(); err != nil {
		SendErrorResponse(w, "GetTasksHandler: Failed to iterate over rows", http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			SendErrorResponse(w, "GetTasksHandler: Error scanning data from the database", http.StatusInternalServerError)
			return
		}
		task.Id = fmt.Sprint(id)
		taskList = append(taskList, task)
	}

	// provide empty array if task list is empty
	if len(taskList) == 0 {
		taskList = []tasks.Task{}
	}

	// sort tasks
	tasks.SortTaskSlice(taskList)

	// get JSON response
	responseMap := map[string][]tasks.Task{"tasks": taskList}
	response, err := json.Marshal(responseMap)
	if err != nil {
		SendErrorResponse(w, "GetTasksHandler: response JSON creation error", http.StatusInternalServerError)
		return
	}

	// send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		SendErrorResponse(w, "GetTasksHandler: Error sending response", http.StatusInternalServerError)
	}
}
