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
		DisableOnDefault: true,
		Help:             "随机群友怪话\n- 看看群友",
	})
	en.OnMessagePrefix("看看群友").Limit(ctxext.LimitByGroup).Handle(func(ctx *nano.Ctx) {
		prompt := ctx.State["args"].(string)
		sb := strings.Builder{}
		errsb := strings.Builder{}
		cmd := exec.Cmd{
			Path:   "/usr/local/bin/llama2.run",
			Args:   []string{"/usr/local/bin/llama2.run", "model.bin"},
			Dir:    "/usr/local/src/llama2.c",
			Stdout: &sb,
			Stderr: &errsb,
		}
		if prompt != "" {
			cmd.Args = append(cmd.Args, "-i", prompt)
		}
		err := cmd.Run()
		if err != nil {
			ctx.SendChain(nano.Text("ERROR: ", err, errsb.String()))
			return
		}
		ctx.SendChain(nano.Text(sb.String()))
	})
}
