package responseglobal

import "nrs_customer_module_backend/internal/types"

func GenerateResponseBody(success bool, message string, data interface{}) *types.SuccessResponse {
	resp := &types.SuccessResponse{
		Success: success,
		Message: message,
		Data:    data,
	}

	return resp
}
