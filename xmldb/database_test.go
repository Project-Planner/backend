package xmldb

import (
	"fmt"
	"github.com/Project-Planner/backend/model"
	"os"
	"testing"
)

//DONE
func TestNew(t *testing.T) {
	//1. Step: Construct a database.
	//――――――――――――――――――――――――――――――――――
	var db = GetDatabase(t)
	t.Cleanup(func() { DeleteDatabase(db, t) })

	//2. Step: Check for entries.
	//――――――――――――――――――――――――――――――――
	if len(db.users) != 0 || len(db.logins) != 0 || len(db.calendars) != 0 {
		t.Fatal(fmt.Sprintf("Falsely identified user and calendar files.\n"+
			"len(users): %d\nlen(logins): %d\nlen(calendars): %d",
			len(db.users), len(db.logins), len(db.calendars)))
	}
}

//DONE
func TestAddUser(t *testing.T) {
	//1. Step: Construct a database.
	//――――――――――――――――――――――――――――――――――
	var db = GetDatabase(t)
	t.Cleanup(func() { DeleteDatabase(db, t) })

	//2. Step: Add user and check whether the user
	//		   now actually exists by looking into
	//		   respective collections and checking
	//		   if the files are present.
	//―――――――――――――――――――――――――――――――――――――――――――――
	var userID = "f5932068"
	var hash = "hash"
	if err := db.AddUser(userID, hash); err != nil {
		t.Fatal(err)
	}

	//Check if user file exists and if user is
	//registered in user collection.
	var path = fmt.Sprintf("%s/%s.xml", db.config.UserDir, userID)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal(fmt.Sprintf("User file for user '%s' doesn't exist.", userID))
	}

	if _, ok := db.users[userID]; !ok {
		t.Fatal(fmt.Sprintf("User '%s' is not registered in user collection.", userID))
	}

	//Check if user authentication file exists and
	//if user is registered in authentication collection.
	path = fmt.Sprintf("%s/%s.xml", db.config.AuthDir, userID)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal(fmt.Sprintf("Authentication file for user '%s' doesn't exist.", userID))
	}

	if _, ok := db.logins[userID]; !ok {
		t.Fatal(fmt.Sprintf("User '%s' is not registered in authentication collection.", userID))
	}

	//Check if initial calendar and its parent folder exist
	//and if it is registered in calendar collection.
	var calID = fmt.Sprintf("%s/%s", userID, userID)
	path = fmt.Sprintf("%s/%s", db.config.CalendarDir, userID)
	if stat, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			t.Fatal(fmt.Sprintf("Calendar folder for user '%s' doesn't exist.", userID))
		} else if !stat.IsDir() {
			t.Fatal(fmt.Sprintf("Calendar 'folder' for user '%s' isn't a folder.", userID))
		}
	}

	path = fmt.Sprintf("%s/%s.xml", path, userID)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal(fmt.Sprintf("Initial calendar for user '%s' doesn't exist.", userID))
	}

	if _, ok := db.calendars[calID]; !ok {
		t.Fatal(fmt.Sprintf("Calendar for user '%s' is not registered in calendar collection", userID))
	}

	//3. Step: Check if parsing has gone wrong.
	//――――――――――――――――――――――――――――――――――――――――――
	var user = db.users[userID]
	if user.Name.Val != userID || len(user.Items.Calendars) != 1 || user.Items.Calendars[0].Link != calID {
		t.Fatal(fmt.Sprintf("User struct for user '%s' contains invalid data.", userID))
	}

	var login = db.logins[userID]
	if login.Name.Val != userID || login.Hash.Val != hash {
		t.Fatal(fmt.Sprintf("Login struct for user '%s' contains invalid data.", userID))
	}

	var calendar = db.calendars[calID]
	if calendar.Name.Val != userID || calendar.Owner.Val != userID {
		t.Fatal(fmt.Sprintf("Calendar struct for user '%s' contains invalid data.", userID))
	}

	//4. Step: Check if the right amount of collection
	//		   entries is present.
	//――――――――――――――――――――――――――――――――――――――――――――――――――
	if len(db.users) != 1 || len(db.logins) != 1 || len(db.calendars) != 1 {
		t.Fatal("Invalid amount of entries in any of the data collections.")
	}
}

