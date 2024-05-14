package terms_n_condition

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
	"nrs_customer_module_backend/internal/model/card"
	"nrs_customer_module_backend/internal/model/tnc"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	otherInfo map[string]interface{}
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext, otherInfo map[string]interface{}) *AddLogic {
	return &AddLogic{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		otherInfo: otherInfo,
	}
}

func (l *AddLogic) Add(req *types.PostTNCAddReq) (resp *types.SuccessResponse, err error) {
	oauthData, err := user.GetOauthDataFromAfterLoginMiddleware(l.ctx)
	if err != nil {
		return nil, err
	}

	// Create an instance of TncModel
	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	cardModel := card.NewCardModel(conn)
	// Get card data based on guid
	cardResult, err := cardModel.FindOneGuid(context.Background(), req.Guid)
	if err != nil {
		return nil, err
	}

	tncModel := tnc.NewTncModel(conn)
	currentTime, err := global.TimeInSingapore()
	if err != nil {
		return nil, err
	}

	// Prepare data for insertion
	data := &tnc.Tnc{
		Guid:        global.GenerateGuid(),
		ReferTable:  "card",
		ReferId:     cardResult.Id,
		Title:       req.Title,
		Description: req.Content,
		Created:     *currentTime,
		Active:      1,
	}

	// Insert data into the database
	result, err := tncModel.Insert(context.Background(), data)

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

	dataResp := map[string]interface{}{
		"guid":    data.Guid,
		"title":   req.Title,
		"content": req.Content,
	}
	changes, err := ConvertChangesActivityLog(data)
	if err != nil {
		return nil, err
	}

	dataActivityLog := &activity_log.ActivityLog{
		StaffId:    oauthData.StaffId.Int64,
		ReferTable: "tnc",
		Action:     "add",
		Changes:    string(changes),
		Created:    fmt.Sprint(*currentTime),
	}

	err = activitylog.AddActivityLog(dataActivityLog, conn)

	if err != nil {
		return nil, err
	}

	return responseglobal.GenerateResponseBody(true, "Successfully add record.", dataResp), nil
}

func ConvertChangesActivityLog(data *tnc.Tnc) (jsonData []byte, err error) {

	// Create a map to hold the changes for each field
	changes := map[string]activity_log.Change{
		"guid": {
			Before: "",        // Before value for Guid
			After:  data.Guid, // After value for Guid
		},
		"refer_table": {
			Before: "",              // Before value for ReferTable
			After:  data.ReferTable, // After value for ReferTable
		},
		"refer_id": {
			Before: "",                       // Before value for ReferId
			After:  fmt.Sprint(data.ReferId), // After value for ReferId
		},
		"title": {
			Before: "",         // Before value for Title
			After:  data.Title, // After value for Title
		},
		"description": {
			Before: "",               // Before value for Description
			After:  data.Description, // After value for Description
		},
		"created": {
			Before: "",                       // Before value for Created
			After:  fmt.Sprint(data.Created), // After value for Created
		},
		"active": {
			Before: "",                      // Before value for Active
			After:  fmt.Sprint(data.Active), // After value for Active
		},
	}

	// Marshal the ChangeData object into JSON format
	jsonData, err = json.Marshal(changes)
	if err != nil {
		return nil, err
	}

	return jsonData, nil

}
