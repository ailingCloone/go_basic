package otp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/zeromicro/go-zero/rest/httpx"

	"nrs_customer_module_backend/internal/global"
	"nrs_customer_module_backend/internal/global/createfile"
	"nrs_customer_module_backend/internal/global/errorglobal"
	httperror "nrs_customer_module_backend/internal/global/http_error"
	validatorFunc "nrs_customer_module_backend/internal/global/validator"
	"nrs_customer_module_backend/internal/logic/app/otp"
	"nrs_customer_module_backend/internal/middleware"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"
)

var (
	logFile   = "otp_request"
	appLogger = createfile.New(logFile)
	v         = validator.New(validator.WithRequiredStructEnabled())

	beforeLoginMiddleware = middleware.NewBeforeLoginTokenCheckingMiddleware()
	afterLoginMiddleware  = middleware.NewAfterLoginTokenCheckingMiddleware()
)

func OtpHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	v.RegisterValidation("typeValidator", validatorFunc.TypeValidator)
	v.RegisterValidation("fromValidator", validatorFunc.FromValidator)
	v.RegisterValidation("pageValidator", validatorFunc.PageValidator)

	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OtpReq

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

		jsonReq, err := json.Marshal(req)
		if err != nil {
			appLogger.Error(logFile).Printf("[x] Request:  %+v -> Error json.Marshal: %v", req, err) // Log
			return
		}

		appLogger.Info(logFile).Printf("Incoming Request: %s", jsonReq)

		// Determine which middleware to apply based on 'Page' parameter
		var selectedMiddleware http.HandlerFunc
		switch req.Page {
		case "login", "register":
			selectedMiddleware = beforeLoginMiddleware.Handle(mainHandler(svcCtx, req, jsonReq))
		case "profile_update":
			selectedMiddleware = afterLoginMiddleware.Handle(mainHandler(svcCtx, req, jsonReq))
		default:
			http.Error(w, "Invalid 'Page' value", http.StatusBadRequest)
			return
		}

		// Apply the selected middleware to the request
		selectedMiddleware.ServeHTTP(w, r)
	}
}

func mainHandler(svcCtx *svc.ServiceContext, req types.OtpReq, jsonReq []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := otp.NewOtpLogic(r.Context(), svcCtx)
		resp, err := l.Otp(&req)
		if err != nil {
			// Log the request
			appLogger.Error(logFile).Printf("[x] Request:  %+v -> Error: %v", req, err)
			data := map[string]interface{}{
				"contact": req.Contact,
			}
			info := errorglobal.ErrorChecking(err, data)
			fmt.Println("info from handler", info)
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