//DONE
func TestAddCalendar(t *testing.T) {
	//1. Step: Construct a database.
	//――――――――――――――――――――――――――――――――――
	var db = GetDatabase(t)
	t.Cleanup(func() { DeleteDatabase(db, t) })

	//2. Step: Add user with initial and
	//		   additional calendar.
	//―――――――――――――――――――――――――――――――――――――――――
	var userID = "a"
	var hash = "hash"
	db.AddUser(userID, hash)

	var calName = "test"
	var calID = fmt.Sprintf("%s/%s", userID, calName)
	db.AddCalendar(userID, calName)

	//3. Step: Check if calendar is present
	//		   in collection and on disk.
	//―――――――――――――――――――――――――――――――――――――――――
	if _, ok := db.calendars[calID]; !ok {
		t.Fatal(fmt.Sprintf("Calendar for user '%s' is not registered in calendar collection.", userID))
	}

	var path = fmt.Sprintf("%s/%s", db.config.CalendarDir, userID)
	if stat, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			t.Fatal(fmt.Sprintf("Calendar folder for user '%s' doesn't exist.", userID))
		} else if !stat.IsDir() {
			t.Fatal(fmt.Sprintf("Calendar 'folder' for user '%s' isn't a folder.", userID))
		}
	}

	path = fmt.Sprintf("%s/%s.xml", path, calName)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal(fmt.Sprintf("File for calendar with id '%s' doesn't exist.", calID))
	}

	//4. Step: Ensuring that parsing hasn't
	//		   gone wrong.
	//―――――――――――――――――――――――――――――――――――――――――
	var cal, _ = db.GetCalendar(calID)
	if cal.Owner.Val != userID || cal.Name.Val != calName {
		t.Fatal(fmt.Sprintf("Calendar struct with id '%s' contains invalid data.", calID))
	}
}

//DONE
func TestGetUser(t *testing.T) {
	//1. Step: Construct a database.
	//――――――――――――――――――――――――――――――――――
	var db = GetDatabase(t)
	t.Cleanup(func() { DeleteDatabase(db, t) })

	//2. Step: Add user, retrieve it, compare the data
	//		   within and check for non-existent user.
	//―――――――――――――――――――――――――――――――――――――――――――――――――
	var userID1 = "f5932068"
	var hash = "hash"
	db.AddUser(userID1, hash)

	//Check for first user.
	//It should be successful, since it
	//has just been added.
	var user, err = db.GetUser(userID1)
	if err != nil {
		t.Fatal(err)
	}

	var calID = fmt.Sprintf("%s/%s", userID1, userID1)
	if user.Name.Val != userID1 || len(user.Items.Calendars) != 1 || user.Items.Calendars[0].Link != calID {
		t.Fatal(fmt.Sprintf("User struct for user '%s' contains invalid data.", userID1))
	}

	//Check for second user.
	//It should fail, since no user with this
	//id has been registered yet.
	var userID2 = "Notch"
	user, err = db.GetUser(userID2)
	if err == nil {
		t.Fatal(fmt.Sprintf("No error thrown for user '%s' although the user doesn't exist.", userID2))
	}
}

