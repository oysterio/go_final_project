// package handlers provides API handlers
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go_final_project/constants"
	"go_final_project/database"
	"go_final_project/tasks"
)

// EditTaskHandler takes request and edit task by ID
func EditTaskHandler(w http.ResponseWriter, r *http.Request, db *database.Database) {
	if r.Method != http.MethodPut {
		SendErrorResponse(w, "EditTaskHandler: Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// decode JSON
	var task tasks.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		SendErrorResponse(w, "EditTaskHandler: JSON deserialization error", http.StatusBadRequest)
		return
	}

	// check task ID
	if task.Id == "" {
		SendErrorResponse(w, "EditTaskHandler: Task ID not found", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(task.Id)
	if err != nil || id <= 0 {
		SendErrorResponse(w, "EditTaskHandler: Invalid task ID", http.StatusBadRequest)
		return
	}

	// check task title
	if task.Title == "" {
		SendErrorResponse(w, "EditTaskHandler: Task title must be specified", http.StatusBadRequest)
		return
	}

	if task.Date == "" {
		task.Date = time.Now().Format(constants.DateFormat)
	}

	// parse task Date
	_, err = time.Parse(constants.DateFormat, task.Date)
	if err != nil {
		SendErrorResponse(w, "EditTaskHandler: Invalid date format", http.StatusBadRequest)
		return
	}

	// check task repetition rule
	if task.Repeat != "" {
		if _, err := strconv.Atoi(task.Repeat[2:]); err != nil || (task.Repeat[0] != 'd' && task.Repeat[0] != 'y') {
			SendErrorResponse(w, "EditTaskHandler: Invalid task repetition format", http.StatusBadRequest)
			return
		}
	}

	// check task existence
	_, statusCode, err := db.GetTaskByID(id)
	if err != nil {
		SendErrorResponse(w, fmt.Sprintf("EditTaskHandler: failed to check task existence: %v", err), statusCode)
		return
	}
	// update task
	err = db.EditTask(task)
	if err != nil {
		SendErrorResponse(w, fmt.Errorf("EditTaskHandler: failed to update task: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	// get JSON response
	response, err := json.Marshal(struct{}{})
	if err != nil {
		SendErrorResponse(w, "EditTaskHandler: response JSON creation error", http.StatusInternalServerError)
		return
	}
	// send response
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		SendErrorResponse(w, "EditTaskHandler: Error sending response", http.StatusInternalServerError)
	}
}
