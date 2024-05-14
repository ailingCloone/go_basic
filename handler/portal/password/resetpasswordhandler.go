package password

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"nrs_customer_module_backend/internal/global"
	"nrs_customer_module_backend/internal/global/createfile"
	"nrs_customer_module_backend/internal/global/errorglobal"
	"nrs_customer_module_backend/internal/logic/portal/password"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ResetPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	var otherInfo map[string]interface{} = map[string]interface{}{} // TODO: When have record in db, please assign it

	logFile := "reset_password"
	appLogger := createfile.New(logFile)

	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PostPortalResetPasswordReq
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

		l := password.NewResetPasswordLogic(r.Context(), svcCtx, otherInfo)
		resp, err := l.ResetPassword(&req)
		if err != nil {
			// Log the request
			appLogger.Error(logFile).Printf("[x] Request:  %+v -> Error: %v", req, err)
			data := map[string]interface{}{
				"code": req.Code,
				"txid": req.TxId,
			}
			info := errorglobal.ErrorChecking(err, data)
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