//DONE
func TestGetLogin(t *testing.T) {
	//1. Step: Construct a database.
	//――――――――――――――――――――――――――――――――――
	var db = GetDatabase(t)
	t.Cleanup(func() { DeleteDatabase(db, t) })

	//2. Step: Add user, retrieve it's login,
	//		   compare the data within and
	//		   check for a non-existent login.
	//―――――――――――――――――――――――――――――――――――――――――
	var userID1 = "f5932068"
	var hash = "hash"
	db.AddUser(userID1, hash)

	//Check for first login.
	//It should be successful, since it has
	//been added, when the user has been added.
	var login, err = db.GetLogin(userID1)
	if err != nil {
		t.Fatal(err)
	}

	if login.Name.Val != userID1 || login.Hash.Val != hash {
		t.Fatal(fmt.Sprintf("Login struct for user '%s' contains invalid data.", userID1))
	}

	//Check for second login.
	//It should fail, since no login for
	//this user has been registered yet.
	var userID2 = "Notch"
	login, err = db.GetLogin(userID2)
	if err == nil {
		t.Fatal(fmt.Sprintf("No error thrown for user '%s' although the user doesn't exist.", userID2))
	}
}

//DONE
func TestGetCalendar(t *testing.T) {
	//1. Step: Construct a database.
	//――――――――――――――――――――――――――――――――――
	var db = GetDatabase(t)
	t.Cleanup(func() { DeleteDatabase(db, t) })

	//2. Step: Add user, retrieve it's initial,
	//		   compare the data within and
	//		   check for a non-existent calendar.
	//――――――――――――――――――――――――――――――――――――――――――――
	var userID = "f5932068"
	var hash = "hash"
	db.AddUser(userID, hash)

	//Check for first calendar.
	//It should be successful, since the initial
	//calendar is generated when the user is added.
	var calName1 = userID
	var calID1 = fmt.Sprintf("%s/%s", userID, calName1)
	var cal, err = db.GetCalendar(calID1)
	if err != nil {
		t.Fatal(err)
	}

	if cal.Name.Val != calName1 || cal.Owner.Val != userID {
		t.Fatal(fmt.Sprintf("Calendar struct for calendar with id '%s' contains invalid data.", calID1))
	}

	//Check for second calendar.
	//It should fail, since no calendar with this
	//it has been added yet.
	var calID2 = fmt.Sprintf("%s/test", userID)
	cal, err = db.GetCalendar(calID2)
	if err == nil {
		t.Fatal(fmt.Sprintf("No error thrown for calendar with id '%s' although the calendar doesn't exist.", calID2))
	}
}

//DONE
func TestSetUser(t *testing.T) {
	//1. Step: Construct a database.
	//――――――――――――――――――――――――――――――――――
	var db = GetDatabase(t)
	t.Cleanup(func() { DeleteDatabase(db, t) })

	//2. Step: Add user, retrieve its struct, modify
	//		   and set it again.
	//――――――――――――――――――――――――――――――――――――――――――――――――――
	var userID = "f5932068"
	var hash = "hash"
	db.AddUser(userID, hash)

	//Retrieve user, modify it and write it back
	var user, _ = db.GetUser(userID)
	var replacement = "Notch"

	user.Name.Val = replacement
	db.SetUser(userID, user)

	//3. Step: Retrieve the user from the collection again
	//		   and compare its values to ensure that the modified
	//		   copy has indeed been written back.
	//――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	user, _ = db.GetUser(userID)
	if user.Name.Val != replacement {
		t.Fatal(fmt.Sprintf("Modified user '%s' has not been written to the collection.", userID))
	}

	//4. Step: Parse from the user file directly to check
	//		   whether the changes have also been written to disk.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	var path = fmt.Sprintf("%s/%s.xml", db.config.UserDir, userID)
	parse(path, &user)
	if user.Name.Val != replacement {
		t.Fatal(fmt.Sprintf("Modified user '%s' has not been written to disk.", userID))
	}
}

