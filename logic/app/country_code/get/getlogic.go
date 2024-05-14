package get

import (
	"context"

	"nrs_customer_module_backend/internal/global/responseglobal"
	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/gs_country_code"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogic {
	return &GetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLogic) Get(req *types.PostCountryCodeGetReq) (resp *types.SuccessResponse, err error) {
	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	countryCodeModel := gs_country_code.NewGsCountryCodeModel(conn)
	countryCode, err := countryCodeModel.FindAll(l.ctx)
	if err != nil {
		return nil, err
	}
	var content []gs_country_code.CountryCodeContent
	for _, cc := range *countryCode {
		content = append(content, gs_country_code.CountryCodeContent{
			Guid:        cc.Guid,
			ImageUrl:    cc.ImageUrl,
			WebImageUrl: cc.WebImageUrl,
			Title:       cc.Name,
			Code:        cc.Code,
		})
	}
	dataResp := map[string]interface{}{
		"list": content,
	}
	return responseglobal.GenerateResponseBody(true, "Successfully retrieved records.", dataResp), nil

}
