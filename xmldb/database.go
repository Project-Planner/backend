package xmldb

import (
	"fmt"
	"github.com/Project-Planner/backend/model"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//The struct implementing the web.Database interface
type Database struct {
	config    DBConfig
	logins    map[string]model.Login
	users     map[string]model.User
	calendars map[string]model.Calendar
}

//New configures and parses a new database struct.
//Using the parent folder (database folder) path in the passed config, it is ensured that
//necessary folders actually exists before parsing their content into the database struct.
//After this process, the returned struct is a 1-by-1 depiction of the current system status.
func New(config DBConfig) (Database, error) {

	//1. Step: Expanding configuration file (e.g. constructing absolute paths
	//		   from relative paths)
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	config.AuthDir = fmt.Sprintf("%s%s", config.DBDir, config.AuthRelDir)
	config.UserDir = fmt.Sprintf("%s%s", config.DBDir, config.UserRelDir)
	config.CalendarDir = fmt.Sprintf("%s%s", config.DBDir, config.CalendarRelDir)

	//2. Step: Ensuring that parent folders (auth, user, calendars) exist.
	//――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	ensureDir(config.DBDir)
	ensureDir(config.AuthDir)
	ensureDir(config.UserDir)
	ensureDir(config.CalendarDir)

	//3. Step: Parsing source files into corresponding structures or collections.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	// Authentication: Each user has his own authentication file containing his
	//				   login data.
	logins := make(map[string]model.Login)
	filepath.Walk(config.AuthDir, func(file string, info os.FileInfo, err error) error {
		var name = info.Name()
		var index = strings.LastIndex(name, ".xml")
		if index >= 0 {
			var login model.Login
			parse(file, &login)
			logins[login.Name.Val] = login
		}
		return nil
	})

	// Users: Each user also has its own user file linking to his calendars.
	users := make(map[string]model.User)
	filepath.Walk(config.UserDir, func(file string, info os.FileInfo, err error) error {
		var name = info.Name()
		var index = strings.LastIndex(name, ".xml")
		if index >= 0 {
			var user model.User
			parse(file, &user)
			users[user.Name.Val] = user
		}
		return nil
	})

	// Calendars: Each calendar has an owner. Hence, it is placed into a folder
	//			  named after its owner, along with other calendars.
	calendars := make(map[string]model.Calendar)
	filepath.Walk(config.CalendarDir, func(folder string, folderInfo os.FileInfo, err error) error {

		if folderInfo.IsDir() && folder != config.CalendarDir {
			filepath.Walk(folder, func(file string, fileInfo os.FileInfo, err error) error {

				if !fileInfo.IsDir() {
					var name = fileInfo.Name()
					var index = strings.LastIndex(name, ".xml")
					if index >= 0 {
						var calendar model.Calendar
						parse(file, &calendar)

						var key = fmt.Sprintf("%s/%s", calendar.Owner.Val, name[:index])
						calendars[key] = calendar
					}
				}
				return nil

			})

		}
		return nil
	})

	return Database{config, logins, users, calendars}, nil
}

//GetUser retrieves the user to a given @userID.
//If the user doesn't exist, an error is thrown.
func (db Database) GetUser(userID string) (model.User, error) {
	val, ok := db.users[userID]
	if !ok {
		return model.User{}, model.ErrNotFound
	}
	return val, nil
}

//SetUser sets the given user to the given @userID.
//This overwrites any existing user or creates a new one,
//on the disk as well as in the collection.
func (db Database) SetUser(userID string, user model.User) {
	var path = fmt.Sprintf("%s/%s.xml", db.config.UserDir, userID)
	write(path, user.String())
	db.users[userID] = user
}

//AddUser creates a new user by creating the user file itself, the user's authentication file
//and an initial calendar file.
func (db Database) AddUser(userID, hash string) error {
	//1. Step: Checking whether user is already registered.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――
	if _, ok := db.users[userID]; ok {
		return model.ErrAlreadyExists
	}

	//2. Step: Ensuring that target folders actually exists
	//		   before creating the new user.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――――
	ensureDir(db.config.DBDir)
	ensureDir(db.config.AuthDir)
	ensureDir(db.config.UserDir)
	ensureDir(db.config.CalendarDir)

	//2. Step: Creating authentication file.
	//―――――――――――――――――――――――――――――――――――――――――
	var path = fmt.Sprintf("%s/%s.xml", db.config.AuthDir, userID)
	var login = model.NewLogin(userID, hash)
	write(path, login.String())
	db.logins[userID] = login

	//3. Step: Creating initial calendar file. The initial calendar is
	//		   placed in a folder dedicated for this user.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	path = fmt.Sprintf("%s/%s", db.config.CalendarDir, userID)
	ensureDir(path)

	var calID = fmt.Sprintf("%s/initial", userID)
	var calendar = model.NewCalendar("initial", userID)
	db.SetCalendar(calID, calendar)

	//4. Step: Creating user file itself, linking initial calendar to it
	//		   and writing to disk.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	var user = model.NewUser(userID)
	if err := user.AssociateCalendar(model.OWNER, calID, db); err != nil {
		log.Fatal(err)
	}

	return nil
}

//DeleteUser not only deletes a user itself, but also his authentication file
//and his calendars. Since the calendar can be referenced by multiple users,
//these users must also be found and disassociated from the calendar.
func (db Database) DeleteUser(userID string) error {
	//1. Step: Checking whether user is already registered.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――――
	_, ok := db.users[userID]
	if !ok {
		return model.ErrNotFound
	}

	//2. Step: Deleting user file itself from disk and
	//		   from the user collection.
	//―――――――――――――――――――――――――――――――――――――
	delete(db.users, userID)
	var path = fmt.Sprintf("%s/%s.xml", db.config.UserDir, userID)
	if err := os.Remove(path); err != nil {
		log.Fatal(err)
	}

	//3. Step: Deleting authentication file from disk and
	//         from the authentication collection.
	//―――――――――――――――――――――――――――――――――――――――
	delete(db.logins, userID)
	path = fmt.Sprintf("%s/%s.xml", db.config.AuthDir, userID)
	if err := os.Remove(path); err != nil {
		log.Fatal(err)
	}

	//4. Step: Delete calendars folder of user to be deleted.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	path = fmt.Sprintf("%s/%s", db.config.CalendarDir, userID)
	if err := os.RemoveAll(path); err != nil {
		log.Fatal(err)
	}

	//4. Step: Finding referenced users, so that the calendar
	//		   can be disassociated in their user files.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	// Users affected from disassociation are stored in @affected.
	// It allows only one occurrence to be recorded, so that each
	// changed user struct is only written to disk once including
	// all changes at once.
	affected := make(map[string]model.User)

	//Each calendar stores his owner's identifier. Only calendars of
	//the user to be deleted, so the owner, are disassociated and deleted.
	for name, cal := range db.calendars {
		if cal.Owner.Val != userID {
			continue
		}

		//Going through all users with VIEW and EDIT permissions for
		//this calendar, and editing user struct on the fly.
		for _, userID := range append(cal.Permissions.View.User, cal.Permissions.Edit.User...) {
			//There could be two sources of the user struct to be modified. First, it could
			//be the @affected collection itself, because an already modified user must be
			//modified again or, second, the struct can be found in the plain user collection.
			//This is basically a fallback - if no match, the userID is probably invalid.
			user, exists := affected[userID.Val]
			if !exists {
				user, exists = db.users[userID.Val]
				if !exists {
					continue
				}
			}

			user.DisassociateCalendar(db.config.CalendarDir, name, cal)
			affected[userID.Val] = user
		}

		delete(db.calendars, name)

	}

	//5. Step: Writing affected users back.
	//―――――――――――――――――――――――――――――――――――――――
	for id, user := range affected {
		db.SetUser(id, user)
	}

	return nil

}

//GetLogin retrieves the login to a given @userID.
//If the user doesn't exist, an error is thrown.
func (db Database) GetLogin(userID string) (model.Login, error) {
	val, ok := db.logins[userID]
	if !ok {
		return model.Login{}, model.ErrNotFound
	}
	return val, nil
}

//GetCalendar retrieves the calendar to a given @calID.
//If the calendar doesn't exist, an error is thrown.
//Note: IDs of calendars are made of several parts.
//Each calendar has an owner (with his unique userID), hence the scheme:
//	<userID>/<unique calender name>.xml
func (db Database) GetCalendar(calID string) (model.Calendar, error) {
	val, ok := db.calendars[calID]
	if !ok {
		return model.Calendar{}, model.ErrNotFound
	}
	return val, nil
}

//AddCalendar creates a new calendar and appends it to the owner's
//collection of calendars.
func (db Database) AddCalendar(ownerID, calName string) error {
	//1. Step: Checking whether calendar already exists.
	//――――――――――――――――――――――――――――――――――――――――――――――――――――
	var calID = fmt.Sprintf("%s/%s", ownerID, calName)
	if _, ok := db.calendars[calID]; ok {
		return model.ErrAlreadyExists
	}

	//2. Step: Ensuring that target folders actually exists
	//		   before creating the new calendar.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――――
	ensureDir(fmt.Sprintf("%s/%s", db.config.CalendarDir, ownerID))

	//3. Step: Creating calendar file and registering it
	//		   in the corresponding collection.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――――
	var calendar = model.NewCalendar(calName, ownerID)
	db.SetCalendar(calID, calendar)

	//4. Step: Associating calendar to issuing owner
	//		   and writing back.
	//―――――――――――――――――――――――――――――――――――――――――――――――――
	var owner = db.users[ownerID]
	if err := owner.AssociateCalendar(model.OWNER, calID, db); err != nil {
		log.Fatal(err)
	}

	return nil
}

//SetCalendar sets the given calendar to the given @calendarID.
//This overrides any existing calendar or creates a new one,
//on the disk as well as in the collection.
func (db Database) SetCalendar(calID string, cal model.Calendar) {
	var path = fmt.Sprintf("%s/%s.xml", db.config.CalendarDir, calID)
	write(path, cal.String())
	db.calendars[calID] = cal
}

//DeleteCalendar not only deletes the calendar file behind @calendarID, but
//also removes references to this file.
func (db Database) DeleteCalendar(calID string) error {
	//1. Step: Checking whether calendar actually exists.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――
	cal, ok := db.calendars[calID]
	if !ok {
		return model.ErrNotFound
	}

	//2. Step: Finding referenced users and disassociating them from
	//		   the calendar to be deleted and deleting the calendar
	//――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	for _, userID := range append(cal.Permissions.View.User, cal.Permissions.Edit.User...) {
		user, exists := db.users[userID.Val]
		if exists {
			user.DisassociateCalendar(db.config.CalendarDir, calID, cal)
			db.SetUser(user.Name.Val, user)
		}
	}

	//The owner itself is not part of the permission list, since there is no need -
	//he automatically has all permissions. He must be called separately, so that
	//the calendar file can be deleted.
	var ownerID = cal.Owner.Val
	var owner = db.users[ownerID]
	owner.DisassociateCalendar(db.config.CalendarDir, calID, cal)
	db.SetUser(ownerID, owner)
	delete(db.calendars, calID)

	return nil
}