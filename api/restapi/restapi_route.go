package restapi

import (
	"net/http"

	"github.com/pjmessi/go-database-usage/internal/service/auth"
	"github.com/pjmessi/go-database-usage/internal/service/user"
	"github.com/pjmessi/go-database-usage/pkg/logger"
	"github.com/pjmessi/go-database-usage/pkg/uuid"

	"github.com/gorilla/mux"
)

func RegisterRoutes(loggerUtil logger.Util, authFacade auth.Facade, userFacade user.Facade, uuidUtil uuid.Util) http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	handler := NewRouteHandler(loggerUtil, authFacade, userFacade, uuidUtil)

	router.HandleFunc("/auth/login", handler.handlePanic(handler.handleLoginApi)).Methods("POST")
	router.HandleFunc("/user/registration", handler.handlePanic(handler.handleUserRegApi)).Methods("POST")

	return router
}
