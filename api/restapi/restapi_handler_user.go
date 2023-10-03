package restapi

import (
	"io"
	"net/http"

	"github.com/pjmessi/go-database-usage/pkg/ctxutil"
)

func (rh *RouteHandler) handleUserRegApi(w http.ResponseWriter, r *http.Request) {
	traceId, err := rh.uuidUtil.GenUuidV4()
	if err != nil {
		rh.handleErr(w, err)
	}
	ctx := ctxutil.NewCtx(traceId)

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

	resByte, err := rh.userFacade.RegisterUser(ctx, reqBytes)
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
