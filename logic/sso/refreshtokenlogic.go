package sso

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"nrs_customer_module_backend/internal/config"
	"nrs_customer_module_backend/internal/global"
	"nrs_customer_module_backend/internal/global/activitylog"
	"nrs_customer_module_backend/internal/global/responseglobal"
	"nrs_customer_module_backend/internal/global/tokenglobal"
	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/activity_log"
	"nrs_customer_module_backend/internal/model/oauth"

	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/core/conf"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	otherInfo map[string]interface{}
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext, otherInfo map[string]interface{}) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		otherInfo: otherInfo,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.PostRefreshTokenReq) (resp *types.SuccessResponse, err error) {
	otherInfo := l.otherInfo
	refreshToken, ok := otherInfo["refresh_token"].(string)
	if !ok {
		// Handle the case where the conversion fails
		// For example, return an error or set a default value
		return nil, errors.New("unable to convert refresh_token to string")
	}

	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}
	oauthModel := oauth.NewOauthModel(conn)
	currentTime, err := global.TimeInSingapore()
	if err != nil {
		return nil, err
	}

	// Get token
	tokenData, err := oauthModel.FindOneByRefreshToken(l.ctx, refreshToken)

	if err != nil {
		if err == sql.ErrNoRows {
			return responseglobal.GenerateResponseBody(false, "Refresh Token is not found.", map[string]interface{}{}), nil
		}

		return nil, err
	}

	// Generate new access token
	var c config.Config
	configFile := "etc/api.yaml"
	if err := conf.LoadConfig(configFile, &c); err != nil {
		return nil, err
	}

	tokenValiditySeconds := int64(c.TokenValiditySecondsAuthorized)
	if tokenData.CustomerId.Int64 > 0 || tokenData.StaffId.Int64 > 0 {
		tokenValiditySeconds = int64(c.TokenValiditySecondsLogin)
	}

	tokenResp, err := tokenglobal.RefreshAccessToken(tokenValiditySeconds, refreshToken)
	if err != nil {
		return nil, err
	}

	// Prepare data for update
	previousData := *tokenData

	tokenData.RefreshToken = tokenResp.RefreshToken
	tokenData.AccessToken = tokenResp.AccessToken
	tokenData.ExpiresIn = tokenResp.ExpiresIn
	tokenData.Updated = *currentTime
	data := tokenData

	// Update data in the database
	if err := oauthModel.Update(context.Background(), data); err != nil {
		return nil, err
	}

	// Log changes
	changes, err := EditConvertChangesActivityLog(data, &previousData)
	if err != nil {
		return nil, err
	}

	dataActivityLog := &activity_log.ActivityLog{
		ReferTable: "oauth",
		Action:     "Edit",
		Changes:    string(changes),
		Created:    fmt.Sprint(*currentTime),
	}

	err = activitylog.AddActivityLog(dataActivityLog, conn)

	if err != nil {
		return nil, err
	}

	return responseglobal.GenerateResponseBody(true, "Successfully authorized.", tokenResp), nil
}

func EditConvertChangesActivityLog(data *oauth.Oauth, previousData *oauth.Oauth) (jsonData []byte, err error) {

	// Create a map to hold the changes for each field
	changes := map[string]activity_log.Change{
		"customer_id": {
			Before: fmt.Sprint(previousData.CustomerId), // Before value for CustomerId
			After:  fmt.Sprint(data.CustomerId),         // After value for CustomerId
		},
		"staff_id": {
			Before: fmt.Sprint(previousData.StaffId), // Before value for StaffId
			After:  fmt.Sprint(data.StaffId),         // After value for StaffId
		},
		"tenant_id": {
			Before: fmt.Sprint(previousData.TenantId), // Before value for TenantId
			After:  fmt.Sprint(data.TenantId),         // After value for TenantId
		},
		"login_by": {
			Before: fmt.Sprint(previousData.LoginBy), // Before value for LoginBy
			After:  fmt.Sprint(data.LoginBy),         // After value for LoginBy
		},
		"access_token": {
			Before: previousData.AccessToken, // Before value for AccessToken
			After:  data.AccessToken,         // After value for AccessToken
		},
		"refresh_token": {
			Before: previousData.RefreshToken, // Before value for RefreshToken
			After:  data.RefreshToken,         // After value for RefreshToken
		},
		"expires_in": {
			Before: fmt.Sprint(previousData.ExpiresIn), // Before value for ExpiresIn
			After:  fmt.Sprint(data.ExpiresIn),         // After value for ExpiresIn
		},
		"scope": {
			Before: fmt.Sprint(previousData.Scope), // Before value for Scope
			After:  fmt.Sprint(data.Scope),         // After value for Scope
		},
		"token_type": {
			Before: previousData.TokenType, // Before value for TokenType
			After:  data.TokenType,         // After value for TokenType
		},
		"created": {
			Before: fmt.Sprint(previousData.Created), // Before value for Created
			After:  fmt.Sprint(data.Created),         // After value for Created
		},
		"updated": {
			Before: fmt.Sprint(previousData.Updated), // Before value for Created
			After:  fmt.Sprint(data.Updated),         // After value for Created
		},
		"active": {
			Before: fmt.Sprint(previousData.Active), // Before value for Active
			After:  fmt.Sprint(data.Active),         // After value for Active
		},
	}

	// Marshal the ChangeData object into JSON format
	jsonData, err = json.Marshal(changes)
	if err != nil {
		return nil, err
	}

	return jsonData, nil

}
