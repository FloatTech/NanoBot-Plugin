package qqwife

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/FloatTech/AnimeAPI/wallet"
	"github.com/FloatTech/NanoBot-Plugin/utils/ctxext"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/math"
	"github.com/FloatTech/gg"
	"github.com/FloatTech/imgfactory"
	control "github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/img/text"
	nano "github.com/fumiama/NanoBot"
	log "github.com/sirupsen/logrus"
)

var sendtext = [...][]string{
	{ // 表白成功
		"是个勇敢的孩子(*/ω＼*) 今天的运气都降临在你的身边~\n\n",
		"(´･ω･`)对方答应了你 并表示愿意当今天的CP\n\n",
	},
	{ // 表白失败
		"今天的运气有一点背哦~明天再试试叭",
		"_(:з」∠)_下次还有机会 咱抱抱你w",
		"今天失败了惹. 摸摸头~咱明天还有机会",
	},
	{ // ntr成功
		"因为你的个人魅力~~今天他就是你的了w\n\n",
	},
	{ // 离婚失败
		"打是情,骂是爱,不打不亲不相爱。答应我不要分手。",
		"床头打架床尾和，夫妻没有隔夜仇。安啦安啦，不要闹变扭。",
	},
	{ // 离婚成功
		"离婚成功力\n话说你不考虑当个1？",
		"离婚成功力\n天涯何处无芳草，何必单恋一枝花？不如再摘一支（bushi",
	},
}

