package web

import (
	"github.com/Project-Planner/backend/model"
	"net/http"
)

func postTaskHandler(w http.ResponseWriter, r *http.Request) {
	i, err := model.NewTask(r)
	c, err := preparePostItem(w, r, i, err)
	if err != nil {
		return
	}

	c.Items.Tasks.Task = append(c.Items.Tasks.Task, i)

	finishItem(w, c, i, http.StatusCreated)
}

func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Parse data for put
	a, err := model.NewTask(r)
	c, err := preparePutItem(w, r, err)
	if err != nil {
		return
	}

	items := c.Items.Tasks.Task

	// find idx of item to be edited
	ids := make([]model.Identifier, len(items))
	for i, v := range items {
		ids[i] = v
	}
	idx, err := itemIdx(w, r, ids...)
	if err != nil {
		return // err reporting already done by method call
	}

	items[idx].Update(a)

	finishItem(w, c, items[idx], http.StatusOK)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	c, err := getCalendarIfPermission(w, r, model.Edit)
	if err != nil {
		// err reporting already done by method call
		return
	}

	items := c.Items.Tasks.Task

	ids := make([]model.Identifier, len(items))
	for i, v := range items {
		ids[i] = v
	}
	idx, err := itemIdx(w, r, ids...)
	if err != nil {
		return // err reporting already done by method call
	}

	items[idx] = items[len(items)-1]
	c.Items.Tasks.Task = items[:len(items)-1]

	finishItem(w, c, nil, http.StatusNoContent)
}
