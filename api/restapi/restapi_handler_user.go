package restapi

import (
	"io"
	"net/http"
)

func (rh *RouteHandler) handleUserRegApi(w http.ResponseWriter, r *http.Request) {
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

	resByte, err := rh.userFacade.RegisterUser(reqBytes)
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
