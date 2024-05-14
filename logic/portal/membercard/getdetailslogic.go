package membercard

import (
	"context"

	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDetailsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDetailsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDetailsLogic {
	return &GetDetailsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDetailsLogic) GetDetails(req *types.GetCardDetailsReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
