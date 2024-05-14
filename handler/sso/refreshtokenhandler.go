package sso

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"nrs_customer_module_backend/internal/global"
	"nrs_customer_module_backend/internal/global/createfile"
	"nrs_customer_module_backend/internal/global/errorglobal"
	"nrs_customer_module_backend/internal/logic/sso"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func RefreshTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	logFile := "refresh_token"
	appLogger := createfile.New(logFile)

	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PostRefreshTokenReq
		
		// Validate data
		if err := httpx.Parse(r, &req); err != nil {
			if strings.Contains(fmt.Sprint(err), "is not set") {
				fieldName := global.ExtractFieldName(fmt.Sprint(err))
				httpx.WriteJsonCtx(r.Context(), w, 400, map[string]interface{}{"error": errorglobal.BadRequest(fieldName, fieldName+" is required")})
				return
			}
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// Log error
		jsonReq, err := json.Marshal(req)
		if err != nil {
			appLogger.Error(logFile).Printf("[x] Request:  %+v -> Error json.Marshal: %v", req, err) // Log
		}

		appLogger.Info(logFile).Printf("Incoming Request: %s", jsonReq)

		// Authentication
		user, _, _ := r.BasicAuth()
		otherInfo := map[string]interface{}{
			"refresh_token":  req.RefreshToken,
		}
		
		l := sso.NewRefreshTokenLogic(r.Context(), svcCtx, otherInfo)
		resp, err := l.RefreshToken(&req)
		if err != nil {
			// Log the request
			appLogger.Info(logFile).Printf("[x] Request: user: %s -> Error:  %v", user, err) // Log the r
			info := map[string]interface{}{
				"source_name": "authorized",
				"source_id":   user,
			}
			errorglobal.ReturnError(err, r, w, info)
			return
		}
		respJSON, err := json.Marshal(resp)
		if err != nil {
			// Log the error
			appLogger.Error(logFile).Printf("[x] Request:   user: %s -> Error json.Marshal: %v", user, err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		appLogger.Info(logFile).Printf("[x] Request: user: %s -> Response:  %s", user, respJSON) // Log the r
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
