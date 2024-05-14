package password

import (
	"context"
	"fmt"
	"time"

	"nrs_customer_module_backend/internal/global"
	"nrs_customer_module_backend/internal/global/errorglobal"
	otpFunc "nrs_customer_module_backend/internal/global/otp"
	"nrs_customer_module_backend/internal/global/responseglobal"
	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/model/customer"
	otps "nrs_customer_module_backend/internal/model/otp"
	"nrs_customer_module_backend/internal/model/staff"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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

func (l *ResetPasswordLogic) ResetPassword(req *types.PostResetPasswordAppReq) (resp *types.SuccessResponse, err error) {
	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	currentTime, err := global.TimeInSingapore()
	if err != nil {
		return nil, err
	}

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

	resp, err = UpdateUserPassword(l.ctx, conn, req.From, otpList, hashPass, currentTime)
	if err != nil {
		return resp, nil
	}

	return responseglobal.GenerateResponseBody(true, "Successfully Reset Password.", map[string]interface{}{}), nil
}

func UpdateUserPassword(ctx context.Context, conn sqlx.SqlConn, from string, otpList *otps.Otp, hashPass string, currentTime *time.Time) (resp *types.SuccessResponse, err error) {
	if from == "staff" {
		data := &staff.Staff{
			Email:    otpList.Value,
			Password: hashPass,
			Updated:  *currentTime,
		}

		staffModel := staff.NewStaffModel(conn)
		if err := staffModel.UpdatePassword(ctx, data); err != nil {
			return responseglobal.GenerateResponseBody(false, "Failed Reset Password.", map[string]interface{}{}), nil
		}
	}

	if from == "customer" {
		data := &customer.Customer{
			Email:    otpList.Value,
			Password: hashPass,
			Updated:  *currentTime,
		}

		customerModel := customer.NewCustomerModel(conn)
		if err := customerModel.UpdatePassword(ctx, data); err != nil {
			return responseglobal.GenerateResponseBody(false, "Failed Reset Password.", map[string]interface{}{}), nil
		}
	}
	return resp, nil
}
