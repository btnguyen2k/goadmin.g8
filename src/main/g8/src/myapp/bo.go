package myapp

const (
	SystemGroupId      = "system"
	AdminUserUsernname = "admin"
)

const (
	fieldGroupId   = "id"
	fieldGroupName = "name"
)

// Group represents a user group
type Group struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// GroupDao defines API to access user group storage
type GroupDao interface {
	Create(id, name string) (bool, error)
	Get(id string) (*Group, error)
}

const (
	fieldUserUsername = "uname"
	fieldUserPassword = "pwd"
	fieldUserName     = "name"
	fieldUserGroupId  = "gid"
)

// User represents a user account
type User struct {
	Username string `json:"uname"`
	Password string `json:"pwd"`
	Name     string `json:"name"`
	GroupId  string `json:"gid"`
}

// UserDao defines API to access user account storage
type UserDao interface {
	Create(username, encryptedPassword, name, groupId string) (bool, error)
	Get(username string) (*User, error)
}
