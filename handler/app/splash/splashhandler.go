package splash

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nrs_customer_module_backend/internal/global"
	"nrs_customer_module_backend/internal/global/createfile"
	"nrs_customer_module_backend/internal/global/errorglobal"
	"nrs_customer_module_backend/internal/logic/app/splash"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func SplashHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	logFile := "splash"
	appLogger := createfile.New(logFile)

	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SplashReq
		if err := httpx.Parse(r, &req); err != nil {
			if strings.Contains(fmt.Sprint(err), "is not set") {
				fieldName := global.ExtractFieldName(fmt.Sprint(err))
				httpx.WriteJsonCtx(r.Context(), w, 400, map[string]interface{}{"error": errorglobal.BadRequest(fieldName, fieldName+" is required")})
				return
			}

			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		jsonReq, err := json.Marshal(req)
		if err != nil {
			appLogger.Error(logFile).Printf("[x] Request:  %+v -> Error json.Marshal: %v", req, err) // Log
		}

		appLogger.Info(logFile).Printf("Incoming Request: %s", jsonReq)

		info := map[string]interface{}{
			"IsFirst": req.IsFirst,
		}

		l := splash.NewSplashLogic(r.Context(), svcCtx)
		resp, err := l.GetSplashes(info)
		if err != nil {
			appLogger.Error(logFile).Printf("[x] Request:  %+v -> Error: %v", req, err)
			// httpx.ErrorCtx(r.Context(), w, err)
			errorglobal.ReturnError(err, r, w, info)

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
