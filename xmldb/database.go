package xmldb

import (
	"encoding/xml"
	"fmt"
	"github.com/Project-Planner/backend/model"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

//The struct implementing the web.Database interface
type database struct {
	config    DBConfig
	mutexes   map[string]*sync.Mutex
	logins    map[string]model.Login
	users     map[string]model.User
	calendars map[string]model.Calendar
}

//New configures and parses a new database struct.
//Using the parent folder (database folder) path in the passed config, it ensures that
//necessary folders actually exists before parsing their content into the database struct.
//After this process, the returned struct is a 1-by-1 depiction of the current system status.
//Since each parsed element represents one resource, except the authentication files, each
//of them gets equipped with a mutex lock.
func New(config DBConfig) (database, error) {

	// Set the "constants" here to make the config file simpler
	config.AuthRelDir = "/auth"
	config.UserRelDir = "/users"
	config.CalendarRelDir = "/calendars"

	//1. Step: Expanding configuration file (e.g. constructing absolute paths
	//		   from relative paths)
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	config.AuthDir = fmt.Sprintf("%s%s", config.DBDir, config.AuthRelDir)
	config.UserDir = fmt.Sprintf("%s%s", config.DBDir, config.UserRelDir)
	config.CalendarDir = fmt.Sprintf("%s%s", config.DBDir, config.CalendarRelDir)

	//2. Step: Ensure that parent folders (auth, user, calendars) exist.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	if err := ensureDir(config.DBDir); err != nil {
		return database{}, err
	}

	if err := ensureDir(config.AuthDir); err != nil {
		return database{}, err
	}

	if err := ensureDir(config.UserDir); err != nil {
		return database{}, err
	}

	if err := ensureDir(config.CalendarDir); err != nil {
		return database{}, err
	}

	//3. Step: Parse source files into corresponding structures or collections.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	// Authentication: Each user has his own authentication file containing his
	//				   login data.
	logins := make(map[string]model.Login)
	if err := filepath.Walk(config.AuthDir, func(file string, info os.FileInfo, err error) error {
		var name = info.Name()
		var index = strings.LastIndex(name, ".xml")
		if index >= 0 {
			var login model.Login
			if err := parse(file, &login); err != nil {
				return err
			}
			logins[login.Name.Val] = login
		}
		return nil
	}); err != nil {
		return database{}, err
	}

	// Users: Each user also has its own user file linking to his calendars.
	mutexes := make(map[string]*sync.Mutex)
	users := make(map[string]model.User)
	if err := filepath.Walk(config.UserDir, func(file string, info os.FileInfo, err error) error {
		var name = info.Name()
		var index = strings.LastIndex(name, ".xml")
		if index >= 0 {
			var user model.User
			if err := parse(file, &user); err != nil {
				return err
			}

			var userID = user.Name.Val
			var lock sync.Mutex
			users[userID] = user
			mutexes[userID] = &lock
		}
		return nil
	}); err != nil {
		return database{}, err
	}

	// Calendars: Each calendar has an owner. Hence, it is placed into a folder
	//			  named after its owner, along with other calendars.
	calendars := make(map[string]model.Calendar)
	if err := filepath.Walk(config.CalendarDir, func(folder string, folderInfo os.FileInfo, err error) error {
		if folderInfo.IsDir() && folder != config.CalendarDir {
			if err := filepath.Walk(folder, func(file string, fileInfo os.FileInfo, err error) error {

				if !fileInfo.IsDir() {
					var name = fileInfo.Name()
					var index = strings.LastIndex(name, ".xml")
					if index >= 0 {
						var calendar model.Calendar
						if err := parse(file, &calendar); err != nil {
							return err
						}

						var key = fmt.Sprintf("%s/%s", calendar.Owner.Val, name[:index])
						var lock sync.Mutex
						calendars[key] = calendar
						mutexes[key] = &lock
					}
				}
				return nil

			}); err != nil {
				return err
			}

		}
		return nil
	}); err != nil {
		return database{}, err
	}

	return database{config, mutexes, logins, users, calendars}, nil
}

//GetUser retrieves the user to a given @userID.
//If the user doesn't exist, an error is thrown.
func (db database) GetUser(userID string) (model.User, error) {
	var val, ok = db.users[userID]
	if !ok {
		return model.User{}, model.ErrNotFound
	}
	return val, nil
}

