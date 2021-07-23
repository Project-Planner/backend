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

	fmt.Println(users)

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

//SetUser is the synchronized version of setUser used
//for concurrent modification.
//Furthermore, it only executes its internal variant
//if the user yet exists.
func (db database) SetUser(userID string, user model.User) error {
	//Obtain mutex and lock resource
	var mutex, ok = db.mutexes[userID]
	if !ok {
		return model.ErrNotFound
	}

	mutex.Lock()
	defer mutex.Unlock()

	//Overwrite only if the user yet exists
	if _, err := db.GetUser(userID); err != nil {
		return err
	}
	return db.setUser(userID, user)
}

//setUser writes the given user data for the user
//with the given @userID. This function overwrites any
//existing user file or creates a new one; on the disk
//as well as in the collection.
func (db database) setUser(userID string, user model.User) error {
	var path = fmt.Sprintf("%s/%s.xml", db.config.UserDir, userID)
	var err = write(path, user.String())
	db.users[userID] = user
	return err
}

//AddUser makes a new user by creating respective files
//(user file, authentication file, initial calendar file)
//and registering its newly created resources along with
//their locks in the collections.
func (db database) AddUser(userID, hash string) error {
	//1. Step: Check whether user is already registered
	//		   before creating resources multiple times.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――
	if _, ok := db.mutexes[userID]; !ok {
		//2. Step: Make mutex and lock resources.
		//――――――――――――――――――――――――――――――――――――――――
		var mutex = new(sync.Mutex)
		db.mutexes[userID] = mutex
		mutex.Lock()
		defer mutex.Unlock()
	} else {
		return model.ErrAlreadyExists
	}

	//3. Step: Ensure that target folders actually
	// 		   exists before creating the new user.
	//――――――――――――――――――――――――――――――――――――――――――――――
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

	//4. Step: Make authentication file and
	//		   register resource in collection.
	//――――――――――――――――――――――――――――――――――――――――――
	var path = fmt.Sprintf("%s/%s.xml", db.config.AuthDir, userID)
	var login = model.NewLogin(userID, hash)
	if err := write(path, login.String()); err != nil {
		return err
	}
	db.logins[userID] = login

	//5. Step: Make user file itself and
	//		   register resource in collection.
	//――――――――――――――――――――――――――――――――――――――――――
	var user = model.NewUser(userID)
	db.users[userID] = user

	//6. Step: Associate owner and initial calendar by adding
	//		   the initial calendar to the user.
	//		   (this implies initially writing to disk and
	//		   collection since included in called function).
	//――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	if err := db.addCalendar(userID, userID); err != nil {
		return err
	}

	return nil
}

