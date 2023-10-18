// Package autowithdraw 触发者撤回时也自动撤回
package autowithdraw

import (
	"github.com/FloatTech/floatbox/process"
	ctrl "github.com/FloatTech/zbpctrl"
	nano "github.com/fumiama/NanoBot"
)

func init() {
	en := nano.Register("autowithdraw", &ctrl.Options[*nano.Ctx]{
		DisableOnDefault: false,
		Brief:            "触发者撤回时也自动撤回",
		Help:             "- 撤回一条消息\n",
	})
	en.OnMessageDelete(nano.OnlyPrivate).SetBlock(false).Handle(func(ctx *nano.Ctx) {
		delmsg := ctx.Value.(*nano.MessageDelete)
		for _, msg := range nano.GetTriggeredMessages(delmsg.Message.ID) {
			process.SleepAbout1sTo2s()
			_ = ctx.DeleteMessageInChannel(delmsg.Message.ChannelID, msg, false)
		}
	})
}
