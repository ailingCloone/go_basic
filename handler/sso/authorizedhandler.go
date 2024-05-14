package sso

import (
	"encoding/json"
	"net/http"

	"nrs_customer_module_backend/internal/config"
	"nrs_customer_module_backend/internal/global/createfile"
	"nrs_customer_module_backend/internal/global/errorglobal"
	"nrs_customer_module_backend/internal/global/tokenglobal"
	"nrs_customer_module_backend/internal/logic/sso"
	"nrs_customer_module_backend/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AuthorizedHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	logFile := "authorized"
	appLogger := createfile.New(logFile)

	return func(w http.ResponseWriter, r *http.Request) {

		user, _, _ := r.BasicAuth()
		// Get value from request context
		tenantID := r.Context().Value("tenant_id")
		// Get value from form data
		roleID := r.Context().Value("role_id")
		otherInfo := map[string]interface{}{
			"tenant_id": tenantID,
			"role_id":   roleID,
		}
		var c config.Config
		configFile := "etc/api.yaml"
		if err := conf.LoadConfig(configFile, &c); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		tokenValiditySeconds := int64(c.TokenValiditySecondsAuthorized)
		tokenResp, err := tokenglobal.GenerateToken(tokenValiditySeconds)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := sso.NewAuthorizedLogic(r.Context(), svcCtx, otherInfo)
		resp, err := l.Authorized(tokenResp)
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
			appLogger.Error(logFile).Printf("[x] Request:  user: %s -> Error json.Marshal: %v", user, err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		appLogger.Info(logFile).Printf("[x] Request: user: %s -> Response:  %s", user, respJSON) // Log the r
		httpx.OkJsonCtx(r.Context(), w, resp)

	}
}