//SetUser sets the given user to the given @userID
//only if the user actually exists.
func (db database) SetUser(userID string, user model.User) error {
	//Retrieve mutex and lock resource
	var mutex, ok = db.mutexes[userID]
	if !ok {
		return model.ErrNotFound
	}

	mutex.Lock()
	defer mutex.Unlock()

	//Write user struct if it
	//actually is registered
	if _, err := db.GetUser(userID); err != nil {
		return err
	}
	return db.setUser(userID, user)
}

//setUser sets the given user to the given @userID.
//This overwrites any existing user or creates a new one,
//on the disk as well as in the collection.
func (db database) setUser(userID string, user model.User) error {
	var path = fmt.Sprintf("%s/%s.xml", db.config.UserDir, userID)
	var err = write(path, user.String())
	db.users[userID] = user
	return err
}

//AddUser creates a new user by creating the user file itself, the
//user's authentication file and an initial calendar file.
func (db database) AddUser(userID, hash string) error {
	//1. Step: Register user mutex and lock the resource.
	//――――――――――――――――――――――――――――――――――――――――――――――――――――
	var userLock sync.Mutex
	db.mutexes[userID] = &userLock
	userLock.Lock()
	defer userLock.Unlock()

	//2. Step: Check whether user is already registered.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――
	if _, err := db.GetUser(userID); err == nil {
		return model.ErrAlreadyExists
	}

	//3. Step: Ensure that target folders actually exists
	//		   before creating the new user.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――
	if err := ensureDir(db.config.DBDir); err != nil {
		return err
	}

	if err := ensureDir(db.config.AuthDir); err != nil {
		return err
	}

	if err := ensureDir(db.config.UserDir); err != nil {
		return err
	}

	if err := ensureDir(db.config.CalendarDir); err != nil {
		return err
	}

	//4. Step: Create authentication file.
	//―――――――――――――――――――――――――――――――――――――
	var path = fmt.Sprintf("%s/%s.xml", db.config.AuthDir, userID)
	var login = model.NewLogin(userID, hash)
	if err := write(path, login.String()); err != nil {
		return err
	}
	db.logins[userID] = login

	//5. Step: Register calendar lock and create initial calendar.
	//		   The initial calendar is placed in a folder
	//   	   dedicated for this user.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	var calID = fmt.Sprintf("%s/%s", userID, userID)

	//Register calendar lock
	var calLock sync.Mutex
	db.mutexes[calID] = &calLock

	//Ensure parent directory of calendar file
	path = fmt.Sprintf("%s/%s", db.config.CalendarDir, userID)
	if err := ensureDir(path); err != nil {
		return err
	}

	//Write calendar struct to disk
	var cal = model.Calendar{
		Name:  model.Attribute{Val: userID},
		Owner: model.Attribute{Val: userID},
	}
	if err := db.setCalendar(calID, cal); err != nil {
		return err
	}

	//6. Step: Create user file itself, linking initial
	//		   calendar to it and writing to disk.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――
	var user = model.NewUser(userID)
	if err := db.AssociateCalendar(user, cal, model.Owner); err != nil {
		return err
	}

	return nil
}

//DeleteUser not only deletes a user itself, but also his authentication file
//and his calendars. Since the calendar can be referenced by multiple users,
//these users must also be found and disassociated from the calendar.
func (db database) DeleteUser(userID string) error {
	//1. Step: Retrieve mutex and lock resource and
	//		   defer call to also delete lock.
	//――――――――――――――――――――――――――――――――――――――――――――――
	var lock = db.mutexes[userID]
	lock.Lock()
	defer lock.Unlock()
	defer delete(db.mutexes, userID)

	//2. Step: Check whether user is already registered.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――
	owner, ok := db.users[userID]
	if !ok {
		return model.ErrNotFound
	}

	//3. Step: Delete user file itself from disk and
	//		   from the user collection.
	//――――――――――――――――――――――――――――――――――――――――――――――――
	delete(db.users, userID)
	var path = fmt.Sprintf("%s/%s.xml", db.config.UserDir, userID)
	if err := os.Remove(path); err != nil {
		return err
	}

	//4. Step: Delete authentication file from disk and
	//         from the authentication collection.
	//――――――――――――――――――――――――――――――――――――――――――――――――――
	delete(db.logins, userID)
	path = fmt.Sprintf("%s/%s.xml", db.config.AuthDir, userID)
	if err := os.Remove(path); err != nil {
		return err
	}

	//5. Step: Find referenced users, so that the calendar
	//		   can be disassociated in their user files.
	//――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	// The user to be deleted has some calendars referenced.
	// These are checked whether the user is their owner before
	// also deleting them.
	for _, reference := range owner.Items.Calendars {
		var calID = reference.Link
		var cal, ok = db.calendars[calID]
		if !ok || cal.Owner.Val != userID {
			continue
		}

		//Going through all users with VIEW and EDIT permissions for
		//this calendar, and editing user struct on the fly.
		for _, userID := range append(cal.Permissions.View.User, cal.Permissions.Edit.User...) {
			var user, exists = db.users[userID.Val]
			if !exists {
				continue
			}

			if err := db.DisassociateCalendar(user, cal); err != nil {
				return err
			}
		}

		delete(db.calendars, calID)
	}

	//6. Step: Delete calendars folder of user to be deleted.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	path = fmt.Sprintf("%s/%s", db.config.CalendarDir, userID)
	if err := os.RemoveAll(path); err != nil {
		return err
	}

	return nil

}

