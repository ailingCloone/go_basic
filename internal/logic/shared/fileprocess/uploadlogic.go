package fileprocess

import (
	"context"

	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadLogic) Upload(req *types.PostFileUploadReq) (resp *types.Response, err error) {

	// todo: add your logic here and delete this line
	return resp, err
}
