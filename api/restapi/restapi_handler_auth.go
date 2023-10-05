package restapi

import (
	"context"
	"io"
	"net/http"
)

func (rh *RouteHandler) handleLoginApi(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		rh.handleErr(w, err)
		return
	}

	// TODO: TEST
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			rh.handleErr(w, err)
		}
	}(r.Body)

	resByte, err := rh.authFacade.Login(ctx, reqBytes)
	if err != nil {
		rh.handleErr(w, err)
		return
	}

	_, err = w.Write(resByte)
	if err != nil {
		rh.handleErr(w, err)
		return
	}
}
