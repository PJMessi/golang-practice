package ctxutil

import "context"

type contextKey string

func NewCtxWithTraceId(traceId string) context.Context {
	return context.WithValue(context.Background(), contextKey("TraceId"), traceId)
}

func GetTraceIdFromCtx(ctx context.Context) string {
	traceIdVal := ctx.Value(contextKey("TraceId"))
	traceId, ok := traceIdVal.(string)
	if !ok {
		return ""
	}

	return traceId
}

func AddValue(ctx context.Context, key string, value interface{}) context.Context {
	return context.WithValue(ctx, contextKey(key), value)
}

func GetValue(ctx context.Context, key string) interface{} {
	value := ctx.Value(contextKey(key))
	return value
}
