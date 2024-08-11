// package tasks provides tools for working on scheduler tasks
package tasks

import (
	"time"

	"go_final_project/dates"
)

// GetSearchQuery get search query and provide SQL query for search execution
func GetSearchQuery(searchStr string) (searchParam, query string) {
	searchDate, err := time.Parse("02.01.2006", searchStr)
	if err == nil {
		// get tasks by date
		searchParam = searchDate.Format(dates.DateFormat)
		query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = $1 ORDER BY date LIMIT $2"
	} else {
		// get tasks by title/comment
		searchParam = "%" + searchStr + "%"
		query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE $1 OR comment LIKE $1 ORDER BY date LIMIT $2"
	}
	return searchParam, query
}