//GetLogin retrieves the login to a given @userID.
//If the user doesn't exist, an error is thrown.
func (db database) GetLogin(userID string) (model.Login, error) {
	var val, ok = db.logins[userID]
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
func (db database) GetCalendar(calID string) (model.Calendar, error) {
	var val, ok = db.calendars[calID]
	if !ok {
		return model.Calendar{}, model.ErrNotFound
	}
	return val, nil
}

//AddCalendar creates a new calendar and appends it to the owner's
//collection of calendars.
func (db database) AddCalendar(ownerID, calName string) error {
	//1. Step: Register mutex for this resource and lock it.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――――
	var calID = fmt.Sprintf("%s/%s", ownerID, calName)
	var lock sync.Mutex
	db.mutexes[calID] = &lock
	lock.Lock()
	defer lock.Unlock()

	//2. Step: Check whether calendar already exists.
	//――――――――――――――――――――――――――――――――――――――――――――――――
	if _, ok := db.calendars[calID]; ok {
		return model.ErrAlreadyExists
	}

	//3. Step: Ensure that target folders actually exists
	//		   before creating the new calendar.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――
	if err := ensureDir(fmt.Sprintf("%s/%s", db.config.CalendarDir, ownerID)); err != nil {
		return err
	}

	//3. Step: Create calendar file and registering it
	//		   in the corresponding collection.
	//――――――――――――――――――――――――――――――――――――――――――――――――――
	var cal = model.Calendar{
		Name:  model.Attribute{Val: calName},
		Owner: model.Attribute{Val: ownerID},
	}
	if err := db.setCalendar(calID, cal); err != nil {
		return err
	}

	//4. Step: Associate calendar to issuing owner
	//		   and writing back.
	//――――――――――――――――――――――――――――――――――――――――――――――
	var owner = db.users[ownerID]
	if err := db.AssociateCalendar(owner, cal, model.Owner); err != nil {
		return err
	}

	return nil
}

//SetCalendar sets the given calendar to the given @calID
//only if the calendar already exists.
func (db database) SetCalendar(calID string, cal model.Calendar) error {
	//Retrieve mutex and lock resource
	var mutex, ok = db.mutexes[calID]
	if !ok {
		return model.ErrNotFound
	}

	mutex.Lock()
	defer mutex.Unlock()

	//Write calendar struct if it
	//actually is registered
	if _, err := db.GetCalendar(calID); err != nil {
		return err
	}
	return db.setCalendar(calID, cal)
}

//setCalendar sets the given calendar to the given @calID.
//This overrides any existing calendar or creates a new one,
//on the disk as well as in the collection.
func (db database) setCalendar(calID string, cal model.Calendar) error {
	var path = fmt.Sprintf("%s/%s.xml", db.config.CalendarDir, calID)
	var err = write(path, cal.String())
	db.calendars[calID] = cal
	return err
}

//DeleteCalendar not only deletes the calendar file behind @calendarID, but
//also removes references to this file.
func (db database) DeleteCalendar(calID string) error {
	//1. Step: Retrieve mutex and lock resource and
	//		   defer call to also delete lock.
	//――――――――――――――――――――――――――――――――――――――――――――――
	var lock = db.mutexes[calID]
	lock.Lock()
	defer lock.Unlock()
	defer delete(db.mutexes, calID)

	//2. Step: Check whether calendar actually exists.
	//―――――――――――――――――――――――――――――――――――――――――――――――――
	cal, ok := db.calendars[calID]
	if !ok {
		return model.ErrNotFound
	}

	//3. Step: Find referenced users and disassociating them from
	//		   the calendar to be deleted and deleting the calendar.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	for _, userID := range append(cal.Permissions.View.User, cal.Permissions.Edit.User...) {
		user, exists := db.users[userID.Val]
		if exists {
			if err := db.DisassociateCalendar(user, cal); err != nil {
				return err
			}
			if err := db.SetUser(user.Name.Val, user); err != nil {
				return err
			}
		}
	}

	//The owner itself is not part of the permission list, since there is no need -
	//he automatically has all permissions. He must be called separately, so that
	//the calendar file can be deleted.
	var ownerID = cal.Owner.Val
	var owner = db.users[ownerID]
	if err := db.DisassociateCalendar(owner, cal); err != nil {
		return err
	}
	delete(db.calendars, calID)

	return nil
}

