// package handlers provides API handlers
package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"go_final_project/database"
	"go_final_project/dates"
	"go_final_project/tasks"
)

// DoneTaskHandler takes request and mark task as done by ID
func DoneTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		SendErrorResponse(w, "DoneTaskHandler: Method not allowed", http.StatusBadRequest)
		return
	}

	idTask := r.FormValue("id")
	if idTask == "" {
		SendErrorResponse(w, "DoneTaskHandler: No ID provided", http.StatusBadRequest)
		return
	}

	idTaskParsed, _ := strconv.Atoi(idTask)

	var task tasks.Task
	errText, statusCode, err := database.GetTaskByID(idTaskParsed, &task, db)
	errMsg := "DoneTaskHandler: " + errText
	if err != nil {
		SendErrorResponse(w, errMsg, statusCode)
		return
	}

	now := time.Now()
	// get task repetition date
	if task.Repeat != "" {
		newTaskDate, err := dates.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			SendErrorResponse(w, "DoneTaskHandler: Invalid repeat pattern", http.StatusBadRequest)
			return
		}

		task.Date = newTaskDate
		// update task if repeat rule set

		errText, err := database.EditTask(&task, db)
		errMsg := "DoneTaskHandler: " + errText
		if err != nil {
			SendErrorResponse(w, errMsg, http.StatusInternalServerError)
			return
		}
	} else {
		// delete task if repeat rule not set
		errText, err := database.DeleteTask(idTaskParsed, db)
		errMsg := "DoneTaskHandler: " + errText
		if err != nil {
			SendErrorResponse(w, errMsg, http.StatusInternalServerError)
			return
		}
	}

	// send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{}`))
	if err != nil {
		SendErrorResponse(w, "DoneTaskHandler: Error sending empty response", http.StatusInternalServerError)
	}
}
