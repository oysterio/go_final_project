// package tasks provides tools for working on scheduler tasks
package tasks

import "sort"

func SortTaskSlice(taskList []Task) {
	sort.Slice(taskList, func(i, j int) bool {
		return taskList[i].Date < taskList[j].Date
	})
}
