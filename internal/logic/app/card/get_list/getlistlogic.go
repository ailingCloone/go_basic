package get_list

import (
	"context"
	"fmt"

	"nrs_customer_module_backend/internal/global/responseglobal"
	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/card"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

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

func (l *GetListLogic) GetList(req *types.PostCardGetListReq) (resp *types.SuccessResponse, err error) {
	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	cardModel := card.NewCardModel(conn)
	cardList, err := cardModel.FindAllCardList(l.ctx)
	if err != nil {
		return nil, err
	}

	var list []card.CardGetList
	for _, cList := range *cardList {
		list = append(list, card.CardGetList{
			Guid:        cList.Guid,
			ImageUrl:    cList.ImageUrl,
			WebImageUrl: cList.WebImageUrl,
			Title:       cList.Title,
			Code:        cList.Code,
			NeedPayment: cList.Payment,
			Price:       fmt.Sprint(cList.Price),
		})
	}
	dataResp := map[string]interface{}{
		"list": list,
	}

	return responseglobal.GenerateResponseBody(true, "Successfully retrieve record.", dataResp), nil
}
