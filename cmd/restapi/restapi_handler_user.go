package restapi

import (
	"context"
	"io"
	"net/http"

	"github.com/pjmessi/golang-practice/internal/pkg/jwt"
)

func (rh *RouteHandler) handleUserRegApi(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		rh.handleErr(ctx, w, err)
		return
	}

	// TODO: TEST
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			rh.handleErr(ctx, w, err)
		}
	}(r.Body)

	resByte, err := rh.userFacade.RegisterUser(ctx, reqBytes)
	if err != nil {
		rh.handleErr(ctx, w, err)
		return
	}

	_, err = w.Write(resByte)
	if err != nil {
		rh.handleErr(ctx, w, err)
		return
	}
}

func (rh *RouteHandler) handleGetProfileApi(ctx context.Context, jwtPayload jwt.JwtPayload, w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		rh.handleErr(ctx, w, err)
		return
	}

	// TODO: TEST
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			rh.handleErr(ctx, w, err)
		}
	}(r.Body)

	resByte, err := rh.userFacade.GetProfile(ctx, reqBytes, jwtPayload)
	if err != nil {
		rh.handleErr(ctx, w, err)
		return
	}

	_, err = w.Write(resByte)
	if err != nil {
		rh.handleErr(ctx, w, err)
		return
	}
}
