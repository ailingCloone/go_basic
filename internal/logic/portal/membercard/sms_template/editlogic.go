package sms_template

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
	"nrs_customer_module_backend/internal/model/sms_template"
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

func (l *EditLogic) Edit(req *types.PostSTEditReq) (resp *types.SuccessResponse, err error) {
	oauthData, err := user.GetOauthDataFromAfterLoginMiddleware(l.ctx)
	if err != nil {
		return nil, err
	}

	status, err := global.ConvertStringToInt64(req.Status)
	if err != nil {
		return nil, err
	}

	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	smsTemplateModel := sms_template.NewSmsTemplateModel(conn)
	currentTime, err := global.TimeInSingapore()
	if err != nil {
		return nil, err
	}

	smsTemplate, err := smsTemplateModel.FindOneGuid(l.ctx, req.Guid)
	if err != nil {
		if err == sql.ErrNoRows {
			return responseglobal.GenerateResponseBody(false, "SMS Template is not found.", map[string]interface{}{}), nil
		}

		return nil, err
	}

	// Prepare data for insertion
	previousData := *smsTemplate

	smsTemplate.Title = req.Title
	smsTemplate.Description = req.Content
	smsTemplate.Status = status
	smsTemplate.Updated = *currentTime
	smsTemplate.Active = 1
	data := smsTemplate
	dataResp := map[string]interface{}{
		"title":   req.Title,
		"content": req.Content,
		"status":  req.Status,
	}
	if err := smsTemplateModel.UpdateByGuid(l.ctx, data); err != nil {
		return responseglobal.GenerateResponseBody(false, "Failed to update Sms Template.", dataResp), nil
	}

	changes, err := EditConvertChangesActivityLog(data, &previousData)
	if err != nil {
		return nil, err
	}

	dataActivityLog := &activity_log.ActivityLog{
		StaffId:    oauthData.StaffId.Int64,
		ReferTable: "sms_template",
		Action:     "edit",
		Changes:    string(changes),
		Created:    fmt.Sprint(*currentTime),
	}

	err = activitylog.AddActivityLog(dataActivityLog, conn)

	if err != nil {
		return nil, err
	}

	dataResponse := map[string]interface{}{
		"title":   req.Title,
		"content": req.Content,
	}
	return responseglobal.GenerateResponseBody(true, "SMS Template updated successfully.", dataResponse), nil
}

func EditConvertChangesActivityLog(data *sms_template.SmsTemplate, previousData *sms_template.SmsTemplate) (jsonData []byte, err error) {

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
		"status": {
			Before: fmt.Sprint(previousData.Status), // Before value for Description
			After:  fmt.Sprint(data.Status),         // After value for Description
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
