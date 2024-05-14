package sso

import (
	"context"
	"database/sql"
	"errors"

	"nrs_customer_module_backend/internal/global"
	"nrs_customer_module_backend/internal/global/responseglobal"
	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/oauth"

	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthorizedLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	otherInfo map[string]interface{}
}

func NewAuthorizedLogic(ctx context.Context, svcCtx *svc.ServiceContext, otherInfo map[string]interface{}) *AuthorizedLogic {
	return &AuthorizedLogic{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		otherInfo: otherInfo,
	}
}

func (l *AuthorizedLogic) Authorized(req *types.TokenResponse) (resp *types.SuccessResponse, err error) {
	otherInfo := l.otherInfo
	tenantId, ok := otherInfo["tenant_id"].(int64)
	if !ok {
		// Handle the case where the conversion fails
		// For example, return an error or set a default value
		return nil, errors.New("unable to convert tenant_id to int64")
	}
	roleId, ok := otherInfo["role_id"].(int64)
	if !ok {
		// Handle the case where the conversion fails
		// For example, return an error or set a default value
		return nil, errors.New("unable to convert role_id to int64")
	}
	// todo: add your logic here and delete this line
	// Create an instance of TncModel
	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}
	oauthModel := oauth.NewOauthModel(conn)
	currentTime, err := global.TimeInSingapore()
	if err != nil {
		return nil, err
	}
	tenantIdNull := sql.NullInt64{Int64: tenantId, Valid: true}
	// Prepare data for insertion
	data := &oauth.Oauth{
		AccessToken:  req.AccessToken,
		RefreshToken: req.RefreshToken,
		TenantId:     tenantIdNull,
		ExpiresIn:    req.ExpiresIn,
		Scope:        roleId,
		TokenType:    req.TokenType,
		Created:      *currentTime,
		Active:       1,
	}

	// Insert data into the database
	result, err := oauthModel.Insert(context.Background(), data)
	// insert record
	if err != nil {
		return nil, err
	}
	// Check the result of the insertion
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return responseglobal.GenerateResponseBody(false, "Failed to add record.", map[string]interface{}{}), nil
	}
	return responseglobal.GenerateResponseBody(true, "Successfullt authorized.", req), nil
}
