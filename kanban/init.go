// Package kanban 打印版本信息
package kanban

import "fmt"

//go:generate go run github.com/FloatTech/NanoBot-Plugin/kanban/gen

func init() {
	PrintBanner()
}

func PrintBanner() {
	fmt.Print(
		"\n======================[NanoBot-Plugin]======================",
		"\n", Banner, "\n",
		"===========================================================\n\n",
	)
}
