package otp

import (
	"context"
	"fmt"
	"strings"

	"nrs_customer_module_backend/internal/config"
	"nrs_customer_module_backend/internal/global"
	"nrs_customer_module_backend/internal/global/createfile"
	"nrs_customer_module_backend/internal/global/errorglobal"
	otpFunc "nrs_customer_module_backend/internal/global/otp"
	"nrs_customer_module_backend/internal/global/responseglobal"

	"nrs_customer_module_backend/internal/model/sms_template"

	"nrs_customer_module_backend/internal/global/user"

	"nrs_customer_module_backend/internal/model"
	"nrs_customer_module_backend/internal/svc"
	"nrs_customer_module_backend/internal/types"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type OtpLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOtpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtpLogic {
	return &OtpLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OtpLogic) Otp(req *types.OtpReq) (resp *types.SuccessResponse, err error) {
	logFile := "otp_request_sms"
	appLogger := createfile.New(logFile)
	//20240417 focus for login only
	var c config.Config
	configFile := "etc/api.yaml"
	if err := conf.LoadConfig(configFile, &c); err != nil {
		return nil, err
	}

	conn, err := model.InitializeDatabase()
	if err != nil {
		return nil, err
	}
	var contactNo string
	var sendStatus string
	var from int64
	var referName string //sms template [refer table name]
	var categoryId int64 //for sms template [refer table = category;refer id ]

	if req.Page == "login" {
		contactNo, from, err = LoginPage(req, l.ctx, conn)
		if err != nil {
			return nil, err
		}
		referName = "category"
		categoryId = 33
	}

	if req.Page == "register" {
		contactNo, from, err = RegisterPage(req, l.ctx, conn)
		if err != nil {
			return nil, err
		}
		referName = "category"
		categoryId = 33
	}

	var auth int64 = 1 //1- SMS, 2- EMAIL
	smsModel := sms_template.NewSmsTemplateModel(conn)
	smsTemplate, err := smsModel.FindOneReferId(l.ctx, referName, categoryId)
	if err != nil {
		return responseglobal.GenerateResponseBody(false, "Failed to get SMS Template", map[string]interface{}{}), err
	}

	smsContent := smsTemplate.Description

	responseTemplate := smsTemplate.Title
	responseMessage := strings.Replace(responseTemplate, "{{contact}}", contactNo, -1)

	otpRecord, sendStatus, err := otpFunc.GenerateSMSOtpRecord(l.ctx, conn, contactNo, auth, from, c.SMSOtpExpire, appLogger, logFile, smsContent, responseMessage)
	if err != nil {
		return nil, err
	}

	phoneNo := global.MaskPhoneNumber(contactNo)
	uiOtpVerification := map[string]interface{}{
		"info_verification": phoneNo,
		"count_down":        "60",
	}

	sendTimeStr := global.FormatMDYTime(otpRecord.SendTime)
	sendTimePtr := &sendTimeStr

	phoneDelivery := &types.PhoneDelivery{
		Contact:    &otpRecord.Value,
		SendStatus: &sendStatus,
		SendTime:   sendTimePtr,
	}
	messageTemplate := " The OTP has been sent to {{contact}}. Please enter the OTP you received to Validate."

	message := strings.Replace(messageTemplate, "{{contact}}", phoneNo, -1)

	content := types.OtpDetailsContent{
		AuthType:      req.Type,
		EmailDelivery: &types.EmailDelivery{}, // sms otp request -> no email
		Message:       message,
		PhoneDelivery: phoneDelivery,
		TxID:          otpRecord.Guid,
	}

	dataResp := map[string]interface{}{
		"ui_otp_verification": uiOtpVerification,
		"content":             content,
	}

	return responseglobal.GenerateResponseBody(true, "Successfully request otp.", dataResp), nil
}

func LoginPage(req *types.OtpReq, ctx context.Context, conn sqlx.SqlConn) (contactNo string, from int64, err error) {
	typeVal := strings.ToUpper(req.Type)

	if req.From == "staff" {
		exists, err := user.CheckStaff(ctx, conn, typeVal, req.Contact)

		if err != nil {
			return "", 0, err
		}
		contactNo = exists.Contact

	}

	if req.From == "customer" {
		exists, err := user.CheckCustomer(ctx, conn, typeVal, req.Contact)
		if err != nil {
			return "", 0, err
		}

		contactNo = exists.Contact
	}

	if req.Type == "CONTACT" {
		from = 2 //1- Forget Password, 2- Login Contact ,3- Login IC,  4- Register, 5- Profile Update
	}

	if req.Type == "IC" {
		from = 3 //1- Forget Password, 2- Login Contact ,3- Login IC,  4- Register, 5- Profile Update
	}
	return contactNo, from, nil
}

func RegisterPage(req *types.OtpReq, ctx context.Context, conn sqlx.SqlConn) (contactNo string, from int64, err error) {
	typeVal := strings.ToUpper(req.Type)

	if req.From == "staff" {
		exists, err := user.CheckStaff(ctx, conn, typeVal, req.Contact)
		if exists != nil && err == nil {
			err = fmt.Errorf(errorglobal.ExistedUser)
			return "", 0, err
		}
	}

	if req.From == "customer" {
		exists, err := user.CheckCustomer(ctx, conn, typeVal, req.Contact)
		if exists != nil && err == nil {
			err = fmt.Errorf(errorglobal.ExistedUser)
			return "", 0, err
		}
	}
	contactNo = req.Contact

	from = 4 //1- Forget Password, 2- Login Contact ,3- Login IC,  4- Register, 5- Profile Update
	return contactNo, from, nil
}
