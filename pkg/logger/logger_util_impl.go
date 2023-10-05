package logger

import (
	"context"
	"fmt"
	"log"

	"github.com/pjmessi/go-database-usage/pkg/ctxutil"
)

type UtilImpl struct {
}

func NewUtil() Util {
	return &UtilImpl{}
}

func (u *UtilImpl) Debug(msg string) {
	msgToPrint := fmt.Sprintf("DEBUG: %s", msg)
	log.Println(msgToPrint)
}

func (u *UtilImpl) Error(msg string) {
	msgToPrint := fmt.Sprintf("ERROR: %s", msg)
	log.Println(msgToPrint)
}

func (u *UtilImpl) DebugCtx(ctx context.Context, msg string) {
	traceId := ctxutil.GetTraceIdFromCtx(ctx)
	msgToPrint := fmt.Sprintf("DEBUG: [TraceId: %s] %s", traceId, msg)
	log.Println(msgToPrint)
}

func (u *UtilImpl) ErrorCtx(ctx context.Context, msg string) {
	traceId := ctxutil.GetTraceIdFromCtx(ctx)
	msgToPrint := fmt.Sprintf("ERROR: [TraceId: %s] %s", traceId, msg)
	log.Println(msgToPrint)
}
