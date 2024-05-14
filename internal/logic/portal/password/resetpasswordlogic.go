package password

import (
	"context"
	"fmt"

	"nrs_customer_module_backend/internal/global"
	"nrs_customer_module_backend/internal/global/errorglobal"
	otpFunc "nrs_customer_module_backend/internal/global/otp"
	"nrs_customer_module_backend/internal/global/responseglobal"
	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/staff"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPasswordLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	otherInfo map[string]interface{}
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext, otherInfo map[string]interface{}) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		otherInfo: otherInfo,
	}
}

func (l *ResetPasswordLogic) ResetPassword(req *types.PostPortalResetPasswordReq) (resp *types.SuccessResponse, err error) {
	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	currentTime, err := global.TimeInSingapore()
	if err != nil {
		return nil, err
	}

	//check is otp valid to use
	otpList, err := otpFunc.OtpExpirationChecking(l.ctx, conn, req.TxId, currentTime)
	if err != nil {
		return nil, err
	}

	if req.Code != fmt.Sprint(otpList.Code) {
		err = fmt.Errorf(errorglobal.InvalidOtp)
		return nil, err
	}

	//update otp active status if verified otp
	otpList.Active = 2
	err = otpFunc.UpdateOTP(l.ctx, conn, otpList)
	if err != nil {
		return nil, err
	}

	checkPass := global.PasswordValidation(req.NewPassword, req.ConfirmPassword)

	if !checkPass {
		err = fmt.Errorf(errorglobal.InvalidPassword)
		return nil, err
	}

	hashPass := global.GenerateSha256(req.ConfirmPassword)

	data := &staff.Staff{
		Email:    otpList.Value,
		Password: hashPass,
		Updated:  *currentTime,
	}

	staffModel := staff.NewStaffModel(conn)
	if err := staffModel.UpdatePassword(l.ctx, data); err != nil {
		return responseglobal.GenerateResponseBody(false, "Failed Reset Password.", nil), nil
	}

	return responseglobal.GenerateResponseBody(true, "Successfully Reset Password.", nil), nil

}
