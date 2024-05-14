package get_list

import (
	"context"

	"nrs_customer_module_backend/internal/config"
	"nrs_customer_module_backend/internal/global"
	"nrs_customer_module_backend/internal/global/memberinfo"
	"nrs_customer_module_backend/internal/global/responseglobal"

	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/homepage"
	"nrs_customer_module_backend/internal/model/request_status"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/core/conf"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListLogic {
	return &GetListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetListLogic) GetList(req *types.GetListRegistrationReq) (resp *types.SuccessResponse, err error) {
	/*
	   11 - Register
	   12 - Renew
	   13 - Upgrade
	*/
	var from string = "11"

	// get the filter day option
	page := "app_registration_list_get_list"
	filterDay := global.FilterDay(page)

	// load api.yaml
	var c config.Config
	configFile := "etc/api.yaml"
	if err := conf.LoadConfig(configFile, &c); err != nil {
		return nil, err
	}

	// open db connection
	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	// connect to request_status
	requestStatusModel := request_status.NewRequestStatusModel(conn)

	// retrieve record
	list, err := requestStatusModel.FindMemberRequestList(l.ctx, from, req)
	if err != nil {
		return nil, err
	}

	// massage to wanted structure
	record := memberinfo.AppContent(list, page)

	// app bar
	homepageModel := homepage.NewHomepageModel(conn)
	homepageData, err := homepageModel.FindOneByCategoryId(l.ctx, 92)
	if err != nil {
		return nil, err
	}

	appBar := map[string]interface{}{
		"title": homepageData.Title,
	}

	response := map[string]interface{}{
		"filter":  filterDay,
		"tab":     record,
		"app_bar": appBar,
	}
	// Prepare the response
	return responseglobal.GenerateResponseBody(true, "Successfully retrieve record.", response), nil

}