//DONE
func TestSetCalendar(t *testing.T) {
	//1. Step: Construct a database.
	//――――――――――――――――――――――――――――――――――
	var db = GetDatabase(t)
	t.Cleanup(func() { DeleteDatabase(db, t) })

	//2. Step: Add user, retrieve his initial calendar,
	//		   modify and set it.
	//――――――――――――――――――――――――――――――――――――――――――――――――――
	var userID = "f5932068"
	var hash = "hash"
	db.AddUser(userID, hash)

	//Retrieve initial calendar, modify it and write it back
	var calID = fmt.Sprintf("%s/%s", userID, userID)
	var cal, _ = db.GetCalendar(calID)
	var replacement = "Notch"

	cal.Owner.Val = replacement
	db.SetCalendar(calID, cal)

	//3. Step: Retrieve the calendar from the collection again
	//		   and compare its values to ensure that the modified
	//		   copy has indeed been written back.
	//――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	cal, _ = db.GetCalendar(calID)
	if cal.Owner.Val != replacement {
		t.Fatal(fmt.Sprintf("Modified calendar with id '%s' has not been written to the collection.", calID))
	}

	//4. Step: Parse from the calendar file directly to check
	//		   whether the changes have also been written to disk.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	var path = fmt.Sprintf("%s/%s/%s.xml", db.config.CalendarDir, userID, userID)
	parse(path, &cal)
	if cal.Owner.Val != replacement {
		t.Fatal(fmt.Sprintf("Modified calendar with id '%s' has not been written to disk.", userID))
	}
}

//DONE
func TestAssociateCalendar(t *testing.T) {
	//1. Step: Construct a database.
	//――――――――――――――――――――――――――――――――
	var db = GetDatabase(t)
	t.Cleanup(func() { DeleteDatabase(db, t) })

	//2. Step: Add two users that each have
	//		   an initial calendar and link
	//		   that of user 1 to user 2.
	//―――――――――――――――――――――――――――――――――――――――
	var userID1 = "a"
	var userID2 = "b"
	var hash = "hash"
	db.AddUser(userID1, hash)
	db.AddUser(userID2, hash)

	var calID = fmt.Sprintf("%s/%s", userID1, userID1)
	var user2, _ = db.GetUser(userID2)

	user2.AssociateCalendar(model.Edit, calID, db)
	db.SetUser(userID2, user2)

	//3. Step: Check if calendar is listed in
	//		   user's calendar list.
	//―――――――――――――――――――――――――――――――――――――――――
	var foundInUser bool
	for _, cal := range user2.Items.Calendars {
		if cal.Link == calID {
			foundInUser = true
			break
		}
	}

	if !foundInUser {
		t.Fatal(fmt.Sprintf("Calendar with id '%s' has not been found in calendar list of user '%s'.", calID, userID2))
	}

	//4. Step: Check if calendar has user registered.
	//――――――――――――――――――――――――――――――――――――――――――――――――
	var cal, _ = db.GetCalendar(calID)
	var entries = append(cal.Permissions.Edit.User, cal.Permissions.View.User...)
	var foundInCalendar bool
	for _, user := range entries {
		if user.Val == userID2 {
			foundInCalendar = true
			break
		}
	}

	if !foundInCalendar {
		t.Fatal(fmt.Sprintf("Calendar with id '%s' does not contain the user '%s'.", calID, userID2))
	}
}

