package get_list

import (
	"context"

	"nrs_customer_module_backend/internal/config"
	"nrs_customer_module_backend/internal/global/memberinfo"
	"nrs_customer_module_backend/internal/global/responseglobal"

	"nrs_customer_module_backend/internal/model"
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
	// todo: add your logic here and delete this line
	listHeader := []types.ListHeader{
		{"Customer Name", "name"},
		{"Email", "email"},
		{"Type of Member Apply", "card"},
		{"Status", "status"},
		{"Action", "action"},
	}

	var data []memberinfo.DataList = []memberinfo.DataList{}

	var c config.Config
	configFile := "etc/api.yaml"
	if err := conf.LoadConfig(configFile, &c); err != nil {
		return nil, err
	}
	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	requestStatusModel := request_status.NewRequestStatusModel(conn)
	/*
		11 - Register
		12 - Renew
		13 - Upgrade
	*/
	from := "11"
	list, err := requestStatusModel.FindMemberRequestList(l.ctx, from, req)
	if err != nil {
		return nil, err
	}

	for _, v := range *list {
		record := memberinfo.DataList{
			Guid:  v.Guid,
			Name:  v.CusFullname,
			Email: v.CusEmail,
			Card:  v.CusCardTitle,
		}

		cusCardCode := v.CusCardCode
		status := v.Status

		record = memberinfo.PortalStatusAction(cusCardCode, status, record)

		data = append(data, record)
	}

	response := map[string]interface{}{
		"list":        data,
		"list_header": listHeader,
	}

	// Prepare the response
	return responseglobal.GenerateResponseBody(true, "Successfully retrieve record.", response), nil

}
