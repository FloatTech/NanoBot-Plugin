// Package score 签到
package score

import (
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/AnimeAPI/bilibili"
	"github.com/FloatTech/AnimeAPI/wallet"
	"github.com/FloatTech/NanoBot-Plugin/utils/ctxext"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/process"
	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/imgfactory"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/img/text"
	nano "github.com/fumiama/NanoBot"
	"github.com/golang/freetype"
	"github.com/wcharczuk/go-chart/v2"
)

const (
	backgroundURL = "https://iw233.cn/api.php?sort=pc"
	referer       = "https://weibo.com/"
	signinMax     = 1
	// SCOREMAX 分数上限定为1200
	SCOREMAX = 1200
)

var (
	rankArray = [...]int{0, 10, 20, 50, 100, 200, 350, 550, 750, 1000, 1200}
	engine    = nano.Register("score", &ctrl.Options[*nano.Ctx]{
		DisableOnDefault:  false,
		Brief:             "签到",
		Help:              "- 签到\n | 获得签到背景\n- 查看等级排名\n注:为跨群排名\n- 查看我的钱包\n- 查看钱包排名\n注:为本群排行，若群人数太多不建议使用该功能!!!",
		PrivateDataFolder: "score",
	})
)

func init() {
	cachePath := engine.DataFolder() + "cache/"
	go func() {
		ok := file.IsExist(cachePath)
		if !ok {
			err := os.MkdirAll(cachePath, 0777)
			if err != nil {
				panic(err)
			}
			return
		}
		files, err := os.ReadDir(cachePath)
		if err == nil {
			for _, f := range files {
				if !strings.Contains(f.Name(), time.Now().Format("20060102")) {
					_ = os.Remove(cachePath + f.Name())
				}
			}
		}
		sdb = initialize(engine.DataFolder() + "score.db")
	}()
	engine.OnMessageRegex("签到").Limit(ctxext.LimitByUser).SetBlock(true).Handle(func(ctx *nano.Ctx) {
		uid := ctx.Message.Author.ID
		if uid == "" {
			_, _ = ctx.SendPlainMessage(false, "ERROR: 未获取到用户uid")
			return
		}
		uidint, _ := strconv.ParseInt(uid, 10, 64)
		today := time.Now().Format("20060102")
		// 签到图片
		drawedFile := cachePath + uid + today + "signin.png"
		picFile := cachePath + uid + today + ".png"
		// 获取签到时间
		si := sdb.GetSignInByUID(uidint)
		siUpdateTimeStr := si.UpdatedAt.Format("20060102")
		switch {
		case si.Count >= signinMax && siUpdateTimeStr == today:
			// 如果签到时间是今天
			_, err := ctx.SendPlainMessage(true, "今天你已经签到过了！")
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			}
			if file.IsExist(drawedFile) {
				_, err := ctx.SendImage("file:///"+file.BOTPATH+"/"+drawedFile, false)
				if err != nil {
					_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				}
			}
			return
		case siUpdateTimeStr != today:
			// 如果是跨天签到就清数据
			err := sdb.InsertOrUpdateSignInCountByUID(uidint, 0)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return
			}
		}
		// 更新签到次数
		err := sdb.InsertOrUpdateSignInCountByUID(uidint, si.Count+1)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			return
		}
		// 更新经验
		level := sdb.GetScoreByUID(uidint).Score + 1
		if level > SCOREMAX {
			level = SCOREMAX
			_, err := ctx.SendPlainMessage(true, "你的等级已经达到上限")
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			}
		}
		err = sdb.InsertOrUpdateScoreByUID(uidint, level)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			return
		}
		// 更新钱包
		rank := getrank(level)
		add := 1 + rand.Intn(10) + rank*5 // 等级越高获得的钱越高

		go func() {
			err = wallet.InsertWalletOf(uidint, add)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return
			}
		}()
		alldata := &scoredata{
			drawedfile: drawedFile,
			picfile:    picFile,
			avatarurl:  ctx.Message.Author.Avatar,
			nickname:   ctx.Message.Author.Username,
			inc:        add,
			score:      wallet.GetWalletOf(uidint),
			level:      level,
			rank:       rank,
		}
		drawimage, err := floatstyle(alldata)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			return
		}
		// done.
		f, err := os.Create(drawedFile)
		if err != nil {
			data, err := imgfactory.ToBytes(drawimage)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return
			}
			_, err = ctx.SendImage("base64://"+nano.BytesToString(data), false)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			}
			return
		}
		_, err = imgfactory.WriteTo(drawimage, f)
		defer f.Close()
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			return
		}
		_, err = ctx.SendImage("file:///"+file.BOTPATH+"/"+drawedFile, false)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
		}
	})

	engine.OnMessageFullMatch("获得签到背景", nano.OnlyPublic).Limit(ctxext.LimitByGroup).SetBlock(true).
		Handle(func(ctx *nano.Ctx) {
			uidStr := ctx.Message.Author.ID
			picFile := cachePath + uidStr + time.Now().Format("20060102") + ".png"
			if file.IsNotExist(picFile) {
				_, err := ctx.SendPlainMessage(true, "请先签到！")
				if err != nil {
					_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				}
				return
			}
			if _, err := ctx.SendImage("file:///"+file.BOTPATH+"/"+picFile, false); err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			}
		})
	engine.OnMessageFullMatch("查看等级排名", nano.OnlyPublic).Limit(ctxext.LimitByGroup).SetBlock(true).
		Handle(func(ctx *nano.Ctx) {
			today := time.Now().Format("20060102")
			drawedFile := cachePath + today + "scoreRank.png"
			if file.IsExist(drawedFile) {
				_, err := ctx.SendImage("file:///"+file.BOTPATH+"/"+drawedFile, false)
				if err != nil {
					_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)

				}
				return
			}
			st, err := sdb.GetScoreRankByTopN(10)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return
			}
			if len(st) == 0 {
				_, _ = ctx.SendPlainMessage(false, "ERROR: 目前还没有人签到过")
				return
			}
			_, err = file.GetLazyData(text.FontFile, nano.Md5File, true)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return
			}
			b, err := os.ReadFile(text.FontFile)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return
			}
			font, err := freetype.ParseFont(b)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return
			}
			f, err := os.Create(drawedFile)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return
			}
			var bars []chart.Value
			for _, v := range st {
				if v.Score != 0 {
					bars = append(bars, chart.Value{
						Label: ctx.Message.Author.Username,
						Value: float64(v.Score),
					})
				}
			}
			err = chart.BarChart{
				Font:  font,
				Title: "等级排名(1天只刷新1次)",
				Background: chart.Style{
					Padding: chart.Box{
						Top: 40,
					},
				},
				YAxis: chart.YAxis{
					Range: &chart.ContinuousRange{
						Min: 0,
						Max: math.Ceil(bars[0].Value/10) * 10,
					},
				},
				Height:   500,
				BarWidth: 50,
				Bars:     bars,
			}.Render(chart.PNG, f)
			_ = f.Close()
			if err != nil {
				_ = os.Remove(drawedFile)
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return
			}
			if _, err := ctx.SendImage("file:///"+file.BOTPATH+"/"+drawedFile, false); err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			}
		})
}

func getHourWord(t time.Time) string {
	h := t.Hour()
	switch {
	case 6 <= h && h < 12:
		return "早上好"
	case 12 <= h && h < 14:
		return "中午好"
	case 14 <= h && h < 19:
		return "下午好"
	case 19 <= h && h < 24:
		return "晚上好"
	case 0 <= h && h < 6:
		return "凌晨好"
	default:
		return ""
	}
}

func getrank(count int) int {
	for k, v := range rankArray {
		if count == v {
			return k
		} else if count < v {
			return k - 1
		}
	}
	return -1
}

func initPic(picFile string, avatarurl string) (avatar []byte, err error) {
	defer process.SleepAbout1sTo2s()
	avatar, err = web.GetData(avatarurl)
	if err != nil {
		return
	}
	if file.IsExist(picFile) {
		return
	}
	url, err := bilibili.GetRealURL(backgroundURL)
	if err != nil {
		return
	}
	data, err := web.RequestDataWith(web.NewDefaultClient(), url, "", referer, "", nil)
	if err != nil {
		return
	}
	return avatar, os.WriteFile(picFile, data, 0644)
}
