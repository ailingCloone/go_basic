package errorglobal

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
)

var (
	InvalidPassword = "invalid password"
	InvalidOtp      = "invalid otp"
	InvalidTxid     = "invalid txid"
	ExistedUser     = "existed user"
	InvalidContact  = "invalid contact"
)

func BadRequest(field, reason string) (Result map[string]interface{}) {
	return map[string]interface{}{
		"code":    "400",
		"message": "Bad Request - Invalid parameters",
		"details": map[string]interface{}{
			"field":  field,
			"reason": reason,
		},
	}
}

var Unauthorized map[string]interface{} = map[string]interface{}{
	"code":    "401",
	"message": "Unauthorized - Authentication required",
	"details": map[string]interface{}{
		"reason": "Invalid API key or token",
	},
}

var Forbidden map[string]interface{} = map[string]interface{}{
	"code":    "403",
	"message": "Forbidden - Insufficient permissions",
	"details": map[string]interface{}{
		"reason": "User lacks necessary permission",
	},
}

func Notfound(resource, id string, message string) (Result map[string]interface{}) {
	if message == "" {
		message = "Not Found - Resource not found"
	}
	Result = map[string]interface{}{
		"code":    "404",
		"message": message,
		"details": map[string]interface{}{
			"resource": resource,
			"id":       id,
		},
	}

	return Result
}

var InternalServerError map[string]interface{} = map[string]interface{}{
	"code":    "500",
	"message": "Internal Server Error - Something went wrong on the server",
	"details": map[string]interface{}{
		"reason": "Internal server error details",
	},
}

var ConnectionServerError map[string]interface{} = map[string]interface{}{
	"code":    "500",
	"message": "Internal Server Error - Something went wrong on the server connection",
	"details": map[string]interface{}{
		"reason": "Internal server error details",
	},
}

var RequestTimeoutError map[string]interface{} = map[string]interface{}{
	"code":    "503",
	"message": "Internal Server Error - Something went wrong on the server",
	"details": map[string]interface{}{
		"reason": "Request Timeout",
	},
}

func ReturnError(err error, r *http.Request, w http.ResponseWriter, info map[string]interface{}) {

	sourceName := fmt.Sprint(info["source_name"])
	sourceId := fmt.Sprint(info["source_id"])
	message := fmt.Sprint(info["message"])

	switch {
	case err.Error() == "sql: no rows in result set", err.Error() == InvalidPassword, err.Error() == InvalidOtp, err.Error() == InvalidTxid, err.Error() == ExistedUser, err.Error() == InvalidContact:
		httpx.WriteJsonCtx(r.Context(), w, 404, map[string]interface{}{"error": Notfound(sourceName, sourceId, message)})
	case strings.Contains(fmt.Sprint(err), "connection refused"):
		httpx.WriteJsonCtx(r.Context(), w, 500, map[string]interface{}{"error": ConnectionServerError})
	case strings.Contains(fmt.Sprint(err), "Request Timeout"):
		httpx.WriteJsonCtx(r.Context(), w, 503, map[string]interface{}{"error": RequestTimeoutError})
	default:
		httpx.WriteJsonCtx(r.Context(), w, 500, map[string]interface{}{"error": err.Error()})
	}

}

func ErrorChecking(err error, sourceId map[string]interface{}) (info map[string]interface{}) {
	if err.Error() == InvalidPassword {
		info = map[string]interface{}{
			"source_name": "password",
			"source_id":   fmt.Sprint(sourceId["code"]),
			"message":     "The Username/Password incorrect or invalid. Please check again.",
		}

	} else if err.Error() == InvalidOtp {
		info = map[string]interface{}{
			"source_name": "otp",
			"source_id":   fmt.Sprint(sourceId["code"]),
			"message":     "The Username/OTP incorrect or invalid. Please check again.",
		}
	} else if err.Error() == InvalidTxid {
		info = map[string]interface{}{
			"source_name": "txid",
			"source_id":   fmt.Sprint(sourceId["txid"]),
			"message":     "The Username/OTP incorrect or invalid. Please check again.",
		}
	} else if err.Error() == ExistedUser {
		info = map[string]interface{}{
			"source_name": "contact",
			"source_id":   fmt.Sprint(sourceId["contact"]),
			"message":     "User already exists.Please proceed to the login page to access your account.",
		}
	} else if err.Error() == InvalidContact {
		info = map[string]interface{}{
			"source_name": "contact",
			"source_id":   fmt.Sprint(sourceId["contact"]),
			"message":     "The contact number incorrect or invalid. Please check again.",
		}
	} else {
		info = map[string]interface{}{
			"source_name": "contact",
			"source_id":   fmt.Sprint(sourceId["contact"]),
			"message":     "The Username/OTP/Password incorrect or invalid. Please check again.",
		}

	}
	return info
}
