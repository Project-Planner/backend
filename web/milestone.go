package web

import (
	"github.com/Project-Planner/backend/model"
	"net/http"
)

func postMilestoneHandler(w http.ResponseWriter, r *http.Request) {
	i, err := model.NewMilestone(r)
	c, err := preparePostItem(w, r, i, err)
	if err != nil {
		return
	}

	c.Items.Milestones.Milestone = append(c.Items.Milestones.Milestone, i)

	finishItem(w, c, i, http.StatusCreated)
}

func putMilestoneHandler(w http.ResponseWriter, r *http.Request) {
	// Parse data for put
	a, err := model.NewMilestone(r)
	c, err := preparePutItem(w, r, err)
	if err != nil {
		return
	}

	items := c.Items.Milestones.Milestone

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

func deleteMilestoneHandler(w http.ResponseWriter, r *http.Request) {
	c, err := getCalendarIfPermission(w, r, model.Edit)
	if err != nil {
		// err reporting already done by method call
		return
	}

	items := c.Items.Milestones.Milestone

	ids := make([]model.Identifier, len(items))
	for i, v := range items {
		ids[i] = v
	}
	idx, err := itemIdx(w, r, ids...)
	if err != nil {
		return // err reporting already done by method call
	}

	items[idx] = items[len(items)-1]
	c.Items.Milestones.Milestone = items[:len(items)-1]

	finishItem(w, c, nil, http.StatusNoContent)
}
