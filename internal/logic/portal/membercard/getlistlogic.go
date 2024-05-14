package membercard

import (
	"context"
	"fmt"

	"nrs_customer_module_backend/internal/global/responseglobal"
	"nrs_customer_module_backend/internal/global/user"
	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/card"
	"nrs_customer_module_backend/internal/model/user_access_setting"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	otherInfo map[string]interface{}
}

func NewGetListLogic(ctx context.Context, svcCtx *svc.ServiceContext, otherInfo map[string]interface{}) *GetListLogic {
	return &GetListLogic{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		otherInfo: otherInfo,
	}
}

func (l *GetListLogic) GetList(req *types.GetCardListReq) (resp *types.SuccessResponse, err error) {
	userPermission, err := user.GetUserPermission(l.ctx)
	if err != nil {
		return nil, err
	}

	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	cardModel := card.NewCardModel(conn)

	cardInfo, err := cardModel.FindAllCard(l.ctx)
	if err != nil {
		return nil, err
	}

	var cardSummary []card.CardInfo
	for _, info := range *cardInfo {
		statusCount, err := cardModel.FindAllMember(l.ctx, int64(info.Id))
		if err != nil {
			return nil, err
		}
		cardSummary = append(cardSummary, card.CardInfo{
			Guid:        info.Guid,
			Title:       info.Title,
			Value:       fmt.Sprint(statusCount),
			Description: "Total Members",
			ImageUrl:    info.ImageUrl,
			WebImageUrl: info.WebImageUrl,
		})
	}

	var permissionSummaries []user_access_setting.Permission
	for _, permission := range *userPermission {
		summary := user_access_setting.Permission{
			Title: permission.Title,
			Allow: user_access_setting.Allow{
				Edit:   permission.AllowEdit,
				Delete: permission.AllowDelete,
				Add:    permission.AllowAdd,
				View:   permission.AllowView,
			},
		}
		permissionSummaries = append(permissionSummaries, summary)
	}

	response := map[string]interface{}{
		"list":       cardSummary,
		"permission": permissionSummaries,
	}

	// Prepare the response
	return responseglobal.GenerateResponseBody(true, "Successfully retrieve record.", response), nil
}
