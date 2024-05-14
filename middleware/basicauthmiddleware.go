package middleware

import (
	"context"
	"fmt"
	"net/http"
	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/tenant"
	"strings"
)

type BasicAuthMiddleware struct {
}

func NewBasicAuthMiddleware() *BasicAuthMiddleware {
	return &BasicAuthMiddleware{}
}

func (m *BasicAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			unauthorized(w)
			return
		}
		appPrefix := "App"
		portalPrefix := "Portal"
		sharedPrefix := "Shared"
		module := 0
		if strings.HasPrefix(user, appPrefix) || strings.HasPrefix(user, portalPrefix) || strings.HasPrefix(user, sharedPrefix) {
			module = 1
		}
		// Create an instance of TncModel
		conn, err := model.InitializeDatabase()
		if err != nil {
			unauthorized(w)
			return
		}
		tenantModel := tenant.NewTenantModel(conn)
		// Get tenant data based on client_id, client_secret and module
		tenantResult, err := tenantModel.FindTenantWithModule(context.Background(), user, pass, module)
		if err != nil || tenantResult.Id == 0 {
			unauthorized(w)
			return
		}
		fmt.Println("tenantIDKey: ", tenantResult.Id)
		fmt.Println("roleIDKey: ", tenantResult.RoleId)
		ctx := context.WithValue(r.Context(), "tenant_id", tenantResult.Id)
		ctx = context.WithValue(ctx, "role_id", tenantResult.RoleId)
		// Call the next handler with the updated context
		next(w, r.WithContext(ctx))
	}
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Unauthorized\n"))
}
