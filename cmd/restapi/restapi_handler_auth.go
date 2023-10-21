package restapi

import (
	"io"
	"net/http"
)

func (rh *RouteHandler) handleLoginApi(w http.ResponseWriter, r *http.Request) {
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

	resByte, err := rh.authFacade.Login(ctx, reqBytes)
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
