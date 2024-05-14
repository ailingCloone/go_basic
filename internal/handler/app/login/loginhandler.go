package login

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"nrs_customer_module_backend/internal/global"
	"nrs_customer_module_backend/internal/global/createfile"
	"nrs_customer_module_backend/internal/global/errorglobal"
	httperror "nrs_customer_module_backend/internal/global/http_error"
	validatorFunc "nrs_customer_module_backend/internal/global/validator"

	"nrs_customer_module_backend/internal/logic/app/login"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/go-playground/validator/v10"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	var otherInfo map[string]interface{} = map[string]interface{}{} // TODO: When have record in db, please assign it

	logFile := "login"
	appLogger := createfile.New(logFile)

	v := validator.New(validator.WithRequiredStructEnabled())

	v.RegisterValidation("typeValidator", validatorFunc.TypeValidator)
	v.RegisterValidation("fromValidator", validatorFunc.FromValidator)
	v.RegisterValidation("pageValidator", validatorFunc.PageValidator)

	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PostLoginAppReq
		if err := httpx.Parse(r, &req); err != nil {
			if strings.Contains(fmt.Sprint(err), "is not set") {
				fieldName := global.ExtractFieldName(fmt.Sprint(err))
				httpx.WriteJsonCtx(r.Context(), w, 400, map[string]interface{}{"error": errorglobal.BadRequest(fieldName, fieldName+" is required")})
				return
			}
			httpx.ErrorCtx(r.Context(), w, err)

			return
		}

		if err := v.Struct(req); err != nil {
			// Validation failed
			fmt.Printf("Validation error: %s\n", err.Error())
			cusError := err.(validator.ValidationErrors)
			httperror.ResponseErrorWithValidationErrors(w, cusError)
			return

		}

		if req.Type != "EMAIL" && req.TxId == "" {
			fieldName := "txid"
			httpx.WriteJsonCtx(r.Context(), w, 400, map[string]interface{}{"error": errorglobal.BadRequest(fieldName, fieldName+" is required")})
			return
		}

		jsonReq, err := json.Marshal(req)
		if err != nil {
			appLogger.Error(logFile).Printf("[x] Request:  %+v -> Error json.Marshal: %v", req, err) // Log
		}

		appLogger.Info(logFile).Printf("Incoming Request: %s", jsonReq)

		l := login.NewLoginLogic(r.Context(), svcCtx, otherInfo)
		resp, err := l.Login(&req)
		if err != nil {
			// Log the request
			appLogger.Error(logFile).Printf("[x] Request:  %+v -> Error: %v", req, err)
			data := map[string]interface{}{
				"code":    req.Code,
				"txid":    req.TxId,
				"contact": req.Contact,
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
