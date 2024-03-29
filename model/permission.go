package model

type Permission int

const (
	None Permission = iota
	Read
	Edit
	Owner
)

func (p Permission) String() string {
	switch p {
	case None:
		return "none"
	case Read:
		return "view"
	case Edit:
		return "edit"
	case Owner:
		return "owner"
	}
	panic("permission: string: not implemented")
}

// CalendarPermissions takes a calendar c, and a user id and returns which permissions this user has for c.
func CalendarPermissions(c Calendar, userID string) Permission {
	if c.Owner.Val == userID {
		return Owner
	}

	s := func(u []Attribute, n string) bool {
		for _, v := range u {
			if v.Val == n {
				return true
			}
		}
		return false
	}

	if s(c.Permissions.Edit.User, userID) {
		return Edit
	}

	if s(c.Permissions.View.User, userID) {
		return Read
	}

	return None
}
