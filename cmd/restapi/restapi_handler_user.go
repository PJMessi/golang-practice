package restapi

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pjmessi/golang-practice/internal/pkg/jwt"
	"github.com/pjmessi/golang-practice/pkg/ctxutil"
)

func (rh *RouteHandler) handleUserRegApi(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		rh.writeHttpResFromErr(ctx, w, err)
		return
	}

	// TODO: TEST
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			rh.writeHttpResFromErr(ctx, w, err)
		}
	}(r.Body)

	resByte, err := rh.userFacade.RegisterUser(ctx, reqBytes)
	if err != nil {
		rh.writeHttpResFromErr(ctx, w, err)
		return
	}

	_, err = w.Write(resByte)
	if err != nil {
		rh.writeHttpResFromErr(ctx, w, err)
		return
	}
}

func (rh *RouteHandler) handleGetProfileApi(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	jwtPayload, ok := ctxutil.GetValue(ctx, "jwtPayload").(jwt.JwtPayload)
	if !ok {
		rh.writeHttpResFromErr(ctx, w, fmt.Errorf("error getting jwt payload"))
		return
	}

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		rh.writeHttpResFromErr(ctx, w, err)
		return
	}

	// TODO: TEST
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			rh.writeHttpResFromErr(ctx, w, err)
		}
	}(r.Body)

	resByte, err := rh.userFacade.GetProfile(ctx, reqBytes, jwtPayload)
	if err != nil {
		rh.writeHttpResFromErr(ctx, w, err)
		return
	}

	_, err = w.Write(resByte)
	if err != nil {
		rh.writeHttpResFromErr(ctx, w, err)
		return
	}
}
