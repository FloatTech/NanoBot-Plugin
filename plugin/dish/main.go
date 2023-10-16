// Package dish 程序员做饭指南, 数据来源Anduin2017/HowToCook
package dish

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	sql "github.com/FloatTech/sqlite"
	ctrl "github.com/FloatTech/zbpctrl"
	nano "github.com/fumiama/NanoBot"

	"github.com/FloatTech/NanoBot-Plugin/utils/ctxext"
)

type dish struct {
	ID        uint32 `db:"id"`
	Name      string `db:"name"`
	Materials string `db:"materials"`
	Steps     string `db:"steps"`
}

var (
	db          = &sql.Sqlite{}
	initialized = false
)

func init() {
	en := nano.Register("dish", &ctrl.Options[*nano.Ctx]{
		DisableOnDefault: false,
		Brief:            "程序员做饭指南",
		Help:             "-怎么做[xxx]|烹饪[xxx]|随机菜谱|随便做点菜",
		PublicDataFolder: "Dish",
	})

	db.DBPath = en.DataFolder() + "dishes.db"

	if _, err := en.GetLazyData("dishes.db", true); err != nil {
		logrus.Warnln("[dish]获取菜谱数据库文件失败")
	} else if err = db.Open(time.Hour); err != nil {
		logrus.Warnln("[dish]连接菜谱数据库失败")
	} else if err = db.Create("dish", &dish{}); err != nil {
		logrus.Warnln("[dish]同步菜谱数据表失败")
	} else if count, err := db.Count("dish"); err != nil {
		logrus.Warnln("[dish]统计菜谱数据失败")
	} else {
		logrus.Infoln("[dish]加载", count, "条菜谱")
		initialized = true
	}

	if !initialized {
		logrus.Warnln("[dish]插件未能成功初始化")
	}

	en.OnMessagePrefixGroup([]string{"怎么做", "烹饪"}).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *nano.Ctx) {
		if !initialized {
			_, err := ctx.SendPlainMessage(false, "客官，本店暂未开业")
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			}
			return
		}

		name := ctx.Message.Author.Username
		dishName := ctx.State["args"].(string)

		if dishName == "" {
			return
		}

		if strings.Contains(dishName, "'") ||
			strings.Contains(dishName, "\"") ||
			strings.Contains(dishName, "\\") ||
			strings.Contains(dishName, ";") {
			return
		}

		var d dish
		if err := db.Find("dish", &d, fmt.Sprintf("WHERE name like %%%s%%", dishName)); err != nil {
			return
		}

		_, err := ctx.SendPlainMessage(false, fmt.Sprintf(
			"已为客官%s找到%s的做法辣！\n"+
				"原材料：%s\n"+
				"步骤：\n"+
				"%s",
			name, dishName, d.Materials, d.Steps),
		)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
		}
	})

	en.OnMessagePrefixGroup([]string{"随机菜谱", "随便做点菜"}).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *nano.Ctx) {
		if !initialized {
			_, err := ctx.SendPlainMessage(false, "客官，本店暂未开业")
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			}
			return
		}

		name := ctx.Message.Author.Username
		var d dish
		if err := db.Pick("dish", &d); err != nil {
			_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			return
		}

		_, err := ctx.SendPlainMessage(false, fmt.Sprintf(
			"已为客官%s送上%s的做法\n"+
				"原材料：%s\n"+
				"步骤：\n"+
				"%s",
			name, d.Name, d.Materials, d.Steps),
		)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
		}
	})
}
