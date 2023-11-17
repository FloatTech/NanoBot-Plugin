// Package ctxext ctx扩展
package ctxext

import (
	"strconv"
	"time"
	"unsafe"

	nano "github.com/fumiama/NanoBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
)

// DefaultSingle 默认反并发处理
//
//	按 发送者 反并发
//	并发时返回 "您有操作正在执行, 请稍后再试!"
var DefaultSingle = nano.NewSingle(
	nano.WithKeyFn(func(ctx *nano.Ctx) int64 {
		switch ctx.Value.(type) {
		case *nano.Message:
			return int64(ctx.UserID())
		}
		return 0
	}),
	nano.WithPostFn[int64](func(ctx *nano.Ctx) {
		_, _ = ctx.SendPlainMessage(false, "您有操作正在执行, 请稍后再试!")
	}),
)

// defaultLimiterManager 默认限速器管理
//
//	每 10s 5次触发
var defaultLimiterManager = rate.NewManager[int64](time.Second*10, 5)

//nolint:structcheck
type fakeLM struct {
	limiters unsafe.Pointer
	interval time.Duration
	burst    int
}

// SetDefaultLimiterManagerParam 设置默认限速器参数
//
//	每 interval 时间 burst 次触发
func SetDefaultLimiterManagerParam(interval time.Duration, burst int) {
	f := (*fakeLM)(unsafe.Pointer(defaultLimiterManager))
	f.interval = interval
	f.burst = burst
}

// LimitByUser 默认限速器 每 10s 5次触发
//
//	按 发送者 限制
func LimitByUser(ctx *nano.Ctx) *rate.Limiter {
	if _, ok := ctx.Value.(*nano.Message); ok {
		return defaultLimiterManager.Load(int64(ctx.UserID()))
	}
	return defaultLimiterManager.Load(0)
}

// LimitByGroup 默认限速器 每 10s 5次触发
//
//	按 group 限制
func LimitByGroup(ctx *nano.Ctx) *rate.Limiter {
	if _, ok := ctx.Value.(*nano.Message); ok {
		return defaultLimiterManager.Load(int64(ctx.GroupID()))
	}
	return defaultLimiterManager.Load(0)
}

// LimitByGuild 默认限速器 每 10s 5次触发
//
//	按 guild 限制
func LimitByGuild(ctx *nano.Ctx) *rate.Limiter {
	if msg, ok := ctx.Value.(*nano.Message); ok {
		id, _ := strconv.ParseUint(msg.GuildID, 10, 64)
		return defaultLimiterManager.Load(int64(id))
	}
	return defaultLimiterManager.Load(0)
}

// LimitByChannel 默认限速器 每 10s 5次触发
//
//	按 channel 限制
func LimitByChannel(ctx *nano.Ctx) *rate.Limiter {
	if _, ok := ctx.Value.(*nano.Message); ok {
		return defaultLimiterManager.Load(int64(ctx.GroupID()))
	}
	return defaultLimiterManager.Load(0)
}

// LimiterManager 自定义限速器管理
type LimiterManager struct {
	m *rate.LimiterManager[int64]
}

// NewLimiterManager 新限速器管理
func NewLimiterManager(interval time.Duration, burst int) (m LimiterManager) {
	m.m = rate.NewManager[int64](interval, burst)
	return
}

// LimitByUser 自定义限速器
//
//	按 发送者 限制
func (m LimiterManager) LimitByUser(ctx *nano.Ctx) *rate.Limiter {
	if _, ok := ctx.Value.(*nano.Message); ok {
		return defaultLimiterManager.Load(int64(ctx.UserID()))
	}
	return defaultLimiterManager.Load(0)
}

// LimitByGuild 自定义限速器
//
//	按 guild 限制
func (m LimiterManager) LimitByGuild(ctx *nano.Ctx) *rate.Limiter {
	if msg, ok := ctx.Value.(*nano.Message); ok {
		id, _ := strconv.ParseUint(msg.GuildID, 10, 64)
		return defaultLimiterManager.Load(int64(id))
	}
	return defaultLimiterManager.Load(0)
}

// LimitByGroup 自定义限速器
//
//	按 group 限制
func (m LimiterManager) LimitByGroup(ctx *nano.Ctx) *rate.Limiter {
	if _, ok := ctx.Value.(*nano.Message); ok {
		return defaultLimiterManager.Load(int64(ctx.GroupID()))
	}
	return defaultLimiterManager.Load(0)
}

// LimitByChannel 自定义限速器
//
//	按 channel 限制
func (m LimiterManager) LimitByChannel(ctx *nano.Ctx) *rate.Limiter {
	if _, ok := ctx.Value.(*nano.Message); ok {
		return defaultLimiterManager.Load(int64(ctx.GroupID()))
	}
	return defaultLimiterManager.Load(0)
}

// MustMessageNotNil 消息是否不为空
func MustMessageNotNil(ctx *nano.Ctx) bool {
	return ctx.Message != nil
}
