// Package tarot 塔罗牌
package tarot

import (
	"encoding/json"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/FloatTech/floatbox/binary"
	fcext "github.com/FloatTech/floatbox/ctxext"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/process"
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"

	"github.com/FloatTech/NanoBot-Plugin/utils/ctxext"
	"github.com/FloatTech/zbputils/img/text"
	nano "github.com/fumiama/NanoBot"
	"github.com/sirupsen/logrus"
)

const bed = "https://gitcode.net/shudorcl/zbp-tarot/-/raw/master/"

type cardInfo struct {
	Description        string `json:"description"`
	ReverseDescription string `json:"reverseDescription"`
	ImgURL             string `json:"imgUrl"`
}
type card struct {
	Name     string `json:"name"`
	cardInfo `json:"info"`
}

type formation struct {
	CardsNum  int        `json:"cards_num"`
	IsCut     bool       `json:"is_cut"`
	Represent [][]string `json:"represent"`
}
type cardSet = map[string]card

var (
	cardMap         = make(cardSet, 80)
	infoMap         = make(map[string]cardInfo, 80)
	formationMap    = make(map[string]formation, 10)
	majorArcanaName = make([]string, 0, 80)
	formationName   = make([]string, 0, 10)
)

func init() {
	engine := nano.Register("tarot", &ctrl.Options[*nano.Ctx]{
		DisableOnDefault: false,
		Brief:            "塔罗牌",
		Help: "- 抽[塔罗牌|大阿卡纳|小阿卡纳]\n" +
			"- 解塔罗牌[牌名]",
		PublicDataFolder: "Tarot",
	}).ApplySingle(ctxext.DefaultSingle)

	cache := engine.DataFolder() + "cache"
	_ = os.RemoveAll(cache)
	err := os.MkdirAll(cache, 0755)
	if err != nil {
		panic(err)
	}

	getTarot := fcext.DoOnceOnSuccess(func(ctx *nano.Ctx) bool {
		data, err := engine.GetLazyData("tarots.json", true)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			return false
		}
		err = json.Unmarshal(data, &cardMap)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			return false
		}
		for _, card := range cardMap {
			infoMap[card.Name] = card.cardInfo
		}
		for i := 0; i < 22; i++ {
			majorArcanaName = append(majorArcanaName, cardMap[strconv.Itoa(i)].Name)
		}
		logrus.Infof("[tarot]读取%d张塔罗牌", len(cardMap))
		formation, err := engine.GetLazyData("formation.json", true)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			return false
		}
		err = json.Unmarshal(formation, &formationMap)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			return false
		}
		for k := range formationMap {
			formationName = append(formationName, k)
		}
		logrus.Infof("[tarot]读取%d组塔罗牌阵", len(formationMap))
		return true
	})
	engine.OnMessageRegex(`^抽((塔罗牌|大阿(尔)?卡纳)|小阿(尔)?卡纳)$`, getTarot).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *nano.Ctx) {
		cardType := ctx.State["regex_matched"].([]string)[1]

		reasons := [...]string{"您抽到的是~\n", "锵锵锵，塔罗牌的预言是~\n", "诶，让我看看您抽到了~\n"}
		position := [...]string{"『正位』", "『逆位』"}
		reverse := [...]string{"", "Reverse/"}
		start := 0
		length := 22
		if strings.Contains(cardType, "小") {
			start = 22
			length = 55
		}

		i := rand.Intn(length) + start
		p := rand.Intn(2)
		card := cardMap[strconv.Itoa(i)]
		name := card.Name
		description := card.Description
		if p == 1 {
			description = card.ReverseDescription
		}
		imgurl := bed + reverse[p] + card.ImgURL
		imgname := ""
		if p == 1 {
			imgname = reverse[p][:len(reverse[p])-1] + name
		} else {
			imgname = name
		}
		imgpath := cache + "/" + imgname + ".png"
		data, err := web.RequestDataWith(web.NewTLS12Client(), imgurl, "GET", "gitcode.net", web.RandUA(), nil)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			return
		}
		f, err := os.Create(imgpath)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			return
		}
		defer f.Close()
		err = os.WriteFile(f.Name(), data, 0755)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			return
		}
		process.SleepAbout1sTo2s()
		_, err = ctx.SendImage("file:///"+file.BOTPATH+"/"+imgpath, false, reasons[rand.Intn(len(reasons))], position[p], "的『", name, "』\n其释义为: ", description)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
		}

	})

	engine.OnMessageRegex(`^解塔罗牌\s?(.*)`, getTarot).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *nano.Ctx) {
		match := ctx.State["regex_matched"].([]string)[1]
		info, ok := infoMap[match]
		if ok {
			card := cardMap[match]
			imgname := card.Name
			imgpath := cache + "/" + imgname + ".png"
			if file.IsNotExist(imgpath) {
				imgurl := bed + info.ImgURL
				data, err := web.RequestDataWith(web.NewTLS12Client(), imgurl, "GET", "gitcode.net", web.RandUA(), nil)
				if err != nil {
					_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
					return
				}
				f, err := os.Create(imgpath)
				if err != nil {
					_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
					return
				}
				defer f.Close()
				err = os.WriteFile(f.Name(), data, 0755)
				if err != nil {
					_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
					return
				}
				process.SleepAbout1sTo2s()
			}
			_, err = ctx.SendImage("file:///"+file.BOTPATH+"/"+imgpath, false, "\n", match, "的含义是~\n『正位』:", info.Description, "\n『逆位』:", info.ReverseDescription)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			}
			return
		}
		var build strings.Builder
		build.WriteString("塔罗牌列表\n大阿尔卡纳:\n")
		build.WriteString(strings.Join(majorArcanaName[:7], " "))
		build.WriteString("\n")
		build.WriteString(strings.Join(majorArcanaName[7:14], " "))
		build.WriteString("\n")
		build.WriteString(strings.Join(majorArcanaName[14:22], " "))
		build.WriteString("\n小阿尔卡纳:\n[圣杯|星币|宝剑|权杖] [0-10|侍从|骑士|王后|国王]")
		txt := build.String()
		cardList, err := text.RenderToBase64(txt, text.FontFile, 420, 20)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			return
		}
		_, err = ctx.SendImage("base64://"+binary.BytesToString(cardList), false, "没有找到", match, "噢~")
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
		}
	})
}
