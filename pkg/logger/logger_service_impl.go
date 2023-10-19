package logger

import (
	"context"
	"fmt"
	"log"

	"github.com/pjmessi/golang-practice/pkg/ctxutil"
)

type ServiceImpl struct {
}

func NewService() Service {
	return &ServiceImpl{}
}

func (s *ServiceImpl) Debug(msg string) {
	msgToPrint := fmt.Sprintf("DEBUG: %s", msg)
	log.Println(msgToPrint)
}

func (s *ServiceImpl) Error(msg string) {
	msgToPrint := fmt.Sprintf("ERROR: %s", msg)
	log.Println(msgToPrint)
}

func (s *ServiceImpl) DebugCtx(ctx context.Context, msg string) {
	traceId := ctxutil.GetTraceIdFromCtx(ctx)
	if traceId == "" {
		s.Debug(msg)
		return
	}
	msgToPrint := fmt.Sprintf("DEBUG: [TraceId: %s] %s", traceId, msg)
	log.Println(msgToPrint)
}

func (s *ServiceImpl) ErrorCtx(ctx context.Context, msg string) {
	traceId := ctxutil.GetTraceIdFromCtx(ctx)
	if traceId == "" {
		s.Error(msg)
		return
	}
	msgToPrint := fmt.Sprintf("ERROR: [TraceId: %s] %s", traceId, msg)
	log.Println(msgToPrint)
}
