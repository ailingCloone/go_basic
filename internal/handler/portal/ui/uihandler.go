package ui

import (
	"encoding/json"
	"net/http"

	"nrs_customer_module_backend/internal/global/createfile"
	"nrs_customer_module_backend/internal/global/errorglobal"
	"nrs_customer_module_backend/internal/logic/portal/ui"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {

	logFile := "get_login_ui"
	appLogger := createfile.New(logFile)

	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PostLoginUIReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		jsonReq, err := json.Marshal(req)
		if err != nil {
			appLogger.Error(logFile).Printf("[x] Request:  %+v -> Error json.Marshal: %v", req, err) // Log
		}

		appLogger.Info(logFile).Printf("Incoming Request: %s", jsonReq)

		l := ui.NewUiLogic(r.Context(), svcCtx)
		resp, err := l.Ui(&req)
		if err != nil {
			// Log the request
			appLogger.Error(logFile).Printf("[x] Request:  %+v -> Error: %v", req, err)
			info := map[string]interface{}{
				"source_name": "version",
				"source_id":   req.Version,
			}
			errorglobal.ReturnError(err, r, w, info)
			return
		}
		respJSON, err := json.Marshal(resp)
		if err != nil {
			// Log the error
			appLogger.Error(logFile).Printf("[x] Request:  %+v -> Error json.Marshal: %v", req, err)
		}

		appLogger.Info(logFile).Printf("[x] Request:  %s -> Response:  %s", jsonReq, respJSON) // Log the request and response
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
