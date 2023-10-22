// Package qqwife 娶群友
package qqwife

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/FloatTech/NanoBot-Plugin/utils/ctxext"
	"github.com/FloatTech/imgfactory"
	ctrl "github.com/FloatTech/zbpctrl"
	nano "github.com/fumiama/NanoBot"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/gg"
	"github.com/FloatTech/zbputils/img/text"
)

var (
	engine = nano.Register("qqwife", &ctrl.Options[*nano.Ctx]{
		DisableOnDefault: false,
		Brief:            "娶群友",
		Help: "- 娶群友\n- 群老婆列表\n" +
			"- [允许|禁止]自由恋爱\n- [允许|禁止]牛头人\n" +
			"- 设置CD为xx小时    →(默认12小时)\n" +
			"- 查好感度@对方QQ\n" +
			"- 好感度列表\n" +
			"--------------------------------\n以下指令存在CD,频道共用,不跨天刷新,前两个受指令开关\n--------------------------------\n" +
			"- (娶|嫁)@对方QQ\n    (好感度越高成功率越高,保底30%概率)\n" +
			"- 牛@对方QQ\n    (好感度越高成功率越高,保底10%概率)\n" +
			"- 闹离婚\n    (好感度越高成功率越低)\n" +
			"- 买礼物给@对方QQ\n    (使用bot钱包插件的金额获取好感度)\n" +
			"- 做媒 @攻方QQ @受方QQ\n    (攻受双方好感度越高成功率越高,保底30%概率)\n" +
			"--------------------------------\n好感度规则\n--------------------------------\n" +
			"\"娶群友\"指令好感度随机增加1~5。\n\"A牛B的C\"会导致C恨A, 好感度-5;\nB为了报复A, 好感度+5(什么柜子play)\nA为BC做媒,成功B、C对A好感度+1反之-1\n做媒成功BC好感度+1" +
			"\nTips: 群老婆列表每天4点刷新",
		PrivateDataFolder: "qqwife",
	}).ApplySingle(nano.NewSingle(
		nano.WithKeyFn(func(ctx *nano.Ctx) int64 {
			gid, _ := strconv.ParseUint(ctx.Message.ChannelID, 10, 64)
			return int64(gid)
		}),
		nano.WithPostFn[int64](func(ctx *nano.Ctx) {
			_, _ = ctx.SendPlainMessage(true, "别着急，民政局门口排长队了！")
		}),
	))
)