//DONE
func TestDeleteUser(t *testing.T) {
	//1. Step: Construct a database.
	//――――――――――――――――――――――――――――――――――
	var db = GetDatabase(t)
	t.Cleanup(func() { DeleteDatabase(db, t) })

	//2. Step: Adding two users and check entry count.
	//―――――――――――――――――――――――――――――――――――――――――――――――――
	var userID1 = "a"
	var userID2 = "b"
	var hash = "hash"
	db.AddUser(userID1, hash)
	db.AddUser(userID2, hash)

	if len(db.users) != 2 || len(db.logins) != 2 || len(db.calendars) != 2 {
		t.Fatal("Invalid amount of entries in any of the data collections.")
	}

	//3. Step: Associate initial calendar of first
	//		   user to second one.
	//――――――――――――――――――――――――――――――――――――――――――――――――
	var calID = fmt.Sprintf("%s/%s", userID1, userID1)
	var user2, _ = db.GetUser(userID2)

	user2.AssociateCalendar(model.Edit, calID, db)
	db.SetUser(userID2, user2)

	//4. Step: Delete user, check for entry count
	//		   again, check for files and for
	//		   association with second user.
	//―――――――――――――――――――――――――――――――――――――――――――――
	db.DeleteUser(userID1)
	if len(db.users) != 1 || len(db.logins) != 1 || len(db.calendars) != 1 {
		t.Fatal("Invalid amount of entries in any of the data collections.")
	}

	//Check if user file still exists.
	var path = fmt.Sprintf("%s/%s.xml", db.config.UserDir, userID1)
	if _, err := os.Stat(path); err == nil {
		t.Fatal(fmt.Sprintf("User file for user '%s' still exists.", userID1))
	}

	//Check if authentication file still exists,
	path = fmt.Sprintf("%s/%s.xml", db.config.AuthDir, userID1)
	if _, err := os.Stat(path); err == nil {
		t.Fatal(fmt.Sprintf("Authentication file for user '%s' still exists.", userID1))
	}

	//Check if calendar folder for still exists.
	path = fmt.Sprintf("%s/%s", db.config.CalendarDir, userID1)
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		t.Fatal(fmt.Sprintf("Calendar folder for user '%s' still exists.", userID1))
	}

	//Check if second user still contains reference to initial
	//calendar of first user.
	user2, _ = db.GetUser(userID2)
	for _, cal := range user2.Items.Calendars {
		if cal.Link == calID {
			t.Fatal(fmt.Sprintf("Reference to calendar '%s' can still be found at user '%s'.", calID, userID2))
		}
	}
}

//DONE
func TestDeleteCalendar(t *testing.T) {
	//1. Step: Construct a database.
	//――――――――――――――――――――――――――――――――――
	var db = GetDatabase(t)
	t.Cleanup(func() { DeleteDatabase(db, t) })

	//2. Step: Adding two users and associate the initial
	//		   calendar of the first user to the second one.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――――
	var userID1 = "a"
	var userID2 = "b"
	var hash = "hash"
	db.AddUser(userID1, hash)
	db.AddUser(userID2, hash)

	//Associate initial calendar of first user to second one.
	var calID = fmt.Sprintf("%s/%s", userID1, userID1)
	var user2, _ = db.GetUser(userID2)

	user2.AssociateCalendar(model.Edit, calID, db)
	db.SetUser(userID2, user2)

	//3. Step: Delete calendar, check for entry count and calendar
	//		   file and if it is still associated to the second user.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	db.DeleteCalendar(calID)
	if len(db.users) != 2 || len(db.logins) != 2 || len(db.calendars) != 1 {
		t.Fatal("Invalid amount of entries in any of the data collections.")
	}

	//Check if calendar still exists in
	//collection and on disk.
	if _, ok := db.calendars[calID]; ok {
		t.Fatal(fmt.Sprintf("Calendar with id '%s' still exists in calendar collection.", calID))
	}

	var path = fmt.Sprintf("%s/%s.xml", db.config.CalendarDir, calID)
	if _, err := os.Stat(path); err == nil {
		t.Fatal(fmt.Sprintf("File for calendar with id '%s' still exists.", calID))
	}

	//Check if it is still associated to the second user.
	user2, _ = db.GetUser(userID2)
	for _, cal := range user2.Items.Calendars {
		if cal.Link == calID {
			t.Fatal(fmt.Sprintf("Calendar with id '%s' is still registered in user '%s'.", calID, userID2))
		}
	}
}

//GetDatabase loads and constructs a new database struct for testing.
func GetDatabase(t *testing.T) database {
	var config = DBConfig{
		DBDir:          "./xmldb",
		AuthRelDir:     "/auth",
		AuthDir:        "",
		UserRelDir:     "/users",
		UserDir:        "",
		CalendarRelDir: "/calendars",
		CalendarDir:    "",
		CacheSize:      0,
	}
	db, err := New(config)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func DeleteDatabase(db database, t *testing.T) {
	if err := os.RemoveAll(db.config.DBDir); err != nil {
		t.Fatal(err)
	}
}