func init() {
	engine.OnMessageRegex(`^设置CD为(\d+)小时`, nano.OnlyChannel, nano.AdminPermission, getdb).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *nano.Ctx) {
		cdTime, err := strconv.ParseInt(ctx.State["regex_matched"].([]string)[1], 10, 64)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.47 ->ERROR]:请设置纯数字\n", err)
			return
		}
		groupInfo, err := wifeData.getSet(ctx.Message.ChannelID)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.52 ->ERROR]:", err)
			return
		}
		groupInfo.CDtime = cdTime
		err = wifeData.updateSet(groupInfo)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.58 ->ERROR]设置CD时长失败\n", err)
			return
		}
		_, _ = ctx.SendPlainMessage(true, "设置成功")
	})
	engine.OnMessageRegex(`^(允许|禁止)(自由恋爱|牛头人)$`, nano.OnlyChannel, nano.AdminPermission, getdb).SetBlock(true).Handle(func(ctx *nano.Ctx) {
		status := ctx.State["regex_matched"].([]string)[1]
		mode := ctx.State["regex_matched"].([]string)[2]
		groupInfo, err := wifeData.getSet(ctx.Message.ChannelID)
		switch {
		case err != nil:
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.69 ->ERROR]:", err)
			return
		case mode == "自由恋爱":
			if status == "允许" {
				groupInfo.CanMatch = 1
			} else {
				groupInfo.CanMatch = 0
			}
		case mode == "牛头人":
			if status == "允许" {
				groupInfo.CanNtr = 1
			} else {
				groupInfo.CanNtr = 0
			}
		}
		err = wifeData.updateSet(groupInfo)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.86 ->ERROR]:", err)
			return
		}
		_, _ = ctx.SendPlainMessage(true, "设置成功")
	})
	// 单身技能
	engine.OnMessageRegex(`^(娶|嫁)\s*<@!(\d+)>$`, nano.OnlyChannel, getdb).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *nano.Ctx) {
		gid := ctx.Message.ChannelID
		setting, err := wifeData.getSet(gid)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.96 ->ERROR]:", err)
			return
		}
		if setting.CanMatch == 0 {
			_, _ = ctx.SendPlainMessage(true, "该频道已发布了禁止自由恋爱,请认真水群")
			return
		}
		uid := ctx.Message.Author.ID
		choice := ctx.State["regex_matched"].([]string)[1]
		cdTime, err := wifeData.checkCD(gid, uid, choice)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.107 ->ERROR]:", err)
			return
		}
		if cdTime > 0 {
			_, _ = ctx.SendPlainMessage(true, "你的技能CD还有", cdTime)
			return
		}
		fiance := ctx.State["regex_matched"].([]string)[2]
		uInfo, err := wifeData.checkUser(gid, uid)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.117 ->ERROR]:", err)
			return
		}
		if uInfo.Users != "" {
			info := strings.Split(uInfo.Users, " & ")
			switch {
			case info[0] == "" || info[1] == "":
				_, _ = ctx.SendPlainMessage(true, "今天的你是单身贵族噢")
				return
			case info[0] == fiance || info[1] == fiance:
				_, _ = ctx.SendPlainMessage(true, "笨蛋！你们已经在一起了！")
				return
			case info[0] == uid: // 如果如为攻
				_, _ = ctx.SendPlainMessage(true, "笨蛋~你家里还有个吃白饭的w")
				return
			case info[1] == uid: // 如果为受
				_, _ = ctx.SendPlainMessage(true, "该是0就是0,当0有什么不好")
				return
			}
		}
		fInfo, err := wifeData.checkUser(gid, uid)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.139 ->ERROR]:", err)
			return
		}
		if fInfo.Users != "" {
			info := strings.Split(fInfo.Users, " & ")
			switch {
			case info[0] == "" || info[1] == "":
				_, _ = ctx.SendPlainMessage(true, "今天的ta是单身贵族噢")
				return
			case info[0] == uid: // 如果如为攻
				_, _ = ctx.SendPlainMessage(true, "他有别的女人了，你该放下了")
				return
			case info[1] == uid: // 如果为受
				_, _ = ctx.SendPlainMessage(true, "ta被别人娶了,你来晚力")
				return
			}
		}
		// 写入CD
		err = wifeData.setCD(uid, choice)
		if err != nil {
			log.Warnln("[qqwife]你的技能CD记录失败,", err)
		}
		uBook, err := getUserInfoIn(ctx, gid, uid)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.163 ->ERROR]:", err)
			return
		}
		if uid == fiance { // 如果是自己
			switch rand.Intn(3) {
			case 1:
				err := wifeData.register(gid, uBook, userInfo{})
				if err != nil {
					_, _ = ctx.SendPlainMessage(false, "[happyplay.go.171 ->ERROR]:", err)
					return
				}
				_, _ = ctx.SendPlainMessage(true, "今日获得成就：单身贵族")
			default:
				_, _ = ctx.SendPlainMessage(true, "今日获得成就：自恋狂")
			}
			return
		}
		fBook, err := getUserInfoIn(ctx, gid, fiance)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.182 ->ERROR]:", err)
			return
		}
		favor, err := wifeData.favorFor(uid, fiance, 0)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.187 ->ERROR]:", err)
			return
		}
		if favor < 30 {
			favor = 30 // 保底30%概率
		}
		if rand.Intn(101) >= favor {
			_, _ = ctx.SendPlainMessage(true, sendtext[1][rand.Intn(len(sendtext[1]))])
			return
		}
		// 去民政局登记
		var choicetext string
		switch choice {
		case "娶":
			err = wifeData.register(gid, uBook, fBook)
			choicetext = "\n今天你的群老婆是"
		default:
			err = wifeData.register(gid, fBook, uBook)
			choicetext = "\n今天你的群老公是"
		}
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.208 ->ERROR]:", err)
			return
		}
		// 请大家吃席
		_, err = ctx.SendChain(nano.ReplyTo(ctx.Message.ID),
			nano.Text(sendtext[0][rand.Intn(len(sendtext[0]))], "\n",
				choicetext, "[", fBook.Nick, "](", fiance, ")\n",
				"当前你们好感度为", favor), nano.Image(fBook.Avatar))
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "happyplay.go.217 ->ERROR: ", err)
		}
	})
	// NTR技能
	engine.OnMessageRegex(`^牛\s*<@!(\d+)>$`, nano.OnlyChannel, getdb).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *nano.Ctx) {
		gid := ctx.Message.ChannelID
		setting, err := wifeData.getSet(gid)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.225 ->ERROR]:", err)
			return
		}
		if setting.CanNtr == 0 {
			_, _ = ctx.SendPlainMessage(true, "该频道已发布了禁止牛头人,请认真水群")
			return
		}
		uid := ctx.Message.Author.ID
		cdTime, err := wifeData.checkCD(gid, uid, "牛")
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.235 ->ERROR]:", err)
			return
		}
		if cdTime > 0 {
			_, _ = ctx.SendPlainMessage(true, "你的技能CD还有", cdTime)
			return
		}
		fiance := ctx.State["regex_matched"].([]string)[1]
		uInfo, err := wifeData.checkUser(gid, uid)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.245 ->ERROR]:", err)
			return
		}
		if uInfo.Users != "" {
			info := strings.Split(uInfo.Users, " & ")
			switch {
			case info[0] == "" || info[1] == "":
				_, _ = ctx.SendPlainMessage(true, "今天的你是单身贵族噢")
				return
			case info[0] == fiance || info[1] == fiance:
				_, _ = ctx.SendPlainMessage(true, "笨蛋！你们已经在一起了！")
				return
			case info[0] == uid: // 如果如为攻
				_, _ = ctx.SendPlainMessage(true, "笨蛋~你家里还有个吃白饭的w")
				return
			case info[1] == uid: // 如果为受
				_, _ = ctx.SendPlainMessage(true, "该是0就是0,当0有什么不好")
				return
			}
		}
		fInfo, err := wifeData.checkUser(gid, fiance)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.267 ->ERROR]:", err)
			return
		}
		if fInfo.Users == "" {
			_, _ = ctx.SendPlainMessage(true, "今天的ta是单身噢,快去明媒正娶吧!")
			return
		}
		// 写入CD
		err = wifeData.setCD(uid, "牛")
		if err != nil {
			log.Warnln("[qqwife]你的技能CD记录失败,", err)
		}
		if fiance == uid {
			_, _ = ctx.SendPlainMessage(true, "今日获得成就：自我攻略")
			return
		}
		favor, err := wifeData.favorFor(uid, fiance, 0)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.285 ->ERROR]:", err)
			return
		}
		if favor < 30 {
			favor = 30 // 保底10%概率
		}
		if rand.Intn(101) >= favor/3 {
			_, _ = ctx.SendPlainMessage(true, "失败了！可惜")
			return
		}
		// 判断target是老公还是老婆
		choicetext := "老公"
		ntrID := uid
		targetID := fiance
		greenID := "" // 被牛的

		err = wifeData.divorce(gid, fiance)
		if err != nil {
			_, _ = ctx.SendPlainMessage(true, "ta不想和原来的对象分手...\n[error]", err)
			return
		}
		user := strings.Split(fInfo.Users, " & ")
		switch {
		case user[0] == fiance: // 是1
			ntrID = fiance
			targetID = uid
			greenID = user[1]
			choicetext = "老公"
		case user[1] == fiance: // 是0
			greenID = user[0]
			choicetext = "老婆"
		default:
			_, _ = ctx.SendPlainMessage(true, "数据库发生问题力")
			return
		}
		userInfo, err := getUserInfoIn(ctx, gid, ntrID)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.322 ->ERROR]:", err)
			return
		}
		fianceInfo, err := getUserInfoIn(ctx, gid, targetID)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.327 ->ERROR]:", err)
			return
		}
		err = wifeData.register(gid, userInfo, fianceInfo)
		if err != nil {
			_, _ = ctx.SendPlainMessage(true, "[qqwife]复婚登记失败力\n", err)
			return
		}
		favor, err = wifeData.favorFor(uid, fiance, -5)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.337 ->ERROR]:", err)
		}
		_, err = wifeData.favorFor(uid, greenID, 5)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.341 ->ERROR]:", err)
		}
		// 输出结果
		_, err = ctx.SendChain(nano.ReplyTo(ctx.Message.ID),
			nano.Text(sendtext[2][rand.Intn(len(sendtext[2]))], "\n",
				choicetext, "[", fianceInfo.Nick, "](", fianceInfo.ID, ")\n",
				"当前你们好感度为", favor), nano.Image(fianceInfo.Avatar))
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.349 ->ERROR]: ", err)
		}
	})
	// 做媒技能
	engine.OnMessageRegex(`^做媒\s*<@!(\d+)>\s*<@!(\d+)>`, nano.OnlyChannel, nano.AdminPermission, getdb).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *nano.Ctx) {
		gid := ctx.Message.ChannelID
		uid := ctx.Message.Author.ID
		cdTime, err := wifeData.checkCD(gid, uid, "媒")
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.358 ->ERROR]:", err)
			return
		}
		if cdTime > 0 {
			_, _ = ctx.SendPlainMessage(true, "你的技能CD还有", cdTime)
			return
		}
		gayOne := ctx.State["regex_matched"].([]string)[1]
		gaynano := ctx.State["regex_matched"].([]string)[2]
		uInfo, err := wifeData.checkUser(gid, gayOne)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.369 ->ERROR]:", err)
			return
		}
		if uInfo.Users != "" {
			_, _ = ctx.SendChain(nano.ReplyTo(ctx.Message.ID), nano.At(gayOne), nano.Text("已有家妻"))
			return
		}
		fInfo, err := wifeData.checkUser(gid, gaynano)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.378 ->ERROR]:", err)
			return
		}
		if fInfo.Users != "" {
			_, _ = ctx.SendChain(nano.ReplyTo(ctx.Message.ID), nano.At(gaynano), nano.Text("已有所属"))
			return
		}
		// 写入CD
		err = wifeData.setCD(uid, "媒")
		if err != nil {
			log.Warnln("[qqwife]你的技能CD记录失败,", err)
		}
		favor, err := wifeData.favorFor(gayOne, gaynano, 0)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.392 ->ERROR]:", err)
			return
		}
		if favor < 30 {
			favor = 30 // 保底30%概率
		}
		if rand.Intn(101) >= favor {
			_, err = wifeData.favorFor(uid, gayOne, -1)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "[happyplay.go.401 ->ERROR]:", err)
			}
			_, err = wifeData.favorFor(uid, gaynano, -1)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "[happyplay.go.64052 ->ERROR]:", err)
			}
			_, _ = ctx.SendPlainMessage(true, sendtext[1][rand.Intn(len(sendtext[1]))])
			return
		}
		// 去民政局登记
		userInfo, err := getUserInfoIn(ctx, gid, gayOne)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.413 ->ERROR]:", err)
			return
		}
		fianceInfo, err := getUserInfoIn(ctx, gid, gaynano)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.418 ->ERROR]:", err)
			return
		}
		err = wifeData.register(gid, userInfo, fianceInfo)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.423 ->ERROR]:", err)
			return
		}
		_, err = wifeData.favorFor(uid, gayOne, 1)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.428 ->ERROR]:", err)
		}
		_, err = wifeData.favorFor(uid, gaynano, 1)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.432 ->ERROR]:", err)
		}
		_, err = wifeData.favorFor(gayOne, gaynano, 1)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.436 ->ERROR]:", err)
		}
		// 请大家吃席
		_, err = ctx.SendChain(nano.ReplyTo(ctx.Message.ID),
			nano.Text("恭喜你成功撮合了一对CP\n\n"), nano.At(gayOne), nano.Text("今天你的群老婆是[", fianceInfo.Nick, "](", fianceInfo.ID, ")"),
			nano.Image(fianceInfo.Avatar))
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "happyplay.go.443 ->ERROR: ", err)
		}
	})
	engine.OnMessageFullMatchGroup([]string{"闹离婚", "办离婚"}, nano.OnlyChannel, getdb).Limit(ctxext.LimitByUser).SetBlock(true).Handle(func(ctx *nano.Ctx) {
		gid := ctx.Message.ChannelID
		uid := ctx.Message.Author.ID
		cdTime, err := wifeData.checkCD(gid, uid, "离")
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.451 ->ERROR]:", err)
			return
		}
		if cdTime > 0 {
			_, _ = ctx.SendPlainMessage(true, "你的技能CD还有", cdTime)
			return
		}
		uInfo, err := wifeData.checkUser(gid, uid)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.460 ->ERROR]:", err)
			return
		}
		if uInfo.Users == "" {
			_, _ = ctx.SendPlainMessage(true, "你还是单身噢,快去娶群友吧!")
			return
		}
		// 写入CD
		err = wifeData.setCD(uid, "离")
		if err != nil {
			_, _ = ctx.SendPlainMessage(true, "[qqwife]你的技能CD记录失败\n", err)
		}
		user := strings.Split(uInfo.Users, " & ")
		mun := 0
		if user[1] == uid {
			mun = 1
		}
		favor, err := wifeData.favorFor(user[0], user[1], 0)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.479 ->ERROR]:", err)
			return
		}
		if favor < 30 {
			favor = 10
		}
		if rand.Intn(101) > 110-favor {
			_, _ = ctx.SendPlainMessage(true, sendtext[3][rand.Intn(len(sendtext[3]))])
			return
		}
		err = wifeData.divorce(gid, uid)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.491 ->ERROR]:", err)
			return
		}
		/*
			if rand.Intn(100) > 50 {
				_, _ = wifeData.favorFor(user[0], user[1], -rand.Intn(favor/2))
			}
		*/
		_, _ = ctx.SendPlainMessage(true, sendtext[4][mun])
	})

	// 好感度系统
	engine.OnMessageRegex(`^查好感度\s*<@!(\d+)>`, nano.OnlyChannel, getdb).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *nano.Ctx) {
		fiance := ctx.State["regex_matched"].([]string)[1]
		uid := ctx.Message.Author.ID
		favor, err := wifeData.favorFor(uid, fiance, 0)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.503 ->ERROR]:", err)
			return
		}
		// 输出结果
		_, _ = ctx.SendPlainMessage(true, "当前你们好感度为", favor)
	})
	// 礼物系统
	engine.OnMessageRegex(`^买礼物给\s*<@!(\d+)>`, getdb).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *nano.Ctx) {
		gid := ctx.Message.ChannelID
		uid := ctx.Message.Author.ID
		fiance := ctx.State["regex_matched"].([]string)[1]
		if fiance == uid {
			_, _ = ctx.SendPlainMessage(true, "你想给自己买什么礼物呢?")
			return
		}
		// 获取CD
		cdTime, err := wifeData.checkCD(gid, uid, "买")
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.521 ->ERROR]:", err)
			return
		}
		if cdTime > 0 {
			_, _ = ctx.SendPlainMessage(true, "你的技能CD还有", cdTime)
			return
		}
		// 获取好感度
		favor, err := wifeData.favorFor(uid, fiance, 0)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.531 ->ERROR]:", err)
			return
		}
		// 对接小熊饼干
		uidint64, _ := strconv.ParseInt(uid, 10, 64)
		walletinfo := wallet.GetWalletOf(uidint64)
		if walletinfo < 1 {
			_, _ = ctx.SendPlainMessage(true, "你钱包没钱啦!")
			return
		}
		moneyToFavor := rand.Intn(math.Min(walletinfo, 100)) + 1
		// 计算钱对应的好感值
		newFavor := 1
		moodMax := 2
		if favor > 50 {
			newFavor = moneyToFavor % 10 // 礼物厌倦
		} else {
			moodMax = 5
			newFavor += rand.Intn(moneyToFavor)
		}
		// 随机对方心情
		mood := rand.Intn(moodMax)
		if mood == 0 {
			newFavor = -newFavor
		}
		// 记录结果
		err = wallet.InsertWalletOf(uidint64, -moneyToFavor)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.559 ->ERROR]:", err)
			return
		}
		lastfavor, err := wifeData.favorFor(uid, fiance, newFavor)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.564 ->ERROR]:", err)
			return
		}
		// 写入CD
		err = wifeData.setCD(uid, "buy")
		if err != nil {
			_, _ = ctx.SendPlainMessage(true, "[qqwife]你的技能CD记录失败\n", err)
		}
		// 输出结果
		if mood == 0 {
			_, _ = ctx.SendPlainMessage(true, "你花了", moneyToFavor, "ATRI币买了一件女装送给了ta,ta很不喜欢,你们的好感度降低至", lastfavor)
		} else {
			_, _ = ctx.SendPlainMessage(true, "你花了", moneyToFavor, "ATRI币买了一件女装送给了ta,ta很喜欢,你们的好感度升至", lastfavor)
		}
	})
	engine.OnMessageFullMatch("好感度列表", nano.OnlyChannel, getdb).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *nano.Ctx) {
		gid := ctx.Message.ChannelID
		uid := ctx.Message.Author.ID
		fianceeInfo, err := wifeData.getGroupFavorability(uid)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.584 ->ERROR]:", err)
			return
		}
		/***********设置图片的大小和底色***********/
		number := len(fianceeInfo)
		if number > 10 {
			number = 10
		}
		fontSize := 50.0
		canvas := gg.NewContext(1150, int(170+(50+70)*float64(number)))
		canvas.SetRGB(1, 1, 1) // 白色
		canvas.Clear()
		/***********下载字体***********/
		data, err := file.GetLazyData(text.BoldFontFile, control.Md5File, true)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.599 ->ERROR]:", err)
		}
		/***********设置字体颜色为黑色***********/
		canvas.SetRGB(0, 0, 0)
		/***********设置字体大小,并获取字体高度用来定位***********/
		if err = canvas.ParseFontFace(data, fontSize*2); err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.605 ->ERROR]:", err)
			return
		}
		sl, h := canvas.MeasureString("你的好感度排行列表")
		/***********绘制标题***********/
		canvas.DrawString("你的好感度排行列表", (1100-sl)/2, 100) // 放置在中间位置
		canvas.DrawString("————————————————————", 0, 160)
		/***********设置字体大小,并获取字体高度用来定位***********/
		if err = canvas.ParseFontFace(data, fontSize); err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.614 ->ERROR]:", err)
			return
		}
		i := 0
		for _, info := range fianceeInfo {
			if info.Favor == 0 {
				break
			}
			if info.Users == "" {
				continue
			}
			user, err := getUserInfoIn(ctx, gid, info.Users)
			if err != nil {
				log.Warnln("[happyplay.go.627 ->ERROR]:", err.Error())
				continue
			}
			canvas.SetRGB255(0, 0, 0)
			canvas.DrawString(user.Nick+"("+user.ID+")", 10, float64(180+(50+70)*i))
			canvas.DrawString(strconv.Itoa(info.Favor), 1020, float64(180+60+(50+70)*i))
			canvas.DrawRectangle(10, float64(180+60+(50+70)*i)-h/2, 1000, 50)
			canvas.SetRGB255(150, 150, 150)
			canvas.Fill()
			canvas.SetRGB255(0, 0, 0)
			canvas.DrawRectangle(10, float64(180+60+(50+70)*i)-h/2, float64(info.Favor)*10, 50)
			canvas.SetRGB255(231, 27, 100)
			canvas.Fill()
			i++
		}
		data, err = imgfactory.ToBytes(canvas.Image())
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "[happyplay.go.62 ->ERROR]:", err)
			return
		}
		ctx.SendImageBytes(data, true)
	})
}
