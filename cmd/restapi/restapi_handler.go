package restapi

import (
	"github.com/pjmessi/golang-practice/internal/service/auth"
	"github.com/pjmessi/golang-practice/internal/service/user"
	"github.com/pjmessi/golang-practice/pkg/logger"
)

type RouteHandler struct {
	authFacade auth.Facade
	userFacade user.Facade
	loggerUtil logger.Util
}

func NewRouteHandler(loggerUtil logger.Util, authFacade auth.Facade, userFacade user.Facade) *RouteHandler {
	return &RouteHandler{
		authFacade: authFacade,
		userFacade: userFacade,
		loggerUtil: loggerUtil,
	}
}
