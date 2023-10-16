// Package emojimix 合成emoji
package emojimix

import (
	"fmt"
	"net/http"

	nano "github.com/fumiama/NanoBot"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/sirupsen/logrus"

	"github.com/FloatTech/NanoBot-Plugin/utils/ctxext"
)

const bed = "https://www.gstatic.com/android/keyboard/emojikitchen/%d/u%x/u%x_u%x.png"

func init() {
	nano.Register("emojimix", &ctrl.Options[*nano.Ctx]{
		DisableOnDefault: false,
		Help: "合成emoji\n" +
			"- [emoji][emoji]",
	}).OnMessage(match).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *nano.Ctx) {
			r := ctx.State["emojimix"].([]rune)
			logrus.Debugln("[emojimix] match:", r)
			r1, r2 := r[0], r[1]
			u1 := fmt.Sprintf(bed, emojis[r1], r1, r1, r2)
			u2 := fmt.Sprintf(bed, emojis[r2], r2, r2, r1)
			resp1, err := http.Head(u1)
			if err == nil {
				resp1.Body.Close()
				if resp1.StatusCode == http.StatusOK {
					_, _ = ctx.SendImage(u1, false)
					return
				}
			}
			resp2, err := http.Head(u2)
			if err == nil {
				resp2.Body.Close()
				if resp2.StatusCode == http.StatusOK {
					_, _ = ctx.SendImage(u2, false)
					return
				}
			}
		})
}

func match(ctx *nano.Ctx) bool {
	r := []rune(ctx.Message.Content)
	if len(r) == 2 {
		if _, ok := emojis[r[0]]; !ok {
			return false
		}
		if _, ok := emojis[r[1]]; !ok {
			return false
		}
		ctx.State["emojimix"] = r
		return true
	}
	return false
}
