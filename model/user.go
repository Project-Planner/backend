package model

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

type User struct {
	XMLName xml.Name  `xml:"user"`
	Name    Attribute `xml:"name"`
	Items   Items     `xml:"items"`
}

type Items struct {
	XMLName   xml.Name            `xml:"items"`
	Calendars []CalendarReference `xml:"calendar"`
}

type CalendarReference struct {
	XMLName xml.Name `xml:"calendar"`
	Link    string   `xml:"href,attr"`
}

func NewUser(name string) User {
	return User{
		Name:  Attribute{Val: name},
		Items: Items{},
	}
}

//DisassociateCalendar removes the calendar with the name @calName from the
//issuing users collection of calendars, so that the updated version can be
//written back to disk. Furthermore, if the user is the owner of the calendar,
//the original file is also deleted.
func (user *User) DisassociateCalendar(path, calID string, cal Calendar) {
	var items = user.Items.Calendars
	for i, item := range items {
		if item.Link == calID {
			//The calendar to be removed has been found. Now, another slice of
			//calendars is constructed that can be assigned to the user.
			user.Items.Calendars = append(items[:i], items[i+1:]...)

		}
	}

	if cal.Owner.Val == user.Name.Val {
		//The issuing user is also the owner of the calendar.
		//Therefore the file must be deleted, because otherwise
		//there would no longer be any reference to the file.
		if err := os.Remove(fmt.Sprintf("%s/%s.xml", path, calID)); err != nil {
			log.Fatal(err)
		}
	}
}

//AssociateCalendar appends the calendar to the collection of the user's calendars,
//if it hasn't been associated to the user yet and also links the user in the
//calendar file itself.
func (user *User) AssociateCalendar(perm Permission, calID string, db Database) error {
	//If any of the iterated items/calendars has the same id as the calendar to
	//be associated, an error is thrown, because the element is already there.
	for _, cal := range user.Items.Calendars {
		if cal.Link == calID {
			return ErrAlreadyExists
		}
	}

	//Append the calendar to the user's collection
	//of calendars.
	var userID = user.Name.Val
	var items = user.Items.Calendars
	var appendix = CalendarReference{
		XMLName: xml.Name{Local: "calendar"},
		Link:    calID,
	}
	user.Items.Calendars = append(items, appendix)
	db.SetUser(userID, *user)

	//Link the user itself to the calendar.
	var cal, _ = db.GetCalendar(calID)
	if perm == Owner {
		cal.Owner.Val = userID
	} else {
		var entry = Attribute{
			Val: userID,
		}
		var users []Attribute

		if perm == Read {
			users = cal.Permissions.View.User
			cal.Permissions.View.User = append(users, entry)
		} else if perm == Edit {
			users = cal.Permissions.Edit.User
			cal.Permissions.Edit.User = append(users, entry)
		}
	}
	db.SetCalendar(calID, cal)

	return nil
}

func (user User) String() string {
	var parsed, _ = xml.MarshalIndent(user, "", "\t")
	return string(parsed)
}
