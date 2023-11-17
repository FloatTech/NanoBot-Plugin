// Package runcode 基于 https://tool.runoob.com 的在线运行代码
package runcode

import (
	"strings"

	nano "github.com/fumiama/NanoBot"

	"github.com/FloatTech/AnimeAPI/runoob"
	ctrl "github.com/FloatTech/zbpctrl"

	"github.com/FloatTech/NanoBot-Plugin/utils/ctxext"
)

var ro = runoob.NewRunOOB("066417defb80d038228de76ec581a50a")

func init() {
	nano.Register("runcode", &ctrl.Options[*nano.Ctx]{
		DisableOnDefault: false,
		Help: "在线代码运行: \n" +
			">runcode [language] [code block]\n" +
			"模板查看: \n" +
			">runcode [language] help\n" +
			"支持语种: \n" +
			"Go || Python || C/C++ || C# || Java || Lua \n" +
			"JavaScript || TypeScript || PHP || Shell \n" +
			"Kotlin  || Rust || Erlang || Ruby || Swift \n" +
			"R || VB || Py2 || Perl || Pascal || Scala",
	}).ApplySingle(ctxext.DefaultSingle).OnMessageRegex(`^\s*[(&gt;)>]runcode(raw)?\s(.+?)\s([\s\S]+)$`).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *nano.Ctx) {
			israw := ctx.State["regex_matched"].([]string)[1] != ""
			language := ctx.State["regex_matched"].([]string)[2]
			language = strings.ToLower(language)
			if _, exist := runoob.LangTable[language]; !exist {
				// 不支持语言
				msg := "> " + ctx.Message.Author.Username + "\n语言" + language + "不是受支持的编程语种呢~"
				if nano.OnlyQQ(ctx) {
					_, _ = ctx.SendPlainMessage(false, msg)
				} else {
					_, _ = ctx.SendPlainMessage(false, nano.MessageEscape(msg))
				}
			} else {
				// 执行运行
				block := ctx.State["regex_matched"].([]string)[3]
				switch block {
				case "help":
					msg := "> " + ctx.Message.Author.Username + "  " + language + "-template:\n>runcode " + language + "\n" + runoob.Templates[language]
					var err error
					if nano.OnlyQQ(ctx) {
						_, err = ctx.SendPlainMessage(false, msg)
					} else {
						_, err = ctx.SendPlainMessage(false, nano.MessageEscape(msg))
					}
					if err != nil {
						_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
					}
				default:
					output, err := ro.Run(block, language, "")
					if err != nil {
						output = "ERROR:\n" + nano.MessageEscape(err.Error())
					}
					output = cutTooLong(strings.Trim(output, "\n"))
					if israw {
						_, err = ctx.SendPlainMessage(false, output)
					} else {
						head := "> " + ctx.Message.Author.Username + "\n"
						if nano.OnlyQQ(ctx) {
							_, err = ctx.SendPlainMessage(false, head+output)
						} else {
							_, err = ctx.SendPlainMessage(false, nano.MessageEscape(head+output))
						}
					}
					if err != nil {
						_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
					}
				}
			}
		})
}

// 截断过长文本
func cutTooLong(text string) string {
	temp := []rune(text)
	count := 0
	for i := range temp {
		switch {
		case temp[i] == 13 && i < len(temp)-1 && temp[i+1] == 10:
			// 匹配 \r\n 跳过，等 \n 自己加
		case temp[i] == 10:
			count++
		case temp[i] == 13:
			count++
		}
		if count > 30 || i > 1000 {
			temp = append(temp[:i-1], []rune("\n............\n............")...)
			break
		}
	}
	return string(temp)
}