//DeleteUser deletes the user itself, its authentication file and his
//calendars. Since calendars can be referenced by multiple other users,
//these must be found in order to remove their references to the calendars
//to be deleted.
func (db database) DeleteUser(userID string) error {
	//1. Step: Check whether user actually exists,
	//		   before deleting the resource and
	//		   lock resource.
	//―――――――――――――――――――――――――――――――――――――――――――――
	var mutex, ok = db.mutexes[userID]
	if !ok {
		return model.ErrNotFound
	}
	mutex.Lock()
	defer mutex.Unlock()
	defer delete(db.mutexes, userID)

	//3. Step: Delete authentication file from disk and
	//         from the authentication collection.
	//――――――――――――――――――――――――――――――――――――――――――――――――――
	delete(db.logins, userID)
	var path = fmt.Sprintf("%s/%s.xml", db.config.AuthDir, userID)
	if err := os.Remove(path); err != nil {
		return err
	}

	//4. Step: Delete the user's calendars and remove
	//		   their references in the other users' files.
	//――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	// The user to be deleted has some calendars referenced.
	// These are checked whether the user is their owner before
	// also deleting them.
	owner, ok := db.users[userID]
	if !ok {
		return model.ErrNotFound
	}

	for _, reference := range owner.Items.Calendars {
		var calID = reference.Link
		var cal, ok = db.calendars[calID]
		if !ok || cal.Owner.Val != userID {
			continue
		}

		//Delete the calendar and disconnect it
		//from referenced users.
		var calLock = db.mutexes[calID]
		calLock.Lock()
		if err := db.deleteCalendar(calID); err != nil {
			return err
		}
		calLock.Unlock()
	}

	//5. Step: Delete calendars folder of user to be deleted.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	path = fmt.Sprintf("%s/%s", db.config.CalendarDir, userID)
	if err := os.RemoveAll(path); err != nil {
		return err
	}

	//2. Step: Delete user file itself from disk and
	//		   from the user collection.
	//――――――――――――――――――――――――――――――――――――――――――――――――
	delete(db.users, userID)
	path = fmt.Sprintf("%s/%s.xml", db.config.UserDir, userID)
	if err := os.Remove(path); err != nil {
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

//AddCalendar is the synchronized version of addCalendar
//used for concurrent modification.
func (db database) AddCalendar(ownerID, calName string) error {
	//Obtain mutexes
	var mutex, ok = db.mutexes[ownerID]
	if !ok {
		return model.ErrNotFound
	}

	//Lock resources
	mutex.Lock()
	defer mutex.Unlock()

	//Call unsafe method
	return db.addCalendar(ownerID, calName)
}

//addCalendar makes a new calendar and appends it to the owner's
//collection of calendars.
func (db database) addCalendar(ownerID, calName string) error {
	//1. Step: Check whether the calendar is already registered
	//		   before creating resources multiple times.
	//――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	var calID = fmt.Sprintf("%s/%s", ownerID, calName)
	if _, ok := db.mutexes[calID]; !ok {
		//2. Step: Make mutex and lock resources.
		//――――――――――――――――――――――――――――――――――――――――
		var mutex = new(sync.Mutex)
		db.mutexes[calID] = mutex
		mutex.Lock()
		defer mutex.Unlock()
	} else {
		return model.ErrAlreadyExists
	}

	//3. Step: Ensure that target folders actually exists
	//		   before creating the new calendar.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――
	if err := ensureDir(fmt.Sprintf("%s/%s", db.config.CalendarDir, ownerID)); err != nil {
		return err
	}

	//4. Step: Make basic calendar and associate
	//		   it to the owner.
	//―――――――――――――――――――――――――――――――――――――――――――
	var cal = model.Calendar{
		Name:  model.Attribute{Val: calName},
		Owner: model.Attribute{Val: ownerID},
		ID: model.Attribute{Val: fmt.Sprintf("%s/%s", ownerID, calName)},
	}

	var owner = db.users[ownerID]
	if err := db.associateCalendar(owner, cal, model.Owner); err != nil {
		return err
	}

	return nil
}

//SetCalendar sets the given calendar to the given @calID
//only if the calendar already exists.
//Furthermore, it only executes its internal variant
//if the calendar yet exists.
func (db database) SetCalendar(calID string, cal model.Calendar) error {
	//Obtain mutex and lock resource
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

//DeleteCalendar is the synchronized version of deleteCalendar
//used for concurrent modification.
func (db database) DeleteCalendar(calID string) error {
	//Obtain mutexes
	var calMutex, ok = db.mutexes[calID]
	if !ok {
		return model.ErrNotFound
	}

	var ownerID = strings.Split(calID, "/")[0]
	ownerMutex, ok := db.mutexes[ownerID]
	if !ok {
		return model.ErrNotFound
	}

	//Lock resources
	calMutex.Lock()
	defer calMutex.Unlock()

	ownerMutex.Lock()
	defer ownerMutex.Unlock()

	//Call unsafe method
	return db.deleteCalendar(calID)
}

//deleteCalendar deletes the calendar file behind @calID and removes
//links in the referenced user files.
func (db database) deleteCalendar(calID string) error {
	//1. Step: Check whether calendar actually exists
	//		   before deleting the resource and lock
	//		   resource.
	//―――――――――――――――――――――――――――――――――――――――――――――――――
	if _, ok := db.mutexes[calID]; !ok {
		return model.ErrNotFound
	}

	//2. Step: Find referenced users and disassociate them from
	//		   from the calendar and then delete the calendar itself.
	//		   Each user resource must also be locked before modification.
	//―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	cal, ok := db.calendars[calID]
	if !ok {
		return model.ErrNotFound
	}

	for _, userID := range append(cal.Permissions.View.User, cal.Permissions.Edit.User...) {
		user, exists := db.users[userID.Val]
		if exists {
			//Locking the user at this point is crucial
			//in order to prevent write errors when
			//calling following unsafe functions
			var userMutex = db.mutexes[userID.Val]
			userMutex.Lock()
			if err := db.disassociateCalendar(user, cal); err != nil {
				return err
			}
			userMutex.Unlock()
		}
	}

	//The owner itself is not part of the permission list, since there is no need;
	//he automatically has all permissions. He must be disassociated separately,
	//so that the calendar file can be deleted.
	var ownerID = cal.Owner.Val
	var owner = db.users[ownerID]
	if err := db.disassociateCalendar(owner, cal); err != nil {
		return err
	}

	delete(db.calendars, calID)
	delete(db.mutexes, calID)

	return nil
}

//DisassociateCalendar is the synchronized version of disassociateCalendar
//used for concurrent modification.
func (db database) DisassociateCalendar(user model.User, cal model.Calendar) error {
	//Obtain mutexes
	userMutex, ok := db.mutexes[user.Name.Val]
	if !ok {
		return model.ErrNotFound
	}

	var calID = fmt.Sprintf("%s/%s", cal.Owner.Val, cal.Name.Val)
	calMutex, ok := db.mutexes[calID]
	if !ok {
		return model.ErrNotFound
	}

	//Lock resources
	userMutex.Lock()
	defer userMutex.Unlock()

	calMutex.Lock()
	defer calMutex.Unlock()

	//Call unsafe method
	return db.disassociateCalendar(user, cal)
}

//disassociateCalendar removes the calendar from the user's collection of
//calendars, so that the updated version can be written back to disk.
//Furthermore, if the user is the owner of the calendar, the original file
//is also deleted.
func (db database) disassociateCalendar(user model.User, cal model.Calendar) error {
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

//AssociateCalendar is the synchronized version of associatedCalendar
//used for concurrent modification.
func (db database) AssociateCalendar(user model.User, cal model.Calendar, perm model.Permission) error {
	//Obtain mutexes
	userMutex, ok := db.mutexes[user.Name.Val]
	if !ok {
		return model.ErrNotFound
	}

	var calID = fmt.Sprintf("%s/%s", cal.Owner.Val, cal.Name.Val)
	calMutex, ok := db.mutexes[calID]
	if !ok {
		return model.ErrNotFound
	}

	//Lock resources
	userMutex.Lock()
	defer userMutex.Unlock()

	calMutex.Lock()
	defer calMutex.Unlock()

	//Call unsafe method
	return db.associateCalendar(user, cal, perm)
}

//associateCalendar appends the calendar to the collection of the user's calendars,
//if it hasn't been associated to this user yet and also links the references the
//user in the calendar file itself.
func (db database) associateCalendar(user model.User, cal model.Calendar, perm model.Permission) error {
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
		Perm:    perm.String(),
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
