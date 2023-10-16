// Package baseamasiro base天城文 与 tea 加解密
package baseamasiro

import (
	nano "github.com/fumiama/NanoBot"

	"github.com/fumiama/unibase2n"

	"github.com/FloatTech/floatbox/crypto"
	ctrl "github.com/FloatTech/zbpctrl"
)

func init() {
	en := nano.Register("baseamasiro", &ctrl.Options[*nano.Ctx]{
		DisableOnDefault: false,
		Help: "base天城文加解密\n" +
			"- 天城文加密xxx\n- 天城文解密xxx\n- 天城文用yyy加密xxx\n- 天城文用yyy解密xxx",
	})
	en.OnMessageRegex(`^天城文加密\s*(.+)$`).SetBlock(true).
		Handle(func(ctx *nano.Ctx) {
			str := ctx.State["regex_matched"].([]string)[1]
			es := unibase2n.BaseDevanagari.EncodeString(str)
			if es != "" {
				_, _ = ctx.SendPlainMessage(false, es)
			} else {
				_, _ = ctx.SendPlainMessage(false, "加密失败!")
			}
		})
	en.OnMessageRegex(`^天城文解密\s*([ऀ-ॿ]+[০-৫]?)$`).SetBlock(true).
		Handle(func(ctx *nano.Ctx) {
			str := ctx.State["regex_matched"].([]string)[1]
			es := unibase2n.BaseDevanagari.DecodeString(str)
			if es != "" {
				_, _ = ctx.SendPlainMessage(false, es)
			} else {
				_, _ = ctx.SendPlainMessage(false, "解密失败!")
			}
		})
	en.OnMessageRegex(`^天城文用(.+)加密\s*(.+)$`).SetBlock(true).
		Handle(func(ctx *nano.Ctx) {
			key, str := ctx.State["regex_matched"].([]string)[1], ctx.State["regex_matched"].([]string)[2]
			t := crypto.GetTEA(key)
			es, err := unibase2n.UTF16BE2UTF8(unibase2n.BaseDevanagari.Encode(t.Encrypt(nano.StringToBytes(str))))
			if err == nil {
				_, _ = ctx.SendPlainMessage(false, nano.BytesToString(es))
			} else {
				_, _ = ctx.SendPlainMessage(false, "加密失败!")
			}
		})
	en.OnMessageRegex(`^天城文用(.+)解密\s*([ऀ-ॿ]+[০-৫]?)$`).SetBlock(true).
		Handle(func(ctx *nano.Ctx) {
			key, str := ctx.State["regex_matched"].([]string)[1], ctx.State["regex_matched"].([]string)[2]
			t := crypto.GetTEA(key)
			es, err := unibase2n.UTF82UTF16BE(nano.StringToBytes(str))
			if err == nil {
				_, _ = ctx.SendPlainMessage(false, nano.BytesToString(t.Decrypt(unibase2n.BaseDevanagari.Decode(es))))
			} else {
				_, _ = ctx.SendPlainMessage(false, "解密失败!")
			}
		})
}
