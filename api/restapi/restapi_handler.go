package restapi

import (
	"github.com/pjmessi/golang-practice/internal/service/auth"
	"github.com/pjmessi/golang-practice/internal/service/user"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/pjmessi/golang-practice/pkg/uuid"
)

type RouteHandler struct {
	authFacade auth.Facade
	userFacade user.Facade
	uuidUtil   uuid.Util
	loggerUtil logger.Util
}

func NewRouteHandler(loggerUtil logger.Util, authFacade auth.Facade, userFacade user.Facade, uuidUtil uuid.Util) *RouteHandler {
	return &RouteHandler{
		authFacade: authFacade,
		userFacade: userFacade,
		uuidUtil:   uuidUtil,
		loggerUtil: loggerUtil,
	}
}
