package otp

import (
	"context"

	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtpVerifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOtpVerifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtpVerifyLogic {
	return &OtpVerifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OtpVerifyLogic) OtpVerify(req *types.OtpVerifyReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
