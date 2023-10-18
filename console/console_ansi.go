//go:build !windows

// Package console sets console's behavior on init
package console

import (
	"fmt"

	"github.com/FloatTech/NanoBot-Plugin/kanban/banner"
)

func init() {
	fmt.Print("\033]0;NanoBot-Blugin " + banner.Version + " " + banner.Copyright + "\007")
}
