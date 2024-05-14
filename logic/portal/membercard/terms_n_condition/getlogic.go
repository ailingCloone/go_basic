package terms_n_condition

import (
	"context"

	"nrs_customer_module_backend/internal/global/responseglobal"
	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/card"
	"nrs_customer_module_backend/internal/model/tnc"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	otherInfo map[string]interface{}
}

func NewGetLogic(ctx context.Context, svcCtx *svc.ServiceContext, otherInfo map[string]interface{}) *GetLogic {
	return &GetLogic{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		otherInfo: otherInfo,
	}
}

func (l *GetLogic) Get(req *types.PostTNCGetReq) (resp *types.SuccessResponse, err error) {
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

	respSlice, err := tncModel.FindOneByCard(l.ctx, cardResult.Id)
	if err != nil {
		return nil, err
	}

	if respSlice == nil {
		return responseglobal.GenerateResponseBody(false, "Failed to retrieve record.", map[string]interface{}{}), nil
	}

	response := map[string]interface{}{
		"guid":    respSlice.Guid,
		"title":   respSlice.Title,
		"content": respSlice.Description,
	}

	return responseglobal.GenerateResponseBody(true, "Successfully retrieve record.", response), nil
}