func init() {
	engine.OnMessageFullMatch("娶群友", nano.OnlyChannel, getdb).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *nano.Ctx) {
		gid := ctx.Message.ChannelID
		uid := ctx.Message.Author.ID

		info, err := wifeData.checkUser(gid, uid)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[main.go.57 ->ERROR]:", err)
			return
		}
		uInfo, err := getUserInfoIn(ctx, gid, uid)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[main.go.62 ->ERROR]:", err)
			return
		}
		if info.Users == "" {
			menbers, err := ctx.GetGuildMembersIn(gid, "0", 1000)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "[main.go.68 ->ERROR]:", err)
				return
			}
			list := make(map[int]userInfo, 1000)
			for i, member := range menbers {
				nick := member.Nick
				if nick == "" {
					nick = member.User.Username
				}
				list[i] = userInfo{
					ID:     member.User.ID,
					Nick:   nick,
					Avatar: member.User.Avatar,
				}
			}
			target := list[rand.Intn(len(list))]
			if target.ID == uid {
				err = wifeData.register(gid, uInfo, userInfo{})
				if err != nil {
					_, _ = ctx.SendPlainMessage(false, "[main.go.87 ->ERROR]:", err)
					return
				}
				_, err = ctx.SendChain(nano.At(uid), nano.Text("今日获得成就：单身贵族"))
				if err != nil {
					_, _ = ctx.SendPlainMessage(false, "main.go.92 ->ERROR: ", err)
				}
			}
			info, err = wifeData.checkUser(gid, target.ID)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "[main.go.97 ->ERROR]:", err)
				return
			}
			if info.Users != "" {
				_, _ = ctx.SendPlainMessage(true, "呜...没娶到，你可以再尝试一次")
				return
			}
			err = wifeData.register(gid, uInfo, target)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "[main.go.106 ->ERROR]:", err)
				return
			}
			favor, err := wifeData.favorFor(uid, target.ID, rand.Intn(5))
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "[main.go.111 ->ERROR]:", err)
				return
			}
			_, err = ctx.SendChain(nano.At(uid), nano.Text("\n今天你的群老婆是\n[", target.Nick, "](", target.ID, ")哒\n当前你们好感度为", favor), nano.Image(target.Avatar))
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "main.go.116 ->ERROR: ", err)
			}
			return
		}
		users := strings.Split(info.Users, " & ")
		switch {
		case (users[0] == uid && users[1] == "") || (users[1] == uid && users[0] == ""): // 如果是单身贵族
			_, _ = ctx.SendPlainMessage(true, "今天你是单身贵族噢")
			return
		case users[0] == uid: // 娶过别人
			favor, err := wifeData.favorFor(uid, users[1], 0)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "[main.go.128 ->ERROR]:", err)
				return
			}
			_, err = ctx.SendChain(nano.At(uid),
				nano.Text("\n今天你在", info.Updatetime, "娶了群友\n[", info.Mname, "](", users[1], ")\n",
					"当前你们好感度为", favor), nano.Image(info.Mpic))
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "main.go.135 ->ERROR: ", err)
			}
			return
		case users[1] == uid: // 嫁给别人
			favor, err := wifeData.favorFor(users[0], uid, 0)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "[main.go.141 ->ERROR]:", err)
				return
			}
			_, err = ctx.SendChain(nano.At(uid),
				nano.Text("\n今天你在", info.Updatetime, "被群友\n[", info.Sname, "](", users[0], ")娶了\n",
					"当前你们好感度为", favor), nano.Image(info.Spic))
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "main.go.148 ->ERROR: ", err)
			}
			return
		}
	})
	engine.OnMessageFullMatch("群老婆列表", nano.OnlyChannel, getdb).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *nano.Ctx) {
		gid := ctx.Message.ChannelID
		list, err := wifeData.getlist(gid)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[main.go.157 ->ERROR]:", err)
			return
		}
		number := len(list)
		if number <= 0 {
			_, _ = ctx.SendPlainMessage(false, "今天没有人结婚哦: ")
			return
		}
		/***********设置图片的大小和底色***********/
		fontSize := 50.0
		if number < 10 {
			number = 10
		}
		canvas := gg.NewContext(1500, int(250+fontSize*float64(number)))
		canvas.SetRGB(1, 1, 1) // 白色
		canvas.Clear()
		/***********下载字体，可以注销掉***********/
		data, err := file.GetLazyData(text.BoldFontFile, nano.Md5File, true)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "main.go.176 ->ERROR: ", err)
		}
		/***********设置字体颜色为黑色***********/
		canvas.SetRGB(0, 0, 0)
		/***********设置字体大小,并获取字体高度用来定位***********/
		if err = canvas.ParseFontFace(data, fontSize*2); err != nil {
			_, _ = ctx.SendPlainMessage(false, "main.go.182 ->ERROR: ", err)
			return
		}
		sl, h := canvas.MeasureString("群老婆列表")
		/***********绘制标题***********/
		canvas.DrawString("群老婆列表", (1500-sl)/2, 160-h) // 放置在中间位置
		canvas.DrawString("————————————————————", 0, 250-h)
		/***********设置字体大小,并获取字体高度用来定位***********/
		if err = canvas.ParseFontFace(data, fontSize); err != nil {
			_, _ = ctx.SendPlainMessage(false, "main.go.191 ->ERROR: ", err)
			return
		}
		_, h = canvas.MeasureString("焯")
		for i, info := range list {
			canvas.DrawString(slicename(info[0], canvas), 0, float64(260+50*i)-h)
			canvas.DrawString("("+info[1]+")", 350, float64(260+50*i)-h)
			canvas.DrawString("←→", 700, float64(260+50*i)-h)
			canvas.DrawString(slicename(info[2], canvas), 800, float64(260+50*i)-h)
			canvas.DrawString("("+info[3]+")", 1150, float64(260+50*i)-h)
		}
		data, err = imgfactory.ToBytes(canvas.Image())
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "main.go.204 ->ERROR: ", err)
			return
		}
		_, _ = ctx.SendImageBytes(data, false)
	})
}
