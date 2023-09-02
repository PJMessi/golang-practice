package restapi

import (
	"github.com/pjmessi/go-database-usage/internal/service/auth"
	"github.com/pjmessi/go-database-usage/internal/service/user"
)

type RouteHandler struct {
	authFacade auth.Facade
	userFacade user.Facade
}

func NewRouteHandler(authFacade auth.Facade, userFacade user.Facade) *RouteHandler {
	return &RouteHandler{
		authFacade: authFacade,
		userFacade: userFacade,
	}
}
