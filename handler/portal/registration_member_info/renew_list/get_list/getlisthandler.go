package get_list

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
	"nrs_customer_module_backend/internal/logic/portal/registration_member_info/renew_list/get_list"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/go-playground/validator/v10"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {

	logFile := "registration_member_info_renew_list_get_list"
	appLogger := createfile.New(logFile)
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterValidation("registerListStatusValidator", validatorFunc.RegisterListStatusValidator)
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetListRegistrationReq
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
			cusError := err.(validator.ValidationErrors)
			httperror.ResponseErrorWithValidationErrors(w, cusError)
			return
		}
		jsonReq, err := json.Marshal(req)
		if err != nil {
			appLogger.Error(logFile).Printf("[x] Request:  %+v -> Error json.Marshal: %v", req, err) // Log
		}
		appLogger.Info(logFile).Printf("Incoming Request: %s", jsonReq)
		l := get_list.NewGetListLogic(r.Context(), svcCtx)
		resp, err := l.GetList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}

}
