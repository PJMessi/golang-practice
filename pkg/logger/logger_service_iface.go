package logger

import "context"

type Service interface {
	Debug(msg string)
	Error(msg string)
	DebugCtx(ctx context.Context, msg string)
	ErrorCtx(ctx context.Context, msg string)
}
