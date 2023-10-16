// Package score 签到
package score

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"strconv"
	"sync"
	"time"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/gg"
	"github.com/FloatTech/rendercard"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/disintegration/imaging"
	nano "github.com/fumiama/NanoBot"

	"github.com/FloatTech/NanoBot-Plugin/kanban"
)

func floatstyle(a *scoredata) (img image.Image, err error) {
	fontdata, err := file.GetLazyData(text.GlowSansFontFile, nano.Md5File, false)
	if err != nil {
		return
	}

	getAvatar, err := initPic(a.picfile, a.avatarurl)
	if err != nil {
		return
	}

	back, err := gg.LoadImage(a.picfile)
	if err != nil {
		return
	}

	bx, by := float64(back.Bounds().Dx()), float64(back.Bounds().Dy())

	sc := 1280 / bx

	colors := gg.TakeColor(back, 3)

	canvas := gg.NewContext(1280, 1280*int(by)/int(bx))

	cw, ch := float64(canvas.W()), float64(canvas.H())

	sch := ch * 6 / 10

	var blurback, scbackimg, backshadowimg, avatarimg, avatarbackimg, avatarshadowimg, whitetext, blacktext image.Image
	var wg sync.WaitGroup
	wg.Add(8)

	go func() {
		defer wg.Done()
		scback := gg.NewContext(canvas.W(), canvas.H())
		scback.ScaleAbout(sc, sc, cw/2, ch/2)
		scback.DrawImageAnchored(back, canvas.W()/2, canvas.H()/2, 0.5, 0.5)
		scback.Identity()

		go func() {
			defer wg.Done()
			blurback = imaging.Blur(scback.Image(), 20)
		}()

		scbackimg = rendercard.Fillet(scback.Image(), 12)
	}()

	go func() {
		defer wg.Done()
		pureblack := gg.NewContext(canvas.W(), canvas.H())
		pureblack.SetRGBA255(0, 0, 0, 255)
		pureblack.Clear()

		shadow := gg.NewContext(canvas.W(), canvas.H())
		shadow.ScaleAbout(0.6, 0.6, cw-cw/3, ch/2)
		shadow.DrawImageAnchored(pureblack.Image(), canvas.W()-canvas.W()/3, canvas.H()/2, 0.5, 0.5)
		shadow.Identity()

		backshadowimg = imaging.Blur(shadow.Image(), 8)
	}()

	aw, ah := (ch-sch)/2/2/2*3, (ch-sch)/2/2/2*3

	go func() {
		defer wg.Done()
		avatar, _, err := image.Decode(bytes.NewReader(getAvatar))
		if err != nil {
			return
		}

		isc := (ch - sch) / 2 / 2 / 2 * 3 / float64(avatar.Bounds().Dy())

		scavatar := gg.NewContext(int(aw), int(ah))

		scavatar.ScaleAbout(isc, isc, aw/2, ah/2)
		scavatar.DrawImageAnchored(avatar, scavatar.W()/2, scavatar.H()/2, 0.5, 0.5)
		scavatar.Identity()

		avatarimg = rendercard.Fillet(scavatar.Image(), 8)
	}()

	err = canvas.ParseFontFace(fontdata, (ch-sch)/2/2/2)
	if err != nil {
		return
	}
	namew, _ := canvas.MeasureString(a.nickname)

	go func() {
		defer wg.Done()
		avatarshadowimg = imaging.Blur(customrectangle(cw, ch, aw, ah, namew, color.Black), 8)
	}()

	go func() {
		defer wg.Done()
		avatarbackimg = customrectangle(cw, ch, aw, ah, namew, colors[0])
	}()

	go func() {
		defer wg.Done()
		whitetext, err = customtext(a, fontdata, cw, ch, aw, color.White)
		if err != nil {
			return
		}
	}()

	go func() {
		defer wg.Done()
		blacktext, err = customtext(a, fontdata, cw, ch, aw, color.Black)
		if err != nil {
			return
		}
	}()

	wg.Wait()
	if scbackimg == nil || backshadowimg == nil || avatarimg == nil || avatarbackimg == nil || avatarshadowimg == nil || whitetext == nil || blacktext == nil {
		err = errors.New("图片渲染失败")
		return
	}

	canvas.DrawImageAnchored(blurback, canvas.W()/2, canvas.H()/2, 0.5, 0.5)

	canvas.DrawImage(backshadowimg, 0, 0)

	canvas.ScaleAbout(0.6, 0.6, cw-cw/3, ch/2)
	canvas.DrawImageAnchored(scbackimg, canvas.W()-canvas.W()/3, canvas.H()/2, 0.5, 0.5)
	canvas.Identity()

	canvas.DrawImage(avatarshadowimg, 0, 0)
	canvas.DrawImage(avatarbackimg, 0, 0)
	canvas.DrawImageAnchored(avatarimg, int((ch-sch)/2/2), int((ch-sch)/2/2), 0.5, 0.5)

	canvas.DrawImage(blacktext, 2, 2)
	canvas.DrawImage(whitetext, 0, 0)

	img = canvas.Image()
	return
}

