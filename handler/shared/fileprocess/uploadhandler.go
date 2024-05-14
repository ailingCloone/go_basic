package fileprocess

import (
	"net/http"

	"nrs_customer_module_backend/internal/global/createfile"
	"nrs_customer_module_backend/internal/global/uploadfile"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var logFile string = "shared_uploadhandler"
		appLogger := createfile.New(logFile)

		req, err := uploadfile.UploadFile(w, r, appLogger, logFile)

		if req.FormValue("status") == "success" {

			data := map[string]interface{}{
				"filename": r.FormValue("filename"),
			}

			resp := &types.SuccessResponse{
				Success: true,
				Message: "Success upload.",
				Data:    data, // Pass any relevant data back to the client
			}
			httpx.OkJsonCtx(r.Context(), w, resp)
			return

		}

		appLogger.Error(logFile).Printf("[x] Error UploadFile: %v", err) // Log

		httpx.ErrorCtx(r.Context(), w, err)

	}
}
