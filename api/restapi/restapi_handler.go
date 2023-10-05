package restapi

import (
	"github.com/pjmessi/go-database-usage/internal/service/auth"
	"github.com/pjmessi/go-database-usage/internal/service/user"
	"github.com/pjmessi/go-database-usage/pkg/logger"
	"github.com/pjmessi/go-database-usage/pkg/uuid"
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
