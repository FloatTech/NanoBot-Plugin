// Package base64gua base64卦 与 tea 加解密
package base64gua

import (
	nano "github.com/fumiama/NanoBot"

	"github.com/fumiama/unibase2n"

	"github.com/FloatTech/floatbox/crypto"
	ctrl "github.com/FloatTech/zbpctrl"
)

func init() {
	en := nano.Register("base64gua", &ctrl.Options[*nano.Ctx]{
		DisableOnDefault: false,
		Help: "base64gua加解密\n" +
			"- 六十四卦加密xxx\n- 六十四卦解密xxx\n- 六十四卦用yyy加密xxx\n- 六十四卦用yyy解密xxx",
	})
	en.OnMessageRegex(`^六十四卦加密\s*(.+)$`).SetBlock(true).
		Handle(func(ctx *nano.Ctx) {
			str := ctx.State["regex_matched"].([]string)[1]
			es := unibase2n.Base64Gua.EncodeString(str)
			if es != "" {
				_, _ = ctx.SendPlainMessage(false, es)
			} else {
				_, _ = ctx.SendPlainMessage(false, "加密失败!")
			}
		})
	en.OnMessageRegex(`^六十四卦解密\s*([䷀-䷿]+[☰☱]?)$`).SetBlock(true).
		Handle(func(ctx *nano.Ctx) {
			str := ctx.State["regex_matched"].([]string)[1]
			es := unibase2n.Base64Gua.DecodeString(str)
			if es != "" {
				_, _ = ctx.SendPlainMessage(false, es)
			} else {
				_, _ = ctx.SendPlainMessage(false, "解密失败!")
			}
		})
	en.OnMessageRegex(`^六十四卦用(.+)加密\s*(.+)$`).SetBlock(true).
		Handle(func(ctx *nano.Ctx) {
			key, str := ctx.State["regex_matched"].([]string)[1], ctx.State["regex_matched"].([]string)[2]
			t := crypto.GetTEA(key)
			es, err := unibase2n.UTF16BE2UTF8(unibase2n.Base64Gua.Encode(t.Encrypt(nano.StringToBytes(str))))
			if err == nil {
				_, _ = ctx.SendPlainMessage(false, nano.BytesToString(es))
			} else {
				_, _ = ctx.SendPlainMessage(false, "加密失败!")
			}
		})
	en.OnMessageRegex(`^六十四卦用(.+)解密\s*([䷀-䷿]+[☰☱]?)$`).SetBlock(true).
		Handle(func(ctx *nano.Ctx) {
			key, str := ctx.State["regex_matched"].([]string)[1], ctx.State["regex_matched"].([]string)[2]
			t := crypto.GetTEA(key)
			es, err := unibase2n.UTF82UTF16BE(nano.StringToBytes(str))
			if err == nil {
				_, _ = ctx.SendPlainMessage(false, nano.BytesToString(t.Decrypt(unibase2n.Base64Gua.Decode(es))))
			} else {
				_, _ = ctx.SendPlainMessage(false, "解密失败!")
			}
		})
}
