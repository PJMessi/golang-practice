package restapi

import (
	"github.com/pjmessi/golang-practice/internal/service/auth"
	"github.com/pjmessi/golang-practice/internal/service/user"
	"github.com/pjmessi/golang-practice/pkg/logger"
)

type RouteHandler struct {
	authFacade auth.Facade
	userFacade user.Facade
	logService logger.Service
}

func NewRouteHandler(logService logger.Service, authFacade auth.Facade, userFacade user.Facade) *RouteHandler {
	return &RouteHandler{
		authFacade: authFacade,
		userFacade: userFacade,
		logService: logService,
	}
}
