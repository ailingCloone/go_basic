package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"nrs_customer_module_backend/internal/global"
	"nrs_customer_module_backend/internal/global/errorglobal"
	otpFunc "nrs_customer_module_backend/internal/global/otp"

	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/oauth"
	"nrs_customer_module_backend/internal/model/user_access_setting"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type AfterLoginTokenCheckingMiddleware struct {
}

func NewAfterLoginTokenCheckingMiddleware() *AfterLoginTokenCheckingMiddleware {
	return &AfterLoginTokenCheckingMiddleware{}
}

func (m *AfterLoginTokenCheckingMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		tenantId := CheckTenantId(path)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: Missing Authorization header", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		if !isValidTokenAfterLogin(token, tenantId) {
			http.Error(w, "Unauthorized: Invalid or expired token", http.StatusUnauthorized)
			return
		}
		// Retrieve OAuth information for the valid token
		oauthData, err := getOAuthDataAfterLogin(token, tenantId)
		fmt.Println("oauthData in middleware", oauthData)
		if err != nil {
			http.Error(w, "Unauthorized: Error retrieving OAuth data", http.StatusUnauthorized)
			return
		}

		// check api only access by staff
		err = APIOnlyAccessByStaff(oauthData, r.URL.Path)
		if err != nil {
			httpx.WriteJsonCtx(r.Context(), w, 400, map[string]interface{}{"error": errorglobal.Forbidden})
			return
		}

		// var ctx context.Context
		if oauthData.StaffId.Int64 != 0 {
			fmt.Println("staff in middle", oauthData.StaffId.Int64)

			path := r.URL.Path
			referTable, otherInfo, method := determineTableInfo(path)

			if referTable != "" {
				fmt.Println("sidemenu refer table:::", referTable)
				userPermission, err := getUserPermission(oauthData.StaffId.Int64, referTable, otherInfo, method)
				fmt.Println("user permission from middleware", userPermission)
				if err != nil {
					httpx.WriteJsonCtx(r.Context(), w, 400, map[string]interface{}{"error": errorglobal.Forbidden})
					return
				}
				r = r.WithContext(context.WithValue(r.Context(), "userPermission", userPermission))

			}

		}

		// Store the OAuth data in the request context
		r = r.WithContext(context.WithValue(r.Context(), "oauthData", oauthData))
		// Proceed to the next handler with the updated context
		next.ServeHTTP(w, r)
	}
}

func getOAuthDataAfterLogin(token string, tenantId int64) (*oauth.Oauth, error) {
	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	oauthModel := oauth.NewOauthModel(conn)
	oauthData, err := oauthModel.CheckBearerTokenAfterLogin(context.Background(), token, tenantId)
	if err != nil {
		return nil, err
	}

	return oauthData, nil
}

func isValidTokenAfterLogin(token string, tenantId int64) bool {
	conn, err := model.InitializeDatabase()
	if err != nil {
		fmt.Println("Database initialization error:", err)
		return false
	}

	oauthModel := oauth.NewOauthModel(conn)
	info, err := oauthModel.CheckBearerTokenAfterLogin(context.Background(), token, tenantId)
	fmt.Println("token in after login", token)
	fmt.Println("info", info)
	if err != nil {
		fmt.Println("Error checking token:", err)
		return false
	}

	if info == nil {
		fmt.Println("Token info is nil")
		return false
	}

	currentTime, err := global.TimeInSingapore()
	if err != nil {
		fmt.Println("Error getting current time:", err)
		return false
	}

	expired := otpFunc.ExpiredChecking(info.Updated.Format(global.DefaultTimeFormat), info.ExpiresIn, *currentTime)
	if expired {
		fmt.Println("Token has expired")
		return false
	}

	validToken := token // Replace with the expected valid token value
	return info.AccessToken == validToken
}

