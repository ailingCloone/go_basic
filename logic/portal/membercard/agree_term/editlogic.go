package agree_term

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"nrs_customer_module_backend/internal/global"
	"nrs_customer_module_backend/internal/global/activitylog"
	"nrs_customer_module_backend/internal/global/responseglobal"
	"nrs_customer_module_backend/internal/global/user"
	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/activity_log"
	"nrs_customer_module_backend/internal/model/agree_term"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	otherInfo map[string]interface{}
}

func NewEditLogic(ctx context.Context, svcCtx *svc.ServiceContext, otherInfo map[string]interface{}) *EditLogic {
	return &EditLogic{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		otherInfo: otherInfo,
	}
}

func (l *EditLogic) Edit(req *types.PostATEditReq) (resp *types.SuccessResponse, err error) {
	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	oauthData, err := user.GetOauthDataFromAfterLoginMiddleware(l.ctx)
	if err != nil {
		return nil, err
	}

	agreeTermModel := agree_term.NewAgreeTermModel(conn)
	agreeTerm, err := agreeTermModel.FindOneGuid(l.ctx, req.Guid)
	if err != nil {
		if err == sql.ErrNoRows {
			return responseglobal.GenerateResponseBody(false, "Agree Term is not found.", map[string]interface{}{}), nil
		}

		return nil, err
	}

	currentTime, err := global.TimeInSingapore()
	if err != nil {
		return nil, err
	}

	// Prepare data for insertion
	previousData := *agreeTerm

	agreeTerm.Title = req.Title
	agreeTerm.Description = req.Content
	agreeTerm.Updated = *currentTime
	agreeTerm.Active = 1
	data := agreeTerm

	dataResp := map[string]interface{}{
		"title":   req.Title,
		"content": req.Content,
	}

	if err := agreeTermModel.UpdateByGuid(l.ctx, data); err != nil {
		return responseglobal.GenerateResponseBody(false, "Failed to update Agree Term.", dataResp), nil
	}

	changes, err := EditConvertChangesActivityLog(data, &previousData)
	if err != nil {
		return nil, err
	}

	dataActivityLog := &activity_log.ActivityLog{
		CustomerId: oauthData.CustomerId.Int64,
		StaffId:    oauthData.StaffId.Int64,
		ReferTable: "agree_term",
		Action:     "edit",
		Changes:    string(changes),
		Created:    fmt.Sprint(*currentTime),
	}

	err = activitylog.AddActivityLog(dataActivityLog, conn)

	if err != nil {
		return nil, err
	}

	return responseglobal.GenerateResponseBody(true, "AgreeTerm updated successfully.", dataResp), nil
}

func EditConvertChangesActivityLog(data *agree_term.AgreeTerm, previousData *agree_term.AgreeTerm) (jsonData []byte, err error) {

	// Create a map to hold the changes for each field
	changes := map[string]activity_log.Change{
		"guid": {
			Before: previousData.Guid, // Before value for Guid
			After:  data.Guid,         // After value for Guid
		},
		"refer_table": {
			Before: previousData.ReferTable, // Before value for ReferTable
			After:  data.ReferTable,         // After value for ReferTable
		},
		"refer_id": {
			Before: fmt.Sprint(previousData.ReferId), // Before value for ReferId
			After:  fmt.Sprint(data.ReferId),         // After value for ReferId
		},
		"title": {
			Before: previousData.Title, // Before value for Title
			After:  data.Title,         // After value for Title
		},
		"description": {
			Before: previousData.Description, // Before value for Description
			After:  data.Description,         // After value for Description
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
