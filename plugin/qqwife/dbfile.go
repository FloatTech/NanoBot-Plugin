// Package qqwife 娶群友
package qqwife

import (
	"errors"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	nano "github.com/fumiama/NanoBot"

	fcext "github.com/FloatTech/floatbox/ctxext"
	"github.com/FloatTech/gg"
	sql "github.com/FloatTech/sqlite"
)

type userInfo struct {
	ID     string
	Nick   string
	Avatar string
}

type dbData struct {
	db *sql.Sqlite
	sync.RWMutex
}

// 群设置
type setting struct {
	GID      string
	LastTime int
	CanMatch int   // 嫁婚开关
	CanNtr   int   // Ntr开关
	CDtime   int64 // CD时间
}

// 结婚证信息
type marriage struct {
	Users      string // 双方QQ号
	Sname      string // 户主名称
	Mname      string // 对象名称
	Updatetime string // 登记时间
	Spic       string
	Mpic       string
}

// 预留10个为后续扩展
type cdsheet struct {
	User string `db:"User"`
	Mar  int64  `db:"CD0"` // 娶
	Rob  int64  `db:"CD1"` // 强
	Lef  int64  `db:"CD2"` // 离
	Buy  int64  `db:"CD3"` // 礼
	MMk  int64  `db:"CD4"` // 做媒
	CD5  int64  `db:"CD5"`
	CD6  int64  `db:"CD6"`
	CD7  int64  `db:"CD7"`
	CD8  int64  `db:"CD8"`
	CD9  int64  `db:"CD9"`
}

// 好感度列表
type favor struct {
	Users string // 双方QQ号
	Favor int    // 好感度
}

var (
	wifeData = &dbData{
		db: &sql.Sqlite{},
	}
	getdb = fcext.DoOnceOnSuccess(func(ctx *nano.Ctx) bool {
		wifeData.db.DBPath = engine.DataFolder() + "结婚登记表.db"
		err := wifeData.db.Open(time.Hour)
		if err == nil {
			// 创建群配置表
			err = wifeData.db.Create("setting", &setting{})
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "[dbfile.go.80 ->ERROR]:", err)
				return false
			}
			// 创建CD表
			err = wifeData.db.Create("cdsheet", &cdsheet{})
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "[dbfile.go.86 ->ERROR]:", err)
				return false
			}
			// 创建好感度表
			err = wifeData.db.Create("favor", &favor{})
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "[dbfile.go.92 ->ERROR]:", err)
				return false
			}
			// 刷新列表
			err = wifeData.refresh(ctx.Message.ChannelID)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "[dbfile.go.98 ->ERROR]:", err)
				return false
			}
			return true
		}
		_, _ = ctx.SendPlainMessage(false, "[dbfile.go.103 ->ERROR]:", err)
		return false
	})
)

func getUserInfoIn(ctx *nano.Ctx, gid, uid string) (info userInfo, err error) {
	uInfo, err := ctx.GetGuildMemberOf(gid, uid)
	if err != nil {
		return
	}
	nick := uInfo.Nick
	if nick == "" {
		nick = uInfo.User.Username
	}
	info = userInfo{
		ID:     uInfo.User.ID,
		Nick:   nick,
		Avatar: uInfo.User.Avatar,
	}
	return
}

func (sql *dbData) getSet(gid string) (dbinfo setting, err error) {
	sql.Lock()
	defer sql.Unlock()
	// 创建群表格
	err = sql.db.Create("setting", &dbinfo)
	if err != nil {
		return
	}
	if !sql.db.CanFind("setting", "where gid is "+gid) {
		// 没有记录
		return setting{
			GID:      gid,
			CanMatch: 1,
			CanNtr:   1,
			CDtime:   12,
		}, nil
	}
	_ = sql.db.Find("setting", &dbinfo, "where gid is "+gid)
	return
}

func (sql *dbData) updateSet(dbinfo setting) error {
	sql.Lock()
	defer sql.Unlock()
	return sql.db.Insert("setting", &dbinfo)
}

func (sql *dbData) refresh(gid string) error {
	sql.Lock()
	defer sql.Unlock()
	// 创建群表格
	err := sql.db.Create("setting", &setting{})
	if err != nil {
		return err
	}
	if !sql.db.CanFind("setting", "where gid is "+gid) {
		return nil
	}
	dbinfo := setting{}
	_ = sql.db.Find("setting", &dbinfo, "where gid is "+gid)
	if time.Now().Day() != dbinfo.LastTime && time.Now().Hour() >= 4 {
		_ = sql.db.Drop("group" + gid)
		// 更新数据时间
		dbinfo.GID = gid
		dbinfo.LastTime = time.Now().Day()
		return sql.db.Insert("setting", &dbinfo)
	}
	return nil
}

func (sql *dbData) checkUser(gid, uid string) (userinfo marriage, err error) {
	sql.Lock()
	defer sql.Unlock()
	gidstr := "group" + gid
	// 创建群表格
	err = sql.db.Create(gidstr, &userinfo)
	if err != nil {
		return
	}
	if !sql.db.CanFind(gidstr, "where Users glob '*"+uid+"*'") {
		return
	}
	err = sql.db.Find(gidstr, &userinfo, "where Users glob '*"+uid+"*'")
	return
}

// 民政局登记数据
func (sql *dbData) register(gid string, uid, target userInfo) error {
	sql.Lock()
	defer sql.Unlock()
	gidstr := "group" + gid
	uidinfo := marriage{
		Users:      uid.ID + " & " + target.ID,
		Sname:      uid.Nick,
		Spic:       uid.Avatar,
		Mname:      target.Nick,
		Mpic:       target.Avatar,
		Updatetime: time.Now().Format("15:04:05"),
	}
	return sql.db.Insert(gidstr, &uidinfo)
}

