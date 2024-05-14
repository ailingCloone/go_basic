package terms_n_condition

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"nrs_customer_module_backend/internal/global"
	"nrs_customer_module_backend/internal/global/createfile"
	"nrs_customer_module_backend/internal/global/errorglobal"
	"nrs_customer_module_backend/internal/logic/portal/membercard/terms_n_condition"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func EditHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	var otherInfo map[string]interface{} = map[string]interface{}{} // TODO: When have record in db, please assign it

	logFile := "edit_tnc"
	appLogger := createfile.New(logFile)

	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PostTNCEditReq
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

		l := terms_n_condition.NewEditLogic(r.Context(), svcCtx, otherInfo)
		resp, err := l.Edit(&req)
		if err != nil {
			// Log the request
			appLogger.Error(logFile).Printf("[x] Request:  %+v -> Error: %v", req, err)
			info := map[string]interface{}{
				"source_name": "guid",
				"source_id":   req.Guid,
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
