package myapp

import "github.com/labstack/echo/v4"

func toGroupModel(c echo.Context, g *Group) *GroupModel {
	if g == nil {
		return nil
	}
	return &GroupModel{c: c, Group: g}
}

func toGroupModelList(c echo.Context, groupList []*Group) []*GroupModel {
	result := make([]*GroupModel, 0)
	for _, g := range groupList {
		result = append(result, toGroupModel(c, g))
	}
	return result
}

// GroupModel represents a user group model to be used in view
type GroupModel struct {
	c echo.Context
	*Group
}

func (m *GroupModel) CanDelete() bool {
	// cannot delete system-group
	return m.Id != SystemGroupId
}

func (m *GroupModel) UrlDelete() string {
	return m.c.Echo().Reverse(actionNameCpDeleteGroup) + "?id=" + m.Id
}

func (m *GroupModel) UrlEdit() string {
	return m.c.Echo().Reverse(actionNameCpEditGroup) + "?id=" + m.Id
}

/*----------------------------------------------------------------------*/

func toUserModel(c echo.Context, u *User) *UserModel {
	if u == nil {
		return nil
	}
	return &UserModel{c: c, User: u}
}

func toUserModelList(c echo.Context, userList []*User) []*UserModel {
	result := make([]*UserModel, 0)
	for _, u := range userList {
		result = append(result, toUserModel(c, u))
	}
	return result
}

// UserModel represents a user model to be used in view
type UserModel struct {
	c echo.Context
	*User
}

func (m *UserModel) IsSystemUser() bool {
	return m.GroupId == SystemGroupId
}

func (m *UserModel) CanDelete() bool {
	// cannot delete system-user
	return m.Username != AdminUserUsernname
}

func (m *UserModel) UrlDelete() string {
	return m.c.Echo().Reverse(actionNameCpDeleteUser) + "?u=" + m.Username
}

func (m *UserModel) UrlEdit() string {
	return m.c.Echo().Reverse(actionNameCpEditUser) + "?u=" + m.Username
}
