package middleware

import (
	"context"
	"fmt"
	"net/http"
	"nrs_customer_module_backend/internal/global"
	otpFunc "nrs_customer_module_backend/internal/global/otp"

	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/oauth"
	"strings"
)

type BeforeLoginTokenCheckingMiddleware struct {
}

func NewBeforeLoginTokenCheckingMiddleware() *BeforeLoginTokenCheckingMiddleware {
	return &BeforeLoginTokenCheckingMiddleware{}
}

func (m *BeforeLoginTokenCheckingMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
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

		if !isValidToken(token, tenantId) {
			http.Error(w, "Unauthorized: Invalid or expired token", http.StatusUnauthorized)
			return
		}
		// Retrieve OAuth information for the valid token
		oauthData, err := getOAuthData(token, tenantId)
		if err != nil {
			http.Error(w, "Unauthorized: Error retrieving OAuth data", http.StatusUnauthorized)
			return
		}

		// Get the expiration time of the token
		expirationTime, err := otpFunc.GetTokenExpirationTime(oauthData.Created.Format(global.DefaultTimeFormat), oauthData.ExpiresIn)
		if err != nil {
			http.Error(w, "Unauthorized: Error getting token expiration time", http.StatusUnauthorized)
			return
		}

		// Store the OAuth data in the request context
		ctx := context.WithValue(r.Context(), "oauthData", oauthData)
		ctx = context.WithValue(ctx, "expirationTime", expirationTime)

		// Proceed to the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func getOAuthData(token string, tenantId int64) (*oauth.Oauth, error) {
	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	oauthModel := oauth.NewOauthModel(conn)
	oauthData, err := oauthModel.CheckBearerTokenBeforeLogin(context.Background(), token, tenantId)
	if err != nil {
		return nil, err
	}

	return oauthData, nil
}

func isValidToken(token string, tenantId int64) bool {
	conn, err := model.InitializeDatabase()
	if err != nil {
		fmt.Println("Database initialization error:", err)
		return false
	}

	oauthModel := oauth.NewOauthModel(conn)
	info, err := oauthModel.CheckBearerTokenBeforeLogin(context.Background(), token, tenantId)
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

	expired := otpFunc.ExpiredChecking(info.Created.Format(global.DefaultTimeFormat), info.ExpiresIn, *currentTime)
	if expired {
		fmt.Println("Token has expired")
		return false
	}

	validToken := token // Replace with the expected valid token value
	return info.AccessToken == validToken
}

func CheckTenantId(path string) (tenantId int64) {
	if strings.Contains(path, "/portal/") {
		tenantId = 2
	} else if strings.Contains(path, "/app/") {
		tenantId = 1
	}
	return tenantId
}
