// package handlers provides API handlers
package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go_final_project/database"
	"go_final_project/dates"
	"go_final_project/tasks"
)

// AddTaskHandler
func AddTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		SendErrorResponse(w, "AddTaskHandler: Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// decode JSON
	var task tasks.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		SendErrorResponse(w, "AddTaskHandler: JSON deserialization error", http.StatusBadRequest)
		return
	}

	// check task title existence
	if task.Title == "" {
		SendErrorResponse(w, "AddTaskHandler: Task title not specified", http.StatusBadRequest)
		return
	}

	if task.Date == "" {
		task.Date = time.Now().Format(dates.DateFormat)
	}

	// check date format
	date, err := time.Parse(dates.DateFormat, task.Date)
	if err != nil {
		SendErrorResponse(w, "AddTaskHandler: Invalid date format", http.StatusBadRequest)
		return
	}

	// get task repetition date
	if task.Repeat != "" {
		dateCheck, err := dates.NextDate(time.Now(), task.Date, task.Repeat)
		if dateCheck == "" && err != nil {
			SendErrorResponse(w, "AddTaskHandler: Invalid repeat rule", http.StatusBadRequest)
			return
		}
	}

	task.Date, err = dates.GetTaskRepetitionDate(task.Repeat, date)
	if err != nil {
		SendErrorResponse(w, "AddTaskHandler: Invalid repeat rule", http.StatusBadRequest)
		return
	}

	// add task
	idTask, errText, err := database.AddTask(&task, db)
	if err != nil {
		SendErrorResponse(w, errText, http.StatusInternalServerError)
		return
	}
	id := *idTask
	task.Id = fmt.Sprint(id)

	taskId := map[string]interface{}{"id": id}
	response, err := json.Marshal(taskId)
	if err != nil {
		SendErrorResponse(w, "AddTaskHandler: response JSON creation  error", http.StatusInternalServerError)
		return
	}
	// send response
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		SendErrorResponse(w, "AddTaskHandler: Error sending response", http.StatusInternalServerError)
	}
}
