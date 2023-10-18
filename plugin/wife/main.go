// Package wife 抽老婆
package wife

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/FloatTech/NanoBot-Plugin/utils/ctxext"
	fcext "github.com/FloatTech/floatbox/ctxext"
	ctrl "github.com/FloatTech/zbpctrl"
	nano "github.com/fumiama/NanoBot"
	"github.com/sirupsen/logrus"
)

func init() {
	engine := nano.Register("wife", &ctrl.Options[*nano.Ctx]{
		DisableOnDefault: false,
		Help:             "- 抽老婆",
		Brief:            "从老婆库抽每日老婆",
		PublicDataFolder: "Wife",
	}).ApplySingle(ctxext.DefaultSingle)
	_ = os.MkdirAll(engine.DataFolder()+"wives", 0755)
	cards := []string{}
	engine.OnMessageFullMatch("抽老婆", fcext.DoOnceOnSuccess(
		func(ctx *nano.Ctx) bool {
			data, err := engine.GetLazyData("wife.json", true)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return false
			}
			err = json.Unmarshal(data, &cards)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return false
			}
			logrus.Infof("[wife]加载%d个老婆", len(cards))
			return true
		},
	)).SetBlock(true).
		Handle(func(ctx *nano.Ctx) {
			uid := ctx.Message.Author.ID
			if uid == "" {
				_, _ = ctx.SendPlainMessage(false, "ERROR: 未获取到用户")
				return
			}
			uidint, _ := strconv.ParseInt(uid, 10, 64)
			card := cards[fcext.RandSenderPerDayN(uidint, len(cards))]
			data, err := engine.GetLazyData("wives/"+card, true)
			card, _, _ = strings.Cut(card, ".")
			if err != nil {
				_, err = ctx.SendChain(nano.At(uid), nano.Text("今天的二次元老婆是~【", card, "】哒\n【图片下载失败: ", err, "】"))
				if err != nil {
					_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				}
				return
			}
			_, err = ctx.SendChain(nano.At(uid), nano.Text("今天的二次元老婆是~【", card, "】哒"), nano.ImageBytes(data))
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			}
		})
}
