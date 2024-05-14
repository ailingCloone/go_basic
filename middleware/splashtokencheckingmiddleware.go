package middleware

import (
	"context"
	"fmt"
	"net/http"
	"nrs_customer_module_backend/internal/global"
	otpFunc "nrs_customer_module_backend/internal/global/otp"
	"time"

	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/oauth"
	"strings"
)

type SplashTokenCheckingMiddleware struct {
}

func NewSplashTokenCheckingMiddleware() *SplashTokenCheckingMiddleware {
	return &SplashTokenCheckingMiddleware{}
}

func (m *SplashTokenCheckingMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
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

		if !isValidSplashToken(token, tenantId) {
			http.Error(w, "Unauthorized: Invalid or expired token", http.StatusUnauthorized)
			return
		}
		// Retrieve OAuth information for the valid token
		oauthData, err := getOAuthDataAfterLogin(token, tenantId)
		if err != nil {
			http.Error(w, "Unauthorized: Error retrieving OAuth data", http.StatusUnauthorized)
			return
		}

		// Store the OAuth data in the request context
		ctx := context.WithValue(r.Context(), "oauthData", oauthData)

		// Proceed to the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func isValidSplashToken(token string, tenantId int64) bool {
	conn, err := model.InitializeDatabase()
	if err != nil {
		fmt.Println("Database initialization error:", err)
		return false
	}

	oauthModel := oauth.NewOauthModel(conn)
	info, err := oauthModel.CheckBearerTokenAfterLogin(context.Background(), token, tenantId)
	if err != nil {
		fmt.Println("Error checking token:", err)
		return false
	}

	if info == nil {
		fmt.Println("Token info is nil")
		return false
	}
	fmt.Println("info from middleware::", info)
	currentTime, err := global.TimeInSingapore()
	if err != nil {
		fmt.Println("Error getting current time:", err)
		return false
	}

	var compareTime time.Time
	if fmt.Sprint(info.Updated) != "0001-01-01 00:00:00 +0000 UTC" {
		compareTime = info.Updated
	} else {
		compareTime = info.Created
	}

	expired := otpFunc.ExpiredChecking(compareTime.Format(global.DefaultTimeFormat), info.ExpiresIn, *currentTime)
	if expired {
		fmt.Println("Token has expired")
		return false
	}

	validToken := token // Replace with the expected valid token value
	return info.AccessToken == validToken
}
