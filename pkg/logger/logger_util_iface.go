package logger

import "context"

type Util interface {
	Debug(msg string)
	Error(msg string)
	DebugCtx(ctx context.Context, msg string)
	ErrorCtx(ctx context.Context, msg string)
}
