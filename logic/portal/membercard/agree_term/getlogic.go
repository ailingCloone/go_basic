package agree_term

import (
	"context"
	"database/sql"

	"nrs_customer_module_backend/internal/global/responseglobal"
	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/agree_term"
	"nrs_customer_module_backend/internal/model/card"
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

func (l *GetLogic) Get(req *types.PostATGetReq) (resp *types.SuccessResponse, err error) {
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

	agreeTermModel := agree_term.NewAgreeTermModel(conn)
	agreeTerm, err := agreeTermModel.FindOneByCard(l.ctx, cardResult.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return responseglobal.GenerateResponseBody(false, "Agree Term is not found.", map[string]interface{}{}), nil
		}
		return nil, err
	}

	if agreeTerm == nil {
		return responseglobal.GenerateResponseBody(false, "Failed retrieve record.", map[string]interface{}{}), nil
	}

	data := map[string]interface{}{
		"guid":    agreeTerm.Guid,
		"title":   agreeTerm.Title,
		"content": agreeTerm.Description,
	}
	return responseglobal.GenerateResponseBody(true, "Success retrieve record.", data), nil
}