// 民政局离婚
func (sql *dbData) divorce(gid, uid string) error {
	sql.Lock()
	defer sql.Unlock()
	gidstr := "group" + gid
	// 创建群表格
	userinfo := marriage{}
	err := sql.db.Create(gidstr, &userinfo)
	if err != nil {
		return err
	}
	if !sql.db.CanFind(gidstr, "where Users glob '*"+uid+"*'") {
		return errors.New("user(" + uid + ") not found")
	}
	return sql.db.Del(gidstr, "where Users glob '*"+uid+"*'")
}

func (sql *dbData) getlist(gid string) (list [][4]string, err error) {
	sql.Lock()
	defer sql.Unlock()
	gidstr := "group" + gid
	number, _ := sql.db.Count(gidstr)
	if number <= 0 {
		return
	}
	var info marriage
	err = sql.db.FindFor(gidstr, &info, "GROUP BY Users", func() error {
		users := strings.Split(info.Users, " & ")
		if users[0] == "" || users[1] == "" {
			return nil
		}
		dbinfo := [4]string{
			info.Sname,
			users[0],
			info.Mname,
			users[1],
		}
		list = append(list, dbinfo)
		return nil
	})
	return
}

func slicename(name string, canvas *gg.Context) (resultname string) {
	usermane := []rune(name) // 将每个字符单独放置
	widthlen := 0
	numberlen := 0
	for i, v := range usermane {
		width, _ := canvas.MeasureString(string(v)) // 获取单个字符的宽度
		widthlen += int(width)
		if widthlen > 350 {
			break // 总宽度不能超过350
		}
		numberlen = i
	}
	if widthlen > 350 {
		resultname = string(usermane[:numberlen-1]) + "......" // 名字切片
	} else {
		resultname = name
	}
	return
}

func (sql *dbData) favorFor(uid, target string, add int) (favorValue int, err error) {
	sql.Lock()
	defer sql.Unlock()
	// 创建群表格
	err = sql.db.Create("favor", &favor{})
	if err != nil {
		return
	}
	number, _ := sql.db.Count("favor")
	if number <= 0 {
		return
	}
	key := uid + " & " + target
	uidInt64, _ := strconv.ParseInt(uid, 10, 64)
	targetInt64, _ := strconv.ParseInt(target, 10, 64)
	if uidInt64 < targetInt64 {
		key = target + " & " + uid
	}
	info := favor{}
	err = sql.db.Find("favor", &info, "where Userinfo is '"+key+"'")
	if add > 0 {
		info.Users = key
		info.Favor += add
		err = sql.db.Insert("favor", &info)
	}
	return info.Favor, err
}

func (sql *dbData) getGroupFavorability(uid string) (list []favor, err error) {
	sql.RLock()
	defer sql.RUnlock()
	info := favor{}
	err = sql.db.FindFor("favorability", &info, "where Userinfo glob '*"+uid+"*' AND Favor > 0 ORDER BY DESC", func() error {
		var target string
		userList := strings.Split(info.Users, " & ")
		switch {
		case len(userList) == 0:
			return errors.New("好感度系统数据存在错误")
		case userList[0] == uid:
			target = userList[1]
		default:
			target = userList[0]
		}
		list = append(list, favor{
			Users: target,
			Favor: info.Favor,
		})
		return nil
	})
	return
}

func (sql *dbData) checkCD(gid, uid string, funcType string) (cdTime time.Duration, err error) {
	setting, err := wifeData.getSet(gid)
	if err != nil {
		return
	}
	sql.Lock()
	defer sql.Unlock()
	// 创建群表格
	err = sql.db.Create("cdsheet", &cdsheet{})
	if err != nil {
		return
	}
	number, _ := sql.db.Count("cdsheet")
	if number <= 0 {
		return
	}
	info := cdsheet{}
	if !sql.db.CanFind("cdsheet", "where User is '"+uid+"'") {
		return
	}
	err = sql.db.Find("cdsheet", &info, "where User is '"+uid+"'")
	if err != nil {
		return
	}
	switch funcType {
	case "娶", "嫁":
		cdTime = time.Duration(setting.CDtime)*time.Hour - time.Since(time.Unix(info.Mar, 0))
	case "牛":
		cdTime = time.Duration(setting.CDtime)*time.Hour - time.Since(time.Unix(info.Rob, 0))
	case "离":
		cdTime = time.Duration(setting.CDtime)*time.Hour - time.Since(time.Unix(info.Lef, 0))
	case "媒":
		cdTime = time.Duration(setting.CDtime)*time.Hour - time.Since(time.Unix(info.MMk, 0))
	case "买":
		cdTime = time.Duration(setting.CDtime)*time.Hour - time.Since(time.Unix(info.Buy, 0))
	}
	return
}

func (sql *dbData) setCD(uid string, funcType string) error {
	sql.Lock()
	defer sql.Unlock()
	// 创建群表格
	err := sql.db.Create("cdsheet", &cdsheet{})
	if err != nil {
		return err
	}
	info := cdsheet{}
	_ = sql.db.Find("cdsheet", &info, "where User is '"+uid+"'")
	info.User = uid
	switch funcType {
	case "娶", "嫁":
		info.Mar = time.Now().Unix()
	case "牛":
		info.Rob = time.Now().Unix()
	case "离":
		info.Lef = time.Now().Unix()
	case "媒":
		info.MMk = time.Now().Unix()
	case "买":
		info.Buy = time.Now().Unix()
	}
	return sql.db.Insert("cdsheet", &info)
}

func getLine() string {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		return path.Base(file) + "." + strconv.Itoa(line)
	}
	return ""
}
