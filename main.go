// Package main NanoBot-Plugin main file
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	_ "github.com/FloatTech/NanoBot-Plugin/plugin/b14"
	_ "github.com/FloatTech/NanoBot-Plugin/plugin/base64gua"
	_ "github.com/FloatTech/NanoBot-Plugin/plugin/baseamasiro"
	_ "github.com/FloatTech/NanoBot-Plugin/plugin/chrev"
	_ "github.com/FloatTech/NanoBot-Plugin/plugin/dish"
	_ "github.com/FloatTech/NanoBot-Plugin/plugin/emojimix"
	_ "github.com/FloatTech/NanoBot-Plugin/plugin/fortune"
	_ "github.com/FloatTech/NanoBot-Plugin/plugin/genshin"
	_ "github.com/FloatTech/NanoBot-Plugin/plugin/hyaku"
	_ "github.com/FloatTech/NanoBot-Plugin/plugin/manager"
	_ "github.com/FloatTech/NanoBot-Plugin/plugin/runcode"
	_ "github.com/FloatTech/NanoBot-Plugin/plugin/score"
	_ "github.com/FloatTech/NanoBot-Plugin/plugin/status"
	_ "github.com/FloatTech/NanoBot-Plugin/plugin/tarot"

	// -----------------------以下为内置依赖，勿动------------------------ //
	nano "github.com/fumiama/NanoBot"
	"github.com/sirupsen/logrus"

	"github.com/FloatTech/NanoBot-Plugin/kanban"
	// -----------------------以上为内置依赖，勿动------------------------ //
)

func main() {
	// 全局 seed，其他插件无需再 seed
	rand.Seed(time.Now().UnixNano()) //nolint: staticcheck

	token := flag.String("t", "", "qq api token")
	appid := flag.String("a", "", "qq appid")
	secret := flag.String("s", "", "qq secret")
	debug := flag.Bool("d", false, "enable debug-level log output")
	timeout := flag.Int("T", 60, "api timeout (s)")
	help := flag.Bool("h", false, "print this help")
	sandbox := flag.Bool("b", false, "run in sandbox api")
	flag.Parse()
	if *help {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	sus := make([]string, 0, 16)
	for _, s := range flag.Args() {
		_, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			continue
		}
		sus = append(sus, s)
	}

	if *sandbox {
		nano.OpenAPI = nano.SandboxAPI
	}

	nano.OnMessageCommandGroup([]string{"help", "帮助", "menu", "菜单"}, nano.OnlyToMe).SetBlock(true).
		Handle(func(ctx *nano.Ctx) {
			_, _ = ctx.SendPlainMessage(false, kanban.Banner)
		})
	_ = nano.Run(&nano.Bot{
		AppID:      *appid,
		Token:      *token,
		Secret:     *secret,
		Intents:    nano.IntentPublic,
		Timeout:    time.Duration(*timeout) * time.Second,
		SuperUsers: sus,
	})
}