//DisassociateCalendar removes the calendar from the users collection of
//calendars, so that the updated version can be written back to disk.
//Furthermore, if the user is the owner of the calendar, the original file is also deleted.
func (db database) DisassociateCalendar(user model.User, cal model.Calendar) error {
	var userID = user.Name.Val
	var calID = fmt.Sprintf("%s/%s", cal.Owner.Val, cal.Name.Val)
	var items = user.Items.Calendars
	for i, item := range items {
		if item.Link == calID {
			//The calendar to be removed has been found. Now, another slice of
			//calendars is constructed that can be assigned to the user.
			user.Items.Calendars = append(items[:i], items[i+1:]...)
			if err := db.setUser(userID, user); err != nil {
				return err
			}
			break
		}
	}

	if cal.Owner.Val == userID {
		//The given user also is the owner of the calendar.
		//Therefore the file must be deleted, because otherwise
		//there would no longer be any reference to the file.
		if err := os.Remove(fmt.Sprintf("%s/%s.xml", db.config.CalendarDir, calID)); err != nil {
			return err
		}
	} else {
		//The given user is not the owner; hence, the calendar file
		//mustn't be deleted, but the user has to be removed from the
		//collection of permitted users

		//1. Step: Remove from EDIT permitted users
		var items = cal.Permissions.Edit.User
		for i, item := range items {
			if item.Val == userID {
				cal.Permissions.Edit.User = append(items[:i], items[i+1:]...)
				if err := db.setCalendar(calID, cal); err != nil {
					return err
				}
				return nil
			}
		}

		//2. Step: Remove from VIEW permitted users
		items = cal.Permissions.View.User
		for i, item := range items {
			if item.Val == userID {
				cal.Permissions.View.User = append(items[:i], items[i+1:]...)
				if err := db.setCalendar(calID, cal); err != nil {
					return err
				}
				return nil
			}
		}
	}
	return nil
}

//AssociateCalendar appends the calendar to the collection of the user's calendars,
//if it hasn't been associated to the user yet and also links the user in the
//calendar file itself.
func (db database) AssociateCalendar(user model.User, cal model.Calendar, perm model.Permission) error {
	//If any of the iterated items/calendars has the same id as the calendar to
	//be associated, an error is thrown, because the element is already there.
	var userID = user.Name.Val
	var calID = fmt.Sprintf("%s/%s", cal.Owner.Val, cal.Name.Val)
	for _, reference := range user.Items.Calendars {
		if reference.Link == calID {
			return model.ErrAlreadyExists
		}
	}

	//Append the calendar to the user's
	//collection of calendars.
	var items = user.Items.Calendars
	var appendix = model.CalendarReference{
		XMLName: xml.Name{Local: "calendar"},
		Link:    calID,
		Perm: 	 perm.String(),
	}
	user.Items.Calendars = append(items, appendix)
	if err := db.setUser(userID, user); err != nil {
		return err
	}

	//Link the user itself to the calendar.
	if perm == model.Owner {
		cal.Owner.Val = userID
	} else {
		var entry = model.Attribute{
			Val: userID,
		}
		var users []model.Attribute

		if perm == model.Read {
			users = cal.Permissions.View.User
			cal.Permissions.View.User = append(users, entry)
		} else if perm == model.Edit {
			users = cal.Permissions.Edit.User
			cal.Permissions.Edit.User = append(users, entry)
		}
	}
	if err := db.setCalendar(calID, cal); err != nil {
		return err
	}

	return nil
}
