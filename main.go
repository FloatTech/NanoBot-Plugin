// Package main NanoBot-Plugin main file
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/NanoBot-Plugin/kanban" // 打印 banner

	_ "github.com/FloatTech/NanoBot-Plugin/plugin/autowithdraw"
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
	_ "github.com/FloatTech/NanoBot-Plugin/plugin/qqwife"
	_ "github.com/FloatTech/NanoBot-Plugin/plugin/runcode"
	_ "github.com/FloatTech/NanoBot-Plugin/plugin/score"
	_ "github.com/FloatTech/NanoBot-Plugin/plugin/status"
	_ "github.com/FloatTech/NanoBot-Plugin/plugin/tarot"
	_ "github.com/FloatTech/NanoBot-Plugin/plugin/wife"
	_ "github.com/FloatTech/NanoBot-Plugin/plugin/wordle"

	// -----------------------以下为内置依赖，勿动------------------------ //
	nano "github.com/fumiama/NanoBot"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/FloatTech/floatbox/process"

	"github.com/FloatTech/NanoBot-Plugin/kanban/banner"
	// -----------------------以上为内置依赖，勿动------------------------ //
)

func main() {
	if !strings.Contains(runtime.Version(), "go1.2") { // go1.20之前版本需要全局 seed，其他插件无需再 seed
		rand.Seed(time.Now().UnixNano()) //nolint: staticcheck
	}

	token := flag.String("t", "", "qq api token")
	appid := flag.String("a", "", "qq appid")
	secret := flag.String("s", "", "qq secret")
	debug := flag.Bool("D", false, "enable debug-level log output")
	timeout := flag.Int("T", 60, "api timeout (s)")
	help := flag.Bool("h", false, "print this help")
	loadconfig := flag.String("c", "", "load from config")
	sandbox := flag.Bool("sandbox", false, "run in sandbox api")
	onlypublic := flag.Bool("public", false, "only listen to public intent")
	shardindex := flag.Uint("shardindex", 0, "shard index")
	shardcount := flag.Uint("shardcount", 0, "shard count")
	savecfg := flag.String("save", "", "save bot config to filename (eg. config.yaml)")
	flag.Parse()
	if *help {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	intent := uint32(nano.IntentPrivate)
	if *onlypublic {
		intent = nano.IntentPublic
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

	bot := []*nano.Bot{}
	if *loadconfig == "" {
		bot = append(bot, &nano.Bot{
			AppID:      *appid,
			Token:      *token,
			Secret:     *secret,
			SuperUsers: sus,
			Timeout:    time.Duration(*timeout) * time.Second,
			Intents:    intent,
			ShardIndex: uint8(*shardindex),
			ShardCount: uint8(*shardcount),
		})
	} else {
		f, err := os.Open(*loadconfig)
		if err != nil {
			logrus.Fatal(err)
		}
		dec := yaml.NewDecoder(f)
		dec.KnownFields(true)
		err = dec.Decode(&bot)
		_ = f.Close()
		if err != nil {
			logrus.Fatal(err)
		}
	}
	if *savecfg != "" {
		f, err := os.Create(*savecfg)
		if err != nil {
			logrus.Fatal(err)
		}
		defer f.Close()
		err = yaml.NewEncoder(f).Encode(bot)
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Infoln("已将当前配置保存到", *savecfg)
		return
	}

	nano.OnMessageCommandGroup([]string{"help", "帮助", "menu", "菜单"}, nano.OnlyToMe).SetBlock(true).
		Handle(func(ctx *nano.Ctx) {
			_, _ = ctx.SendChain(nano.Text(banner.Banner))
		})
	nano.OnMessageFullMatch("查看nbp公告", nano.OnlyToMe, nano.AdminPermission).SetBlock(true).
		Handle(func(ctx *nano.Ctx) {
			_, _ = ctx.SendChain(nano.Text(strings.ReplaceAll(kanban.Kanban(), "\t", "")))
		})
	_ = nano.Run(process.GlobalInitMutex.Unlock, bot...)
}
