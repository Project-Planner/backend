package xmldb

import (
	"errors"
	"fmt"
	"github.com/Project-Planner/backend/model"
	"github.com/google/uuid"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// The struct implementing the web.Database interface
type Database struct {
	config DBConfig
	auth   map[string]model.Login
	users  map[string]model.User
}

func NewDatabase(config DBConfig) (Database, error) {
	// 1. Step: Ensuring that the parent folders (auth, user, calendars) exist.
	ensureDir(config.RootDir)
	ensureDir(config.AuthDir)
	ensureDir(config.UserDir)
	ensureDir(config.CalendarDir)

	//2. Step: Parsing source files into corresponding structs or struct collections.
	// #Auth#
	// a) Ensuring that the authentication file actually exists.
	//	  If it doesn't, then it is created, so that parsing
	//	  can go on.
	var authPath = fmt.Sprintf("%s/%s", config.AuthDir, "/auth.xml")
	var authPayload model.Auth
	var auth = make(map[string]model.Login)
	if !exists(authPath) {
		create(authPath, authPayload.ToString())
		log.Print("Authentication file didn't exist - was created now.")
	}

	// b) Parsing authentication file into struct.
	parse(authPath, &authPayload)
	for _, login := range authPayload.Logins {
		auth[login.Name.Value] = login
	}

	// #Users#
	// Users are each represented by a separate file inside the corresponding folder.
	// Each file to be iterated over in this folder yields the users unique id (UUID)
	// and content for the user struct to be parsed with.
	var users = make(map[string]model.User)
	filepath.Walk(config.UserDir, func(path string, info os.FileInfo, err error) error {

		//For completeness, it is tested, whether the file
		//has an "xml" file extension and hence can be seen as such.
		var name = info.Name()
		var matching, _ = regexp.MatchString(".+\\.xml$", name)

		if matching {
			var index = strings.LastIndex(name, ".xml")
			var uuid = name[0:index]

			var user model.User
			parse(path, &user)
			users[uuid] = user
		}

		return nil
	})

	return Database{config, auth, users}, nil
}

func (db Database) AddUser(name, hash string) error {

	//1. Step: Checking whether name is already in use.
	if _, exists := db.auth[name]; !exists {

		//2. Step: Ensuring that the target folder (user directory)
		//		   actually exists before creating the new user file.
		ensureDir(db.config.UserDir)

		//2. Step: Generating UUID and creating corresponding user file.
		var uuid = uuid.New().String()
		var path = fmt.Sprintf("%s/%s.xml", db.config.UserDir, uuid)
		var user = model.NewUser(name)
		create(path, user.ToString())

		//3. Step: Registering authentication for this user.
		var login = model.NewLogin(name, hash)
		db.auth[name] = login

		//4. Step: Transferring logins into struct that can be parsed.
		var logins = []model.Login{}
		for _, value := range db.auth {
			logins = append(logins, value)
		}

		var auth = model.NewAuth()
		auth.Logins = logins
		path = fmt.Sprintf("%s/%s", db.config.AuthDir, db.config.AuthSubpath)
		setFile(path, auth.ToString())
		return nil

	} else {
		return errors.New(fmt.Sprintf("User '%s' is already registered.\n", name))
	}
}

func (db Database) GetLogin(name string) (model.Login, error) {
	var val, exists = db.auth[name]
	if exists {
		return val, nil
	} else {
		return model.Login{}, model.ErrNotFound
	}
}

// Example implementation of a method
// func (db database) GetUser(id string) (model.User, error) {}
