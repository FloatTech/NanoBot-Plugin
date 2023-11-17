// Package hyaku 百人一首
package hyaku

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"unsafe"

	nano "github.com/fumiama/NanoBot"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"

	"github.com/FloatTech/NanoBot-Plugin/utils/ctxext"
)

const bed = "https://gitea.seku.su/fumiama/OguraHyakuninIsshu/raw/branch/master/"

//nolint:asciicheck, structcheck
type line struct {
	番号, 歌人, 上の句, 下の句, 上の句ひらがな, 下の句ひらがな string
}

func (l *line) String() string {
	b := binary.NewWriterF(func(w *binary.Writer) {
		r := reflect.ValueOf(l).Elem().Type()
		for i := 0; i < r.NumField(); i++ {
			switch i {
			case 0:
				w.WriteString("●")
			case 1:
				w.WriteString("◉")
			case 2, 3:
				w.WriteString("○")
			case 4, 5:
				w.WriteString("◎")
			}
			w.WriteString(r.Field(i).Name)
			w.WriteString(": ")
			w.WriteString((*[6]string)(unsafe.Pointer(l))[i])
			w.WriteString("\n")
		}
	})
	return binary.BytesToString(b)
}

var lines [100]*line

func init() {
	engine := nano.Register("hyaku", &ctrl.Options[*nano.Ctx]{
		DisableOnDefault: false,
		Help: "百人一首\n" +
			"- 百人一首(随机发一首)\n" +
			"- 百人一首之n",
		PrivateDataFolder: "hyaku",
	})
	csvfile := engine.DataFolder() + "hyaku.csv"
	go func() {
		var f *os.File
		if file.IsNotExist(csvfile) {
			data, err := web.RequestDataWith(web.NewTLS12Client(), bed+"小倉百人一首.csv", "GET", "gitcode.net", web.RandUA(), nil)
			if err != nil {
				_ = os.Remove(csvfile)
				panic(err)
			}
			f, err = os.Create(csvfile)
			if err != nil {
				panic(err)
			}
			_, _ = f.Write(data)
			_, _ = f.Seek(0, io.SeekStart)
		} else {
			var err error
			f, err = os.Open(csvfile)
			if err != nil {
				panic(err)
			}
		}
		records, err := csv.NewReader(f).ReadAll()
		if err != nil {
			panic(err)
		}
		_ = f.Close()
		records = records[1:] // skip title
		if len(records) != 100 {
			panic("invalid csvfile")
		}
		for j, r := range records {
			if len(r) != 6 {
				panic("invalid csvfile")
			}
			i, err := strconv.Atoi(r[0])
			if err != nil {
				panic(err)
			}
			i--
			if j != i {
				panic("invalid csvfile")
			}
			lines[i] = (*line)(*(*unsafe.Pointer)(unsafe.Pointer(&r)))
		}
	}()
	engine.OnMessageFullMatch("百人一首").SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *nano.Ctx) {
		i := rand.Intn(100)
		_, _ = ctx.SendImage(fmt.Sprintf(bed+"img/%03d.jpg", i+1), false, lines[i].String())
		_, _ = ctx.SendImage(fmt.Sprintf(bed+"img/%03d.png", i+1), false)
	})
	engine.OnMessageRegex(`^百人一首之\s?(\d+)$`).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *nano.Ctx) {
		i, err := strconv.Atoi(ctx.State["regex_matched"].([]string)[1])
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			return
		}
		if i > 100 || i < 1 {
			_, _ = ctx.SendPlainMessage(false, "ERROR:超出范围")
			return
		}
		_, _ = ctx.SendImage(fmt.Sprintf(bed+"img/%03d.jpg", i), false)
		_, _ = ctx.SendImage(fmt.Sprintf(bed+"img/%03d.png", i), false, lines[i-1].String())
	})
}
