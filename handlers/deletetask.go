// package handlers provides API handlers
package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"go_final_project/database"
)

// DeleteTaskHandler takes request and delete task by ID
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodDelete {
		SendErrorResponse(w, "DeleteTaskHandler: Method not allowed", http.StatusBadRequest)
		return
	}

	idTask := r.FormValue("id")
	if idTask == "" {
		SendErrorResponse(w, "DeleteTaskHandler: No ID provided", http.StatusBadRequest)
		return
	}

	idTaskParsed, err := strconv.Atoi(idTask)
	if err != nil {
		SendErrorResponse(w, "DeleteTaskHandler: Invalid ID format", http.StatusInternalServerError)
		return
	}
	// delete task by ID
	errText, err := database.DeleteTask(idTaskParsed, db)
	errMsg := "DeleteTaskHandler: " + errText
	if err != nil {
		SendErrorResponse(w, errMsg, http.StatusInternalServerError)
		return
	}
	// send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{}`))
	if err != nil {
		SendErrorResponse(w, "DeleteTaskHandler: Error sending empty response", http.StatusInternalServerError)
	}
}