func determineTableInfo(path string) (string, string, string) {
	var referTable, otherInfo, method string

	switch {

	case strings.Contains(path, "/sidemenu"):
		referTable = "sidemenu"
		otherInfo = ""
		method = ""

	case strings.Contains(path, "/membercard"):
		referTable = "sub_module"
		if strings.Contains(path, "/terms_n_condition") {
			otherInfo = "membercard_terms_n_condition"
			method = extractMethodFromPath(path)
		}
		if strings.Contains(path, "/application_form") {
			otherInfo = "membercard_application_form"
			method = extractMethodFromPath(path)
		}
		if strings.Contains(path, "/agree_term") {
			otherInfo = "membercard_agree_term"
			method = extractMethodFromPath(path)
		}
		if strings.Contains(path, "/sms_template") {
			otherInfo = "membercard_sms_template"
			method = extractMethodFromPath(path)
		}
		if strings.Contains(path, "/get_list") {
			otherInfo = "membercard"
		}
		if strings.Contains(path, "/setting") {
			otherInfo = "membercard_setting"
			method = extractMethodFromPath(path)
		}
	default:
		referTable, otherInfo, method = "", "", ""
	}
	return referTable, otherInfo, method
}

func getUserPermission(staffId int64, referTable, otherInfo, method string) (userPermission *[]user_access_setting.UserAccessSettings, err error) {
	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	fmt.Println("staffId:", staffId, "\nreferTable:", referTable, "\notherInfo:", otherInfo, "\nmethod", method)

	userAccessSettingModel := user_access_setting.NewUserAccessSettingModel(conn)
	if otherInfo == "" {
		userPermissionSetting, err := userAccessSettingModel.FindAll(context.Background(), staffId, referTable)
		fmt.Println("setting", userPermissionSetting)
		if err != nil {
			return nil, err
		}
		userPermission = convertUserAccessSettingToSettings(userPermissionSetting)

		fmt.Println("user permission in checking is ", userPermission)

	} else {
		userPermission, err = userAccessSettingModel.FindAllPermission(context.Background(), staffId, referTable, otherInfo)
		fmt.Println("user permission in find all permission checking is ", userPermission)

		if err != nil {
			return nil, err
		}
		for _, permission := range *userPermission {
			fmt.Println("matched id", permission)
			for _, permissions := range *userPermission {
				if permissions.AllowView == 0 { //check for get api
					fmt.Println("permissions in allow view checking:", permissions)
					fmt.Println("allow view", permissions.AllowView)
					userPermission = nil
				}
			}
		}
	}

	fmt.Println("user permissions is ", userPermission)

	if userPermission == nil || len(*userPermission) == 0 {
		fmt.Println("user permission is nil or empty")
		return nil, errors.New("user permission is nil or empty")
	}

	if method != "" {
		for _, checkPermission := range *userPermission {
			if method == "add" && checkPermission.AllowAdd == 0 {
				return nil, errors.New("user permission is nil or empty")
			}
			if method == "edit" && checkPermission.AllowEdit == 0 {
				return nil, errors.New("user permission is nil or empty")
			}
			if method == "delete" && checkPermission.AllowDelete == 0 {
				return nil, errors.New("user permission is nil or empty")
			}
		}
	}

	return userPermission, nil
}

func extractMethodFromPath(path string) string {
	if strings.Contains(path, "/add") {
		return "add"
	} else if strings.Contains(path, "/edit") {
		return "edit"
	} else if strings.Contains(path, "/delete") {
		return "delete"
	}
	return "" // Default method if none of the patterns match
}

func convertUserAccessSettingToSettings(settings *[]user_access_setting.UserAccessSetting) *[]user_access_setting.UserAccessSettings {
	if settings == nil {
		return nil
	}

	converted := make([]user_access_setting.UserAccessSettings, len(*settings))
	for i, setting := range *settings {
		converted[i] = user_access_setting.UserAccessSettings{
			Id:                  setting.Id,
			StaffId:             setting.StaffId,
			MainModuleId:        setting.MainModuleId,
			SubModuleReferTable: setting.SubModuleReferTable,
			SubModuleId:         setting.SubModuleId,
			AllowAdd:            setting.AllowAdd,
			AllowEdit:           setting.AllowEdit,
			AllowDelete:         setting.AllowDelete,
			AllowView:           setting.AllowView,
			Updated:             setting.Updated,
			Created:             setting.Created,
			Active:              setting.Active,
			Title:               "", // Title is not available in UserAccessSetting
		}
	}

	return &converted
}

func APIOnlyAccessByStaff(oauthData *oauth.Oauth, path string) error {
	onlyStaffAccess := strings.Contains(path, "/app/registration_list/get_list")

	if oauthData.CustomerId.Int64 > 0 && onlyStaffAccess {
		return errors.New("user permission is nil or empty")
	}

	return nil

}
