package web

import (
	"context"
	"github.com/Project-Planner/backend/model"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestPutAppointmentHandler(t *testing.T) {
	var setCalendar model.Calendar

	myId := "1234"

	app := model.Appointment{
		ID:        myId,
		Name:      model.Attribute{Val: "My Birthday"},
		StartDate: model.Attribute{Val: "15.02.2000"},
		StartTime: model.Attribute{Val: "15:00"},
		EndTime:   model.Attribute{Val: "16:00"},
		EndDate:   model.Attribute{Val: "15.02.2000"},
		Desc:      "my desc",
	}
	cWithApp := defCalendar
	cWithApp.Items.Appointments.Appointment = append(cWithApp.Items.Appointments.Appointment, app)

	tt := []struct {
		path      string
		authed    string
		code      int
		urlValues map[string]string
		db        dbMock
	}{
		// Kosher case
		{
			path:   "/c/appointments/" + testOwner + "/" + testOwner + "/" + myId,
			authed: testOwner,
			code:   http.StatusOK,
			urlValues: map[string]string{
				"name": "My Birthday Party",
				"desc": "I am partying",
			},
			db: dbMock{setCalendar: func(s string, calendar model.Calendar) error {
				setCalendar = calendar
				return nil
			},
				data: map[string]struct {
					d interface{}
					e error
				}{
					"GetCalendar": {d: cWithApp},
				},
			},
		},
		// Kosher case empty description
		{
			path:   "/c/appointments/" + testOwner + "/" + testOwner + "/" + myId,
			authed: testOwner,
			code:   http.StatusOK,
			urlValues: map[string]string{
				"name": "My Birthday Party",
				"desc": " ",
			},
			db: dbMock{setCalendar: func(s string, calendar model.Calendar) error {
				setCalendar = calendar
				return nil
			},
				data: map[string]struct {
					d interface{}
					e error
				}{
					"GetCalendar": {d: cWithApp},
				},
			},
		},
	}

	for _, tc := range tt {
		db = tc.db

		data := url.Values{}
		for k, v := range tc.urlValues {
			data.Set(k, v)
		}

		r, err := http.NewRequest("PUT", tc.path, strings.NewReader(data.Encode()))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(putAppointmentHandler)

		r = mux.SetURLVars(r, map[string]string{itemIDStr: myId})
		ctx := context.WithValue(r.Context(), userIDStr, tc.authed)
		handler.ServeHTTP(rr, r.WithContext(ctx))

		if rr.Code != tc.code {
			t.Fatalf("wrong status code: got: %d want: %d \n%s\n%v", rr.Code, tc.code, rr.Body.String(), tc)
		}

		got := setCalendar.Items.Appointments.Appointment[0]
		if tc.urlValues["name"] != got.Name.Val || tc.urlValues["desc"] != got.Desc {
			t.Error("not correctly parsed")
		}
	}
}

func TestPostAppointmentHandler(t *testing.T) {
	var setCalendar model.Calendar

	urlValues := map[string]string{
		"name":      "My Birthday",
		"startDate": "15.02.2000",
		"endDate":   "15.02.2000",
		"startTime": "14:34",
		"endTime":   "20:34",
		"desc":      "my desc",
	}

	tt := []struct {
		path      string
		authed    string
		code      int
		urlValues map[string]string
		db        dbMock
	}{
		// Kosher case
		{
			path:      "/c/appointments/" + testOwner + "/" + testOwner,
			authed:    testOwner,
			code:      http.StatusCreated,
			urlValues: urlValues,
			db: dbMock{setCalendar: func(s string, calendar model.Calendar) error {
				setCalendar = calendar
				return nil
			},
				data: map[string]struct {
					d interface{}
					e error
				}{
					"GetCalendar": {d: defCalendar},
				},
			},
		},
	}

	for _, tc := range tt {
		db = tc.db

		data := url.Values{}
		for k, v := range tc.urlValues {
			data.Set(k, v)
		}

		r, err := http.NewRequest("POST", tc.path, strings.NewReader(data.Encode()))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(postAppointmentHandler)

		ctx := context.WithValue(r.Context(), userIDStr, tc.authed)
		handler.ServeHTTP(rr, r.WithContext(ctx))

		if rr.Code != tc.code {
			t.Fatalf("wrong status code: got: %d want: %d \n%s\n%v", rr.Code, tc.code, rr.Body.String(), tc)
		}

		got := setCalendar.Items.Appointments.Appointment[0]
		if tc.urlValues["name"] != got.Name.Val || tc.urlValues["startDate"] != got.StartDate.Val {
			t.Error("not correctly parsed")
		}
	}
}
