package myapp

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"math"
	"runtime"
	"strings"

	"github.com/btnguyen2k/consu/reddo"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

const (
	systemGroupId = "system" // reserved id for the "system" group
)

var (
	systemUserUsername = "admin"         // reserved username for the "system" user
	systemUserName     = "Administrator" // reserved name for the "system" user
)

const (
	flashPrefixInfo    = "_I_:"
	flashPrefixWarning = "_W_:"
	flashPrefixError   = "_E_:"
)

func getSession(c echo.Context) *sessions.Session {
	sess, _ := session.Get(namespace, c)
	return sess
}

func setSessionValue(c echo.Context, key string, value interface{}) {
	sess := getSession(c)
	if value == nil {
		delete(sess.Values, key)
	} else {
		sess.Values[key] = value
	}
	sess.Save(c.Request(), c.Response())
}

func addFlashMsg(c echo.Context, msg string) {
	sess := getSession(c)
	sess.AddFlash(msg)
	sess.Save(c.Request(), c.Response())
}

func encryptPassword(salt, rawPassword string) string {
	saltAndPwd := salt + "." + rawPassword
	out := sha1.Sum([]byte(saltAndPwd))
	return strings.ToLower(hex.EncodeToString(out[:]))
}

func getCurrentUser(c echo.Context) (*User, error) {
	sess := getSession(c)
	if uid, has := sess.Values[sessionMyUid]; has {
		uid, _ = reddo.ToString(uid)
		if uid != nil {
			username := uid.(string)
			return userDao.Get(username)
		}
	}
	return nil, nil
}

// available since template-r3
func getCookieString(c echo.Context, cookieName string) string {
	cookie, err := c.Cookie(cookieName)
	if err != nil || cookie == nil {
		return ""
	}
	return cookie.Value
}

// available since template-r3
func getContextString(c echo.Context, key string) string {
	val, err := reddo.ToString(c.Get(key))
	if err != nil {
		return ""
	}
	return val
}

/*----------------------------------------------------------------------*/

type OsUtils struct {
}

func (u *OsUtils) CpuCores() int {
	return runtime.NumCPU()
}

func (u *OsUtils) CpuLoad() float64 {
	stats, err := load.Avg()
	if err != nil || stats == nil {
		return -1
	}
	return math.Floor(stats.Load1*100) / 100
}

func (u *OsUtils) MemoryUsed() uint64 {
	v, err := mem.VirtualMemory()
	if err != nil || v == nil {
		return 0
	}
	return v.Used
}

func (u *OsUtils) MemoryUsedKb() float64 {
	return float64(u.MemoryUsed()) / 1014
}

func (u *OsUtils) MemoryUsedMb() float64 {
	return float64(u.MemoryUsed()) / 1024 / 1024
}

func (u *OsUtils) MemoryUsedGb() float64 {
	return float64(u.MemoryUsed()) / 1024 / 1024 / 1024
}

func (u *OsUtils) MemoryFree() uint64 {
	v, err := mem.VirtualMemory()
	if err != nil || v == nil {
		return 0
	}
	return v.Free
}

func (u *OsUtils) MemoryFreeKb() float64 {
	f := float64(u.MemoryFree()) / 1024.0
	return math.Floor(f*100.0) / 100.0
}

func (u *OsUtils) MemoryFreeMb() float64 {
	f := float64(u.MemoryFree()) / 1024.0 / 1024.0
	return math.Floor(f*100.0) / 100.0
}

func (u *OsUtils) MemoryFreeGb() float64 {
	f := float64(u.MemoryFree()) / 1024.0 / 1024.0 / 1024.0
	return math.Floor(f*100.0) / 100.0
}

func (u *OsUtils) MemoryFreePercent() float64 {
	v, err := mem.VirtualMemory()
	if err != nil || v == nil {
		return -1
	}
	p := float64(v.Free) / float64(v.Total)
	return math.Floor(p * 100.0)
}

func (u *OsUtils) AppMemUsed() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

func (u *OsUtils) AppMemUsedKb() float64 {
	f := float64(u.AppMemUsed()) / 1024.0
	return math.Floor(f*100.0) / 100.0
}

func (u *OsUtils) AppMemUsedMb() float64 {
	f := float64(u.AppMemUsed()) / 1024.0 / 1024.0
	return math.Floor(f*100.0) / 100.0
}

func (u *OsUtils) AppMemUsedGb() float64 {
	f := float64(u.AppMemUsed()) / 1024.0 / 1024.0 / 1024.0
	return math.Floor(f*100.0) / 100.0
}

func (u *OsUtils) GoNumRoutines() int {
	return runtime.NumGoroutine()
}

/*----------------------------------------------------------------------*/

type MyAppUtils struct {
	c echo.Context
}

func (u *MyAppUtils) NumUserGroups() int {
	if groupList, err := groupDao.GetAll(); err != nil {
		log.Printf("error while getting user groups: %e", err)
		return -1
	} else {
		return len(groupList)
	}
}

func (u *MyAppUtils) AllUserGroups() []*GroupModel {
	if groupList, err := groupDao.GetAll(); err != nil {
		log.Printf("error while getting user groups: %e", err)
		return make([]*GroupModel, 0)
	} else {
		return toGroupModelList(u.c, groupList)
	}
}

func (u *MyAppUtils) NumUsers() int {
	if userList, err := userDao.GetAll(); err != nil {
		log.Printf("error while getting users: %e", err)
		return -1
	} else {
		return len(userList)
	}
}

func (u *MyAppUtils) AllUsers() []*UserModel {
	if userList, err := userDao.GetAll(); err != nil {
		log.Printf("error while getting users: %e", err)
		return make([]*UserModel, 0)
	} else {
		return toUserModelList(u.c, userList)
	}
}
