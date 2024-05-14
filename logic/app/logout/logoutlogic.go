package logout

import (
	"context"
	"encoding/json"
	"fmt"

	"nrs_customer_module_backend/internal/global"
	"nrs_customer_module_backend/internal/global/activitylog"
	"nrs_customer_module_backend/internal/global/responseglobal"
	"nrs_customer_module_backend/internal/global/user"
	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/activity_log"
	"nrs_customer_module_backend/internal/model/oauth"
	"nrs_customer_module_backend/internal/model/ui"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	otherInfo map[string]interface{}
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext, otherInfo map[string]interface{}) *LogoutLogic {
	return &LogoutLogic{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		otherInfo: otherInfo,
	}
}

func (l *LogoutLogic) Logout(req *types.PostLogoutReq) (resp *types.SuccessResponse, err error) {
	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	oauthData, err := user.GetOauthDataFromAfterLoginMiddleware(l.ctx)
	if err != nil {
		return nil, err
	}

	currentTime, err := global.TimeInSingapore()
	if err != nil {
		return nil, err
	}

	previousData := *oauthData

	oauthData.Active = 0
	oauthData.Updated = *currentTime
	oauthModel := oauth.NewOauthModel(conn)
	if err := oauthModel.UpdateById(l.ctx, oauthData); err != nil {
		return responseglobal.GenerateResponseBody(false, "Failed to logout", map[string]interface{}{}), nil
	}
	data := oauthData

	// Construct the dialog response
	dialog := ui.Dialog{
		Content: ui.DialogContent{
			ImageUrl:    "http://aaa",
			WebImageUrl: "http://aaa",
			Title:       "Success logout.",
			Subtitle:    "Success logout.",
		},
		Button: []ui.Button{
			{
				Text:   "OK",
				Action: "/",
			},
		},
	}
	response := map[string]interface{}{
		"dialog": dialog,
	}

	// Add activity log
	changes, err := ConvertChangesActivityLog(data, &previousData)
	if err != nil {
		return nil, err
	}

	dataActivityLog := &activity_log.ActivityLog{
		CustomerId: oauthData.CustomerId.Int64,
		StaffId:    oauthData.StaffId.Int64,
		ReferTable: "oauth",
		Action:     "logout",
		Changes:    string(changes),
		Created:    fmt.Sprint(*currentTime),
	}

	err = activitylog.AddActivityLog(dataActivityLog, conn)
	if err != nil {
		return nil, err
	}

	return responseglobal.GenerateResponseBody(true, "Successfully logout", response), nil
}

func ConvertChangesActivityLog(data *oauth.Oauth, previousData *oauth.Oauth) (jsonData []byte, err error) {

	// Create a map to hold the changes for each field
	changes := map[string]activity_log.Change{
		"id": {
			Before: fmt.Sprint(previousData.Id), // Before value for Guid
			After:  fmt.Sprint(data.Id),         // After value for Guid
		},
		"customer_id": {
			Before: fmt.Sprint(previousData.CustomerId), // Before value for Guid
			After:  fmt.Sprint(data.CustomerId),         // After value for Guid
		},
		"staff_id": {
			Before: fmt.Sprint(previousData.StaffId), // Before value for Guid
			After:  fmt.Sprint(data.StaffId),         // After value for Guid
		},
		"tenant_id": {
			Before: fmt.Sprint(previousData.TenantId), // Before value for Guid
			After:  fmt.Sprint(data.TenantId),         // After value for Guid
		},
		"access_token": {
			Before: previousData.AccessToken, // Before value for Created
			After:  data.AccessToken,         // After value for Created
		},
		"refresh_token": {
			Before: previousData.RefreshToken, // Before value for Created
			After:  data.RefreshToken,         // After value for Created
		},
		"expire_in": {
			Before: fmt.Sprint(previousData.ExpiresIn), // Before value for Created
			After:  fmt.Sprint(data.ExpiresIn),         // After value for Created
		},
		"scope": {
			Before: fmt.Sprint(previousData.Scope), // Before value for Created
			After:  fmt.Sprint(data.Scope),         // After value for Created
		},
		"token_type": {
			Before: previousData.TokenType, // Before value for Created
			After:  data.TokenType,         // After value for Created
		},
		"updated": {
			Before: fmt.Sprint(previousData.Updated), // Before value for Created
			After:  fmt.Sprint(data.Updated),         // After value for Created
		},
		"created": {
			Before: fmt.Sprint(previousData.Created), // Before value for Created
			After:  fmt.Sprint(data.Created),         // After value for Created
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
