// Package qunyou ...
package qunyou

import (
	"os/exec"
	"strings"

	nano "github.com/fumiama/NanoBot"

	"github.com/FloatTech/NanoBot-Plugin/utils/ctxext"
	ctrl "github.com/FloatTech/zbpctrl"
)

func init() {
	en := nano.Register("qunyou", &ctrl.Options[*nano.Ctx]{
		DisableOnDefault: false,
		Help:             "随机群友怪话\n- 看看群友",
	})
	en.OnMessagePrefix("看看群友").Limit(ctxext.LimitByGroup).Handle(func(ctx *nano.Ctx) {
		prompt := ctx.State["args"].(string)
		sb := strings.Builder{}
		cmd := exec.Cmd{
			Path:   "/usr/local/bin/llama2.run",
			Args:   []string{"model.bin"},
			Dir:    "/usr/local/src/llama2.c",
			Stdout: &sb,
		}
		if prompt != "" {
			cmd.Args = append(cmd.Args, "-i", prompt)
		}
		err := cmd.Run()
		if err != nil {
			ctx.SendChain(nano.Text("ERROR: ", err))
			return
		}
		ctx.SendChain(nano.Text(sb.String()))
	})
}