func customrectangle(cw, ch, aw, ah, namew float64, rtgcolor color.Color) (img image.Image) {
	canvas := gg.NewContext(int(cw), int(ch))
	sch := ch * 6 / 10
	canvas.DrawRoundedRectangle((ch-sch)/2/2-aw/2-aw/40, (ch-sch)/2/2-aw/2-ah/40, aw+aw/40*2, ah+ah/40*2, 8)
	canvas.SetColor(rtgcolor)
	canvas.Fill()
	canvas.DrawRoundedRectangle((ch-sch)/2/2, (ch-sch)/2/2-ah/4, aw/2+aw/40*5+namew, ah/2, 8)
	canvas.Fill()

	img = canvas.Image()
	return
}

func customtext(a *scoredata, fontdata []byte, cw, ch, aw float64, textcolor color.Color) (img image.Image, err error) {
	canvas := gg.NewContext(int(cw), int(ch))
	canvas.SetColor(textcolor)
	scw, sch := cw*6/10, ch*6/10
	err = canvas.ParseFontFace(fontdata, (ch-sch)/2/2/2)
	if err != nil {
		return
	}
	canvas.DrawStringAnchored(a.nickname, (ch-sch)/2/2+aw/2+aw/40*2, (ch-sch)/2/2, 0, 0.5)
	err = canvas.ParseFontFace(fontdata, (ch-sch)/2/2/3*2)
	if err != nil {
		return
	}
	canvas.DrawStringAnchored(time.Now().Format("2006/01/02"), cw-cw/6, ch/2-sch/2-canvas.FontHeight(), 0.5, 0.5)

	err = canvas.ParseFontFace(fontdata, (ch-sch)/2/2/2)
	if err != nil {
		return
	}
	nextrankScore := 0
	if a.rank < 10 {
		nextrankScore = rankArray[a.rank+1]
	} else {
		nextrankScore = SCOREMAX
	}
	nextLevelStyle := strconv.Itoa(a.level) + "/" + strconv.Itoa(nextrankScore)

	canvas.DrawStringAnchored("Level "+strconv.Itoa(a.rank), cw/3*2-scw/2, ch/2+sch/2+canvas.FontHeight(), 0, 0.5)
	canvas.DrawStringAnchored(nextLevelStyle, cw/3*2+scw/2, ch/2+sch/2+canvas.FontHeight(), 1, 0.5)

	err = canvas.ParseFontFace(fontdata, (ch-sch)/2/2/3)
	if err != nil {
		return
	}

	canvas.DrawStringAnchored("Create By NanoBot-Plugin "+kanban.Version, 0+4, ch, 0, -0.5)

	err = canvas.ParseFontFace(fontdata, (ch-sch)/2/5*3)
	if err != nil {
		return
	}

	tempfh := canvas.FontHeight()

	canvas.DrawStringAnchored(getHourWord(time.Now()), ((cw-scw)-(cw/3-scw/2))/8, (ch-sch)/2+sch/4, 0, 0.5)

	err = canvas.ParseFontFace(fontdata, (ch-sch)/2/5)
	if err != nil {
		return
	}

	canvas.DrawStringAnchored("ATRI币 + "+strconv.Itoa(a.inc), ((cw-scw)-(cw/3-scw/2))/8, (ch-sch)/2+sch/4+tempfh, 0, 0.5)
	canvas.DrawStringAnchored("EXP + 1", ((cw-scw)-(cw/3-scw/2))/8, (ch-sch)/2+sch/4+tempfh+canvas.FontHeight(), 0, 1)

	err = canvas.ParseFontFace(fontdata, (ch-sch)/2/4)
	if err != nil {
		return
	}

	canvas.DrawStringAnchored("你有 "+strconv.Itoa(a.score)+" 枚ATRI币", ((cw-scw)-(cw/3-scw/2))/8, (ch-sch)/2+sch/4*3, 0, 0.5)

	img = canvas.Image()
	return
}
